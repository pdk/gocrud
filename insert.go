package gocrud

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/pdk/gocrud/describe"
)

// Inserter is a function taking a db conn and an item, and it inserts a new row in the DB
type Inserter func(db *sql.DB, item interface{}) error

// NewAutoIncrIDInserter generates a new Inserter. To be called on startup. Will
// terminate program if there is any error. The item value passed to the
// generated function MUST be a pointer to the type of struct passed as example.
// f := NewInserterWithAutoIncrID("foo", "ID", DaStruct{})
// err := f(db, &daStructInstance)
func NewAutoIncrIDInserter(tableName string, idField string, exampleItem interface{}) Inserter {

	desc, err := describe.Describe(exampleItem)
	if err != nil {
		log.Fatalf("cannot build inserter: %v", err)
	}

	insertColumnNames := strings.Join(desc.ColumnsOmitField(idField), ", ")
	insertColumnIndexes := desc.FieldIndexesOmitField(idField)
	valueCount := len(insertColumnNames)

	idColumnName := desc.ColumnNameOfField(idField)
	idColumnIndex := desc.FieldIndex(idField)

	insertStmt := fmt.Sprintf("insert into %s (%s) values (%s) returning %s",
		tableName, insertColumnNames, markers(len(insertColumnIndexes)), idColumnName)

	return func(db *sql.DB, item interface{}) error {

		itemValue := reflect.ValueOf(item).Elem()
		insertValues := make([]interface{}, valueCount, valueCount)
		for i, index := range insertColumnIndexes {
			insertValues[i] = itemValue.Field(index).Interface()
		}

		var newID int64
		err = db.QueryRow(insertStmt, insertValues...).Scan(&newID)

		itemValue.Field(idColumnIndex).SetInt(newID)

		return err
	}
}
