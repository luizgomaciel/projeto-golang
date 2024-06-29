# [POC] Serviço de geração de pool de massa do TDM em `Golang`

## Configurando ambiente

Para rodar a app siga os seguintes passos:

* Execute o `docker-compose up -d`
* Executar o comando `docker-compose exec app bash` para entrar no console do container
* Dentro do container e na raiz do projeto, executar o comando `go mod init encoder`
* Executar o comando `go mod tidy`
* Executar o comando `go run server.go`
* Acesse a administração do rabbitmq com `http://localhost:15672/` (Usuário e senha estão no arquivo docker-compose)
* No rabbit, na queue account-queue criada na inicialização, faça o bind da exchange amq.direct com Routing key jobs
* Escreva a mensagem na exchange

## Executando

* No payload da exchange `amq.direct`, publique a seguinte mensagem modelo com Routing key `jobs`
```
{
   "typeAccount":2,
   "qtd":10,
   "products":[
      "LOAN_PRODUCT"
   ]
}
```

## Validando a execução

* Acesse a administração do Mongo Express com `http://localhost:8081/`
`OBS.: Se a execução não criar as collections via auto migrate, criar o database tdm e dentro do database as collections accounts e loans` 
