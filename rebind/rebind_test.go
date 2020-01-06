package rebind_test

import (
	"testing"

	"github.com/pdk/gocrud/rebind"
)

func TestConvertsToDollar(t *testing.T) {

	result := rebind.ToDollar("select * from foo where x = ? and y <> ?")
	exp := "select * from foo where x = $1 and y <> $2"
	if result != exp {
		t.Errorf("expected %s, but got %s", exp, result)
	}
}

func TestConvertsToAtSign(t *testing.T) {

	result := rebind.ToAtSign("select * from foo where x = ? and y <> ?")
	exp := "select * from foo where x = @p1 and y <> @p2"
	if result != exp {
		t.Errorf("expected %s, but got %s", exp, result)
	}
}

func TestConvertsToNamed(t *testing.T) {

	result := rebind.ToNamed("select * from foo where x = ? and y <> ?")
	exp := "select * from foo where x = :arg1 and y <> :arg2"
	if result != exp {
		t.Errorf("expected %s, but got %s", exp, result)
	}
}
