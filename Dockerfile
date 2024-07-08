FROM golang:1.23rc1-alpine3.20

# Instalar Bash e Git
RUN apk update && apk add bash git

# Instalar Protocol Buffer Compiler (protoc) 
RUN apk add protobuf

# Verificar a instalação
RUN protoc --version

WORKDIR /go/src

# Instalando os plugins do protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN go install github.com/ktr0731/evans@latest

# Definir Bash como ponto de entrada
ENTRYPOINT ["tail", "-f", "/dev/null"]