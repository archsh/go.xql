package xql

import (
    "fmt"
    "time"
)

// Ref:> https://www.postgresql.org/docs/9.6/static/datatype.html

//Numeric Types

//Name	Storage Size	Description	Range
//smallint	2 bytes	small-range integer	-32768 to +32767
//integer	4 bytes	typical choice for integer	-2147483648 to +2147483647
//bigint	8 bytes	large-range integer	-9223372036854775808 to +9223372036854775807
//decimal	variable	user-specified precision, exact	up to 131072 digits before the decimal point; up to 16383 digits after the decimal point
//numeric	variable	user-specified precision, exact	up to 131072 digits before the decimal point; up to 16383 digits after the decimal point
//real	4 bytes	variable-precision, inexact	6 decimal digits precision
//double precision	8 bytes	variable-precision, inexact	15 decimal digits precision
//smallserial	2 bytes	small autoincrementing integer	1 to 32767
//serial	4 bytes	autoincrementing integer	1 to 2147483647
//bigserial	8 bytes	large autoincrementing integer	1 to 9223372036854775807

//Character Types

//Name	Description
//character varying(n), varchar(n)	variable-length with limit
//character(n), char(n)	fixed-length, blank padded
//text	variable unlimited length


//Binary Data Types
//
//Name	Storage Size	Description
//bytea	1 or 4 bytes plus the actual binary string	variable-length binary string

//Date/Time Types
//
//Name	Storage Size	Description	Low Value	High Value	Resolution
//timestamp [ (p) ] [ without time zone ]	8 bytes	both date and time (no time zone)	4713 BC	294276 AD	1 microsecond / 14 digits
//timestamp [ (p) ] with time zone	8 bytes	both date and time, with time zone	4713 BC	294276 AD	1 microsecond / 14 digits
//date	4 bytes	date (no time of day)	4713 BC	5874897 AD	1 day
//time [ (p) ] [ without time zone ]	8 bytes	time of day (no date)	00:00:00	24:00:00	1 microsecond / 14 digits
//time [ (p) ] with time zone	12 bytes	times of day only, with time zone	00:00:00+1459	24:00:00-1459	1 microsecond / 14 digits
//interval [ fields ] [ (p) ]	16 bytes	time interval	-178000000 years	178000000 years	1 microsecond / 14 digits

//Boolean Data Type
//
//Name	Storage Size	Description
//boolean	1 byte	state of true or false


//Declaration of Enumerated Types

//Enum types are created using the CREATE TYPE command, for example:

//CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy');

type String string

func (s String) Declare(props PropertySet) string {
    length, _ := props.GetUInt("size", 32)
    return fmt.Sprintf("VARCHAR(%d)", length)
}

type UUID string

func (s UUID) Declare(props PropertySet) string {
    return "UUID"
}

type Text string

func (s Text) Declare(props PropertySet) string {
    return "TEXT"
}

type TinyText string

func (s TinyText) Declare(props PropertySet) string {
    return "TINYTEXT"
}

type MediumText string

func (s MediumText) Declare(props PropertySet) string {
    return "MEDIUMTEXT"
}

type LongText string

func (s LongText) Declare(props PropertySet) string {
    return "LONGTEXT"
}

type Bolb []byte

func (s Bolb) Declare(props PropertySet) string {
    return "BLOB"
}

type TinyBolb []byte

func (s TinyBolb) Declare(props PropertySet) string {
    return "TINYBOLB"
}

type MediumBolb []byte

func (s MediumBolb) Declare(props PropertySet) string {
    return "MEDIUMBLOB"
}

type LongBolb []byte

func (s LongBolb) Declare(props PropertySet) string {
    return "LONGBLOB"
}

type Integer int

func (s Integer) Declare(props PropertySet) string {
    return "INT"
}

type SmallInteger int16

func (s SmallInteger) Declare(props PropertySet) string {
    return "SMALLINT"
}

type TinyInteger int8

func (s TinyInteger) Declare(props PropertySet) string {
    return "TINYINT"
}

type BigInteger int64

func (s BigInteger) Declare(props PropertySet) string {
    return "BIGINT"
}

type Float float32

func (s Float) Declare(props PropertySet) string {
    return "FLOAT"
}

type Double float64

func (s Double) Declare(props PropertySet) string {
    return "DOUBLE"
}

type Decimal string

func (d Decimal) Declare(props PropertySet) string {
    return "DECIMAL"
}

type Serial uint

func (s Serial) Declare(props PropertySet) string {
    return "SERIAL"
}

type SmallSerial uint16

func (s SmallSerial) Declare(props PropertySet) string {
    return "SMALLSERIAL"
}

type TinySerial uint8

func (s TinySerial) Declare(props PropertySet) string {
    return "TINYSERIAL"
}

type BigSerial uint64

func (s BigSerial) Declare(props PropertySet) string {
    return "BIGSERIAL"
}

type Enum uint16

func (s Enum) Declare(props PropertySet) string {
    opts, _ := props.GetString("options", "")
    if opts == "" {
        panic("Empty options for Enum is not allowed")
    }
    return fmt.Sprintf("ENUM(%s)", opts)
}

type Set []uint8

func (s Set) Declare(props PropertySet) string {
    opts, _ := props.GetString("options", "")
    if opts == "" {
        panic("Empty options for Set is not allowed")
    }
    return fmt.Sprintf("SET(%s)", opts)
}

type Date time.Time

func (s Date) Declare(props PropertySet) string {
    return "DATE"
}

type Time time.Time

func (s Time) Declare(props PropertySet) string {
    return "TIME"
}

type DateTime time.Time

func (s DateTime) Declare(props PropertySet) string {
    return "DATETIME"
}

type TimeStamp time.Time

func (s TimeStamp) Declare(props PropertySet) string {
    return "TIMESTAMP"
}
