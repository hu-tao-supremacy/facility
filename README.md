# Facility Service
[![Go Report Card](https://goreportcard.com/badge/github.com/hu-tao-supremacy/facility)](https://goreportcard.com/report/github.com/hu-tao-supremacy/facility)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Test](https://github.com/hu-tao-supremacy/facility/actions/workflows/test.yml/badge.svg)
![Code Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/new5558/4c5f04edd09de877e2792257f7c98bba/raw/badge.json)

1. Setting up PostgreSQL for dev environment
```
docker-compose -f docker-compose.dev.yaml up -d
```
2. Run Makefile for cloning service proto and symlink
```
make apis
```
3. Run migration for database
```
git clone https://github.com/hu-tao-supremacy/migrations/
cd migrations
yarn
make migrate
```
4. Prepare Go's env
```
source dev-env
```
5. Code
```
code .
```
6. Run
```
go run ./cmd/!(*_test).go
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

## Testing
- Coverage
```
go test  -cover  $(go list ./... | grep -v hts) -coverprofile=coverage.out
```
- View test coverage in terminal
```
go tool cover -func=coverage.out
```
- View test coverage in html format
```
go tool cover -html=coverage.out -o coverage.html
```

## WSL2 guide
- make sure you have setup Docker Desktop and connected to WSL2
- for connecting with WSL2's GRPC port from windows
    - run `wsl hostname -I` on windows cmd to find WSL2 ip address


## Direct connection to PosgresSQL

```
psql -U username -h localhost -p 5432 dbname
```
