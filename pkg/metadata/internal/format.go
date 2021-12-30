package internal

import (
	"regexp"
	"strings"

	"github.com/lib/pq/oid"
)

var (
	variableNameRuleRegexp *regexp.Regexp
)

func init() {
	variableNameRuleRegexp = regexp.MustCompile(`(?i)(\?|(sql)|(id)|(uri)|(url))`)
}

func TypeString(typ oid.Oid) string {
	switch typ {
	case oid.T_char, oid.T_varchar, oid.T_text:
		return "string"
	case oid.T_bytea:
		return "[]byte"
	case oid.T_timestamptz:
		fallthrough
	case oid.T_timestamp, oid.T_date:
		fallthrough
	case oid.T_time:
		fallthrough
	case oid.T_timetz:
		return "time.Time"
	case oid.T_bool:
		return "bool"
	case oid.T_int8, oid.T_int4, oid.T_int2:
		return "int64"
	case oid.T_float4, oid.T_float8:
		return "float64"
	}

	return "[]byte"
}

func VariableNameRule(in string) (out string) {
	return variableNameRuleRegexp.ReplaceAllStringFunc(in, func(w string) string {
		return strings.ToUpper(w)
	})
}

func Camelify(in string) (out string) {
	phrase := strings.Split(in, "_")
	for i, v := range phrase {
		phrase[i] = strings.Join([]string{strings.ToUpper(string(v[0])), string(v[1:])}, "")
	}
	return strings.Join(phrase, "")
}
