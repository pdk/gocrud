// Code generated with github.com/pdk/gocrud DO NOT EDIT

package {{ .packageName }}

import (
	"database/sql"
	"strings"

	"github.com/pdk/gocrud/rebind"
)

// Callers to these functions may pass either a *sql.DB or a *sql.Tx
type dbHandle interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// scanner is either *sql.Row or *sql.Rows
type scanner interface {
	Scan(dest ...interface{}) error
}

type hookFunc func({{ .structName }}) ({{ .structName }}, error)

func noopHook({{ .instanceName }} {{ .structName }}) ({{ .structName }}, error) {
	return {{ .instanceName }}, nil
}

var (
	insertStmt, updateStmt, deleteStmt, selectStmt, selectBy{{ .idFieldName }}Stmt string

	preInsert hookFunc = noopHook
    preUpdate hookFunc = noopHook
    postUpdate hookFunc = noopHook
    postDelete hookFunc = noopHook
)

// scan{{ .structName }} creates {{ .particle }} {{ .structName }} and scans into it, returning the new {{ .structName }}.
func scan{{ .structName }}(row scanner) ({{ .structName }}, error) {
	{{ .instanceName }} := {{ .structName }}{}
	err := row.Scan({{ .idAddress }}, {{ .fieldAddresses }})
	return {{ .instanceName }}, err
}

func init() {
	selectStmt = "select {{ .idColumnName }}, {{ .fieldColumns }} from {{ .tableName }} "

	selectBy{{ .idFieldName }}Stmt = selectStmt + "where {{ .idColumnName }} = ?"
	selectBy{{ .idFieldName }}Stmt = rebind.ToDollar(selectBy{{ .idFieldName }}Stmt)

	insertStmt = "insert into {{ .tableName }} ({{ .fieldColumns }}) values ({{ .fieldBindMarks }}) returning {{ .idColumnName }}"

	updateStmt = "update {{ .tableName }} set {{ .updateFields }} where {{ .idColumnName }} = ?"
	updateStmt = rebind.ToDollar(updateStmt)

	deleteStmt = "delete from {{ .tableName }} where {{ .idColumnName }} = ?"
	deleteStmt = rebind.ToDollar(deleteStmt)

	var 

	preInsert, preUpdate, postUpdate, postDelete = noopHook, noopHook, noopHook, noopHook
}

// Insert creates {{ .particle }} {{ .structName }} record in the database, returning a new {{ .structName }} with a new ID value.
func Insert(db dbHandle, {{ .instanceName }} {{ .structName }}) ({{ .structName }}, error) {

	{{ .instanceName }}, err := preInsert({{ .instanceName }})
	if err != nil {
		return {{ .instanceName }}, err
	}

	err = db.QueryRow(insertStmt, {{ .instanceName }}.CreatedAt, {{ .instanceName }}.UpdatedAt, {{ .instanceName }}.CreatedByUserID, {{ .instanceName }}.UpdatedByUserID, {{ .instanceName }}.Name, {{ .instanceName }}.URLStub).
		Scan(&{{ .instanceName }}.ID)

	return {{ .instanceName }}, err
}

// Update modifies {{ .particle }} {{ .structName }} record in the database.
func Update(db dbHandle, {{ .instanceName }} {{ .structName }}) ({{ .structName }}, error) {

	{{ .instanceName }}, err := preUpdate({{ .instanceName }})
	if err != nil {
		return {{ .instanceName }}, err
	}		

	_, err = db.Exec(updateStmt, {{ .instanceName }}.CreatedAt, {{ .instanceName }}.UpdatedAt, {{ .instanceName }}.CreatedByUserID, {{ .instanceName }}.UpdatedByUserID, {{ .instanceName }}.Name, {{ .instanceName }}.URLStub,	{{ .instanceName }}.ID)
	if err != nil {
		return {{ .instanceName }}, err
	}

	return postUpdate({{ .instanceName }})
}

// Delete removes {{ .particle }} {{ .structName }} record from the database.
func Delete(db dbHandle, {{ .instanceName }} {{ .structName }}) ({{ .structName }}, error) {

	_, err := db.Exec(deleteStmt, {{ .instanceName }}.ID)
	if err != nil {
		return {{ .instanceName }}, nil
	}

	return postDelete({{ .instanceName }})
}

// Query executes a select statement and returns a slice of {{ .structName }}{{ .plural }}.
func Query(db dbHandle, queryExtra string, bindValues ...interface{}) ([]{{ .structName }}, error) {

	results := []{{ .structName }}{}

	query :=  selectStmt + queryExtra
	query = rebind.ToDollar(query)

	rows, err := db.Query(query, bindValues...)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		{{ .instanceName }}, err := scan{{ .structName }}(rows)
		if err != nil {
			return results, err
		}

		results = append(results, {{ .instanceName }})
	}

	err = rows.Err()
	if err != nil {
		return results, err
	}

	return results, rows.Close()
}

// QueryRow executes a select statement expecting at most 1 result.
func QueryRow(db dbHandle, queryExtra string, bindValues ...interface{}) ({{ .structName }}, error) {

	query := selectStmt + queryExtra
	query = rebind.ToDollar(query)

	return scan{{ .structName }}(db.QueryRow(query, bindValues...))
}

// QueryRowBy{{ .idFieldName }} executes a "select ... where {{ .idColumnName }} = ?" statement
func QueryRowBy{{ .idFieldName }}(db dbHandle, val int64) ({{ .structName }}, error) {

	return scan{{ .structName }}(db.QueryRow(selectBy{{ .idFieldName }}Stmt, val))
}

// QueryBy executes a select statement with a series of name, value pairs to construct a where clause.
func QueryBy(db dbHandle, nameValuePairs ...interface{}) ([]{{ .structName }}, error) {

	whereClause := strings.Builder{}
	values := []interface{}{}

	for i := 0; i < len(nameValuePairs); i += 2 {
		if i == 0 {
			whereClause.WriteString("where ")
		} else {
			whereClause.WriteString(" and ")
		}
		whereClause.WriteString(nameValuePairs[i].(string))
		whereClause.WriteString(" = ?")

		values = append(values, nameValuePairs[i+1])
	}

	return Query(db, whereClause.String(), values...)
}

// SetPreInsert sets the function to call before inserting the record.
func SetPreInsert(f hookFunc) {
	preInsert = f
}

// SetPreUpdate sets the function to call before updating the record.
func SetPreUpdate(f hookFunc) {
	preUpdate = f
}

// SetPostDelete sets the function to call after deleting the record.
func SetPostDelete(f hookFunc) {
	postDelete = f
}