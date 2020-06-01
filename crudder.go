package gocrud

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Crudder provides enough metadata to drive CRUD operations.
type Crudder interface {
	Table() string        // the name of the table
	IDColumn() []string   // the auto-incr ID column name
	KeyColumns() []string // alternate to using an auto-incr ID column
	Columns() []string    // the OTHER column names

}

type CrudderItem interface {
}

// Scannable allows loading db query results
type Scannable interface {
	ScanId() []interface{}     // slice of pointer to ID scannable
	ScanKeys() []interface{}   // slice of pointers to key scannables
	ScanFields() []interface{} // slice of pointers for loading query result
}

// CreatedAter allows setting CreatedAt
type CreatedAter interface {
	SetCreatedAt(time.Time)
}

// UpdatedAter allows setting UpdatedAt
type UpdatedAter interface {
	SetUpdatedAt(time.Time)
}

// selectBy{{ .idFieldName }}Stmt = selectStmt + "where {{ .idColumnName }} = ?"
// selectBy{{ .idFieldName }}Stmt = rebind.ToDollar(selectBy{{ .idFieldName }}Stmt)

// insertStmt = "insert into {{ .tableName }} ({{ .fieldColumns }}) values ({{ .fieldBindMarks }}) returning {{ .idColumnName }}"

// updateStmt = "update {{ .tableName }} set {{ .updateFields }} where {{ .idColumnName }} = ?"
// updateStmt = rebind.ToDollar(updateStmt)

// deleteStmt = "delete from {{ .tableName }} where {{ .idColumnName }} = ?"
// deleteStmt = rebind.ToDollar(deleteStmt)

// Insert creates a new row
func Insert(db *sql.DB, item Crudder) {

	insertStmt := fmt.Sprintf("insert into %s (%s) values (%s) returning %s",
		item.Table(), allColumns(item), allMarkers(item), item.IDColumn())

	log.Print(insertStmt)
}

func insertColumns(item Crudder) []string {
	return append(item.KeyColumns(), item.Columns()...)
}

func insertMarkers(item Crudder) string {
	return markers(len(item.KeyColumns()) + len(item.Columns()))
}

func allColumns(item Crudder) []string {
	return append(append(item.IDColumn(), item.KeyColumns()...), item.Columns()...)
}

func allMarkers(item Crudder) string {
	return markers(len(item.IDColumn()) + len(item.KeyColumns()) + len(item.Columns()))

}

// markers returns a string of n bind markers, comma separated
func markers(n int) string {
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("$")
		sb.WriteString(strconv.Itoa(i))
	}

	return sb.String()
}
