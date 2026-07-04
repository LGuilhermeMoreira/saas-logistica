# Documento de Requisitos
## Visão Geral

O sistema tem como objetivo gerenciar entregas, rastreamento de motoristas e faturamento de fretes através de uma arquitetura baseada em microsserviços.

## Requisitos Funcionais### Módulo de Autenticação e Segurança
*   **RF-11 – Cadastro de Usuários:** O sistema deve permitir o cadastro de novos usuários, diferenciando perfis (ex: Administrador, Motorista, Cliente).
*   **RF-12 – Autenticação de Usuários:** O sistema deve permitir que usuários realizem login utilizando e-mail e senha.
*   **RF-13 – Geração de Token de Acesso:** Após o login bem-sucedido, o sistema deve emitir um token seguro (ex: JWT) para autenticar requisições subsequentes.
*   **RF-14 – Controle de Acesso (RBAC):** O sistema deve restringir o acesso às funcionalidades com base no perfil do usuário (ex: apenas motoristas podem atualizar status de entrega).


### Módulo de Entregas
*   **RF-01 – Criar Entrega:** O sistema deve permitir que um cliente cadastre uma nova entrega informando:
    *   Endereço de origem;
    *   Endereço de destino;
    *   Peso do pacote.
*   **RF-02 – Atribuir Motorista:** O sistema deve permitir associar uma entrega a um motorista específico.
*   **RF-03 – Atualizar Status da Entrega:** O motorista deve poder atualizar o status da entrega seguindo o fluxo:
    *   `CRIADA` → `EM_TRANSITO` → `CONCLUIDA` ou `FALHA`.
*   **RF-04 – Consultar Entrega:** O sistema deve permitir consultar os detalhes de uma entrega.
    *   *Nota:* Quando a entrega estiver com status `EM_TRANSITO`, a consulta deverá apresentar também a última localização conhecida do motorista responsável (via integração com Telemetria).
*  **RF-15 – Notificar Conclusão de Entrega por E-mail:** Quando uma entrega atingir o status CONCLUIDA, o sistema deve enviar automaticamente um e-mail de notificação ao cliente responsável pela entrega, contendo:

* Identificador da entrega;
* Endereço de origem e destino;
* Data e hora de conclusão;

### Módulo de Telemetria
*   **RF-05 – Registrar Localização do Motorista:** O sistema deve permitir registrar periodicamente a localização de um motorista contendo:
    *   Identificador do motorista;
    *   Latitude e Longitude;
    *   Timestamp.
*   **RF-06 – Consultar Posição Atual:** O sistema deve permitir consultar a última posição registrada de um motorista.
*   **RF-07 – Consultar Histórico de Rota:** O sistema deve permitir consultar o histórico de coordenadas de um motorista dentro de um intervalo de tempo informado.

### Módulo de Faturamento
*   **RF-08 – Gerar Fatura Automaticamente:** Quando uma entrega for concluída, o sistema deverá gerar uma fatura correspondente ao frete realizado.
*   **RF-09 – Calcular Valor do Frete:** O valor da fatura deve ser calculado utilizando a fórmula:
    > **Valor** = Taxa Base + (Valor por Km × Distância) + (Valor por Kg × Peso)
    
    **Parâmetros iniciais:**
    *   Taxa base: R$ 10,00
    *   Valor por quilômetro: R$ 2,00
    *   Valor por quilograma: R$ 0,50
*   **RF-10 – Consultar Faturas:** O cliente deve poder consultar suas faturas e visualizar:
    *   Valor cobrado, data de emissão e status (`PENDENTE` ou `PAGA`).

## Requisitos Não Funcionais
| ID | Categoria | Descrição |
| :--- | :--- | :--- |
| **RNF-01** | Escalabilidade | O sistema deve suportar crescimento independente dos módulos de entregas, telemetria e faturamento. |
| **RNF-02** | Disponibilidade | Falhas em um módulo não devem interromper completamente a operação dos demais módulos (Isolamento de Falhas). |
| **RNF-03** | Desempenho | O registro de localizações deve suportar altas taxas de escrita sem impactar negativamente as funcionalidades de entrega e faturamento. |
| **RNF-04** | Consistência | As informações de cada domínio devem ser mantidas de forma independente e consistente (Database per Service). |
| **RNF-05** | Comunicação | A comunicação entre módulos deve garantir entrega confiável de informações, utilizando mensageria para eventos de negócio. |

## Regras de Negócio

1.  **RN-01:** Uma entrega só pode ser faturada após atingir o status `CONCLUIDA`.
2.  **RN-02:** O cálculo da distância para o faturamento deve ser fornecido pelo evento de conclusão da entrega.
3.  **RN-03:** A telemetria deve ser ingerida em alta frequência (ex: a cada 10 segundos), mas a persistência não deve bloquear o serviço de entregas.
