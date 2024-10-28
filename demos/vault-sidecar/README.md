# External-Secrets-Operator Example

In this example, we use Vault Agend Sidecar to sync secrets stored in a local instance of Hashicorp Vault into a mounted volume.

## Install prerequisites

### Infrastructure

### Starting minikube

```bash
minikube start --driver=docker --profile=minikube-vault-sc
```

### Installing dependencies

```bash
# add helm chart repos
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

# install hashicorp vault and csi driver
helm install vault hashicorp/vault \
  --set "server.dev.enabled=true" \
  --set "injector.enabled=true" \
```

### Configuring vault

```bash
# enable eso to talk to vault
kubectl apply -f vault-secret.yaml
```

```bash
# start a shell inside vault
kubectl exec -it vault-0 -- /bin/sh

# enable the Kubernetes authentication method
vault auth enable kubernetes

# configure the Kubernetes authentication method with the Kubernetes API address.
# it will automatically use the Vault pod's own service account token.
vault write auth/kubernetes/config \
    kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443"
    
# create a policy named internal-app.
# this will be used to give the app-sa service account permission to read the kv secret created earlier.
vault policy write internal-app - <<EOF
path "secret/data/myapp" {
  capabilities = ["read"]
}
EOF

# create a Kubernetes authentication role named app that binds this policy with a Kubernetes service account named app-sa.
vault write auth/kubernetes/role/app \
    bound_service_account_names=app-sa \
    bound_service_account_namespaces=default \
    policies=internal-app \
    ttl=20m
    
exit
```

```bash
kubectl create serviceaccount app-sa
```

## Configure and Deploy Demo

### Store a secret inside vault

```bash
# Login to Vault using the CLI
kubectl port-forward svc/vault 8200:8200
export VAULT_ADDR='http://127.0.0.1:8200'
vault login token=root

# Adding the Secret
vault kv put secret/myapp ttl=10s SECRET_VALUE=my-secret-vault-sc
```

### Create encrypted Kubernetes Secret

```bash
# example: encode base64 on macos
echo -n "your-text-here" | base64

# example: decode base64 on macos
echo -n "SGVsbG8sIFdvcmxkIQ==" | base64 --decode
```

### Deploy Kubernetes Secret

### Deploy the Application

```bash
kubectl apply -f deployment.yaml
```

### Expose the Application

```bash
kubectl expose deployment vault-sc-secret-app --type=NodePort --port=8080

minikube service vault-sc-secret-app --profile minikube-vault-sc
```

## Clean up

```bash
kubectl delete deployment vault-sc-secret-app
kubectl delete ExternalSecret vault-sc-my-secret
kubectl delete service vault-sc-secret-app
```
