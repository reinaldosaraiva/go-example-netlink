# NFTables Manager em Go

Este projeto contém um script em Go para interagir com o NFTables, permitindo a criação e listagem de tabelas, cadeias e regras.

## Descrição

Este script em Go utiliza a biblioteca `nftables` para criar tabelas e cadeias no NFTables e listar suas configurações de forma hierárquica. É útil para administradores de sistemas e desenvolvedores que precisam gerenciar regras de firewall no Linux de forma programática.

## Funcionalidades

- Criação de tabelas no NFTables.
- Criação de cadeias (chains) dentro das tabelas.
- Listagem de tabelas, cadeias e regras em um formato legível.

## Pré-requisitos

Para executar este script, você precisará:

- Go versão 1.19 ou superior.
- Acesso ao NFTables no seu sistema Linux.
- Privilégios de superusuário para interação com o NFTables.

## Uso
Para executar o script, navegue até o diretório do projeto e use o comando go run:

```bash
cd seu-repositorio
sudo go run cmd/exXX/main.go
```

## Licença
"THE BEER-WARE LICENSE" (Revision 42):
Enquanto você mantiver este aviso, você pode fazer o que quiser com este código. Se nos encontrarmos algum dia, e você achar que este código vale a pena, você pode me comprar uma cerveja em retorno.