package storage

import (
	"log"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/rv-0451/cert-generator/pkg/certs"
	"github.com/rv-0451/cert-generator/pkg/config"
	"github.com/rv-0451/cert-generator/pkg/kubeapi"
)

type SecretStorage struct {
	SecretName      string
	SecretNamespace string
	CertName        string
	KeyName         string
	ClientCAName    string
}

func NewSecretStorage(cfg *config.Config) *SecretStorage {
	return &SecretStorage{
		SecretName:      cfg.Storage.SecretName,
		SecretNamespace: cfg.Storage.SecretNamespace,
		CertName:        cfg.Storage.CertName,
		KeyName:         cfg.Storage.KeyName,
		ClientCAName:    cfg.Storage.ClientCAName,
	}
}

func (s *SecretStorage) Save(c *certs.CertContainer) error {
	kclient := kubeapi.KClient

	data := map[string][]byte{
		s.CertName:     c.ServerCertPEM.Bytes(),
		s.KeyName:      c.ServerKeyPEM.Bytes(),
		s.ClientCAName: c.CaPEM.Bytes(),
	}

	log.Printf("Saving cert data to secret '%s' in namespace '%s'.", s.SecretName, s.SecretNamespace)
	if err := kclient.CreateSecret(s.SecretName, s.SecretNamespace, data); err != nil {
		if apierrors.IsAlreadyExists(err) {
			log.Println("Secret already exists. Updating existing secret.")
			if err := kclient.UpdateSecret(s.SecretName, s.SecretNamespace, data); err != nil {
				log.Printf("Failed to update secret '%s' in namespace '%s'.", s.SecretName, s.SecretNamespace)
				return err
			}
			log.Println("Secret updated. Cert data successfully saved to secret.")
			return nil
		}
		log.Printf("Failed to create secret '%s' in namespace '%s'.", s.SecretName, s.SecretNamespace)
		return err
	}
	log.Println("Secret created. Cert data successfully saved to secret.")
	return nil
}
