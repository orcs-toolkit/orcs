package config

import "os"

type Config struct {
ListenAddr    string
SecretKey     string
ShadowEnabled bool
LegacyAPIBase string
}

func Load() Config {
addr := os.Getenv("GO_API_ADDR")
if addr == "" {
addr = ":4001"
}

secret := os.Getenv("SECRET_KEY")
if secret == "" {
secret = "asdf"
}

return Config{
ListenAddr:    addr,
SecretKey:     secret,
ShadowEnabled: os.Getenv("ORCS_SHADOW_REST") == "1",
LegacyAPIBase: envOrDefault("ORCS_LEGACY_API_BASE", "http://127.0.0.1:4001"),
}
}

func envOrDefault(key, fallback string) string {
v := os.Getenv(key)
if v == "" {
return fallback
}
return v
}
