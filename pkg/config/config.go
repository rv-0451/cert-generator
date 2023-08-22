package config

import (
	"flag"
	"log"
)

const (
	FilesystemStorageType string = "filesystem"
	SecretStorageType     string = "secret"
)

var (
	org                    string
	svc                    string
	storagetype            string
	certdir                string
	secretname             string
	secretnamespace        string
	certname               string
	keyname                string
	clientcaname           string
	validatingwebhooknames string
	mutatingwebhooknames   string
)

var Cfg *Config

type Config struct {
	Certs   certsConfig
	Storage storageConfig
}

type certsConfig struct {
	Svc                    string
	Org                    string
	ValidatingWebhookNames string
	MutatingWebhookNames   string
}

type storageConfig struct {
	Type            string
	CertDir         string
	SecretName      string
	SecretNamespace string
	CertName        string
	KeyName         string
	ClientCAName    string
}

func init() {
	Cfg = NewConfig()
}

func NewConfig() *Config {
	flag.StringVar(&svc, "svc", "", "service DNS")
	flag.StringVar(&org, "org", "organization.com", "organization name")
	flag.StringVar(&storagetype, "storagetype", SecretStorageType, "type of storage where certs will be saved")
	flag.StringVar(&certdir, "certdir", "/tmp/webhook/certs/", "dir for filesystem storage")
	flag.StringVar(&secretname, "secretname", "cert-generator-secret", "secret name for k8s secret storage")
	flag.StringVar(&secretnamespace, "secretnamespace", "default", "namespace for k8s secret")
	flag.StringVar(&certname, "certname", "tls.crt", "certificate name")
	flag.StringVar(&keyname, "keyname", "tls.key", "server key name")
	flag.StringVar(&clientcaname, "clientcaname", "ca.crt", "client ca name")
	flag.StringVar(&validatingwebhooknames, "validatingwebhooknames", "", "comma-separated validating webhook names")
	flag.StringVar(&mutatingwebhooknames, "mutatingwebhooknames", "", "comma-separated mutating webhook names")
	flag.Parse()

	if svc == "" {
		log.Panicln("No service DNS name specified. Unable to generate webhook server certificate without DNS name.")
	}

	return &Config{
		Certs: certsConfig{
			Svc:                    svc,
			Org:                    org,
			MutatingWebhookNames:   mutatingwebhooknames,
			ValidatingWebhookNames: validatingwebhooknames,
		},
		Storage: storageConfig{
			Type:            storagetype,
			CertDir:         certdir,
			SecretName:      secretname,
			SecretNamespace: secretnamespace,
			CertName:        certname,
			KeyName:         keyname,
			ClientCAName:    clientcaname,
		},
	}
}
