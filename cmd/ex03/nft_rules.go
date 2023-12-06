package main

import (
	"fmt"
	"log"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
	"golang.org/x/sys/unix"
)

func main() {
    conn := &nftables.Conn{}

    // Listar tabelas existentes
    tables, err := conn.ListTables()
    if err != nil {
        log.Fatalf("Falha ao listar tabelas: %v", err)
    }

    // Verificar se a tabela inet_filter já existe
    var table *nftables.Table
    exists := false
    for _, t := range tables {
        if t.Name == "filtro_pacotes" && t.Family == nftables.TableFamilyINet {
            table = t
            exists = true
            break
        }
    }

    // Criar a tabela se não existir
    if !exists {
        table = &nftables.Table{
            Family: nftables.TableFamilyINet,
            Name:   "filtro_pacotes",
        }
        if err := conn.AddTable(table); err != nil {
            log.Fatalf("Falha ao adicionar tabela: %v", err)
        }
    }

    // Criar uma nova cadeia para entrada
    chain := &nftables.Chain{
        Name:     "entrada",
        Table:    table,
        Type:     nftables.ChainTypeFilter,
        Hooknum:  nftables.ChainHookInput,
        Priority: nftables.ChainPriorityFilter,
    }
    if err := conn.AddChain(chain); err != nil {
        log.Fatalf("Falha ao adicionar cadeia: %v", err)
    }

    // Adicionar regras para as portas 8080, 22 e 80
    for _, port := range []uint16{8080, 22, 80} {
        rule := &nftables.Rule{
            Table: table,
            Chain: chain,
            Exprs: []expr.Any{
                &expr.Payload{
                    DestRegister: 1,
                    Base:         expr.PayloadBaseTransportHeader,
                    Offset:       2, // Offset para a porta de destino no cabeçalho TCP
                    Len:          2, // Tamanho da porta (2 bytes)
                },
                &expr.Cmp{
                    Op:       expr.CmpOpEq,
                    Register: 1,
                    Data:     []byte{byte(port >> 8), byte(port & 0xff)},
                },
                &expr.Verdict{
                    Kind: expr.VerdictAccept,
                },
            },
        }
        if err := conn.AddRule(rule); err != nil {
            log.Fatalf("Falha ao adicionar regra para porta %d: %v", port, err)
        }
    }

    // Adicionar regra para registrar e descartar outros pacotes
    conn.AddRule(&nftables.Rule{
        Table: table,
        Chain: chain,
        Exprs: []expr.Any{
            &expr.Meta{
                Key:      expr.MetaKeyNFPROTO,
                Register: 1,
            },
            &expr.Cmp{
                Op:       expr.CmpOpNeq,
                Register: 1,
                Data:     []byte{byte(unix.AF_INET), byte(unix.AF_INET6)},
            },
            &expr.Log{
                Level: 7, // Nível de log; 7 é 'debug'
            },
            &expr.Verdict{
                Kind: expr.VerdictDrop,
            },
        },
    })
	if err := conn.AddTable(table); err != nil {
		log.Fatalf("Falha ao adicionar tabela '%s': %v", table.Name, err)
	}
	
    // Aplicar as mudanças
    if err := conn.Flush(); err != nil {
        log.Fatalf("Falha ao aplicar regras: %v", err)
    }

    fmt.Println("Regras do NFTables aplicadas com sucesso.")
}
