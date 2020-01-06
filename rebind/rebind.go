package rebind

import (
	"strconv"
	"strings"
)

// ToDollar converts question marks to $N bindmarkers, suitable for
// postgres and descendents.
func ToDollar(query string) string {

	// Start with enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	for i, j := strings.Index(query, "?"), 1; i != -1; i, j = strings.Index(query, "?"), j+1 {

		rqb = append(rqb, query[:i]...)
		rqb = append(rqb, '$')
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}

// ToNamed converts question marks to :argN bindmarkers, suitable for Oracle.
func ToNamed(query string) string {

	// Start with enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	for i, j := strings.Index(query, "?"), 1; i != -1; i, j = strings.Index(query, "?"), j+1 {

		rqb = append(rqb, query[:i]...)
		rqb = append(rqb, ':', 'a', 'r', 'g')
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}

// ToAtSign converts question marks to @pN bindmarkers, suitable for sqlserver.
func ToAtSign(query string) string {

	// Start with enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	for i, j := strings.Index(query, "?"), 1; i != -1; i, j = strings.Index(query, "?"), j+1 {

		rqb = append(rqb, query[:i]...)
		rqb = append(rqb, '@', 'p')
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}
