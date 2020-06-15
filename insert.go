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

	insertColumnNames := desc.ColumnsOmitField(idField)
	insertColumnNamesStr := strings.Join(insertColumnNames, ", ")
	insertColumnIndexes := desc.FieldIndexesOmitField(idField)
	valueCount := len(insertColumnNames)

	idColumnName := desc.ColumnNameOfField(idField)
	idColumnIndex := desc.FieldIndex(idField)

	insertStmt := fmt.Sprintf("insert into %s (%s) values (%s) returning %s",
		tableName, insertColumnNamesStr, markers(len(insertColumnIndexes)), idColumnName)

	return func(db *sql.DB, item interface{}) error {

		ptr := reflect.ValueOf(item)
		if ptr.Kind() != reflect.Ptr {
			return fmt.Errorf("gocrud.NewAutoIncrIDInserter created func expected a pointer, but got a %T", item)
		}

		itemValue := ptr.Elem()
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

// MakeInserter will create a function to insert, and set the input variable to hold it.
// E.g:
// var insertMyStruct func(*sql.DB, MyStruct) (MyStruct, error)
// MakeInserter(&insertMyStruct, "table_name", MyStruct{})
func MakeInserter(funcVar interface{}, tableName string, exampleItem interface{}) {

	fn := reflect.ValueOf(funcVar).Elem()

	xfunc := NewInserter(tableName, exampleItem)

	inserterFunc := reflect.MakeFunc(fn.Type(), func(args []reflect.Value) []reflect.Value {
		db := args[0].Interface().(*sql.DB)
		data := args[1].Interface()
		err := xfunc(db, data)
		return []reflect.Value{args[1], reflect.ValueOf(err)}
	})

	fn.Set(inserterFunc)
}

// NewInserter generates a new Inserter. To be called on startup. Will
// terminate program if there is any error. The item value passed to the
// generated function MUST be a pointer to the type of struct passed as example.
// f := NewInserter("foo", "ID", DaStruct{})
// err := f(db, &daStructInstance)
func NewInserter(tableName string, exampleItem interface{}) Inserter {

	desc, err := describe.Describe(exampleItem)
	if err != nil {
		log.Fatalf("cannot build inserter: %v", err)
	}

	insertColumnNames := desc.Columns()
	insertColumnNamesStr := strings.Join(insertColumnNames, ", ")
	valueCount := len(insertColumnNames)

	insertStmt := fmt.Sprintf("insert into %s (%s) values (%s)",
		tableName, insertColumnNamesStr, markers(len(insertColumnNames)))

	return func(db *sql.DB, item interface{}) error {

		itemValue := reflect.ValueOf(item).Elem()
		insertValues := make([]interface{}, valueCount, valueCount)
		for i := range insertColumnNames {
			insertValues[i] = itemValue.Field(i).Interface()
		}

		_, err = db.Exec(insertStmt, insertValues...)

		return err
	}
}
