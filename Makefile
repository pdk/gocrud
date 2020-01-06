srcfiles = *.go
templatefiles = templates/*.go.tpl

gocrud : $(srcfiles)
	go build -o gocrud *.go

pkged.go : $(templatefiles)
	pkger

test :
	go test ./...

run : gocrud
	./gocrud -template postgres -an -package example -source example/account.go -struct Account -instance acct > example/account_crud.go
