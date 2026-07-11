# TravelGo — Plataforma de Reservas com Saga Distribuída (AWS Native)

## Contexto Empresarial Realista
Uma plataforma de viagens (similar a Decolar ou Booking) precisa processar reservas de **pacotes combinados** (Voo + Hotel + Pagamento). O sistema deve ser **altamente consistente** (se o voo falhar, o hotel deve ser cancelado) e **escalável** para picos de Black Friday. Operamos inteiramente na AWS, utilizando serviços gerenciados para minimizar custo operacional.

### Stack Completa

| Camada | Tecnologia | Justificativa Técnica |
| :--- | :--- | :--- |
| **API Gateway** | AWS API Gateway (HTTP API) | Throttling, CORS, integração nativa com VPC-Links para ECS. |
| **Compute** | AWS ECS Fargate (Go 1.23) | Serverless containers. Permite rodar Go com binários estáticos, sem cold-start do Lambda. |
| **Orquestração** | AWS SQS (FIFO) + Go Worker | Desacopla os passos do Saga. Fila FIFO garante ordenação e exatamente-uma-entrega para `order_id`. |
| **Database (Transactional)** | AWS RDS Aurora (PostgreSQL) + sqlc | Consistência ACID para a reserva principal e status final. |
| **Database (State)** | AWS DynamoDB (PAY_PER_REQUEST) | Armazena o estado do Saga (`PENDING`, `FLIGHT_OK`, `HOTEL_OK`, `COMPLETED`, `ROLLBACK`). Baixa latência e TTL nativo. |
| **Resiliência** | `sony/gobreaker` + `golang.org/x/sync/semaphore` | Circuit Breaker nas chamadas HTTP para parceiros externos. Bulkhead para limitar concorrência por serviço. |
| **Observabilidade** | OpenTelemetry (OTel) + AWS Distro + X-Ray | Correlação automática entre traces HTTP → SQS → DB → HTTP externo. Logs estruturados com `trace_id`. |
| **IaC** | Terraform (v1.5+) | Módulos separados: Networking (VPC/Subnets), RDS, DynamoDB, SQS, ECS (Fargate), IAM (Least Privilege). |
| **CI/CD** | GitHub Actions | `golangci-lint`, `go test -race`, build de imagem, push para ECR, `terraform apply` no stage/prod. |

---

**Fluxo de Dados:**
1. **API** recebe `{user_id, flight_id, hotel_id, payment_card}` e gera `order_id`. Salva `PENDING` no RDS. Publica mensagem na SQS `saga-start.fifo` (com `MessageGroupId: order_id`).
2. **Orquestrador (Worker)** consome a fila. Para cada passo, publica nas filas respectivas e espera a atualização de estado no DynamoDB (via padrão de *polling* ou *callback*).
3. **Workers isolados** executam as chamadas HTTP para parceiros. Cada um tem **Circuit Breaker** (se parceiro está lento/falhando, abre e falha rápido).
4. **Atualização de Estado**: Cada worker escreve no DynamoDB (`PK: order_id, SK: step`).
5. **Compensação (Saga)**:
   - Se todos os passos OK → Orquestrador atualiza RDS para `CONFIRMED`.
   - Se algum passo falha (ex: Hotel indisponível) → Orquestrador publica mensagens de **rollback** nas filas (Flight cancela, Payment estorna) e atualiza RDS para `FAILED`.
6. **Notificação** é enviada ao final.

---

### Roadmap em Cards (Tasks Detalhadas)

#### 📦 Card 1: Fundação AWS e Bootstrap
- [ ] **T1:** Escrever módulos Terraform para VPC, Subnets (públicas/privadas), Security Groups, Internet Gateway, NAT Gateway (para workers acessarem APIs externas).
- [ ] **T2:** Escrever módulo Terraform para RDS Aurora (PostgreSQL 16) com `db.t4g.small` (dev), backup automático e Multi-AZ (desligado em dev, ligado em prod via variável).
- [ ] **T3:** Escrever módulo Terraform para DynamoDB (PAY_PER_REQUEST) com chave composta `order_id` + `step`.
- [ ] **T4:** Criar SQS FIFO queues: `saga-start.fifo`, `flight-book.fifo`, `hotel-book.fifo`, `payment-process.fifo`, `rollback.fifo`. Configurar **Dead-Letter Queue (DLQ)** para cada uma com `maxReceiveCount = 3`.

#### 📦 Card 2: Core Microservices em Go (Estrutura Hexagonal)
- [ ] **T5:** Implementar **Order Service** (API Gateway → ECS): Handler HTTP (`POST /orders`), validação, geração de UUID, insert no RDS (sqlc), publish no SQS (AWS SDK v2).
- [ ] **T6:** Implementar **Saga Orchestrator** (Worker): Consome `saga-start.fifo`. Carrega o estado inicial do DynamoDB. Publica mensagens em paralelo nas filas de Flight/Hotel/Payment.
- [ ] **T7:** Implementar **Flight/Hotel/Payment Workers**: Cada um consome sua fila, chama um `mock.Client` (com `http.Client` configurado para timeout de 5s), atualiza o estado no DynamoDB com TTL (se falhar, insere estado com erro).

#### 📦 Card 3: Resiliência (Preenche o GAP principal)
- [ ] **T8:** Adicionar **Circuit Breaker** (`sony/gobreaker`) em cada `mock.Client` (Flight, Hotel, Payment). Configurar: `MaxRequests = 3`, `Interval = 5s`, `Timeout = 30s` (meio-aberto). Logar cada transição de estado (Closed → Open → Half-Open).
- [ ] **T9:** Adicionar **Bulkhead** nos workers: usar `golang.org/x/sync/semaphore` (Weighted) limitando a 20 chamadas HTTP concorrentes por worker. Protege a aplicação de sobrecarga de dependências externas.
- [ ] **T10:** Implementar **Retry com Backoff** para filas SQS (já nativo via `visibility_timeout` nas DLQs). Para chamadas HTTP, implementar retry exponencial customizado (1s, 2s, 4s) **antes** do Circuit Breaker abrir.

#### 📦 Card 4: Observabilidade Total com OTel + X-Ray (Preenche o GAP do Tracing)
- [ ] **T11:** Instrumentar cada serviço com **OpenTelemetry Go SDK**. Adicionar:
  - Tracer para HTTP Server (API Gateway).
  - Tracer para SQS Producer/Consumer (criar spans para `SendMessage` e `ReceiveMessage`).
  - Tracer para HTTP Client (chamadas externas).
- [ ] **T12:** Configurar **AWS X-Ray Exporter** (ou OTel Collector sidecar no ECS). Garantir que `trace_id` seja injetado nos logs estruturados (Zap com Field `trace_id`).
- [ ] **T13:** Criar **Custom Metrics** no Prometheus (scrapado via AWS Managed Prometheus ou Self-hosted): `saga_duration_seconds`, `external_call_duration_seconds`, `circuit_breaker_state`.
- [ ] **T14:** Provisionar **Grafana** (via Terraform) e criar dashboard com:
  - Diagrama do Saga em tempo real (número de pedidos por estado).
  - Latência p95 do Circuit Breaker + Taxa de abertura.
  - Fila SQS (ApproximateNumberOfMessagesVisible).

#### 📦 Card 5: CI/CD e Deploy (Fargate)
- [ ] **T15:** Escrever `Dockerfile` multi-stage para cada serviço (`golang:1.23-alpine` → `alpine:latest`). Build com `CGO_ENABLED=0`.
- [ ] **T16:** Escrever módulo Terraform para **ECS Fargate** (Cluster, Task Definitions com sidecar para OTel Collector, Services com load balancer interno).
- [ ] **T17:** Configurar **GitHub Actions**:
  - Job 1: Lint (`golangci-lint`).
  - Job 2: Test (`go test -race -cover ./...`).
  - Job 3: Build & Push Imagens para ECR.
  - Job 4: Terraform Apply (usando `aws s3` backend para remote state).

#### 📦 Card 6: Documentação e Apresentação
- [ ] **T18:** Escrever `README.md` com diagrama Mermaid, pré-requisitos, e **"Run Locally"** usando `docker-compose` para simular localmente (Postgres + DynamoDB Local + SQS Localstack).
- [ ] **T19:** Escrever `ARCHITECTURE.md` detalhando as **decisões técnicas**:
  - *Por que Saga coreografada vs orquestrada?* (Orquestrada mantém a lógica centralizada e facilita rollback).
  - *Por que OTel em vez de X-Ray SDK puro?* (Portabilidade e padronização OTLP).
  - *Por que RDS + DynamoDB juntos?* (RDS para transações ACID do pedido; DynamoDB para alta throughput de estado do Saga, evitando locking).
- [ ] **T20:** Gravar demo de 5 minutos: Mostrar um pedido bem-sucedido, um pedido com falha (Hotel rollback) e o dashboard do Grafana com os traces do X-Ray.
