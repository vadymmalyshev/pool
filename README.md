# Hiveon Pool

Run syncronizer

```go run ./cmd/hasbin/hasbin.go
```

Run Admin:

fill dbs:
```go run ./cmd/hadmin/hadmin.go migrate```
run admin:
```go run ./cmd/hadmin/hadmin.go```

#Use custom config file:
use flags to set config file name from ./config directory by default used config.yaml
note: use configName without ".yaml"
-c=configName      OR:   --c=configName
-c configName            --c configName

Add an admin access to the user which found by email
```go run ./cmd/hadmin/hadmin.go admin add <email>```

Remove an admin access to the user which found by email
```go run ./cmd/hadmin/hadmin.go admin remove <email>```

expose remote influx ports to use it as local service
```ssh root@95.216.199.4 -L 8086:127.0.0.1:8086```


For Debian and Ubuntu based distros, install librdkafka-dev from the standard repositories or using Confluent's Deb repository (https://docs.confluent.io/current/installation/installing_cp/index.html#rpm-packages-via-yum).


For MAC
`brew install librdkafka pkg-config`