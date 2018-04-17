# Note you cannot run golang binaries on Alpine directly
FROM            debian:buster-slim

MAINTAINER      chris.mague@shokunin.co

COPY            consul-vault-exporter /consul-vault-exporter

WORKDIR		/
ENV		GIN_MODE=release

EXPOSE          8080

ENTRYPOINT      [ "/consul-vault-exporter" ]
