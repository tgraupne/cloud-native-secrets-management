# SOPS Example

In this example, we encrypt the manifest file of a kubernetes secret in order to store it in our VCS.
To encrypt the file we used the tool SOPS.

## Install prerequisites

### MacOS

```bash
# install the tool sops for encryption
brew install sops gnupg

# prepare key
gpg --full-generate-key
```

### Using the key in CI

```bash
# get key id
gpg --list-secret-keys --keyid-format LONG

# extract the key
gpg --export-secret-keys --armor <YOUR_GPG_KEY_ID> > private-key.asc
```

### Starting minikube

```bash
minikube start --driver=docker --profile=minikube-sops
```

Store the content of `private-key.asc` as a GitHub secret named `GPG_PRIVATE_KEY`. 

**Example GitHub Actions Configuration**
```yaml
name: Deploy with SOPS

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Import GPG key
      env:
        GPG_PRIVATE_KEY: ${{ secrets.GPG_PRIVATE_KEY }}
      run: |
        echo "$GPG_PRIVATE_KEY" | gpg --batch --import

    - name: Decrypt SOPS secret
      run: |
        sops --decrypt secret.enc.yaml > secret.yaml

    - name: Apply Kubernetes secret
      run: |
        kubectl apply -f secret.yaml
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
# get key id
gpg --list-secret-keys --keyid-format LONG

# encrypt secret 
sops --encrypt --pgp <YOUR_GPG_KEY_ID> secret.yaml > secret.enc.yaml
```

### Deploy Kubernetes Secret

```bash
# with sops
sops --decrypt secret.enc.yaml | kubectl apply -f -
```

### Deploy the Application

```bash
kubectl apply -f deployment.yaml
```

### Expose the Application

```bash
kubectl expose deployment sops-secret-app --type=NodePort --port=8080

minikube service sops-secret-app  --profile minikube-sops
```

## Clean up

```bash
kubectl delete deployment sops-secret-app
kubectl delete service sops-secret-app
kubectl delete secret sops-my-secret
```
