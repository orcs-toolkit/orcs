# ORCS Polyglot Services

- `go-api`: REST APIs for auth, user, policy, logs.
- `go-master`: realtime orchestration service for admin and node events.

Both services support migration safety toggles:
- `ORCS_SHADOW_REST=1` for REST shadow write calls.
- `ORCS_SHADOW_SOCKET=1` for socket path migration toggle.
