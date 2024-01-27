# Desafio Go: Cotação do Dólar 💰

**Objetivo:**
Desenvolver dois sistemas em Go - `client.go` e `server.go` - que interagem para obter e salvar a cotação do dólar.

**Requisitos:**
1. **client.go:**
    - Realizar uma requisição HTTP no `server.go` para obter a cotação do dólar.
    - Utilizar o package "context" com timeout máximo de 300ms.
    - Salvar o valor da cotação em um arquivo "cotacao.txt" no formato: Dólar: {valor}.

2. **server.go:**
    - Consumir a API de câmbio Dólar-Real em https://economia.awesomeapi.com.br/json/last/USD-BRL.
    - Retornar para o cliente, em formato JSON, o resultado da cotação (campo "bid").
    - Usar "context" para registrar no banco SQLite com timeout máximo de 10ms para persistir os dados e 200ms para chamar a API.
    - O endpoint necessário será /cotacao e a porta do servidor HTTP será 8080.

