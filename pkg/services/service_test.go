package services

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/example/proofcraft/internal/zk"
	"github.com/example/proofcraft/pkg/db"
)

func TestServiceGenerateAndProve(t *testing.T) {
	// use a temp sqlite file
	path := "test_service.db"
	defer os.Remove(path)
	dbConn, err := db.Init(path)
	if err != nil {
		t.Fatalf("db init: %v", err)
	}
	svc := NewService(dbConn)

	res, err := svc.GenerateKey("test")
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	privHex := res["private_key"].(string)

	out, err := svc.CreateProof("msg", privHex, "label")
	if err != nil {
		t.Fatalf("create proof: %v", err)
	}
	_ = out

	// quick verify using zk.Verify
	// get public from private
	privBytes, _ := hex.DecodeString(privHex)
	proofMap := out["proof"].(zk.Proof)
	_, pubX, pubY, err := zk.Prove("dummy", privBytes) // new prove just to get pubX/pubY quickly
	if err != nil {
		t.Fatalf("internal: %v", err)
	}
	ok, err := zk.Verify("msg", proofMap, hex.EncodeToString(pubX.Bytes()), hex.EncodeToString(pubY.Bytes()))
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !ok {
		t.Fatalf("expected proof verify true")
	}
}
