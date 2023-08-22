package storage

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/rv-0451/cert-generator/pkg/certs"
	"github.com/rv-0451/cert-generator/pkg/config"
)

type FilesystemStorage struct {
	CertDir      string
	CertName     string
	KeyName      string
	ClientCAName string
}

func NewFilesystemStorage(cfg *config.Config) *FilesystemStorage {
	return &FilesystemStorage{
		CertDir:      cfg.Storage.CertDir,
		CertName:     cfg.Storage.CertName,
		KeyName:      cfg.Storage.KeyName,
		ClientCAName: cfg.Storage.ClientCAName,
	}
}

func (fs *FilesystemStorage) Save(c *certs.CertContainer) error {
	log.Printf("Creating webhook directory %s\n", fs.CertDir)
	err := os.MkdirAll(fs.CertDir, 0775)
	if err != nil {
		return err
	}

	serverCertPath := filepath.Join(fs.CertDir, fs.CertName)
	log.Printf("Saving webhook server cert to %s\n", serverCertPath)
	err = writeFile(serverCertPath, c.ServerCertPEM)
	if err != nil {
		return err
	}

	serverKeyPath := filepath.Join(fs.CertDir, fs.KeyName)
	log.Printf("Saving webhook server key to %s\n", serverKeyPath)
	err = writeFile(serverKeyPath, c.ServerKeyPEM)
	if err != nil {
		return err
	}

	caCertPath := filepath.Join(fs.CertDir, fs.ClientCAName)
	log.Printf("Saving ca client cert to %s\n", caCertPath)
	err = writeFile(caCertPath, c.CaPEM)
	if err != nil {
		return err
	}

	return nil
}

func writeFile(filepath string, sCert *bytes.Buffer) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(sCert.Bytes()); err != nil {
		return err
	}

	return nil
}
