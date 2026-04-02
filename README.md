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

# Собираем образы Minicube
Переключить терминал внутрь Minicube
```bash
# Переключаем терминал на Docker внутри Minikube
eval $(minikube docker-env)
```
```bash
# Собираем бэкенд
cd services/backend
docker build -t mrm-backend:v1 .
```
```bash
# Собираем фронтенд
cd services/frontend
docker build -t mrm-frontend:v1 .
```

Проверить что контейнеры улетели в Minicube
```bash
docker images | grep mrm
```

Применяем манифесты в кластер
```bash
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
```

Посмотреть мои поды из моего пространства имён
```bash
kubectl get pods -n mrm-project
```

## Проверить что frontend видит backend
```bash
 kubectl port-forward svc/frontend-service 8081:8081 -n mrm-project
 ```
 Открываем ``http://localhost:8081/`` и видим 
 ```
 Frontend received: Response from Backend (v1)
 ```

 ## Конфигурации шлюза Istio
```bash
kubectl apply -f istio-configs/gateway.yaml
kubectl apply -f istio-configs/frontend-vs.yaml
```
Получить список адресов для доступа к шлюзу
```bash
minikube service istio-ingressgateway -n istio-system --url
```
Рабочий шлюз Istio
```bash
http://192.168.49.2:30711
```
Такой же ответ
```
Frontend received: Response from Backend (v1)
```





