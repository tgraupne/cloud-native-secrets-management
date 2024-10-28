# External-Secrets-Operator Example

In this example, we use External-Secrets-Operator to sync secrets stored in a local instance of Hashicorp Vault.

## Install prerequisites

### Infrastructure

### Starting minikube

```bash
minikube start --driver=docker --profile=minikube-eso
```

```bash
# add helm chart repos of hashicorp vault and eso
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo add external-secrets https://charts.external-secrets.io
helm repo update

# install hashicorp vault and eso
helm install vault hashicorp/vault --set "server.dev.enabled=true"
helm install external-secrets external-secrets/external-secrets

# enable eso to talk to vault
kubectl apply -f vault-secret.yaml
kubectl apply -f vault-secret-store.yaml
```

## Configure and Deploy Demo

### Store a secret inside vault

```bash
# Login to Vault using the CLI
kubectl port-forward svc/vault 8200:8200
export VAULT_ADDR='http://127.0.0.1:8200'
vault login token=root

# Adding the Secret
vault kv put secret/myapp SECRET_VALUE=my-secret-vault-value
```

### Create encrypted Kubernetes Secret

```bash
# example: encode base64 on macos
echo -n "your-text-here" | base64

# example: decode base64 on macos
echo -n "SGVsbG8sIFdvcmxkIQ==" | base64 --decode
```

### Deploy Kubernetes Secret

```bash
kubectl apply -f external-secret.yaml
```

### Deploy the Application

```bash
kubectl apply -f deployment.yaml
```

### Expose the Application

```bash
kubectl expose deployment eso-secret-app --type=NodePort --port=8080

minikube service eso-secret-app  --profile minikube-eso
```

## Clean up

```bash
kubectl delete deployment eso-secret-app
kubectl delete ExternalSecret eso-my-secret
kubectl delete service eso-secret-app
```
