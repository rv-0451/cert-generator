package main

import (
	"log"
	"strings"

	"github.com/rv-0451/cert-generator/pkg/certs"
	"github.com/rv-0451/cert-generator/pkg/config"
	"github.com/rv-0451/cert-generator/pkg/kubeapi"
	"github.com/rv-0451/cert-generator/pkg/storage"
)

func main() {
	cfg := config.Cfg
	kclient := kubeapi.KClient

	c := certs.NewCertContainer(cfg.Certs.Svc, cfg.Certs.Org)

	s, err := storage.NewStorage(cfg)
	if err != nil {
		log.Panicf("Failed to get Storage instance: %s", err)
	}

	if err := s.Save(c); err != nil {
		log.Panicf("Failed to save cert data to storage: %s", err)
	}

	if cfg.Certs.ValidatingWebhookNames != "" {
		for _, vw := range strings.Split(cfg.Certs.ValidatingWebhookNames, ",") {
			if err := kclient.InjectCAtoValidatingWebhook(strings.TrimSpace(vw), c.CaPEM); err != nil {
				log.Panicf("Failed inject CA into validating webhook: %s", err)
			}
		}
	} else {
		log.Println("No validating webhook names supplied, skipping...")
	}

	if cfg.Certs.MutatingWebhookNames != "" {
		for _, mw := range strings.Split(cfg.Certs.ValidatingWebhookNames, ",") {
			if err := kclient.InjectCAtoMutatingWebhook(strings.TrimSpace(mw), c.CaPEM); err != nil {
				log.Panicf("Failed inject CA into mutating webhook: %s", err)
			}
		}
	} else {
		log.Println("No mutating webhook names supplied, skipping...")
	}
}
