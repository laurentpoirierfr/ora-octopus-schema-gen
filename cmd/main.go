package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/laurentpoirierfr/ora-octopus-schema-gen/internal/format/octopus"
	_ "github.com/sijms/go-ora/v2"
)

func main() {
	schema := octopus.Schema{}
	//db, err := go_ora.NewConnection(os.Getenv("ORACLE_URL"), &configurations.ConnectionConfig{})
	conn, err := sql.Open("oracle", os.Getenv("ORACLE_URL"))
	DieIfError(err)
	defer conn.Close()

	schema.Author = "github.com/laurentpoirierfr/ora-octopus-schema-gen"
	schema.Version = "1.0.0"
	tables, err := listTablesSchema(conn, os.Getenv("ORACLE_SCHEMA_NAME"))
	DieIfError(err)
	schema.Tables = tables

	j, err := json.MarshalIndent(schema, "", "\t")
	DieIfError(err)

	fmt.Println(string(j))
}

func DieIfError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func listTablesSchema(conn *sql.DB, owner string) ([]*octopus.Table, error) {
	var tables []*octopus.Table
	sql := fmt.Sprintf("SELECT table_name FROM all_tables WHERE owner = '%s'", strings.ToUpper(owner))
	rows, err := conn.Query(sql)

	if err != nil {
		return []*octopus.Table{}, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			fmt.Println("Can't close dataset: ", err)
		}
	}()
	var table_name string

	for rows.Next() {
		err = rows.Scan(&table_name)
		if err != nil {
			return []*octopus.Table{}, err
		}

		columns, err := listColumnsTable(conn, owner, table_name)
		DieIfError(err)
		indices, err := listIndecesTable(conn, owner, table_name)
		DieIfError(err)

		table := octopus.Table{
			Name:    table_name,
			Columns: columns,
			Indices: indices,
		}
		tables = append(tables, &table)
	}
	return tables, nil
}

func listColumnsTable(conn *sql.DB, owner, tableName string) ([]*octopus.Column, error) {
	var columns []*octopus.Column
	sql := fmt.Sprintf("select column_name, data_type, nullable, char_length, identity_column from all_tab_cols WHERE owner = '%s' AND table_name = '%s'", strings.ToUpper(owner), tableName)
	rows, err := conn.Query(sql)

	if err != nil {
		return []*octopus.Column{}, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			fmt.Println("Can't close dataset: ", err)
		}
	}()
	var (
		name           string
		dataType       string
		nullable       string
		charLength     uint16
		identityColumn string
	)
	for rows.Next() {
		err = rows.Scan(&name, &dataType, &nullable, &charLength, &identityColumn)
		if err != nil {
			return []*octopus.Column{}, err
		}

		reference, err := listReferencesTable(conn, tableName, name)
		if err != nil {
			return []*octopus.Column{}, err
		}

		columns = append(columns, &octopus.Column{
			Name:            strings.ToLower(name),
			Type:            strings.ToLower(dataType),
			Description:     "",
			Size:            charLength,
			Scale:           0,
			NotNull:         !convertYesNoToBoolean(nullable),
			PrimaryKey:      convertYesNoToBoolean(identityColumn),
			UniqueKey:       false,
			AutoIncremental: false,
			DefaultValue:    "",
			OnUpdate:        "",
			Values:          []string{},
			Ref:             reference,
		})

	}
	return columns, nil
}

func listIndecesTable(conn *sql.DB, owner, tableName string) ([]*octopus.Index, error) {
	var indeces []*octopus.Index
	sql := fmt.Sprintf("SELECT index_name, column_name FROM DBA_IND_COLUMNS WHERE table_name = '%s' AND index_owner = '%s'", strings.ToUpper(tableName), strings.ToUpper(owner))

	rows, err := conn.Query(sql)
	if err != nil {
		return indeces, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			fmt.Println("Can't close dataset: ", err)
		}
	}()
	var (
		indexName  string
		columnName string
	)
	for rows.Next() {
		err = rows.Scan(&indexName, &columnName)
		if err != nil {
			return indeces, err
		}
		columns := []string{strings.ToLower(columnName)}
		indeces = append(indeces, &octopus.Index{
			Name:    indexName,
			Columns: columns,
		})
	}

	return indeces, nil
}

func listReferencesTable(conn *sql.DB, tableName, columnName string) (*octopus.Reference, error) {
	var reference *octopus.Reference
	sql := fmt.Sprintf("SELECT a.table_name child_table, a.column_name child_column, a.constraint_name, b.table_name parent_table, b.column_name parent_column FROM all_cons_columns a JOIN all_constraints c ON a.owner = c.owner AND a.constraint_name = c.constraint_name join all_cons_columns b on c.owner = b.owner and c.r_constraint_name = b.constraint_name WHERE c.constraint_type = 'R' AND a.table_name = '%s' AND a.column_name = '%s'", strings.ToUpper(tableName), strings.ToUpper(columnName))
	rows, err := conn.Query(sql)
	if err != nil {
		return &octopus.Reference{}, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			fmt.Println("Can't close dataset: ", err)
		}
	}()
	var (
		childTable    string
		childColumn   string
		contraintName string
		parentTable   string
		parentColumn  string
	)
	for rows.Next() {
		err = rows.Scan(&childTable, &childColumn, &contraintName, &parentTable, &parentColumn)
		if err != nil {
			return &octopus.Reference{}, err
		}
		reference = &octopus.Reference{
			Table:        parentTable,
			Column:       strings.ToLower(parentColumn),
			Relationship: "1:1",
		}
	}
	return reference, nil
}

func convertYesNoToBoolean(str string) bool {
	if len(str) > 1 {
		return strings.ToUpper(string(str[0])) == "Y"
	}
	return false
}
