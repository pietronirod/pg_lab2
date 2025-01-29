# Sistema de Temperatura por CEP

## Objetivo

Desenvolver um sistema em Go que receba um CEP, identifique a cidade e retorne o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin) juntamente com a cidade. Este sistema implementa OpenTelemetry (OTEL) e Zipkin para rastreamento distribuído.

## Arquitetura

O sistema é composto por dois serviços:

- **Serviço A**: Responsável por receber o input do usuário.
- **Serviço B**: Responsável pela orquestração, realizando a pesquisa do CEP e retornando as temperaturas.

## Requisitos

### Serviço A (responsável pelo input)

- Receber um input de 8 dígitos via POST, através do schema: `{ "cep": "29902555" }`
- Validar se o input é válido (contém 8 dígitos) e é uma STRING
- Caso seja válido, encaminhar para o Serviço B via HTTP
- Caso não seja válido, retornar:
  - Código HTTP: 422
  - Mensagem: `invalid zipcode`

### Serviço B (responsável pela orquestração)

- Receber um CEP válido de 8 dígitos
- Realizar a pesquisa do CEP e encontrar o nome da localização
- Retornar as temperaturas em Celsius, Fahrenheit e Kelvin juntamente com o nome da localização
- Responder adequadamente nos seguintes cenários:
  - Em caso de sucesso:
    - Código HTTP: 200
    - Response Body: `{ "city": "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }`
  - Em caso de falha, caso o CEP não seja válido (com formato correto):
    - Código HTTP: 422
    - Mensagem: `invalid zipcode`
  - Em caso de falha, caso o CEP não seja encontrado:
    - Código HTTP: 404
    - Mensagem: `can not find zipcode`

## Implementação de OTEL e Zipkin

- Implementar tracing distribuído entre Serviço A e Serviço B
- Utilizar spans para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura

## APIs Utilizadas

- [viaCEP](https://viacep.com.br/): Para encontrar a localização a partir do CEP
- [WeatherAPI](https://www.weatherapi.com/): Para consultar as temperaturas

## Fórmulas de Conversão

- Celsius para Fahrenheit: `F = C * 1.8 + 32`
- Celsius para Kelvin: `K = C + 273`

## Como Rodar o Projeto

### Pré-requisitos

- Docker
- Docker Compose

### Passos

1. Clone o repositório:

   ```sh
   git clone https://github.com/seu-usuario/seu-repositorio.git
   cd seu-repositorio
   ```

2.Configure as variáveis de ambiente no arquivo `.env`:

   ```env
   VIACEP_API_URL=https://viacep.com.br/ws/
   WEATHERAPI_URL=http://api.weatherapi.com/v1/current.json
   WEATHERAPI_KEY=your_weatherapi_key
   OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
   OTEL_EXPORTER_OTLP_PROTOCOL=grpc
   ````

3.Inicie os serviços com Docker Compose:

   ```sh
   docker-compose up --build
   ```

4.Acesse o Serviço A na porta 8080 e o Serviço B na porta 8090.

### Endpoints

#### Serviço A

- **POST /cep**
  - Request Body: `{ "cep": "29902555" }`
  - Response: Encaminha a requisição para o Serviço B

#### Serviço B

- **GET /cep/{cep}**
  - Response: `{ "city": "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }`

## Acessando e Visualizando os Logs no Zipkin

1. Certifique-se de que o Zipkin está rodando. O Zipkin é iniciado automaticamente com o Docker Compose.

2. Acesse a interface do Zipkin no navegador: <http://localhost:9411>

3. Utilize a interface do Zipkin para visualizar os traces e spans gerados pelos serviços A e B.

## Referências

- [OpenTelemetry](https://opentelemetry.io/)
- [Zipkin](https://zipkin.io/)
- [viaCEP](https://viacep.com.br/)
- [WeatherAPI](https://www.weatherapi.com/)
