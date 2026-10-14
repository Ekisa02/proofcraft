package zk

import (
	"encoding/hex"
	"testing"
)

func TestSchnorrProveVerify(t *testing.T) {
	// Generate keypair
	priv, pubX, pubY, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("keygen err: %v", err)
	}

	message := "hello zeroknowledge"
	proof, _, _, err := Prove(message, priv.Bytes())
	if err != nil {
		t.Fatalf("prove err: %v", err)
	}

	ok, err := Verify(message, proof, hex.EncodeToString(pubX.Bytes()), hex.EncodeToString(pubY.Bytes()))
	if err != nil {
		t.Fatalf("verify err: %v", err)
	}
	if !ok {
		t.Fatalf("expected verification to succeed")
	}
}
