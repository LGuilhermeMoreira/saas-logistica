# Decisões Tomadas (ADR)

Este documento registra as decisões arquiteturais do projeto, seus motivos e trade-offs assumidos.

---

## ADR-001 — Trocar Casbin por OPA

**Contexto:** O controle de acesso (RBAC) inicialmente seria implementado com Casbin, acoplado diretamente à aplicação.

**Decisão:** Utilizar OPA (Open Policy Agent) como serviço independente, rodando como um container "sidecar" colado a cada microsserviço.

**Consequências:**
- (+) Desacopla a lógica de autorização do código da aplicação, centralizando políticas.
- (+) Permite distribuir e escalar a decisão de autorização de forma independente por serviço.
- (-) Adiciona um container extra por serviço, aumentando a complexidade de deploy e orquestração.
- (-) Introduz uma nova linguagem de políticas (Rego) que a equipe precisa aprender.

---

## ADR-002 — Deploy na AWS

**Contexto:** O sistema precisa de uma infraestrutura de deploy que suporte múltiplos microsserviços de forma escalável e gerenciada (RNF-01 — Escalabilidade).

**Decisão:** Utilizar as seguintes tecnologias AWS para o deploy:

| Tecnologia | Papel |
| :--- | :--- |
| **ECR** | Registro dos containers dos microsserviços e do gateway |
| **Fargate** | Provisionamento e gerenciamento dos containers, sem gerenciar servidores |
| **ECS** | Orquestração dos containers, rodando sobre o Fargate |
| **AWS Cloud Map** | Service discovery entre os microsserviços |

**Consequências:**
- (+) Deploy serverless de containers, sem gerenciamento de instâncias EC2.
- (+) Service discovery nativo, reduzindo a necessidade de ferramentas externas (ex: Consul).
- (-) Maior acoplamento ao ecossistema AWS (vendor lock-in).
- (-) Curva de aprendizado adicional para configuração de ECS/Fargate/Cloud Map em conjunto.

---

## ADR-003 — Trocar Postgres por MongoDB em Delivery

**Contexto:** Uma entrega pode carregar metadados variáveis conforme o tipo de carga (ex: temperatura para carga refrigerada, instruções de manuseio para carga frágil). Modelar isso em schema relacional exigiria colunas frequentemente nulas ou uma tabela EAV, aumentando a complexidade.

**Decisão:** Utilizar MongoDB no serviço de Delivery, aproveitando o schema flexível orientado a documentos para acomodar esses metadados sem alterações estruturais recorrentes.

**Consequências:**
- (+) Novos tipos de metadados podem ser adicionados sem migration.
- (+) O histórico de status pode ser embutido no próprio documento, evitando joins em consultas frequentes (RF-04).
- (-) Perde-se a consistência transacional forte do Postgres para o fluxo de status (RF-03), exigindo cuidado a nível de aplicação (ex: updates atômicos).
- (-) Introduz uma terceira tecnologia de persistência no projeto (Postgres, Cassandra, MongoDB), aumentando o custo operacional.

---

## ADR-004 — Arquitetura Orientada a Eventos para Notificações

**Contexto:** O sistema precisa notificar usuários sobre múltiplos tipos de eventos (entrega criada, status atualizado, fatura gerada, etc.). Uma abordagem acoplada, onde cada serviço chama diretamente a API de notificações, violaria o princípio de baixo acoplamento e criaria dependências frágeis (RNF-02 — Isolamento de Falhas).

**Decisão:** Implementar um microsserviço independente de Notificações que consome eventos publicados por outros serviços através de um message broker (AWS SQS / SNS ou RabbitMQ). A comunicação será baseada em eventos (Event-Driven Architecture).

**Arquitetura:**
1. Serviços publicam eventos em tópicos (ex: `delivery.events`, `invoice.events`) no message broker.
2. O serviço de Notificações consome esses eventos de forma assíncrona.
3. Cada evento dispara a geração de uma notificação, respeitando preferências do usuário.
4. Notificações são armazenadas em banco de dados próprio (Postgres) e enfileiradas para múltiplos canais (E-mail, SMS, Push, In-app).
5. Adaptadores de canal executam em paralelo, permitindo falha de um canal sem afetar outros.

**Padrões Utilizados:**
- **Event Sourcing:** Eventos são imutáveis e servem como fonte de verdade do estado de notificações.
- **CQRS:** Separação entre escrita (evento consumido) e leitura (consultar histórico).
- **Dead Letter Queue (DLQ):** Eventos não processáveis são desviados para análise posterior.
- **Idempotência via Deduplicação:** Cada evento possui um `event_id` único para evitar duplicação ao reenviar.

**Consequências:**
- (+) Desacopla completamente o serviço de Notificações dos demais serviços, permitindo deploy e escala independentes.
- (+) Falhas no envio de notificações não impactam o fluxo de negócio (entregas, faturas, etc.).
- (+) Fácil adicionar novos tipos de eventos ou canais de notificação sem modificar código existente.
- (+) Permite replay de eventos para reenvio de notificações falhadas.
- (-) Introduz complexidade operacional (message broker, DLQ, monitoramento de latência de processamento).
- (-) Notificações são entregar de forma assíncrona, não síncrona (latência de alguns segundos é esperada).
- (-) Exige implementação de retry logic, idempotência e garantias de entrega (At-Least-Once).

**Message Broker Recomendado:**
- **AWS SQS + SNS** (alinhado com ADR-002): SNS publica eventos, SQS consome. Simplicidade operacional, integração nativa com ECS/Fargate.
- Alternativa: **RabbitMQ** em container, oferecendo mais controle e features como priority queues, mas maior overhead operacional.

---

## ADR-005 — Persistência de Notificações com Postgres

**Contexto:** O serviço de Notificações precisa manter registro de todas as notificações enviadas para auditoria (RNF-13) e permitir queries complexas (histórico, filtros, estatísticas).

**Decisão:** Usar Postgres como banco de dados principal do serviço de Notificações, mantendo schema relacional simples com tabelas para: `notifications`, `notification_events`, `user_preferences`, `notification_attempts`.

**Schema Simplificado:**
```sql
-- Armazena notificações e seu estado
CREATE TABLE notifications (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  event_id VARCHAR(255) UNIQUE,
  event_type VARCHAR(100),
  title VARCHAR(255),
  body TEXT,
  status ENUM('PENDING', 'SENT', 'FAILED', 'READ'),
  channel VARCHAR(50),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Armazena preferências do usuário
CREATE TABLE user_preferences (
  user_id UUID PRIMARY KEY,
  email_enabled BOOLEAN DEFAULT TRUE,
  sms_enabled BOOLEAN DEFAULT FALSE,
  push_enabled BOOLEAN DEFAULT TRUE,
  inapp_enabled BOOLEAN DEFAULT TRUE,
  quiet_hours_start TIME,
  quiet_hours_end TIME,
  updated_at TIMESTAMP
);

-- Registra tentativas de envio (para retry logic)
CREATE TABLE notification_attempts (
  id UUID PRIMARY KEY,
  notification_id UUID REFERENCES notifications(id),
  attempt_number INT,
  status VARCHAR(50),
  error_message TEXT,
  created_at TIMESTAMP
);
```

**Consequências:**
- (+) Queries simples e eficientes para histórico e auditoria.
- (+) Índices em `user_id` e `event_id` garantem performance mesmo com alto volume de notificações.
- (+) Postgres suporta JSONB para armazenar payload flexível de cada evento.
- (-) Replicação entre serviços (Postgres para Delivery/Auth, MongoDB para Delivery, Cassandra para Telemetria, agora Postgres novamente para Notificações) aumenta complexidade operacional.
- (-) Não é o melhor para time-series de eventos de alta frequência (melhor seria Cassandra, mas usamos para Telemetria apenas).

**Alternativa Considerada:**
- Usar apenas Redis em memória com TTL: mais rápido, mas perde persistência para auditoria. Rejeitada por violar RNF-13.