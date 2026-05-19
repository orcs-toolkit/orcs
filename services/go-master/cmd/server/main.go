package main

import (
"log"
"net/http"

"github.com/orcs-toolkit/orcs/services/go-master/internal/config"
"github.com/orcs-toolkit/orcs/services/go-master/internal/realtime"
)

func main() {
cfg := config.Load()
hub := realtime.NewHub(cfg)
mux := http.NewServeMux()
mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`{"ok":true}`)) })
mux.HandleFunc("/ws", hub.HandleWS)

log.Printf("go-master listening on %s", cfg.ListenAddr)
if err := http.ListenAndServe(cfg.ListenAddr, mux); err != nil {
log.Fatalf("go-master failed: %v", err)
}
}
