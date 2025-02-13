# 🚀 HETIC CDN Project 🚀

L’objectif de ce projet est de construire un prototype de CDN en Go, qui permettra d’accélérer la distribution de contenu en réduisant la charge sur les serveurs d’origine.

## How to run the project
1. Clone the repository
```bash
git clone git@github.com:DFarau/CDN-project.git
```
2. Run Docker
```bash
docker-compose up -d
```
3. Install the dependencies
```bash
go mod tidy
```
4. Run the project
```bash
go run main.go
```
App is now running on [`http://localhost:8080/`](http://localhost:8080/)

## Running with Kurbenetes
1. Launch minikube
```bash
minikube start --driver=docker
minikube addons enable ingress
```
2. Create docker image
```bash
docker build -t go-backend .
```
3. Load the image into minikube
```bash
minikube image load go-backend
```
4. Create the secret
```bash
kubectl create secret generic secret --from-env-file=k8s/.env
```
5. Apply the deployment and service
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```
6. Launch the minikube tunnel
```bash
minikube tunnel
```
App is now running on [`http://127.0.0.1/`](http://127.0.0.1/)

## 👤️ Authors 👤

- Maxence BREUILLES ([@MisterAzix](https://github.com/MisterAzix))<br />
- Doriane FARAU ([@DFarau](https://github.com/DFarau))<br />
- Charles LAMBRET ([@CharlesLambret](https://github.com/CharlesLambret))<br />
- Antonin CHARPENTIER ([@toutouff](https://github.com/toutouff))<br />
- Louis FORTRIE ([@louisFortrie](https://github.com/louisFortrie))
