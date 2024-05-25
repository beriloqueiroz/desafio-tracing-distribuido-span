# OpenTelemetry + zipkin with clean architecture

- service A -> service B -> external http apis -> service A  

- to run tests

  ```bash
    go test ./...
  ```

- to run applications, service A: 8080 and service B: 8081

  ```bash
    docker compose up -d
  ```

- to run request: POST to <http://locahost:8080>

 ```bash
  curl  -X POST \
  'http://localhost:8080/get-cep' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "cep": "01550020"
  }'
 ```

- to see zpkin:
  - <http://localhost:9411/>
