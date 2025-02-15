# Exoplanet

## Introduction
This program downloads a catalog of exoplanet data and displays the following information:
1. The number of orphan planets (no star).
2. The name (planet identifier) of the planet orbiting the hottest star.
3. A timeline of the number of planets discovered per year grouped by size. The following are the size groupings: “small” is less than 1 Jupiter radii, “medium” is less than 2 Jupiter radii, and anything bigger is considered “large”. For example, in 2004 we discovered 2 small planets, 5 medium planets, and 0 large planets.

## Prerequisites
Before starting, ensure you have the following software installed on your computer:
* [Go](https://go.dev/) (Golang version 1.24)
* [Docker Desktop](https://www.docker.com/products/docker-desktop/)
* [Git](https://git-scm.com/downloads?form=MG0AV3)

## Project Setup

### Clone the Repository
1. Open your terminal (Command Prompt, PowerShell, or Terminal).
2. Clone the project repository using the following command:
```
git clone https://github.com/curtlengemann/exoplanet.git
```
3. Navigate to the project directory:
```
cd exoplanet
```

### Run the Project
To run the project, use the go run command:
```
go run .
```
The output will be displayed in the terminal.

### Run Tests
To run the tests, use the go test command:
```
go test ./...
```
This will run all the tests in the project.

### Building a Docker Image
Use the following command to build the Docker image:
```
docker build -t exoplanet:latest .
```

### Run the Docker Container
Use the following command to run the Docker container:
```
docker run -d exoplanet:latest
```
The container id will be output.

Since this program displays output in the terminal run the following command to see the output:
```
docker logs <container_id>
```

## Troubleshooting

### Problem: Go command not found
* Make sure Go is installed correctly. You should be able to run `go version` in your terminal and your current go version should print.
* Verify your `GOPATH` and `GOROOT` environment variables are set up correctly.
* Add your Go binary path to your system path.

### Problem: Docker build errors
* Ensure Docker Desktop is running
* Use `docker build --no-cache` to avoid using any caching



