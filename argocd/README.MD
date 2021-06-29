> Ref: https://artifacthub.io/packages/helm/argo/argo-cd

### Add repository
```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update
```

### Install chart
```bash
kubectl create ns argocd
helm install argocd argo/argo-cd --version 3.6.11 -n argocd
helm upgrade argocd argo/argo-cd -n argocd
```