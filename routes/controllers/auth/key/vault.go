package key

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
)

type VaultStore struct {
	client *vault.Client
	kvPath string
}

type Config struct {
	Address   string
	Token     string
	MountPath string
	Timeout   time.Duration
}

func NewVaultStore(cfg Config) (*VaultStore, error) {
	config := &vault.Config{Address: cfg.Address, Timeout: cfg.Timeout}
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(cfg.Token)
	return &VaultStore{client: client, kvPath: cfg.MountPath}, nil
}

func (v *VaultStore) StoreFile(ctx context.Context, localPath string, key string) error {
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	fullPath := fmt.Sprintf("%s/%s", v.kvPath, key)
	_, err = v.client.KVv2(fullPath).Put(ctx, key, map[string]interface{}{
		"file": encoded,
	})
	return err
}

func (v *VaultStore) RetrieveFile(ctx context.Context, key string) ([]byte, error) {
	secret, err := v.client.KVv2(v.kvPath).Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("retrieving secret: %w", err)
	}
	encoded, ok := secret.Data["file"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'file' key in Vault")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}

	return data, nil
}
