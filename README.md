# Readme

- para rodar os testes da aplicação, entre em cada service e rode:

  ```bash
    go test ./...
  ```

- para rodar a aplicação localmente, service A na porta 8080 e service B na porta 8081

  ```bash
    docker compose up -d
  ```

- para executar a request basta enviar um POST para <http://locahost:8080>

 ```bash
  curl  -X POST \
  'http://localhost:8080' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "cep": "60541646"
  }'
 ```

- para acessar o zpkin:
  - <http://localhost:9411/>
