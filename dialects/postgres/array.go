package postgres

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lib/pq"

	"github.com/archsh/go.xql"
)

type Elemented interface {
	Elem2Strings() []string
	Strings2Elem(...string) error
}

type StringArray []string

func (h StringArray) Declare(props xql.PropertySet) string {
	size, _ := props.GetInt("size", 32)
	return fmt.Sprintf("varchar(%d)[]", size)
}
//
//func (a StringArray) Elem2Strings() []string {
//	var ss []string
//	for _, x := range a {
//		ss = append(ss,
//			strings.Replace(
//				strings.Replace(x, `'`, `\'`, -1), `"`, `\"`, -1))
//	}
//	return ss
//}
//
//func (a *StringArray) Strings2Elem(ss ...string) error {
//	(*a) = ss
//	return nil
//}

func (p *StringArray) Scan(src interface{}) error {
	return pq.Array(p).Scan(src)
	//return Array_Scan(src, p)
}

func (p StringArray) Value() (driver.Value, error) {
	return pq.Array(p).Value()
	//return Array_Value(&p)
}

type IntegerArray []int

func (h IntegerArray) Declare(props xql.PropertySet) string {
	return "integer[]"
}

//func (a IntegerArray) Elem2Strings() []string {
//	ss := []string{}
//	for _, x := range a {
//		ss = append(ss, fmt.Sprintf("%d", x))
//	}
//	return ss
//}
//
//func (a *IntegerArray) Strings2Elem(ss ...string) error {
//	for _, s := range ss {
//		n, e := strconv.ParseInt(s, 10, 32)
//		if nil != e {
//			return e
//		} else {
//			(*a) = append(*a, int(n))
//		}
//	}
//	return nil
//}

func (p *IntegerArray) Scan(src interface{}) error {
	return pq.Array(p).Scan(src)
	//return Array_Scan(src, p)
}

func (p IntegerArray) Value() (driver.Value, error) {
	return pq.Array(p).Value()
	//return Array_Value(&p)
}

type SmallIntegerArray []int16

func (h SmallIntegerArray) Declare(props xql.PropertySet) string {
	return "smallint[]"
}

//func (a SmallIntegerArray) Elem2Strings() []string {
//	ss := []string{}
//	for _, x := range a {
//		ss = append(ss, fmt.Sprintf("%d", x))
//	}
//	return ss
//}
//
//func (a *SmallIntegerArray) Strings2Elem(ss ...string) error {
//	for _, s := range ss {
//		n, e := strconv.ParseInt(s, 10, 16)
//		if nil != e {
//			return e
//		} else {
//			(*a) = append(*a, int16(n))
//		}
//	}
//	return nil
//}

func (p *SmallIntegerArray) Scan(src interface{}) error {
	return pq.Array(p).Scan(src)
	//return Array_Scan(src, p)
}

func (p SmallIntegerArray) Value() (driver.Value, error) {
	return pq.Array(p).Value()
	//return Array_Value(&p)
}

type BigIntegerArray []int64

func (h BigIntegerArray) Declare(props xql.PropertySet) string {
	return "bigint[]"
}

//func (a BigIntegerArray) Elem2Strings() []string {
//	ss := []string{}
//	for _, x := range a {
//		ss = append(ss, fmt.Sprintf("%d", x))
//	}
//	return ss
//}
//
//func (a *BigIntegerArray) Strings2Elem(ss ...string) error {
//	for _, s := range ss {
//		n, e := strconv.ParseInt(s, 10, 64)
//		if nil != e {
//			return e
//		} else {
//			(*a) = append(*a, n)
//		}
//	}
//	return nil
//}

func (p *BigIntegerArray) Scan(src interface{}) error {
	return pq.Array(p).Scan(src)
	//return Array_Scan(src, p)
}

func (p BigIntegerArray) Value() (driver.Value, error) {
	return pq.Array(p).Value()
	//return Array_Value(&p)
}

type RealArray []float32

func (h RealArray) Declare(props xql.PropertySet) string {
	return "real[]"
}

//func (a RealArray) Elem2Strings() []string {
//	ss := []string{}
//	for _, x := range a {
//		ss = append(ss, fmt.Sprintf("%f", x))
//	}
//	return ss
//}
//
//func (a *RealArray) Strings2Elem(ss ...string) error {
//	for _, s := range ss {
//		n, e := strconv.ParseFloat(s, 32)
//		if nil != e {
//			return e
//		} else {
//			(*a) = append(*a, float32(n))
//		}
//	}
//	return nil
//}

func (p *RealArray) Scan(src interface{}) error {
	return pq.Array(p).Scan(src)
	//return Array_Scan(src, p)
}

func (p RealArray) Value() (driver.Value, error) {
	return pq.Array(p).Value()
	//return Array_Value(&p)
}

type DoubleArray []float64

func (h DoubleArray) Declare(props xql.PropertySet) string {
	return "double[]"
}

//func (a DoubleArray) Elem2Strings() []string {
//	ss := []string{}
//	for _, x := range a {
//		ss = append(ss, fmt.Sprintf("%f", x))
//	}
//	return ss
//}
//
//func (a *DoubleArray) Strings2Elem(ss ...string) error {
//	for _, s := range ss {
//		n, e := strconv.ParseFloat(s, 64)
//		if nil != e {
//			return e
//		} else {
//			(*a) = append(*a, n)
//		}
//	}
//	return nil
//}

func (p *DoubleArray) Scan(src interface{}) error {
	return pq.Array(p).Scan(src)
	//return Array_Scan(src, p)
}

func (p DoubleArray) Value() (driver.Value, error) {
	return pq.Array(p).Value()
	//return Array_Value(&p)
}

type BoolArray []bool

func (h BoolArray) Declare(props xql.PropertySet) string {
	return "bool[]"
}

//func (a BoolArray) Elem2Strings() []string {
//	ss := []string{}
//	for _, x := range a {
//		ss = append(ss, fmt.Sprintf("%s", x))
//	}
//	return ss
//}
//
//func (a *BoolArray) Strings2Elem(ss ...string) error {
//	for _, s := range ss {
//		switch strings.ToLower(s) {
//		case "y", "yes", "t", "true", "ok":
//			(*a) = append(*a, true)
//		default:
//			(*a) = append(*a, false)
//		}
//	}
//	return nil
//}

func (p *BoolArray) Scan(src interface{}) error {
	return pq.Array(p).Scan(src)
	//return Array_Scan(src, p)
}

func (p BoolArray) Value() (driver.Value, error) {
	return pq.Array(p).Value()
	//return Array_Value(&p)
}

// PARSING ARRAYS
// SEE http://www.postgresql.org/docs/9.1/static/arrays.html#ARRAYS-IO
// Arrays are output within {} and a delimiter, which is a comma for most
// postgres types (; for box)
//
// Individual values are surrounded by quotes:
// The array output routine will put double quotes around element values if
// they are empty strings, contain curly braces, delimiter characters,
// double quotes, backslashes, or white space, or match the word NULL.
// Double quotes and backslashes embedded in element values will be
// backslash-escaped. For numeric data types it is safe to assume that double
// quotes will never appear, but for textual data types one should be prepared
// to cope with either the presence or absence of quotes.

// construct a regexp to extract values:
var (
	// unquoted array values must not contain: (" , \ { } whitespace NULL)
	// and must be at least one char
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	// quoted array values are surrounded by double quotes, can be any
	// character except " or \, which must be backslash escaped:
	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	// an array value may be either quoted or unquoted:
	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	// Array values are separated with a comma IF there is more than one value:
	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	valueIndex int
)

// Find the index of the 'value' named expression
func init() {
	for i, subexp := range arrayExp.SubexpNames() {
		if subexp == "value" {
			valueIndex = i
			break
		}
	}
}

// Parse the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func parseArray(array string) []string {
	results := make([]string, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		results = append(results, s)
	}
	return results
}

//func (p *StringArray) Scan(src interface{}) error {
//    asBytes, ok := src.([]byte)
//    if !ok {
//        return error(errors.New("Scan source was not []bytes."))
//    }
//
//    asString := string(asBytes)
//    parsed := parseArray(asString)
//    (*p) = StringArray(parsed)
//
//    return nil
//}
//
//func (p StringArray) Value() (driver.Value, error) {
//    var ss []string
//    for _, s := range p {
//        ss = append(ss, fmt.Sprintf(`"%s"`, s))
//    }
//    return strings.Join([]string{"{", strings.Join(ss, ","), "}"}, ""), nil
//}

func Array_Scan(src interface{}, dest interface{}) error {
	if nil == src || dest == nil {
		return nil
	}
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes."))
	}

	asString := string(asBytes)
	parsed := parseArray(asString)
	if vv, ok := dest.(Elemented); ok {
		return vv.Strings2Elem(parsed...)
	}
	return errors.New("Elemented should be implemented.")
}

func Array_Value(v interface{}) (driver.Value, error) {
	if nil == v {
		return nil, nil
	}
	if vv, ok := v.(Elemented); ok {
		return strings.Join([]string{"{", strings.Join(vv.Elem2Strings(), ","), "}"}, ""), nil
	}
	return nil, errors.New("Elemented should be implemented.")
}
