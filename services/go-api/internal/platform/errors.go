package platform

import (
"encoding/json"
"net/http"
)

type HTTPError struct {
Code    int
Message string
Err     error
}

func (h HTTPError) Error() string {
if h.Err != nil {
return h.Err.Error()
}
return h.Message
}

func WriteError(w http.ResponseWriter, code int, payload map[string]any) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(code)
_ = json.NewEncoder(w).Encode(payload)
}
