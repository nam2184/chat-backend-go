package key

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nam2184/mymy/util"
)

type KeyManager struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

var (
	instance *KeyManager
	once     sync.Once // Ensures KeyManager is initialized only once
)

func InitializeKeyManager(vault *VaultStore, key string) {
	data, err := vault.RetrieveFile(context.Background(), key)
	if err != nil {
		log.Fatalf(err.Error())
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(data)
	if err != nil {
		log.Fatalf(err.Error())
	}

	once.Do(func() {
		instance = &KeyManager{
			privateKey: privateKey,
			publicKey:  &privateKey.PublicKey,
		}
	})
}

func GetKeyManager() (*KeyManager, error) {
	if instance == util.GetZero[*KeyManager]() {
		return nil, fmt.Errorf("Keys not initialised")
	}
	return instance, nil
}

func GetKeyManagerTest(vault VaultStore, key string) (*KeyManager, error) {
	data, err := vault.RetrieveFile(context.Background(), key)
	if err != nil {
		log.Fatalf(err.Error())
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(data)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return &KeyManager{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

func (k *KeyManager) IssueToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, err := token.SignedString(k.privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (k *KeyManager) VerifyToken(tokenString string) (*CustomClaims, error) {
	// Parse the token with the public key
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is using the correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return k.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
