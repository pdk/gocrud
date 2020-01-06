package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/markbates/pkger"
	"github.com/pdk/struct2json"
)

var (
	goFile       string
	structName   string
	instanceName string
	packageName  string
	methodPrefix string
	methodSuffix string
	templateName string
	tableName    string
	useAn        bool
	useEs        bool
	particle     = "a"
	plural       = "s"
)

func init() {
	pkger.Include("/templates/")

	flag.StringVar(&goFile, "source", "", "name of .go file to read (required)")
	flag.StringVar(&structName, "struct", "", "name of struct for which to generate CRUD (required)")
	flag.StringVar(&instanceName, "instance", "", "name of instance vars of the struct")
	flag.StringVar(&packageName, "package", "", "package name (required)")
	flag.StringVar(&methodPrefix, "prefix", "", "string to prefix method names")
	flag.StringVar(&methodSuffix, "suffix", "", "string to suffix method names")
	flag.StringVar(&templateName, "template", "postgres", "name of which template to use (postgres or sqlite)")
	flag.StringVar(&tableName, "table", "", "name of database table")
	flag.BoolVar(&useAn, "an", false, "use 'an' instead of 'a'")
	flag.BoolVar(&useEs, "es", false, "use 'es' instead of 's' for plurals")
}

func main() {

	flag.Parse()

	if goFile == "" || structName == "" || packageName == "" {
		fmt.Fprintf(os.Stderr, "usage: gocrud -flag ...\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if useAn {
		particle = "an"
	}
	if useEs {
		plural = "es"
	}

	structInfo, ok := struct2json.GetStructs(goFile).Get(structName)
	if !ok {
		log.Fatalf("struct %s not found in %s", structName, goFile)
	}

	data := map[string]interface{}{
		"packageName":    packageName,
		"structName":     structName,
		"instanceName":   getInstanceName(),
		"struct":         structInfo,
		"idFieldName":    structInfo.Fields[0].Name,
		"idColumnName":   columnName(structInfo.Fields[0]),
		"idAddress":      "&" + getInstanceName() + "." + structInfo.Fields[0].Name,
		"fieldNames":     getFieldNames(structInfo.Fields),
		"fieldColumns":   getFieldColumns(structInfo.Fields),
		"fieldValues":    getFieldValues(structInfo.Fields),
		"fieldAddresses": getFieldAddressses(structInfo.Fields),
		"fieldBindMarks": getFieldBindMarks(structInfo.Fields),
		"updateFields":   getUpdateFields(structInfo.Fields),
		"particle":       particle,
		"plural":         plural,
		"tableName":      getTableName(),
	}

	t := readTemplate(templateName)
	err := t.Execute(os.Stdout, data)
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}
}

func getFieldNames(flds []struct2json.Field) string {

	names := []string{}
	for _, f := range flds[1:] {
		names = append(names, f.Name)
	}

	return strings.Join(names, ", ")
}

func getFieldColumns(flds []struct2json.Field) string {

	names := []string{}
	for _, f := range flds[1:] {
		names = append(names, columnName(f))
	}

	return strings.Join(names, ", ")
}

func getFieldValues(flds []struct2json.Field) string {

	instanceName = getInstanceName()

	names := []string{}
	for _, f := range flds[1:] {
		names = append(names, instanceName+"."+f.Name)
	}

	return strings.Join(names, ", ")
}

func getFieldBindMarks(flds []struct2json.Field) string {

	marks := []string{}
	for range flds[1:] {
		marks = append(marks, "?")
	}

	return strings.Join(marks, ", ")
}

func getFieldAddressses(flds []struct2json.Field) string {

	instanceName = getInstanceName()

	addrs := []string{}
	for _, f := range flds[1:] {
		addrs = append(addrs, "&"+instanceName+"."+f.Name)
	}

	return strings.Join(addrs, ", ")
}

func getUpdateFields(flds []struct2json.Field) string {

	upds := []string{}
	for _, f := range flds[1:] {
		upds = append(upds, columnName(f)+" = ?")
	}

	return strings.Join(upds, ", ")
}

func getInstanceName() string {
	if instanceName != "" {
		return instanceName
	}

	return strcase.ToLowerCamel(structName)
}

func getTableName() string {
	if tableName != "" {
		return tableName
	}

	return strcase.ToSnake(structName) + plural
}

func columnName(fld struct2json.Field) string {

	if colName, ok := fld.Tags["db"]; ok {
		return colName
	}
	return strcase.ToSnake(fld.Name)
}

func readTemplate(name string) *template.Template {

	f, err := pkger.Open("/templates/" + name + ".go.tpl")
	if err != nil {
		log.Fatalf("can't read template %s: %v", name, err)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f)
	if err != nil {
		log.Fatalf("can't read template %s: %v", name, err)
	}

	t := template.New("template")
	t, err = t.Parse(buf.String())
	if err != nil {
		log.Fatalf("can't parse template %s: %v", name, err)
	}

	return t
}
