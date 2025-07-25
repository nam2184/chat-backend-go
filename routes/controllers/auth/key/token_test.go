package key

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"testing"
	"time"

	util "github.com/nam2184/generic-queries/utils"
)

var (
	username = "hello"
	id       = int64(65)
	keyID    = "hello"
)

func generateECDSAPrivateKey() (*ecdsa.PrivateKey, error) {
	// Generate an ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func encodePrivateKey(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	}
	return pem.EncodeToMemory(pemBlock), nil
}

func TestGenerateAndVerifyToken(t *testing.T) {
	fmt.Println("[Start] TestGenerateAndVerifyToken")

	// Step 1: Create Vault connection
	vaultAddr := "http://127.0.0.1:8200"
	vaultToken := os.Getenv("VAULT_TOKEN")
	fmt.Println("[Info] Vault address:", vaultAddr)

	cfg := Config{
		Address:   vaultAddr,
		Token:     vaultToken,
		MountPath: "secret",
	}

	vaultStore, err := NewVaultStore(cfg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("[OK] Connected to Vault")

	// Step 2: Generate two key pairs (PEM format)
	fmt.Println("[Info] Generating ECDSA key pairs")
	secret1Path := "secret1"
	secret2Path := "secret2"

	privateKey, err := generateECDSAPrivateKey()
	privateKey2, err := generateECDSAPrivateKey()
	if err != nil {
		t.Fatal("Failed to generate ECDSA private key:", err)
	}

	keyByte, err := encodePrivateKey(privateKey)
	if err != nil {
		t.Fatal("Failed to encode ECDSA private key:", err)
	}

	keyByte2, err := encodePrivateKey(privateKey2)
	if err != nil {
		t.Fatal("Failed to encode ECDSA private key:", err)
	}

	fmt.Println("[OK] Keys generated and encoded")

	fmt.Println("[Info] Writing keys to Vault")
	_, err = vaultStore.client.KVv2(vaultStore.kvPath).Put(context.Background(), secret1Path, map[string]interface{}{
		"file": base64.StdEncoding.EncodeToString(keyByte),
	})
	if err != nil {
		t.Fatal("Failed to write key 1 to Vault:", err)
	}

	_, err = vaultStore.client.KVv2(vaultStore.kvPath).Put(context.Background(), secret2Path, map[string]interface{}{
		"file": base64.StdEncoding.EncodeToString(keyByte2),
	})
	if err != nil {
		t.Fatal("Failed to write key 2 to Vault:", err)
	}
	fmt.Println("[OK] Keys written to Vault")

	// Step 3: Create claims
	fmt.Println("[Info] Creating claims")
	accessCfg := TokenConfig{
		Username: "john",
		ID:       123,
		Type:     AccessToken,
		Expiry:   24 * time.Hour,
		Issuer:   "auth-service",
	}

	claims := NewClaims(accessCfg)

	// Step 4: Issue token with secret1
	fmt.Println("[Info] Initializing KeyManager with secret1")
	manager, err := GetKeyManagerTest(*vaultStore, secret1Path)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tokenString, err := manager.IssueToken(claims)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("[OK] Token issued with secret1")

	// Step 5: Issue token with secret2 (different key)
	fmt.Println("[Info] Initializing KeyManager with secret2")
	InitializeKeyManager(vaultStore, secret2Path)
	managerFalse, err := GetKeyManagerTest(*vaultStore, secret2Path)
	if err != nil {
		t.Fatalf(err.Error())
	}

	falseToken, err := managerFalse.IssueToken(claims)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("[OK] Token issued with secret2")

	fmt.Println("Correct Token : ", tokenString)
	fmt.Println("False Token: ", falseToken)

	// Step 6: Verify correct token
	fmt.Println("[Info] Verifying correct token")
	newClaims, err := manager.VerifyToken(tokenString)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("[OK] Correct token verified")

	util.PrintStructAttributes(newClaims)

	// Step 7: Attempt to verify incorrect token
	fmt.Println("[Info] Verifying incorrect token")
	_, err = manager.VerifyToken(falseToken)
	if err == nil {
		t.Fatalf("False token but verified as true")
	}

	fmt.Println("[OK] Verification failed as expected:", err.Error())
	fmt.Println("[End] TestGenerateAndVerifyToken")
}
