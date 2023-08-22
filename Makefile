# Current init image version
VERSION ?= 0.0.1
# Image URL to use all building/pushing image targets
IMG ?= 192.168.33.16:5000/automation/cert-generator:$(VERSION)

##@ General

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

##@ Build

docker-build: fmt vet  ## Build docker image with the cert-generator.
	docker build . -t ${IMG}

# Push the docker image
docker-push: ## Push docker image with the cert-generator.
	docker push ${IMG}

run: fmt vet ## Run the cert-generator from your host.
	go run main.go \
		-svc="myname.mynamespace.svc" \
		-org="organization.com" \
		-storagetype="secret" \
		-secretname="cert-generator-secret" \
		-secretnamespace="mynamespace" \
		-certname="tls.crt" \
		-keyname="tls.key" \
		-clientcaname="ca.crt" \
		-validatingwebhooknames=""
