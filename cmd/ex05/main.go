package main

import (
	"fmt"
	"log"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func formatRule(rule *nftables.Rule) string {
	var ruleDesc string
	for _, e := range rule.Exprs {
		switch v := e.(type) {
		case *expr.Payload:
			// Simplificando para este caso específico
			if v.Base == expr.PayloadBaseTransportHeader && v.Offset == 2 {
				ruleDesc += "tcp dport "
			}
		case *expr.Cmp:
			// Simplificando para este caso específico
			if len(v.Data) == 2 {
				port := int(v.Data[0])<<8 + int(v.Data[1])
				ruleDesc += fmt.Sprintf("%d ", port)
			}
		case *expr.Verdict:
			if v.Kind == expr.VerdictAccept {
				ruleDesc += "accept"
			}
		}
	}
	return ruleDesc
}

func main() {
	conn, err := nftables.New()
	if err != nil {
		log.Fatalf("Erro ao criar conexão NFTables: %v", err)
	}

	tables, err := conn.ListTables()
	if err != nil {
		log.Fatalf("Erro ao listar tabelas: %v", err)
	}

	for _, table := range tables {
		fmt.Printf("table ip %s {\n", table.Name)

		chains, err := conn.ListChains()
		if err != nil {
			log.Fatalf("Erro ao listar cadeias: %v", err)
		}

		for _, chain := range chains {
			if chain.Table.Name == table.Name {
				fmt.Printf("\tchain %s {\n", chain.Name)
				// fmt.Printf("\t\ttype %s hook %s priority %d; policy %s;\n", chain.Type, chain.Hooknum, chain.Priority, chain.Policy)

				rules, err := conn.GetRules(table, chain)
				if err != nil {
					log.Fatalf("Erro ao listar regras: %v", err)
				}

				for _, rule := range rules {
					formattedRule := formatRule(rule)
					if formattedRule != "" {
						fmt.Printf("\t\t%s\n", formattedRule)
					}
				}
				fmt.Printf("\t}\n")
			}
		}
		fmt.Printf("}\n")
	}
}
