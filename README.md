# Facility Service
[![Go Report Card](https://goreportcard.com/badge/github.com/hu-tao-supremacy/facility)](https://goreportcard.com/report/github.com/hu-tao-supremacy/facility)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

1. Setting up PostgreSQL for dev environment
```
docker-compose -f docker-compose.dev.yaml up -d
```
2. Run Makefile for cloning service proto and symlink
```
make apis
```
3. Prepare Go's env
```
source dev-env
```
4. Code
```
code .
```
5. Run
```
go run ./cmd/*.go
```

## Build binary file
1. Run go build command
```
go build -o main ./cmd/*.go
```
2. Execute binary file
```
./main
```

## WSL2 guide
- make sure you have setup Docker Desktop and connected to WSL2
- for connecting with WSL2's GRPC port from windows
    - run `wsl hostname -I` on windows cmd to find WSL2 ip address


## Direct connection to PosgresSQL

```
psql -U username -h localhost -p 5432 dbname
```
