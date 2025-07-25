package key

import (
	"context"
	"crypto/ecdsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type FileKeyProvider struct {
	path string
}

func (f *FileKeyProvider) retrieveKeyPair(ctx context.Context) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return nil, nil, err
	}
	key, err := jwt.ParseECPrivateKeyFromPEM(data)
	if err != nil {
		return nil, nil, err
	}
	return key, &key.PublicKey, nil
}
