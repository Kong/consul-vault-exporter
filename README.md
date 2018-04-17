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

Takes the environment variable: CONSUL_ADDRESS which is set to consul:8500 by default

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

## Outputs
```
# HELP is vault initialized?
# TYPE vault_initialized gauge
vault_initialized{instance="vault-8665717464.consul:8200",cluster="kongvault",version="0.9.5"} 1
# HELP is vault sealed?
# TYPE vault_sealed gauge
vault_sealed{instance="vault-8665717464.consul:8200",cluster="kongvault",version="0.9.5"} 0
# HELP is vault standby?
# TYPE vault_standby gauge
vault_standby{instance="vault-8665717464.consul:8200",cluster="kongvault",version="0.9.5"} 1
```
