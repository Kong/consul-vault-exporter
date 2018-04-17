# consul-vault-exporter

Uses Consul to find all vault servers and exports data

## Building

### Mac/Linux

0) set GOROOT environment variable
1) Install Go and Make
2) make

### Docker

0) set GOROOT environment variable
1) Install Docker, Go and Make
2) make docker


## Running

### Mac/Linux

```
./consul-vault-exporter
```

### Docker

```
docker pull maguec/consul-vault-exporter:latest
docker run -i -t -p 8080:8080 maguec/consul-vault-exporter
```

## Testing

run either the docker container or the raw application binary

```
curl http://localhost:8080/health
```

---
Copyright © 2018, Chris Mague
