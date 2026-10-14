package services

import (
	"encoding/hex"
	"errors"

	"gorm.io/gorm"
	"github.com/example/proofcraft/internal/zk"
	"github.com/example/proofcraft/pkg/models"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GenerateKey creates a new keypair, stores the public key and returns the private key to caller
func (s *Service) GenerateKey(label string) (map[string]interface{}, error) {
	priv, pubX, pubY, err := zk.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	k := models.Key{
		Label:   label,
		PubXHex: hex.EncodeToString(pubX.Bytes()),
		PubYHex: hex.EncodeToString(pubY.Bytes()),
	}
	if err := s.db.Create(&k).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          k.ID,
		"label":       k.Label,
		"private_key": hex.EncodeToString(priv.Bytes()),
		"public_x":    k.PubXHex,
		"public_y":    k.PubYHex,
	}, nil
}

// CreateProof creates a proof using the provided private key and stores it
func (s *Service) CreateProof(message, privateHex, label string) (map[string]interface{}, error) {
	if message == "" || privateHex == "" {
		return nil, errors.New("message and private_key are required")
	}
	privBytes, err := hex.DecodeString(privateHex)
	if err != nil {
		return nil, err
	}

	proof, pubX, pubY, err := zk.Prove(message, privBytes)
	if err != nil {
		return nil, err
	}

	p := models.Proof{
		Label:   label,
		Message: message,
		RxHex:   proof.RxHex,
		RyHex:   proof.RyHex,
		SHex:    proof.SHex,
	}
	if err := s.db.Create(&p).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"proof":  proof,
		"stored": true,
		"public": map[string]string{"public_x": hex.EncodeToString(pubX.Bytes()), "public_y": hex.EncodeToString(pubY.Bytes())},
	}, nil
}

// VerifyProof verifies a proof for message and public key
func (s *Service) VerifyProof(message, rxHex, ryHex, sHex, pubXHex, pubYHex string) (bool, error) {
	if message == "" {
		return false, errors.New("message required")
	}
	proof := zk.Proof{RxHex: rxHex, RyHex: ryHex, SHex: sHex}
	return zk.Verify(message, proof, pubXHex, pubYHex)
}
