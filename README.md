#pr-gobs

Projeto de observabilidade utilizando
- Prometheus;
- OpenTelemetry;
- Grafana;
- Jaeger e 
- Zipkin.

## Execução
Na raiz do projeto, execute o comando:
```bash
docker-compose up -d
```
e veja a mágica acontecer.

## Acesso
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
  - _Usuário_: admin
  - _Senha_: admin
- **Jaeger**: http://localhost:16686
- **Zipkin**: http://localhost:9411
- **Aplicação** weather: http://localhost:8080
- **Aplicação** reader: http://localhost:8081

## Endpoints de teste
### weather
```bash
curl http://localhost:8080/temperature/{cep}
```
Substitua `{cep}` pelo CEP da cidade que deseja consultar.
### reader
```bash
curl http://localhost:8080/temperature/{cep}
```
Substitua `{cep}` pelo CEP da cidade que deseja consultar.

> [!WARNING]
> Ambas aplicações funcionam com mesmo endpoint, atente-se para amudança de porta.

## Descrição
O projeto é composto por duas aplicações:
- **weather**: aplicação que retorna a temperatura de uma cidade a partir do CEP informado.
- **reader**: aplicação que consome a aplicação weather e retorna a temperatura da cidade.
