package xql

import (
    "errors"
    "strings"
    "fmt"
    "regexp"
    "database/sql"
    "database/sql/driver"
    "encoding/json"
)

type JSONDictionary map[string]interface{}
type JSONDictionaryArray []JSONDictionary
type HSTOREDictionary map[string]interface{}
type StringArray []string
type IntegerArray []int
type SmallIntegerArray []int16
type BoolArray []bool


func (p *JSONDictionary) Scan(src interface{}) error {
    if nil == src {
        *p = nil
        return nil
    }
    source, ok := src.([]byte)
    if !ok {
        return errors.New("Type assertion .([]byte) failed.")
    }
    var i JSONDictionary
    err := json.Unmarshal(source, &i)
    if err != nil {
        return err
    }
    *p = i
    return nil
}

func (p JSONDictionary) Value() (driver.Value, error) {
    j, err := json.Marshal(p)
    return j, err
}


func (p *JSONDictionaryArray) Scan(src interface{}) error {
    source, ok := src.([]byte)
    if !ok {
        return errors.New("Type assertion .([]byte) failed.")
    }
    var i JSONDictionaryArray
    err := json.Unmarshal(source, &i)
    if err != nil {
        return err
    }
    *p = i
    return nil
}

func (p JSONDictionaryArray) Value() (driver.Value, error) {
    j, err := json.Marshal(p)
    return j, err
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
    unquotedChar = `[^",\\{}\s(NULL)]`
    unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

    // quoted array values are surrounded by double quotes, can be any
    // character except " or \, which must be backslash escaped:
    quotedChar = `[^"\\]|\\"|\\\\`
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

func (p *StringArray) Scan(src interface{}) error {
    asBytes, ok := src.([]byte)
    if !ok {
        return error(errors.New("Scan source was not []bytes"))
    }

    asString := string(asBytes)
    parsed := parseArray(asString)
    (*p) = StringArray(parsed)

    return nil
}

func (p StringArray) Value() (driver.Value, error) {
    j, err := json.Marshal(p)
    return j, err
}



// escapes and quotes hstore keys/values
// s should be a sql.NullString or string
func hQuote(s interface{}) string {
    var str string
    switch v := s.(type) {
    case sql.NullString:
        if !v.Valid {
            return "NULL"
        }
        str = v.String
    case string:
        str = v
    default:
        panic("not a string or sql.NullString")
    }

    str = strings.Replace(str, "\\", "\\\\", -1)
    return `"` + strings.Replace(str, "\"", "\\\"", -1) + `"`
}

// Scan implements the Scanner interface.
//
// Note h.Map is reallocated before the scan to clear existing values. If the
// hstore column's database value is NULL, then h.Map is set to nil instead.
func (h *HSTOREDictionary) Scan(value interface{}) error {
    if value == nil {
        *h = nil
        return nil
    }
    m := make(map[string]interface{})
    var b byte
    pair := [][]byte{{}, {}}
    pi := 0
    inQuote := false
    didQuote := false
    sawSlash := false
    bindex := 0
    for bindex, b = range value.([]byte) {
        if sawSlash {
            pair[pi] = append(pair[pi], b)
            sawSlash = false
            continue
        }

        switch b {
        case '\\':
            sawSlash = true
            continue
        case '"':
            inQuote = !inQuote
            if !didQuote {
                didQuote = true
            }
            continue
        default:
            if !inQuote {
                switch b {
                case ' ', '\t', '\n', '\r':
                    continue
                case '=':
                    continue
                case '>':
                    pi = 1
                    didQuote = false
                    continue
                case ',':
                    s := string(pair[1])
                    if !didQuote && len(s) == 4 && strings.ToLower(s) == "null" {
                        m[string(pair[0])] = nil // sql.NullString{String: "", Valid: false}
                    } else {
                        m[string(pair[0])] = string(pair[1]) // sql.NullString{String: string(pair[1]), Valid: true}
                    }
                    pair[0] = []byte{}
                    pair[1] = []byte{}
                    pi = 0
                    continue
                }
            }
        }
        pair[pi] = append(pair[pi], b)
    }
    if bindex > 0 {
        s := string(pair[1])
        if !didQuote && len(s) == 4 && strings.ToLower(s) == "null" {
            m[string(pair[0])] = nil // sql.NullString{String: "", Valid: false}
        } else {
            m[string(pair[0])] = string(pair[1]) // sql.NullString{String: string(pair[1]), Valid: true}
        }
    }
    *h = m
    return nil
}

// Value implements the driver Valuer interface. Note if h.Map is nil, the
// database column value will be set to NULL.
func (h HSTOREDictionary) Value() (driver.Value, error) {
    if h == nil {
        return nil, nil
    }
    parts := []string{}
    for key, val := range h {
        thispart := hQuote(key) + "=>" + hQuote(val)
        parts = append(parts, thispart)
    }
    return []byte(strings.Join(parts, ",")), nil
}

