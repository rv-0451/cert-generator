package storage

import (
	"fmt"

	"github.com/rv-0451/cert-generator/pkg/certs"
	"github.com/rv-0451/cert-generator/pkg/config"
)

type Storage interface {
	Save(*certs.CertContainer) error
}

func NewStorage(cfg *config.Config) (Storage, error) {
	switch cfg.Storage.Type {
	case config.FilesystemStorageType:
		return NewFilesystemStorage(cfg), nil
	case config.SecretStorageType:
		return NewSecretStorage(cfg), nil
	default:
		return nil, fmt.Errorf("unknown storage type: '%s'", cfg.Storage.Type)
	}
}
