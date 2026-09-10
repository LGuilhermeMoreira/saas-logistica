# Documento de Requisitos

## Visão Geral
O sistema tem como objetivo gerenciar entregas, rastreamento de motoristas e faturamento de fretes através de uma arquitetura baseada em microsserviços. Adicionalmente, o sistema contará com interfaces frontend (Web e Mobile) para o consumo das APIs pelos usuários.

## Requisitos Funcionais

### Módulo de Autenticação e Segurança
*   **RF-11 – Cadastro de Usuários:** O sistema deve permitir o cadastro de novos usuários, diferenciando perfis (ex: Administrador, Motorista, Cliente).
*   **RF-12 – Autenticação de Usuários:** O sistema deve permitir que usuários realizem login utilizando e-mail e senha.
*   **RF-13 – Geração de Token de Acesso:** Após o login bem-sucedido, o sistema deve emitir um token seguro (ex: JWT) para autenticar requisições subsequentes.
*   **RF-14 – Controle de Acesso (RBAC):** O sistema deve restringir o acesso às funcionalidades com base no perfil do usuário (ex: apenas motoristas podem atualizar status de entrega). O frontend deve refletir essas restrições bloqueando rotas não autorizadas.

### Módulo de Entregas
*   **RF-01 – Criar Entrega:** O sistema deve permitir que um cliente cadastre uma nova entrega informando obrigatoriamente: Endereço de origem, Endereço de destino e Peso do pacote. Opcionalmente, a entrega poderá conter metadados adicionais.
*   **RF-02 – Atribuir Motorista:** O sistema deve permitir associar uma entrega a um motorista específico. O sistema operará com um modelo de *Marketplace*, onde motoristas escolhem e aceitam entregas pendentes através do aplicativo.
*   **RF-03 – Atualizar Status da Entrega:** O motorista deve poder atualizar o status da entrega seguindo o fluxo: `CRIADA` → `EM_TRANSITO` → `CONCLUIDA` ou `FALHA`.
*   **RF-04 – Consultar Entrega:** O sistema deve permitir consultar os detalhes de uma entrega. Quando a entrega estiver com status `EM_TRANSITO`, a consulta deverá apresentar também a última localização conhecida do motorista responsável (via integração com Telemetria).
*   **RF-15 – Notificar Conclusão de Entrega por E-mail:** Quando uma entrega atingir o status CONCLUIDA, o sistema deve enviar automaticamente um e-mail de notificação ao cliente responsável pela entrega, contendo: Identificador da entrega, Endereço de origem e destino, e Data e hora de conclusão.

### Módulo de Telemetria
*   **RF-05 – Registrar Localização do Motorista:** O sistema deve permitir registrar periodicamente a localização de um motorista contendo: Identificador do motorista, Latitude e Longitude, e Timestamp. No app, isso deve rodar como um processo em background.
*   **RF-06 – Consultar Posição Atual:** O sistema deve permitir consultar a última posição registrada de um motorista. O painel web deverá refletir isso em tempo real no mapa.
*   **RF-07 – Consultar Histórico de Rota:** O sistema deve permitir consultar o histórico de coordenadas de um motorista dentro de um intervalo de tempo informado.

### Módulo de Faturamento
*   **RF-08 – Gerar Fatura Automaticamente:** Quando uma entrega for concluída, o sistema deverá gerar uma fatura correspondente ao frete realizado.
*   **RF-09 – Calcular Valor do Frete:** O valor da fatura deve ser calculado utilizando a fórmula: Valor = Taxa Base + (Valor por Km × Distância) + (Valor por Kg × Peso). Parâmetros iniciais: Taxa base: R$ 10,00, Valor por quilômetro: R$ 2,00, Valor por quilograma: R$ 0,50.
*   **RF-10 – Consultar Faturas:** O cliente deve poder consultar suas faturas e visualizar: Valor cobrado, data de emissão e status (`PENDENTE` ou `PAGA`).

### Módulo de Notificações (Orientado a Eventos)
*   **RF-16 – Notificar Criação de Entrega:** Quando uma entrega for criada (evento: `delivery.created`), o cliente responsável deve receber uma notificação contendo: Identificador da entrega, Endereço de origem e destino, Data de criação, e Valor estimado do frete.
*   **RF-17 – Notificar Atualização de Status de Entrega:** Quando o status de uma entrega mudar (evento: `delivery.status_changed`), o cliente deve ser notificado com: Identificador da entrega, Novo status, e Timestamp da mudança. Adicionalmente, se o status for `EM_TRANSITO`, a notificação deve incluir dados do motorista responsável (nome, telefone e localização aproximada).
*   **RF-18 – Notificar Conclusão de Entrega:** Quando uma entrega atingir o status `CONCLUIDA` (evento: `delivery.completed`), o cliente deve receber uma notificação contendo: Identificador da entrega, Endereço de origem e destino, Data e hora de conclusão, e Distância percorrida.
*   **RF-19 – Notificar Geração de Fatura:** Quando uma fatura for gerada (evento: `invoice.created`), o cliente responsável deve receber uma notificação contendo: Identificador da fatura, Identificador da entrega associada, Valor total, Data de emissão, e Link para visualizar detalhes.
*   **RF-20 – Notificar Pagamento de Fatura:** Quando uma fatura for paga (evento: `invoice.paid`), o cliente deve receber uma notificação de confirmação contendo: Identificador da fatura, Data do pagamento, Valor pago, e Método de pagamento.
*   **RF-21 – Notificar Falha em Entrega:** Quando uma entrega mudar para status `FALHA` (evento: `delivery.failed`), o cliente e o motorista devem ser notificados com: Identificador da entrega, Motivo da falha, Data e hora da falha, e Instruções para próximos passos.
*   **RF-22 – Configurar Preferências de Notificação:** O usuário deve poder configurar suas preferências de notificação, incluindo: Canais de comunicação (E-mail, SMS, Push, In-app), Tipos de evento de interesse (entregas, faturas, etc.), Frequência de notificações (em tempo real, resumo diário, etc.), e Horários de silêncio (períodos em que não deseja receber notificações).
*   **RF-23 – Consultar Histórico de Notificações:** O usuário deve poder consultar seu histórico de notificações enviadas, filtrando por tipo de evento, data, status de leitura, e canal de entrega.
*   **RF-24 – Marcar Notificação como Lida:** O usuário deve poder marcar notificações como lidas no painel In-app ou via API.
*   **RF-25 – Reenviar Notificação Falhada:** O sistema deve automaticamente tentar reenviar notificações que falharam, com política de retry exponencial (máximo de 3 tentativas).

## Requisitos Não Funcionais

| ID | Categoria | Descrição |
| :--- | :--- | :--- |
| **RNF-01** | Escalabilidade | O sistema deve suportar crescimento independente dos módulos de entregas, telemetria, faturamento e notificações. |
| **RNF-02** | Disponibilidade | Falhas em um módulo não devem interromper completamente a operação dos demais módulos (Isolamento de Falhas). |
| **RNF-03** | Desempenho | O registro de localizações deve suportar altas taxas de escrita sem impactar negativamente as funcionalidades de entrega e faturamento. |
| **RNF-04** | Consistência | As informações de cada domínio devem ser mantidas de forma independente e consistente (Database per Service). |
| **RNF-05** | Comunicação | A comunicação entre módulos deve garantir entrega confiável de informações, utilizando mensageria para eventos de negócio. |
| **RNF-06** | UX/UI Mobile | O aplicativo do motorista deve focar em acessibilidade e botões de ação rápidos para uso em trânsito. |
| **RNF-07** | Resiliência Offline| O aplicativo deve armazenar localmente as posições de telemetria se perder conexão, enviando os dados em lote quando a rede voltar. |
| **RNF-08** | API Gateway | Utilizar um padrão de API Gateway ou BFF para gerenciar e unificar as chamadas do cliente para os microsserviços do backend. |
| **RNF-09** | Latência de Notificação | O sistema de notificações deve manter latência máxima de 5 segundos entre a ocorrência do evento e o envio da notificação ao usuário para canais em tempo real (E-mail e Push). |
| **RNF-10** | Confiabilidade de Entrega | O sistema de notificações deve garantir entrega de pelo menos uma cópia de cada notificação (At-Least-Once Delivery), com deduplicação em caso de reenvio. |
| **RNF-11** | Suporte a Múltiplos Canais | O sistema de notificações deve suportar envio através de múltiplos canais (E-mail, SMS, Push, In-app) de forma independente e escalável. |
| **RNF-12** | Idempotência | As operações de envio de notificação devem ser idempotentes, permitindo replay de eventos sem duplicação de notificações para o usuário. |
| **RNF-13** | Rastreabilidade | O sistema deve manter logs de auditoria de todas as notificações enviadas, incluindo: timestamp, usuário destinatário, canal, status de entrega, e motivos de falha. |

## Regras de Negócio
1.  **RN-01:** Uma entrega só pode ser faturada após atingir o status `CONCLUIDA`.
2.  **RN-02:** O cálculo da distância para o faturamento deve ser fornecido pelo evento de conclusão da entrega.
3.  **RN-03:** A telemetria deve ser ingerida em alta frequência (ex: a cada 10 segundos), mas a persistência não deve bloquear o serviço de entregas.
4.  **RN-04:** As notificações são geradas exclusivamente através de eventos publicados por outros microsserviços (Delivery, Faturamento, etc.), não devem existir APIs diretas de criação de notificações.
5.  **RN-05:** Se um usuário desabilitar notificações para um tipo de evento, o sistema não deve enviar notificações daquele tipo, mesmo que o evento ocorra.
6.  **RN-06:** Notificações não entregues devem ser automaticamente retentadas com backoff exponencial (1s, 2s, 4s) por no máximo 3 tentativas.
7.  **RN-07:** Se um usuário não possui um canal de notificação configurado (ex: sem email ou sem número de celular), o sistema não deve tentar enviar notificação por esse canal e deve registrar este fato.
8.  **RN-08:** A marca de "lida" de uma notificação deve ser sincronizada apenas entre o painel In-app e a API, não afetando outras cópias já entregues (ex: e-mail já recebido).
9.  **RN-09:** O sistema deve respeitar horários de silêncio configurados pelo usuário, enfileirando notificações para entrega após o horário permitido, exceto para eventos críticos (ex: falha de entrega).