# Facility Service

1. Setting up PostgreSQL for dev environment
```
docker-compose -f docker-compose.dev up -d
```
2. Run Makefile for cloning service proto and symlink
```
make apis
```
3. Code
```
code .
```
4. Run
```
go run main.go
```