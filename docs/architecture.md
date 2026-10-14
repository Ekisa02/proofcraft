# Architecture

ProofCraft is organized to separate concerns and allow independent development of the cryptographic primitives, application services, and HTTP controllers.

- internal/zk: "internal" package for cryptographic code. Keeping cryptography internal prevents accidental reuse outside the module without explicit import path.
- pkg/services: business logic that composes zk primitives and persistence
- pkg/controllers: HTTP layer (Gin). Minimal, thin controllers that delegate to services and return JSON.
- pkg/models: GORM models describing persisted entities.
- pkg/db: DB initialization wrapper.

Flow examples

- Key generation: Controller -> Service.GenerateKey -> zk.GenerateKeyPair -> models.Key stored
- Prove flow: Controller receives message + private key -> Service.CreateProof -> zk.Prove -> models.Proof stored
- Verify: Controller -> Service.VerifyProof -> zk.Verify

Notes

- The zk implementation is an instructive Schnorr proof and is intentionally simple. Replace or augment with vetted, high-performance ZK systems for production use.
