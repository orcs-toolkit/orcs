package realtime

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"github.com/orcs-toolkit/orcs/services/go-master/internal/config"
"github.com/orcs-toolkit/orcs/services/go-master/internal/contracts"
)

func TestEventContracts(t *testing.T) {
h := NewHub(config.Config{AuthSecret: "admin"})
server := httptest.NewServer(http.HandlerFunc(h.HandleWS))
defer server.Close()

payload, _ := json.Marshal(contracts.Envelope{Event: contracts.ClientAuth, Data: map[string]any{"key": "admin"}})
resp, err := http.Post(server.URL, "application/json", bytes.NewReader(payload))
if err != nil || resp.StatusCode != http.StatusOK {
t.Fatalf("clientAuth failed: %v status=%d", err, resp.StatusCode)
}

payload, _ = json.Marshal(contracts.Envelope{Event: contracts.NodeLogs, Data: map[string]any{"type": "connected"}})
resp, err = http.Post(server.URL, "application/json", bytes.NewReader(payload))
if err != nil || resp.StatusCode != http.StatusOK {
t.Fatalf("node logs failed: %v status=%d", err, resp.StatusCode)
}
}
