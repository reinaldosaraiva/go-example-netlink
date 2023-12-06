package main

import (
	"fmt"
	"log"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
	"golang.org/x/sys/unix"
)

func main() {
	conn, err := nftables.New()
	if err != nil {
		log.Fatalf("Erro ao criar conexão NFTables: %v", err)
	}

	tableName := "test_nft_table"

	// Verificar se a tabela já existe
	tables, err := conn.ListTables()
	if err != nil {
		log.Fatalf("Erro ao listar tabelas: %v", err)
	}
	exists := false
	for _, t := range tables {
		if t.Name == tableName {
			exists = true
			break
		}
	}

	table := &nftables.Table{
		Family: nftables.TableFamilyIPv4,
		Name:   tableName,
	}

	if exists {
		fmt.Printf("A tabela '%s' já existe.\n", tableName)
	} else {
		// Criar a tabela
		conn.AddTable(table)
		fmt.Printf("Tabela '%s' criada com sucesso.\n", tableName)
	}

	// Criar a cadeia
	chain := &nftables.Chain{
		Name:     "test_nft_chain",
		Table:    table,
		Type:     nftables.ChainTypeFilter,
		Hooknum:  nftables.ChainHookInput,
		Priority: nftables.ChainPriorityFilter,
	}
	conn.AddChain(chain)
	fmt.Printf("Cadeia '%s' criada com sucesso.\n", chain.Name)

	// Adicionar regras para portas 80 e 22
	for _, port := range []uint16{80, 22} {
		rule := &nftables.Rule{
			Table: table,
			Chain: chain,
			Exprs: []expr.Any{
				&expr.Meta{Key: expr.MetaKeyL4PROTO, Register: 1},
				&expr.Cmp{Op: expr.CmpOpEq, Register: 1, Data: []byte{unix.IPPROTO_TCP}},
				&expr.Payload{DestRegister: 1, Base: expr.PayloadBaseTransportHeader, Offset: 2, Len: 2},
				&expr.Cmp{Op: expr.CmpOpEq, Register: 1, Data: []byte{(byte)(port >> 8), (byte)(port & 0xff)}},
				&expr.Verdict{Kind: expr.VerdictAccept},
			},
		}
		conn.AddRule(rule)
		fmt.Printf("Regra para a porta '%d' criada com sucesso.\n", port)
	}

	// Aplicar as mudanças
	if err := conn.Flush(); err != nil {
		log.Fatalf("Erro ao aplicar mudanças: %v", err)
	}
}
