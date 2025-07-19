# kim

项目由kubebuilder自动生成
```bash
kubebuilder init --domain cts.io

kubebuilder edit --multigroup=true

kubebuilder create api --group kim.io --version v1 --kind User --resource --controller
```

helm chart 由kubebuilder自动生成
```bash
kubebuilder edit --plugins=helm/v1-alpha
```




