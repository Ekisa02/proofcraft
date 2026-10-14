package zk

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
)

// Proof holds hex-encoded components
type Proof struct {
	RxHex string `json:"rx"`
	RyHex string `json:"ry"`
	SHex  string `json:"s"`
}

// GenerateKeyPair returns a private scalar and public point coordinates
func GenerateKeyPair() (*big.Int, *big.Int, *big.Int, error) {
	curve := elliptic.P256()
	priv, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, nil, err
	}
	// priv from GenerateKey is a []byte; convert to big.Int
	d := new(big.Int).SetBytes(priv)
	return d, x, y, nil
}

// Prove produces a non-interactive Schnorr proof of knowledge of private key d for point (pubX,pubY)
// Returns the proof and the public key coordinates used
func Prove(message string, privBytes []byte) (Proof, *big.Int, *big.Int, error) {
	curve := elliptic.P256()
	d := new(big.Int).SetBytes(privBytes)
	if d.Sign() == 0 {
		return Proof{}, nil, nil, errors.New("invalid private key bytes")
	}

	// Compute public key
	pubX, pubY := curve.ScalarBaseMult(d.Bytes())

	// Commitment: R = k*G
	k, _, _, err := GenerateKeyPair()
	if err != nil {
		return Proof{}, nil, nil, err
	}
	Rx, Ry := curve.ScalarBaseMult(k.Bytes())

	// Fiat-Shamir: c = H(R || Pub || msg)
	c := hashChallenge(Rx, Ry, pubX, pubY, []byte(message))

	// s = k + c * d mod n
	n := curve.Params().N
	cd := new(big.Int).Mul(c, d)
	s := new(big.Int).Add(k, cd)
	s.Mod(s, n)

	return Proof{RxHex: hex.EncodeToString(Rx.Bytes()), RyHex: hex.EncodeToString(Ry.Bytes()), SHex: hex.EncodeToString(s.Bytes())}, pubX, pubY, nil
}

// Verify checks the proof
func Verify(message string, proof Proof, pubXHex, pubYHex string) (bool, error) {
	curve := elliptic.P256()

	rxBytes, err := hex.DecodeString(proof.RxHex)
	if err != nil {
		return false, err
	}
	ryBytes, err := hex.DecodeString(proof.RyHex)
	if err != nil {
		return false, err
	}
	sBytes, err := hex.DecodeString(proof.SHex)
	if err != nil {
		return false, err
	}
	pubXBytes, err := hex.DecodeString(pubXHex)
	if err != nil {
		return false, err
	}
	pubYBytes, err := hex.DecodeString(pubYHex)
	if err != nil {
		return false, err
	}

	rx := new(big.Int).SetBytes(rxBytes)
	ry := new(big.Int).SetBytes(ryBytes)
	s := new(big.Int).SetBytes(sBytes)
	pubX := new(big.Int).SetBytes(pubXBytes)
	pubY := new(big.Int).SetBytes(pubYBytes)

	// c = H(R || Pub || msg)
	c := hashChallenge(rx, ry, pubX, pubY, []byte(message))

	// Compute s*G and R + c*Pub
	// left = s*G
	leftX, leftY := curve.ScalarBaseMult(s.Bytes())
	// right = R + c*Pub
	cx, cy := curve.ScalarMult(pubX, pubY, c.Bytes())
	rx2, ry2 := curve.Add(rx, ry, cx, cy)

	if leftX.Cmp(rx2) == 0 && leftY.Cmp(ry2) == 0 {
		return true, nil
	}
	return false, nil
}

func hashChallenge(rx, ry, pubX, pubY *big.Int, message []byte) *big.Int {
	h := sha256.New()
	h.Write(rx.Bytes())
	h.Write(ry.Bytes())
	h.Write(pubX.Bytes())
	h.Write(pubY.Bytes())
	h.Write(message)
	sum := h.Sum(nil)
	c := new(big.Int).SetBytes(sum)
	// reduce modulo curve order
	n := elliptic.P256().Params().N
	c.Mod(c, n)
	return c
}
