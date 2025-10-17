## Client-Server API (Go) — Cotação USD/BRL

### O que é
Aplicação simples cliente-servidor escrita em Go para consultar a cotação do dólar (USD → BRL), persistir as consultas em SQLite e salvar o último valor de `bid` em arquivo de texto no cliente.

### Requisitos
- Go 1.20+
- Acesso à internet (para consumir `https://economia.awesomeapi.com.br`)
- Ambiente com CGO habilitado (driver `github.com/mattn/go-sqlite3`); em Linux/macOS, é recomendável ter um compilador C instalado (ex.: `build-essential` no Ubuntu)

### Estrutura de pastas
```
client-server-api/
  client/
    main.go        # cliente HTTP: consome o server e grava cotacao.txt
    cotacao.txt    # gerado pelo cliente com o valor do bid
  server/
    main.go        # servidor HTTP: rota /cotacao
    db             # arquivo SQLite criado automaticamente (nome de arquivo "db")
```

### Como funciona
- O servidor (`server/main.go`):
  - Expõe `GET /cotacao` em `:8080`.
  - Busca a cotação USD/BRL em `https://economia.awesomeapi.com.br/json/last/USD-BRL` (timeout de 200ms).
  - Persiste a cotação na tabela `quotations` em SQLite (timeout de 10ms para a inserção).
  - Retorna JSON ao cliente com os campos da cotação.
- O cliente (`client/main.go`):
  - Faz `GET http://localhost:8080/cotacao` (timeout de 300ms).
  - Salva o valor de `bid` no arquivo `cotacao.txt` no diretório `client`.

### Executando em desenvolvimento
1. Inicie o servidor (terminal 1):
```bash
cd server
go run main.go
```

2. Em outro terminal, execute o cliente (terminal 2):
```bash
cd client
go run main.go
```

3. (Opcional) Teste o endpoint diretamente via cURL:
```bash
curl http://localhost:8080/cotacao
```

Após executar o cliente, verifique o arquivo `client/cotacao.txt` com o conteúdo no formato `Dólar: <bid>`.

### Endpoint
- `GET /cotacao`
  - Resposta 200 (exemplo simplificado):
    ```json
    {
      "code": "USD",
      "codein": "BRL",
      "name": "Dólar Americano/Real Brasileiro",
      "high": "5.10",
      "low": "5.00",
      "varBid": "0.01",
      "pctChange": "0.20",
      "bid": "5.05",
      "ask": "5.06",
      "timestamp": "1729190000",
      "create_date": "2025-10-17 12:00:00"
    }
    ```

### Banco de dados
- Driver: `github.com/mattn/go-sqlite3`.
- O arquivo do banco é criado no diretório `server` com nome `db` (sem extensão).
- A tabela `quotations` é criada automaticamente na inicialização, caso não exista.
- Cada chamada a `/cotacao` insere um novo registro.

### Build e execução dos binários
- Server:
  ```bash
  cd server
  go build -o server
  ./server
  ```
- Client:
  ```bash
  cd client
  go build -o client
  ./client
  ```

### Timeouts e mensagens
- Server:
  - Chamada à API externa: 200ms. Em estouro, loga aviso e retorna erro 500.
  - Inserção no SQLite: 10ms. Em estouro, loga aviso e retorna erro 500.
- Client:
  - Requisição ao servidor: 300ms. Em estouro, loga aviso e encerra.

### Dicas e erros comuns
- Certifique-se de iniciar o servidor antes do cliente.
- Se houver erro relacionado ao SQLite, verifique a presença de um compilador C (CGO).
- O arquivo do banco (`server/db`) e o `client/cotacao.txt` são criados automaticamente nas respectivas pastas.

---
Se algo não funcionar como esperado, verifique sua versão do Go (`go version`) e se você possui conectividade com a internet para acessar a API externa.


