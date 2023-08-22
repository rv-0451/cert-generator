# Cert-generator

## Description

The main purpose of the cert-generator is to provide Kubernetes admission webhooks with valid SSL configuration.  
It makes sense to use it as an init-container for your deployment if your deployment runs a webhook server.  
Cert-generator generates a certificate, private key, and CA, saves them into the storage for the webhook server to use, and injects webhooks with CA.  
Storage can be either filesystem (to mount via pv/emptydir) or secret (preferred).

## Flags

See the example in the Usage section below.

| Name | Requirement | Defaults | Description |
| :---: | :---: | :---: | :---: |
| svc | mandatory | - | DNS name of the webhook service |
| org | optional | organization.com | Organization name in the generated certificate |
| storagetype | optional | secret | Type of Storage where certs will be saved. Possible values are "secret" and "filesystem" |
| certdir | optional | /tmp/webhook/certs/ | Dir for the filesystem storage where data will be stored |
| secretname | optional | cert-generator-secret | Secret name for the k8s secret storage |
| secretnamespace | optional | default | Namespace for the k8s secret |
| certname | optional | tls.crt | Certificate name |
| keyname | optional | tls.key | Key name |
| clientcaname | optional | ca.crt | CA name |
| validatingwebhooknames | optional | - | Comma-separated list of validating webhook names |
| mutatingwebhooknames | optional | - | Comma-separated list of mutating webhook names |

## Permissions

The cert-generator should be able to create/update secrets at the namespace scope (if secret storage is chosen) and get/update webhooks at the cluster scope:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cert-generator
rules:
  - apiGroups:
      - admissionregistration.k8s.io
    resources:
      - validatingwebhookconfigurations # you may only need one of those
      - mutatingwebhookconfigurations # depending on your scenario
    resourceNames:
      - my-webhook-for-configmaps
      - my-webhook-for-pvs
    verbs:
      - get
      - update
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cert-generator
  namespace: default # update with the namespace the contianer is deployed to
rules:
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - secrets
    resourceNames:
      - cert-generator-secret # update the default secret name if necessary
    verbs:
      - update
```

## Usage

For secret storage update your deployment as follows:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  # ...
spec:
  # ...
  template:
    # ...
    spec:
      # ...
      initContainers:
        - name: certgen-webhooks
          args:
            - -svc=myname.mynamespace.svc
            - -org=organization.com
            - -storagetype=secret
            - -secretname=cert-generator-secret
            - -secretnamespace=mynamespace
            - -certname=tls.crt
            - -keyname=tls.key
            - -clientcaname=ca.crt
            - -validatingwebhooknames=my-webhook-for-configmaps, my-webhook-for-pvs
         # ...
      containers:
         # ...
          volumeMounts:
            - name: webhook-certs
              readOnly: true
              mountPath: /tmp/k8s-webhook-server/serving-certs # update with the dir your webhook server is using
      # ...
      volumes:
        - name: webhook-certs
          secret:
            secretName: cert-generator-secret # update the default secret name if necessary
```

And for filesystem storage (emptyDir example):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  # ...
spec:
  # ...
  template:
    # ...
    spec:
      # ...
      initContainers:
        - name: certgen-webhooks
          args:
            - -svc=myname.mynamespace.svc
            - -org=organization.com
            - -storagetype=filesystem
            - -certdir=/tmp/webhook/certs/
            - -certname=tls.crt
            - -keyname=tls.key
            - -clientcaname=ca.crt
            - -validatingwebhooknames=my-webhook-for-configmaps, my-webhook-for-pvs
          volumeMounts:
            - mountPath: /tmp/webhook/certs
              name: webhook-certs
         # ...
      containers:
         # ...
          volumeMounts:
            - name: webhook-certs
              readOnly: true
              mountPath: /tmp/k8s-webhook-server/serving-certs # update with the dir your webhook server is using
      # ...
      volumes:
        - name: webhook-certs
          emptyDir: {}
```

## Build

Clone this repo, change Makefile (or update corresponding env vars) as necessary and:

```bash
make docker-build docker-push
```

You may run it locally for testing (update flags in Makefile accordingly), just specify kubeconfig via env var `export KUBECONFIG=/path/to/kubeconfig` and:

```bash
make run
```

For full actions see:

```bash
make help
```
