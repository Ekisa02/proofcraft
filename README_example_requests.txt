Examples (curl)

# Generate key
curl -X POST -H "Content-Type: application/json" -d '{"label":"demo"}' http://localhost:8080/keys/generate

# Create proof (supply private key returned by previous call)
curl -X POST -H "Content-Type: application/json" -d '{"message":"hello","private_key":"<hexprivate>"}' http://localhost:8080/proofs/prove

# Verify proof
curl -X POST -H "Content-Type: application/json" -d '{"message":"hello","proof":{"rx":"...","ry":"...","s":"..."},"public_x":"...","public_y":"..."}' http://localhost:8080/proofs/verify
