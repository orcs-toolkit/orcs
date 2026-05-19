# ORCS Internal Package Ownership and Versioning

- Core backend logic must live in this repository under `/services` and `/agents`.
- Shared contracts must live under `/contracts` and act as the source of truth.
- No external npm package is allowed for core backend logic.
- Go internal packages are versioned with semantic tags from this repository.
- C++ reusable libraries (`orcs_process`, `orcs_telemetry`, `orcs_platform`) are versioned together with repository releases.
- Contract-breaking changes require a new contract version and compatibility tests before release.
