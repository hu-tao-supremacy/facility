# Facility Service

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
go run main.go
```