FROM golang:1.23rc1-alpine3.20

# Instalar Bash
RUN apk update && apk add bash

WORKDIR /go/src

# Definir Bash como ponto de entrada
ENTRYPOINT ["tail", "-f", "/dev/null"]