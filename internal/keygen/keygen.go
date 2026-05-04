package keygen

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

type Pair struct {
	PrivatePath string
	PublicPath  string
	Authorized  []byte
	Signer      ssh.Signer
}

func EnsureEd25519(privatePath string) (Pair, error) {
	pubPath := privatePath + ".pub"
	if _, err := os.Stat(privatePath); err == nil {
		return load(privatePath, pubPath)
	}
	if err := os.MkdirAll(filepath.Dir(privatePath), 0o700); err != nil {
		return Pair{}, err
	}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return Pair{}, err
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return Pair{}, err
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if err := os.WriteFile(privatePath, pemBytes, 0o600); err != nil {
		return Pair{}, err
	}
	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		return Pair{}, err
	}
	authorized := ssh.MarshalAuthorizedKey(sshPub)
	if err := os.WriteFile(pubPath, authorized, 0o644); err != nil {
		return Pair{}, err
	}
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		return Pair{}, err
	}
	return Pair{PrivatePath: privatePath, PublicPath: pubPath, Authorized: authorized, Signer: signer}, nil
}

func load(privatePath, pubPath string) (Pair, error) {
	b, err := os.ReadFile(privatePath)
	if err != nil {
		return Pair{}, err
	}
	signer, err := ssh.ParsePrivateKey(b)
	if err != nil {
		return Pair{}, fmt.Errorf("parse private key: %w", err)
	}
	pub := ssh.MarshalAuthorizedKey(signer.PublicKey())
	_ = os.WriteFile(pubPath, pub, 0o644)
	return Pair{PrivatePath: privatePath, PublicPath: pubPath, Authorized: pub, Signer: signer}, nil
}
