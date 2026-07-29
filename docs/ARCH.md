# Documento de Arquitetura

## Visão Arquitetural

O sistema será composto por quatro microsserviços independentes:

1. **Serviço de Entregas**
2. **Serviço de Telemetria**
3. **Serviço de Faturamento**
4. **Serviço de Autenticação e Segurança**

Cada serviço é responsável por seu próprio domínio de negócio.

---

## Microsserviços

### 2.1 Serviço de Entregas

**Responsabilidades:**
- Gerenciamento de entregas;
- Associação de motoristas;
- Atualização de status;
- Consulta de entregas.

**Banco de dados próprio:** `deliveries_db`

---

### 2.2 Serviço de Telemetria

**Responsabilidades:**
- Receber coordenadas dos motoristas;
- Armazenar histórico de localização;
- Disponibilizar posição atual;
- Disponibilizar histórico de rotas.

**Banco de dados próprio:** `telemetry_db`

---

### 2.3 Serviço de Faturamento

**Responsabilidades:**
- Receber notificações de entregas concluídas;
- Calcular fretes;
- Gerar faturas;
- Gerenciar pagamentos.

**Banco de dados próprio:** `billing_db`

---

### 2.4 Serviço de Autenticação e Segurança

**Responsabilidades:**
- Cadastro e gerenciamento de usuários;
- Autenticação de usuários (Login);
- Emissão e validação de tokens JWT;
- Gestão de permissões e perfis (RBAC).

**Banco de dados próprio:** `auth_db`

---

## Persistência

### Banco de Dados por Serviço

Cada microsserviço possui banco de dados isolado.

---

## Comunicação entre Serviços

### Comunicação Síncrona

Para consultas em tempo real:

**Serviço de Entregas → Serviço de Telemetria**
- Exemplo: `GET /drivers/{id}/current-location`

Utilizado para obter a última localização conhecida durante a consulta de uma entrega.

### Comunicação Assíncrona

Para eventos de negócio, como o evento **DeliveryCompleted**, publicado pelo Serviço de Entregas quando uma entrega é concluída.

**Exemplo de payload:**
```json
{
  "delivery_id": "123",
  "driver_id": "456",
  "distance_km": 15.7,
  "weight_kg": 12.5,
  "completed_at": "2026-06-15T10:00:00Z"
}
```

Consumido pelo Serviço de Faturamento para geração automática da fatura.

---

## Mensageria

A integração assíncrona deve utilizar um broker de mensagens.

**Fluxo:**

```
Entrega Concluída
       │
       ▼
  DeliveryCompleted
       │
       ▼
Message Broker
       │
       ▼
Serviço de Faturamento
       │
       ▼
  Geração da Fatura
```

---

## Resiliência

### Serviço de Telemetria

Por possuir alto volume de escrita:
- Deve ser escalável horizontalmente;
- Deve suportar picos de requisições;
- Lentidão ou indisponibilidade do banco de telemetria não deve afetar diretamente os demais serviços.

**Estratégias sugeridas:**
- Filas de buffer;
- Retry;
- Circuit Breaker;
- Cache para consultas frequentes.

---

## API Gateway (Opcional)

Um API Gateway pode atuar como ponto único de entrada.

**Exemplo:**
```text
api.empresa.com/
├── deliveries/
├── telemetry/
├── billing/
└── auth/
```

**Responsabilidades:**
- Roteamento;
- Autenticação;
- Rate Limiting;
- Observabilidade.