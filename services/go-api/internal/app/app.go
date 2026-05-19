package app

import (
"net/http"

"github.com/orcs-toolkit/orcs/services/go-api/internal/config"
"github.com/orcs-toolkit/orcs/services/go-api/internal/http/handlers"
"github.com/orcs-toolkit/orcs/services/go-api/internal/http/middleware"
"github.com/orcs-toolkit/orcs/services/go-api/internal/shadow"
"github.com/orcs-toolkit/orcs/services/go-api/internal/store"
"github.com/orcs-toolkit/orcs/services/go-api/internal/store/memory"
)

type App struct {
Config config.Config
Router http.Handler
}

func New() (*App, error) {
cfg := config.Load()
var db store.Store = memory.New()
sh := shadow.NewRESTShadow(cfg.ShadowEnabled, cfg.LegacyAPIBase)
h := handlers.New(db, cfg.SecretKey, sh)
router := h.Router()

return &App{Config: cfg, Router: middleware.WithCORS(router)}, nil
}
