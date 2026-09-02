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