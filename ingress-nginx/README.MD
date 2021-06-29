> https://artifacthub.io/packages/helm/ingress-nginx/ingress-nginx
### Add repository
```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
```

### Install chart
```bash
kubectl create ns ingress-nginx
helm install ingress-nginx ingress-nginx/ingress-nginx --version 3.34.0 -n ingress-nginx
helm upgrade ingress-nginx ingress-nginx/ingress-nginx --install -n ingress-nginx
```