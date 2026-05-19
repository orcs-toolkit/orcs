package realtime

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/orcs-toolkit/orcs/services/go-master/internal/config"
	"github.com/orcs-toolkit/orcs/services/go-master/internal/contracts"
)

type Hub struct {
	cfg      config.Config
	mu       sync.RWMutex
	admins   map[*client]struct{}
	machines map[*client]string
}

type client struct {
	out chan contracts.Envelope
}

func NewHub(cfg config.Config) *Hub {
	return &Hub{cfg: cfg, admins: map[*client]struct{}{}, machines: map[*client]string{}}
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var env contracts.Envelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch env.Event {
	case contracts.ClientAuth:
		payload, _ := env.Data.(map[string]any)
		key, _ := payload["key"].(string)
		if key != h.cfg.AuthSecret {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, contracts.Envelope{Event: contracts.Logs, Data: map[string]any{"type": "connected", "data": "admin joined"}})
	case contracts.PerfData:
		h.broadcastAdmins(contracts.Envelope{Event: contracts.Data, Data: env.Data})
		writeJSON(w, http.StatusOK, contracts.Envelope{Event: contracts.Data, Data: env.Data})
	case contracts.UpdatedBan:
		writeJSON(w, http.StatusOK, contracts.Envelope{Event: contracts.UpdatedBan, Data: env.Data})
	case contracts.NodeLogs:
		h.broadcastAdmins(contracts.Envelope{Event: contracts.Logs, Data: env.Data})
		writeJSON(w, http.StatusOK, contracts.Envelope{Event: contracts.Logs, Data: env.Data})
	default:
		log.Printf("unsupported event: %s", env.Event)
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *Hub) broadcastAdmins(env contracts.Envelope) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.admins {
		select {
		case c.out <- env:
		default:
		}
	}
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
