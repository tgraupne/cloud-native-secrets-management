# SOPS Example

In this example, we use the Sealed Secrets and `kubeseal` to encrypt a kubernetes secret to prepared it for storage in our VCS.

## Install prerequisites

### MacOS

```bash
# install the tool kubeseal for encryption
brew install kubeseal
```

### Infrastructure

```bash
# install the sealed secrets operator
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.23.0/controller.yaml
```

## Configure and Deploy Demo

### Create encrypted Kubernetes Secret

```bash
# example: encode base64 on macos
echo -n "your-text-here" | base64

# example: decode base64 on macos
echo -n "SGVsbG8sIFdvcmxkIQ==" | base64 --decode
```

```bash
# encrypt secret 
kubeseal --format yaml < secret.yaml > sealed-secret.yaml
```

### Deploy Kubernetes Secret

```bash
kubectl apply -f sealed-secret.yaml
```

### Deploy the Application

```bash
kubectl apply -f deployment.yaml
```

### Expose the Application

```bash
kubectl expose deployment sealed-secret-app --type=NodePort --port=8080

minikube service sealed-secret-app
```

## Clean up

```bash
kubectl delete deployment sealed-secret-app
kubectl delete sealedsecret sealed-my-secret
kubectl delete service sealed-secret-app
```
