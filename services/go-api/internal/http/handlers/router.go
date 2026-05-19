package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/orcs-toolkit/orcs/services/go-api/internal/domain"
	"github.com/orcs-toolkit/orcs/services/go-api/internal/http/middleware"
	"github.com/orcs-toolkit/orcs/services/go-api/internal/platform"
	"github.com/orcs-toolkit/orcs/services/go-api/internal/shadow"
	"github.com/orcs-toolkit/orcs/services/go-api/internal/store"
)

type Handler struct {
	store  store.Store
	secret string
	shadow *shadow.RESTShadow
}

func New(s store.Store, secret string, shadowClient *shadow.RESTShadow) *Handler {
	return &Handler{store: s, secret: secret, shadow: shadowClient}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.health)

	mux.HandleFunc("/auth/register", h.signup)
	mux.HandleFunc("/auth/login", h.signin)
	mux.HandleFunc("/auth/currentuser", h.currentUser)
	mux.HandleFunc("/auth/logout", h.signout)

	mux.HandleFunc("/user/getAllUsers", h.getAllUsers)
	mux.HandleFunc("/user/getUser/", h.getUserByID)
	mux.HandleFunc("/user/updateUser/", h.requireAdmin(h.updateUser))
	mux.HandleFunc("/user/deleteUser/", h.requireAdmin(h.deleteUser))

	mux.HandleFunc("/policy/getPolicies", h.getPolicies)
	mux.HandleFunc("/policy/getPolicy/", h.getPolicyByID)
	mux.HandleFunc("/policy/getRolePolicy/", h.getRolePolicy)
	mux.HandleFunc("/policy/getRoleWisePolicy", h.requireAdmin(h.getRoleWisePolicy))
	mux.HandleFunc("/policy/getFavoriteProcesses", h.requireAdmin(h.getFavoriteProcesses))
	mux.HandleFunc("/policy/setPolicy", h.requireAdmin(h.setPolicy))
	mux.HandleFunc("/policy/updatePolicy", h.requireAdmin(h.updatePolicy))
	mux.HandleFunc("/policy/updateSinglePolicy", h.requireAdmin(h.updateSinglePolicy))
	mux.HandleFunc("/policy/addSingleProcess", h.requireAdmin(h.addSingleProcess))
	mux.HandleFunc("/policy/deletePolicy/", h.requireAdmin(h.deletePolicy))

	mux.HandleFunc("/logs", h.getLogs)
	return mux
}

func (h *Handler) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	wrapped := middleware.RequireLogin(h.secret, true)(next)
	return func(w http.ResponseWriter, r *http.Request) {
		wrapped.ServeHTTP(w, r)
	}
}

func decodeJSON(r *http.Request, dst any) error { return json.NewDecoder(r.Body).Decode(dst) }

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func pathParam(r *http.Request) string {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) signout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "Successfully logged out.", "error": nil})
}

func (h *Handler) currentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 {
		writeJSON(w, http.StatusOK, map[string]any{"currentUser": nil})
		return
	}
	payload, err := platform.VerifyJWT(parts[1], h.secret)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"currentUser": nil})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"currentUser": payload})
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
		Role     string `json:"role"`
	}
	if err := decodeJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"message": "Invalid request", "error": err.Error()})
		return
	}
	if !strings.Contains(req.Email, "@") || len(req.Password) < 6 || len(req.Password) > 24 || strings.TrimSpace(req.Name) == "" {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"message": "Validation failed"})
		return
	}
	hashed, err := platform.HashPassword(req.Password)
	if err != nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"success": false, "message": err.Error()})
		return
	}
	u, err := h.store.CreateUser(r.Context(), domain.User{Name: req.Name, Email: req.Email, Password: hashed, IsAdmin: req.IsAdmin, Role: strings.ToLower(req.Role)})
	if err != nil {
		platform.WriteError(w, http.StatusUnauthorized, map[string]any{"success": false, "message": "Email already registered!"})
		return
	}
	token, _ := platform.SignJWT(map[string]any{"id": u.ID, "email": u.Email, "name": u.Name, "admin": u.IsAdmin, "role": u.Role}, h.secret)
	h.shadow.Fire(http.MethodPost, "/auth/register", req)
	writeJSON(w, http.StatusCreated, map[string]any{"user": map[string]any{"id": u.ID, "email": u.Email, "name": u.Name, "admin": u.IsAdmin, "role": u.Role}, "success": true, "token": token})
}

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"success": false, "message": "Invalid request"})
		return
	}
	u, _ := h.store.FindUserByEmail(r.Context(), req.Email)
	if u == nil || !platform.ComparePassword(u.Password, req.Password) {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"success": false, "message": "Invalid Credentials"})
		return
	}
	token, _ := platform.SignJWT(map[string]any{"id": u.ID, "email": u.Email, "name": u.Name, "isAdmin": u.IsAdmin, "role": u.Role}, h.secret)
	h.shadow.Fire(http.MethodPost, "/auth/login", req)
	writeJSON(w, http.StatusOK, map[string]any{"user": map[string]any{"id": u.ID, "email": u.Email, "name": u.Name, "admin": u.IsAdmin, "role": u.Role}, "success": true, "token": token})
}

func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	users, _ := h.store.GetUsers(r.Context())
	writeJSON(w, http.StatusOK, users)
}

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	u, _ := h.store.FindUserByID(r.Context(), pathParam(r))
	if u == nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"message": "User not found!"})
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
		Role     string `json:"role"`
	}
	if err := decodeJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"message": "User update failed!"})
		return
	}
	id := pathParam(r)
	pwd := ""
	if strings.TrimSpace(req.Password) != "" {
		hashed, _ := platform.HashPassword(req.Password)
		pwd = hashed
	}
	u, _ := h.store.UpdateUser(r.Context(), domain.User{ID: id, Name: req.Name, Email: req.Email, Password: pwd, IsAdmin: req.IsAdmin, Role: req.Role})
	if u == nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"message": "User not found!"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "User updated successfully", "user": u})
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	u, _ := h.store.DeleteUser(r.Context(), pathParam(r))
	if u == nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"message": "User not found!"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Successfully deleted user", "user": u})
}

func (h *Handler) setPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Role    string   `json:"role"`
		BanList []string `json:"banList"`
	}
	if err := decodeJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"message": "Failed to save Policy"})
		return
	}
	p, _ := h.store.CreatePolicy(r.Context(), domain.Policy{Role: req.Role, BanList: req.BanList})
	writeJSON(w, http.StatusCreated, map[string]any{"role": p.Role, "banList": p.BanList, "createdAt": p.CreatedAt})
}

func (h *Handler) getPolicies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	policies, _ := h.store.GetPolicies(r.Context())
	writeJSON(w, http.StatusOK, policies)
}

func (h *Handler) getPolicyByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p, _ := h.store.GetPolicy(r.Context(), pathParam(r))
	if p == nil {
		platform.WriteError(w, http.StatusUnauthorized, map[string]any{"message": "Policy not found!"})
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) getRolePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p, _ := h.store.GetPolicyByRole(r.Context(), pathParam(r))
	if p == nil {
		platform.WriteError(w, http.StatusUnauthorized, map[string]any{"message": "Policy not found!"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"policy": p})
}

func (h *Handler) updatePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Role    string   `json:"role"`
		BanList []string `json:"banList"`
	}
	if err := decodeJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"message": "Policy not found!"})
		return
	}
	p, _ := h.store.UpdatePolicyBanList(r.Context(), req.Role, req.BanList)
	if p == nil {
		platform.WriteError(w, http.StatusUnauthorized, map[string]any{"message": "Policy not found!"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"message": "Successfully updated policy", "policy": p})
}

func (h *Handler) updateSinglePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Role        string `json:"role"`
		ProcessName string `json:"processName"`
	}
	if err := decodeJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusInternalServerError, map[string]any{"message": "Couldn't update data"})
		return
	}
	_ = h.store.UpdateSinglePolicy(r.Context(), req.Role, req.ProcessName)
	writeJSON(w, http.StatusOK, map[string]any{"acknowledged": true})
}

func (h *Handler) addSingleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Role          []string `json:"role"`
		ProcessName   string   `json:"processName"`
		AddToFavorite bool     `json:"addToFavorite"`
	}
	if err := decodeJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, map[string]any{"message": "Policy not found!"})
		return
	}
	all := len(req.Role) > 0 && strings.EqualFold(req.Role[0], "all")
	_ = h.store.AddSingleProcess(r.Context(), req.Role, req.ProcessName, all, req.AddToFavorite)
	writeJSON(w, http.StatusCreated, map[string]any{"message": "Successfully updated policy", "policy": map[string]any{"acknowledged": true}})
}

func (h *Handler) getFavoriteProcesses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, _ := h.store.GetFavoriteProcesses(r.Context())
	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) getRoleWisePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, _ := h.store.GetRoleWisePolicy(r.Context())
	writeJSON(w, http.StatusOK, data)
}

func (h *Handler) deletePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p, _ := h.store.DeletePolicy(r.Context(), pathParam(r))
	if p == nil {
		platform.WriteError(w, http.StatusUnauthorized, map[string]any{"message": "Policy not found!"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"message": "Successfully deleted Policy", "policy": p})
}

func (h *Handler) getLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit > 500 {
		limit = 500
	}
	logs, _ := h.store.GetLogs(r.Context(), limit)
	writeJSON(w, http.StatusOK, logs)
}
