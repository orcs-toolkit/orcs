package contracts

const (
ClientAuth   = "clientAuth"
Admin        = "admin"
Logout       = "logout"
Disconnect   = "disconnect"
InitPerfData = "initPerfData"
PerfData     = "perfData"
UpdatedBan   = "updated:Ban"
NodeLogs     = "node:logs"
Logs         = "logs"
Data         = "data"
)

type Envelope struct {
Event string `json:"event"`
Data  any    `json:"data"`
}
