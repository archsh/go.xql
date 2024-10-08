package xql

import (
	"fmt"
	"time"
)

// Ref:> https://www.postgresql.org/docs/9.6/static/datatype.html

//Numeric Types

// SmallInteger	2 bytes	small-range integer	-32768 to +32767
type SmallInteger int16

func (s SmallInteger) Declare(props PropertySet) string {
	return "smallint"
}

// Integer	4 bytes	typical choice for integer	-2147483648 to +2147483647
type Integer int

func (s Integer) Declare(props PropertySet) string {
	return "integer"
}

// BigInteger	8 bytes	large-range integer	-9223372036854775808 to +9223372036854775807
type BigInteger int64

func (s BigInteger) Declare(props PropertySet) string {
	return "bigint"
}

// Decimal	variable	user-specified precision, exact	up to 131072 digits before the decimal point; up to 16383 digits after the decimal point
type Decimal string

func (d Decimal) Declare(props PropertySet) string {
	return "decimal"
}

// Numeric	variable	user-specified precision, exact	up to 131072 digits before the decimal point; up to 16383 digits after the decimal point
type Numeric Decimal

// Real	4 bytes	variable-precision, inexact	6 decimal digits precision
type Real float32

func (s Real) Declare(props PropertySet) string {
	return "real"
}

// Double precision	8 bytes	variable-precision, inexact	15 decimal digits precision
type Double float64

func (s Double) Declare(props PropertySet) string {
	return "double"
}

// SmallSerial	2 bytes	small autoincrementing integer	1 to 32767
type SmallSerial uint16

func (s SmallSerial) Declare(props PropertySet) string {
	return "smallserial"
}

// Serial	4 bytes	autoincrementing integer	1 to 2147483647
type Serial uint

func (s Serial) Declare(props PropertySet) string {
	return "serial"
}

// BigSerial	8 bytes	large autoincrementing integer	1 to 9223372036854775807
type BigSerial uint64

func (s BigSerial) Declare(props PropertySet) string {
	return "bigserial"
}

//Character Types

// Varchar character varying(n), varchar(n)	variable-length with limit
type Varchar string

func (s Varchar) Declare(props PropertySet) string {
	length, _ := props.GetUInt("size", 32)
	return fmt.Sprintf("character varying(%d)", length)
}

// Char character(n), char(n)	fixed-length, blank padded
type Char string

func (s Char) Declare(props PropertySet) string {
	length, _ := props.GetUInt("size", 32)
	return fmt.Sprintf("character(%d)", length)
}

// Text	variable unlimited length
type Text string

func (s Text) Declare(props PropertySet) string {
	return "text"
}

// Bit String Types
type Bit string

func (b Bit) Declare(props PropertySet) string {
	length, _ := props.GetUInt("size", 1)
	return fmt.Sprintf("bit(%d)", length)
}

// Bitvar bit var
type Bitvar string

func (b Bitvar) Declare(props PropertySet) string {
	length, _ := props.GetUInt("size", 1)
	return fmt.Sprintf("bit varying(%d)", length)
}

// Binary Data Types

// Bytea	1 or 4 bytes plus the actual binary string	variable-length binary string
type Bytea []byte

func (b Bytea) Declare(props PropertySet) string {
	return "bytea"
}

// Date/Time Types

// TimeStamp [ (p) ] [ without time zone ]	8 bytes	both date and time (no time zone)	4713 BC	294276 AD	1 microsecond / 14 digits
// TimeStamp [ (p) ] with time zone	8 bytes	both date and time, with time zone	4713 BC	294276 AD	1 microsecond / 14 digits
type TimeStamp time.Time

func (s TimeStamp) Declare(props PropertySet) string {
	return "timestamp"
}

// Date	4 bytes	date (no time of day)	4713 BC	5874897 AD	1 day
type Date time.Time

func (s Date) Declare(props PropertySet) string {
	return "date"
}

// Time [ (p) ] [ without time zone ]	8 bytes	time of day (no date)	00:00:00	24:00:00	1 microsecond / 14 digits
// Time [ (p) ] with time zone	12 bytes	times of day only, with time zone	00:00:00+1459	24:00:00-1459	1 microsecond / 14 digits
type Time time.Time

func (s Time) Declare(props PropertySet) string {
	return "time"
}

// Interval [ columns ] [ (p) ]	16 bytes	time interval	-178000000 years	178000000 years	1 microsecond / 14 digits
type Interval time.Duration

func (s Interval) Declare(props PropertySet) string {
	return "interval"
}

// Boolean Data Type
//
// Name	Storage Size	Description
// boolean	1 byte	state of true or false
type Boolean bool

func (s Boolean) Declare(props PropertySet) string {
	return "boolean"
}

//Declaration of Enumerated Types

//Enum types are created using the CREATE TYPE command, for example:

// Enum CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy');
type Enum string

func (s Enum) Declare(props PropertySet) string {
	return "boolean"
}

// UUID Type
type UUID string

func (s UUID) Declare(props PropertySet) string {
	return "uuid"
}
