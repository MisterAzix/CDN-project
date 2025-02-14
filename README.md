# 🚀 Projet CDN HETIC 🚀

L’objectif de ce projet est de construire un prototype de CDN en Go, qui permettra d’accélérer la distribution de contenu en réduisant la charge sur les serveurs d’origine.

## Comment exécuter le projet
1. Cloner le dépôt
```bash
git clone git@github.com:DFarau/CDN-project.git
```
2. Lancer Docker
```bash
docker-compose up -d
```
3. Installer les dépendances
```bash
go mod tidy
```
4. Exécuter le projet
```bash
go run main.go
```
L'application est maintenant en cours d'exécution sur [`http://localhost:8080/`](http://localhost:8080/)

## Exécution avec Kubernetes
1. Lancer minikube
```bash
minikube start --driver=docker
minikube addons enable ingress
```
2. Créer l'image Docker
```bash
docker build -t go-backend .
```
3. Charger l'image dans minikube
```bash
minikube image load go-backend
```
4. Créer le secret
```bash
kubectl create secret generic secret --from-env-file=k8s/.env
```
5. Appliquer le déploiement et le service
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```
6. Lancer le tunnel minikube
```bash
minikube tunnel
```
L'application est maintenant en cours d'exécution sur [`http://127.0.0.1/`](http://127.0.0.1/)


## Endpoints de l'API

Vous trouverez des requêtes de test dans la collection postman du repo.

### Upload File
- **URL**: `POST /upload`
- **Description**: Télécharge un fichier.
- **Paramètres**:
  - `file` (file): Le fichier à télécharger.
  - `parent_id` (text): L'ID du dossier parent.

### Delete File
- **URL**: `DELETE /delete/{file_id}`
- **Description**: Supprime un fichier.
- **Paramètres**:
  - `file_id` (path): L'ID du fichier à supprimer.

### List Files
- **URL**: `GET /list?parent_id=parent_folder_id`
- **Description**: Liste les fichiers dans un dossier.
- **Paramètres**:
  - `parent_id` (query): L'ID du dossier parent.

### Create Folder
- **URL**: `POST /folder/upload`
- **Description**: Crée un nouveau dossier.
- **Paramètres**:
  - `name` (text): Le nom du nouveau dossier.
  - `uploader_id` (text): L'ID de l'uploader.
  - `parent_id` (text): L'ID du dossier parent.

### Serve File
- **URL**: `GET /serve-file?id=file_id&uploader_id=uploader_id`
- **Description**: Sert un fichier.
- **Paramètres**:
  - `id` (query): L'ID du fichier.
  - `uploader_id` (query): L'ID de l'uploader.

### Delete Folder
- **URL**: `DELETE /folder/delete`
- **Description**: Supprime un dossier.
- **Paramètres**:
  - `id` (text): L'ID du dossier.
  - `uploader_id` (text): L'ID de l'uploader.

### Login
- **URL**: `POST /login`
- **Description**: Connecte un utilisateur.
- **Paramètres**:
  - `username` (text): Le nom d'utilisateur.
  - `password` (text): Le mot de passe.

### Register
- **URL**: `POST /register`
- **Description**: Enregistre un nouvel utilisateur.
- **Paramètres**:
  - `username` (text): Le nom d'utilisateur.
  - `password` (text): Le mot de passe.
  - `email` (text): L'adresse email.

## Vidéo de démonstration

[Vidéo de démonstration](https://youtu.be/F_M_L0AyD_I)

## 👤️ Auteurs 👤

- Maxence BREUILLES ([@MisterAzix](https://github.com/MisterAzix))<br />
- Doriane FARAU ([@DFarau](https://github.com/DFarau))<br />
- Charles LAMBRET ([@CharlesLambret](https://github.com/CharlesLambret))<br />
- Antonin CHARPENTIER ([@toutouff](https://github.com/toutouff))<br />
- Louis FORTRIE ([@louisFortrie](https://github.com/louisFortrie))
