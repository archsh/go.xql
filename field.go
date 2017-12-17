package xql

import (
    "time"
    "reflect"
)

type Field struct {
    FieldName  string
    ElemName   string
    Type       reflect.Type
    Indexed    bool
    Nullable   bool
    Unique     bool
    ForeignKey interface{}
    table      interface{}
    extras     map[string]string
}

type Declarable interface {
    Declare(props *PropertySet) string
}

type String string

func (s String) Declare(props *PropertySet) string {
    return ""
}

type UUID string

func (s UUID) Declare(props *PropertySet) string {
    return ""
}

type Text string

func (s Text) Declare(props *PropertySet) string {
    return ""
}

type TinyText string

func (s TinyText) Declare(props *PropertySet) string {
    return ""
}

type MediumText string

func (s MediumText) Declare(props *PropertySet) string {
    return ""
}

type LongText string

func (s LongText) Declare(props *PropertySet) string {
    return ""
}

type Bolb []byte

func (s Bolb) Declare(props *PropertySet) string {
    return ""
}

type TinyBolb []byte

func (s TinyBolb) Declare(props *PropertySet) string {
    return ""
}

type MediumBolb []byte

func (s MediumBolb) Declare(props *PropertySet) string {
    return ""
}

type LongBolb []byte

func (s LongBolb) Declare(props *PropertySet) string {
    return ""
}

type Integer int

func (s Integer) Declare(props *PropertySet) string {
    return ""
}

type SmallInteger int16

func (s SmallInteger) Declare(props *PropertySet) string {
    return ""
}

type TinyInteger int8

func (s TinyInteger) Declare(props *PropertySet) string {
    return ""
}

type BigInteger int64

func (s BigInteger) Declare(props *PropertySet) string {
    return ""
}

type Float float32

func (s Float) Declare(props *PropertySet) string {
    return ""
}

type Double float64

func (s Double) Declare(props *PropertySet) string {
    return ""
}

type Serial uint

func (s Serial) Declare(props *PropertySet) string {
    return ""
}

type SmallSerial uint16

func (s SmallSerial) Declare(props *PropertySet) string {
    return ""
}

type TinySerial uint8

func (s TinySerial) Declare(props *PropertySet) string {
    return ""
}

type BigSerial uint64

func (s BigSerial) Declare(props *PropertySet) string {
    return ""
}

type Enum uint16

func (s Enum) Declare(props *PropertySet) string {
    return ""
}

type Set []uint8

func (s Set) Declare(props *PropertySet) string {
    return ""
}

type Date time.Time

func (s Date) Declare(props *PropertySet) string {
    return ""
}

type Time time.Time

func (s Time) Declare(props *PropertySet) string {
    return ""
}

type DateTime time.Time

func (s DateTime) Declare(props *PropertySet) string {
    return ""
}

type TimeStamp time.Time

func (s TimeStamp) Declare(props *PropertySet) string {
    return ""
}
