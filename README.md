# Encurtador de URL

Um encurtador de URL é um serviço que transforma endereços longos em códigos curtos. Quando o código é acessado, o usuário é redirecionado para o endereço original. Esse mecanismo facilita o compartilhamento de links, reduz erros de digitação e pode ser usado para monitorar o número de acessos.

Este projeto implementa um encurtador simples em Go. Ele oferece um endpoint `POST /api/shorten` que recebe um JSON com o campo `url` e retorna um código gerado aleatoriamente. Ao acessar `/{code}`, o servidor consulta o código armazenado em memória e redireciona para a URL associada.

## Executando o projeto

```bash
go run main.go
```

## Proximas Features

Integrar a um banco de dados noSQL
