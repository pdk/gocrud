// Code generated with github.com/pdk/gocrud DO NOT EDIT

package example

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

type hookFunc func(Account) (Account, error)

func noopHook(acct Account) (Account, error) {
	return acct, nil
}

var (
	insertStmt, updateStmt, deleteStmt, selectStmt, selectByIDStmt string

	preInsert hookFunc = noopHook
    preUpdate hookFunc = noopHook
    postUpdate hookFunc = noopHook
    postDelete hookFunc = noopHook
)

// scanAccount creates an Account and scans into it, returning the new Account.
func scanAccount(row scanner) (Account, error) {
	acct := Account{}
	err := row.Scan(&acct.ID, &acct.CreatedAt, &acct.UpdatedAt, &acct.CreatedByUserID, &acct.UpdatedByUserID, &acct.Name, &acct.URLStub)
	return acct, err
}

func init() {
	selectStmt = "select account_id, created_at, updated_at, created_by_user_id, updated_by_user_id, name, url_stub from accounts "

	selectByIDStmt = selectStmt + "where account_id = ?"
	selectByIDStmt = rebind.ToDollar(selectByIDStmt)

	insertStmt = "insert into accounts (created_at, updated_at, created_by_user_id, updated_by_user_id, name, url_stub) values (?, ?, ?, ?, ?, ?) returning account_id"

	updateStmt = "update accounts set created_at = ?, updated_at = ?, created_by_user_id = ?, updated_by_user_id = ?, name = ?, url_stub = ? where account_id = ?"
	updateStmt = rebind.ToDollar(updateStmt)

	deleteStmt = "delete from accounts where account_id = ?"
	deleteStmt = rebind.ToDollar(deleteStmt)

	var 

	preInsert, preUpdate, postUpdate, postDelete = noopHook, noopHook, noopHook, noopHook
}

// Insert creates an Account record in the database, returning a new Account with a new ID value.
func Insert(db dbHandle, acct Account) (Account, error) {

	acct, err := preInsert(acct)
	if err != nil {
		return acct, err
	}

	err = db.QueryRow(insertStmt, acct.CreatedAt, acct.UpdatedAt, acct.CreatedByUserID, acct.UpdatedByUserID, acct.Name, acct.URLStub).
		Scan(&acct.ID)

	return acct, err
}

// Update modifies an Account record in the database.
func Update(db dbHandle, acct Account) (Account, error) {

	acct, err := preUpdate(acct)
	if err != nil {
		return acct, err
	}		

	_, err = db.Exec(updateStmt, acct.CreatedAt, acct.UpdatedAt, acct.CreatedByUserID, acct.UpdatedByUserID, acct.Name, acct.URLStub,	acct.ID)
	if err != nil {
		return acct, err
	}

	return postUpdate(acct)
}

// Delete removes an Account record from the database.
func Delete(db dbHandle, acct Account) (Account, error) {

	_, err := db.Exec(deleteStmt, acct.ID)
	if err != nil {
		return acct, nil
	}

	return postDelete(acct)
}

// Query executes a select statement and returns a slice of Accounts.
func Query(db dbHandle, queryExtra string, bindValues ...interface{}) ([]Account, error) {

	results := []Account{}

	query :=  selectStmt + queryExtra
	query = rebind.ToDollar(query)

	rows, err := db.Query(query, bindValues...)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		acct, err := scanAccount(rows)
		if err != nil {
			return results, err
		}

		results = append(results, acct)
	}

	err = rows.Err()
	if err != nil {
		return results, err
	}

	return results, rows.Close()
}

// QueryRow executes a select statement expecting at most 1 result.
func QueryRow(db dbHandle, queryExtra string, bindValues ...interface{}) (Account, error) {

	query := selectStmt + queryExtra
	query = rebind.ToDollar(query)

	return scanAccount(db.QueryRow(query, bindValues...))
}

// QueryRowByID executes a "select ... where account_id = ?" statement
func QueryRowByID(db dbHandle, val int64) (Account, error) {

	return scanAccount(db.QueryRow(selectByIDStmt, val))
}

// QueryBy executes a select statement with a series of name, value pairs to construct a where clause.
func QueryBy(db dbHandle, nameValuePairs ...interface{}) ([]Account, error) {

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