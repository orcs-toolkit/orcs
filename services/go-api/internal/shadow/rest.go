package shadow

import (
"bytes"
"encoding/json"
"log"
"net/http"
"time"
)

type RESTShadow struct {
Enabled bool
BaseURL string
Client  *http.Client
}

func NewRESTShadow(enabled bool, baseURL string) *RESTShadow {
return &RESTShadow{Enabled: enabled, BaseURL: baseURL, Client: &http.Client{Timeout: 2 * time.Second}}
}

func (s *RESTShadow) Fire(method, path string, body any) {
if !s.Enabled {
return
}
go func() {
payload, _ := json.Marshal(body)
req, err := http.NewRequest(method, s.BaseURL+path, bytes.NewReader(payload))
if err != nil {
log.Printf("shadow request build failed: %v", err)
return
}
req.Header.Set("Content-Type", "application/json")
resp, err := s.Client.Do(req)
if err != nil {
log.Printf("shadow request failed: %v", err)
return
}
_ = resp.Body.Close()
}()
}
