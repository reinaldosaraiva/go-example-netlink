package main

import (
	"fmt"
	"log"

	"github.com/google/nftables"
)

func main() {
    conn, err := nftables.New()
    if err != nil {
        log.Fatalf("Erro ao criar conexão NFTables: %v", err)
    }

    // Tentando criar uma nova tabela com um nome único
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

    if exists {
        fmt.Printf("A tabela '%s' já existe.\n", tableName)
    } else {
        // Criar a tabela
        table := &nftables.Table{
            Family: nftables.TableFamilyIPv4,
            Name:   tableName,
        }
        conn.AddTable(table) // Removido a verificação de erro
        fmt.Printf("Tabela '%s' criada com sucesso.\n", tableName)
    }

    // Aplicar as mudanças
    if err := conn.Flush(); err != nil {
        log.Fatalf("Erro ao aplicar mudanças: %v", err)
    }
}
