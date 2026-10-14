# ProofCraft

ProofCraft is a modular Go project demonstrating simple, educational zero-knowledge proof implementations (Schnorr non-interactive proofs using Fiat-Shamir) with a Gin HTTP API and GORM-based persistence. It is designed as a starting point for experimentation, learning, and extension into production-grade ZK systems.

Key features

- Toy Schnorr ZK implementation on P-256 (educational) with non-interactive Fiat-Shamir transform
- Gin-powered HTTP API for key generation, proof creation, and proof verification
- GORM-backed SQLite storage for public keys and recorded proofs
- Clean, modular project layout suitable for expansion (internal packages, services, controllers)
- Unit tests for core cryptography and integration-style tests for handlers

Architecture

- cmd/ or root main.go: bootstraps the server and dependencies
- pkg/db: GORM DB initialization and migrations
- pkg/models: GORM models (Key, Proof)
- pkg/controllers: HTTP handlers for keys and proofs
- pkg/services: application logic (interact with zk primitives and DB)
- internal/zk: cryptographic proofs and helpers (Schnorr toy implementation)
- docs/: supplementary documentation and architecture notes

Security note

This repository contains an educational/ toy implementation of ZK primitives. Do NOT use this code in production without careful cryptographic review. Private keys are sensitive; in this demo the private key is returned on generation to illustrate workflow — in a real system private keys must be securely managed by the client or a secure KMS.

Getting started (local)

Prerequisites

- Go 1.20+ installed
- git (optional)

Clone and run

1. Clone repository

   git clone https://github.com/example/proofcraft.git
   cd proofcraft

2. Build and run

   go run main.go

   The server will start on :8080 and create a local SQLite DB file proofcraft.db.

API (basic)

- POST /keys/generate
  - Request: { "label": "optional label" }
  - Response: { "private_key": "<hex>", "public_x": "<hex>", "public_y": "<hex>", "id": <id> }

- POST /proofs/prove
  - Request: { "message": "msg", "private_key": "<hex>", "label": "optional label" }
  - Response: { "proof": {"rx": "<hex>", "ry": "<hex>", "s": "<hex>"}, "stored": true }

- POST /proofs/verify
  - Request: { "message": "msg", "proof": { "rx": "<hex>", "ry": "<hex>", "s": "<hex>" }, "public_x": "<hex>", "public_y": "<hex>" }
  - Response: { "valid": true }

Project layout

- main.go: server bootstrap
- go.mod: Go module file
- internal/zk: cryptographic primitives and tests
- pkg/controllers, pkg/services, pkg/models, pkg/db: app structure
- docs/architecture.md: textual architecture overview

Development & tests

- Unit tests: go test ./... (includes zk tests and service tests)

Extending

- Replace SQLite + file storage with an enterprise DB: change pkg/db/driver
- Replace toy Schnorr with established ZK libraries (bulletproofs, PLONK, Groth16) as needed

License & contributions

This is an example educational project. Contributions are welcome — please open issues or PRs with improvements, tests, or realistic ZK integrations.
