package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"time"
)

type CertContainer struct {
	ServerCertPEM *bytes.Buffer
	ServerKeyPEM  *bytes.Buffer
	CaPEM         *bytes.Buffer
}

func NewCertContainer(svcDnsName string, org string) *CertContainer {
	// CA config
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2020),
		Subject: pkix.Name{
			Organization: []string{org},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// CA private key
	log.Println("Generating CA private key")
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Println(err)
	}

	// Self signed CA certificate
	log.Println("Creating CA self signed certificate")
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
	if err != nil {
		log.Println(err)
	}

	// PEM encode CA cert
	caPEM := new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	// server cert config
	cert := &x509.Certificate{
		DNSNames: []string{
			svcDnsName,
			svcDnsName + ".cluster.local",
		},
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName:   "project-operator.project-operator.svc",
			Organization: []string{org},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// server private key
	log.Println("Generating webhook server private key")
	serverKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Println(err)
	}

	// sign the server cert
	log.Println("Signing webhook server private key using CA certificate")
	serverCertBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &serverKey.PublicKey, caKey)
	if err != nil {
		log.Println(err)
	}

	// PEM encode the  server cert and key
	serverCertPEM := new(bytes.Buffer)
	_ = pem.Encode(serverCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertBytes,
	})

	serverKeyPEM := new(bytes.Buffer)
	_ = pem.Encode(serverKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverKey),
	})

	return &CertContainer{
		CaPEM:         caPEM,
		ServerCertPEM: serverCertPEM,
		ServerKeyPEM:  serverKeyPEM,
	}
}
