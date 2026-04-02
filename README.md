Запуск kubernetes
``bash
minikube start --cpus=4 --memory=8192 --disk-size=20g --driver=docker
``

Информация о кластере
``bash
kubectl cluster-info
Kubernetes control plane is running at https://192.168.49.2:8443
CoreDNS is running at https://192.168.49.2:8443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

Установка istio
```bash
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.29.1
export PATH=$PWD/bin:$PATH
```

