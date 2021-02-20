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
go run ./cmd/*.go
```

## WSL2 guide
- make sure you have setup docker desktop and connected to wsl2
- for connecting with wls2's grpc port  from windows
    - run `wsl hostname -I` on windows cmd to find wsl2 ip address


## Direct connection to PostgresSQL

```
psql -U username -h localhost -p 5432 dbname
```