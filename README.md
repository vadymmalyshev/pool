# Hiveon Pool

Run syncronizer

```go run ./cmd/hasbin/hasbin.go
```

Run Admin
```go run ./cmd/hasbin/hadmin.go```

TODO: Add an admin access to the user which found by email
```go run ./cmd/hasbin/hadmin.go admin add <email>```

TODO: Remove an admin access to the user which found by email
```go run ./cmd/hasbin/hadmin.go admin remove <email>```

expose remote influx ports to use it as local service
```ssh root@95.216.199.4 -L 8086:127.0.0.1:8086```
