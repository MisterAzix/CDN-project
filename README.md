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

## 👤️ Authors 👤

- Maxence BREUILLES ([@MisterAzix](https://github.com/MisterAzix))<br />
- Doriane FARAU ([@DFarau](https://github.com/DFarau))<br />
- Charles LAMBRET ([@CharlesLambret](https://github.com/CharlesLambret))<br />
- Antonin CHARPENTIER ([@toutouff](https://github.com/toutouff))<br />
- Louis FORTRIE ([@louisFortrie](https://github.com/louisFortrie))
