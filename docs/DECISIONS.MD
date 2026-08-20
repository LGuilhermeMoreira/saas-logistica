# DECISÕES TOMADAS 

#### Trocar Casbin por OPA

A troca faz sentido, pois o OPA roda como um container colado ao microsserviço, o que permite distribuir melhor. 

#### Deploy na AWS

Vamos utilizar as seguintes tecnologias para fazer deploy
* ECR: registrar os containers dos microsserviços e gateway
* fargate: Gerencia os containers
* ECS: serviço de containers que vai rodar junto com fargate
* AWS Cloud Map: Vai gerenciar a parte de service discorery na aplicação.