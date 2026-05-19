package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orcs-toolkit/orcs/services/go-api/internal/platform"
	"github.com/orcs-toolkit/orcs/services/go-api/internal/shadow"
	"github.com/orcs-toolkit/orcs/services/go-api/internal/store/memory"
)

func TestAuthFlowAndPolicyContract(t *testing.T) {
	h := New(memory.New(), "test-secret", shadow.NewRESTShadow(false, ""))
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	register := map[string]any{"name": "Admin", "email": "admin@example.com", "password": "secret12", "isAdmin": true, "role": "admin"}
	body, _ := json.Marshal(register)
	resp, err := http.Post(srv.URL+"/auth/register", "application/json", bytes.NewReader(body))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("register failed: %v status=%d", err, resp.StatusCode)
	}

	token, _ := platform.SignJWT(map[string]any{"role": "admin"}, "test-secret")
	policy := []byte(`{"role":"student","banList":["chrome"]}`)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/policy/setPolicy", bytes.NewReader(policy))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("set policy failed: %v status=%d", err, resp.StatusCode)
	}

	resp, err = http.Get(srv.URL + "/policy/getRolePolicy/student")
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("get role policy failed: %v status=%d", err, resp.StatusCode)
	}
}
