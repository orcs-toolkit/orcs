package config

import "os"

type Config struct {
	ListenAddr      string
	AuthSecret      string
	ShadowEnabled   bool
	LegacySocketURL string
}

func Load() Config {
	addr := os.Getenv("GO_MASTER_ADDR")
	if addr == "" {
		addr = ":4000"
	}
	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		secret = "admin"
	}
	legacy := os.Getenv("ORCS_LEGACY_SOCKET_URL")
	if legacy == "" {
		legacy = "http://127.0.0.1:4000"
	}
	return Config{ListenAddr: addr, AuthSecret: secret, ShadowEnabled: os.Getenv("ORCS_SHADOW_SOCKET") == "1", LegacySocketURL: legacy}
}
