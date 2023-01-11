# Verathread Modules SDK & Module Runner

## Building the image locally

```shell
goreleaser release --rm-dist --snapshot 
```

## Re-Tag image for local 

```shell
docker tag k3d-local-registry:5000/module-runner:v1.19.0-linux-amd64 vth-module-runner:v1.19.1-linux-amd64
```
