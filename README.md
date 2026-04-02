## Mesh-Ready Microservices (MRM) 🚀
Проект по развертыванию отказоустойчивой микросервисной архитектуры в Kubernetes с использованием Istio Service Mesh.
## 🏗 Архитектура
Проект состоит из двух микросервисов на Golang:

   1. Frontend Service: Принимает внешний трафик и обращается к Backend.
   2. Backend Service (v1/v2): Возвращает данные с указанием версии.

Весь трафик между сервисами управляется через Istio Envoy Sidecars.
------------------------------
## 🛠 1. Подготовка окружения## Запуск кластера Minikube

minikube start --cpus=4 --memory=8192 --disk-size=20g --driver=docker
kubectl cluster-info

## Установка Istio

curl -L https://istio.io/downloadIstio | sh -
cd istio-1.29.1
export PATH=$PWD/bin:$PATH
istioctl install --set profile=demo -y

------------------------------
## 📦 2. Сборка и Деплой## Подготовка пространства

kubectl create namespace mrm-project
kubectl label namespace mrm-project istio-injection=enabled

## Сборка образов внутри Minikube

eval $(minikube docker-env)
# Сборка v1 и v2
docker build -t mrm-backend:v1 ./services/backend
docker build -t mrm-backend:v2 ./services/backend
docker build -t mrm-frontend:v1 ./services/frontend

## Применение манифестов K8s

kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
kubectl get pods -n mrm-project

------------------------------
## 🌐 3. Трафик и Service Mesh## Настройка Ingress Gateway

kubectl apply -f istio-configs/gateway.yaml
kubectl apply -f istio-configs/frontend-vs.yaml

Доступ к приложению:

minikube service istio-ingressgateway -n istio-system --url

## Canary Deployment (80/20)
Применяем разделение трафика между версиями v1 и v2:

kubectl apply -f istio-configs/backend-split.yaml

Тест распределения:

while true; do curl -s $(minikube service istio-ingressgateway -n istio-system --url | head -n 1) | grep "Response"; sleep 0.5; done

------------------------------
## 🛡 4. Resilience & Security (Отказоустойчивость и Безопасность)## Rate Limiting & HPA
Ограничение частоты запросов и автоscaling:

kubectl apply -f k8s/hpa.yaml
kubectl apply -f istio-configs/rate-limit.yaml

## mTLS (Шифрование)
Включение строгого шифрования между сервисами:

kubectl apply -f istio-configs/security.yaml

------------------------------
## 📊 5. Наблюдаемость (Observability)## Визуализация графа трафика (Kiali)

kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.29/samples/addons/kiali.yaml
istioctl dashboard kiali

## Трассировка запросов (Jaeger)

kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.29/samples/addons/jaeger.yaml
istioctl dashboard jaeger

------------------------------
## Основные возможности проекта:

* ✅ Traffic Management: Умная маршрутизация и Canary-релизы.
* ✅ Observability: Визуальный мониторинг и трассировка каждого запроса.
* ✅ Security: Автоматическое mTLS шифрование.
* ✅ Resilience: Повторные попытки (Retries) и Circuit Breaker.

## 🧠 Описание архитектурных решений (Service Mesh в деталях)
В этом проекте каждый конфигурационный файл решал конкретную проблему распределенных систем:
## 📡 Трафик и Маршрутизация

* gateway.yaml: Заменяет стандартный Ingress. Позволяет управлять входящим трафиком на уровне L7 (HTTP), обеспечивая единую точку входа в меш.
* backend-split.yaml (VirtualService): Решает проблему Canary-релизов. Мы можем направлять фиксированный процент трафика на новую версию (v2) без создания отдельных балансировщиков.
* backend-split.yaml (DestinationRule): Группирует поды по версиям (subsets). Без него Kubernetes видел бы v1 и v2 как одну серую массу подов.

## 🛡 Отказоустойчивость (Resilience)

* Retries (Повторы): Решают проблему «мигающей» сети. Если запрос к бэкенду сорвался, прокси-сервер сам переспросит его. Пользователь даже не узнает о сбое.
* Circuit Breaker (Размыкатель): Защищает систему от «эффекта домино». Если один сервис начинает тормозить, Istio отключает его, чтобы он не перегрузил остальные части системы.

## 🔒 Безопасность

* security.yaml (mTLS): Решает проблему безопасности внутри периметра. Даже если злоумышленник получит доступ к кластеру, он не сможет прочитать трафик между микросервисами, так как он зашифрован сертификатами, которые обновляются каждые несколько часов автоматически.

## 🚦 Ограничение ресурсов

* rate-limit.yaml: Защищает от перегрузки (DDoS или ошибок в коде). Мы ограничиваем количество запросов в секунду на уровне прокси, не нагружая само приложение проверками.

## 👁 Наблюдаемость (Observability)

* Kiali & Jaeger: Решают проблему «черного ящика». В микросервисах сложно понять, где именно застрял запрос. Эти инструменты дают полную визуальную карту и таймлайны каждого вызова.




