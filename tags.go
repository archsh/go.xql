package xql

import (
	"strconv"
	"strings"
)

type PropertySet map[string]string

const (
	SingleQuoteOpened uint = 0x01
	DoubleQuoteOpened uint = 0x02
	SBraceOpened      uint = 0x04
	MBraceOpened      uint = 0x08
	BBraceOpened      uint = 0x10
)

func ParseDottedArgs(s string) (ret []string) {
	var opened uint
	var chars []rune
	for _, c := range s {
		switch c {
		case '\'':
			opened ^= SingleQuoteOpened
			chars = append(chars, c)
		case '"':
			opened ^= DoubleQuoteOpened
			chars = append(chars, c)
		case '(':
			opened |= SBraceOpened
			chars = append(chars, c)
		case ')':
			opened &= ^SBraceOpened
			chars = append(chars, c)
		case '{':
			opened |= BBraceOpened
			chars = append(chars, c)
		case '}':
			opened &= ^BBraceOpened
			chars = append(chars, c)
		case '[':
			opened |= MBraceOpened
			chars = append(chars, c)
		case ']':
			opened &= ^MBraceOpened
			chars = append(chars, c)
		case ',':
			if opened == 0 && len(chars) > 0 {
				ret = append(ret, string(chars))
			}
			chars = []rune{}
			/*else {
				chars = append(chars, c)
			}*/
		default:
			chars = append(chars, c)
		}
	}
	if len(chars) > 0 {
		ret = append(ret, string(chars))
	}
	return
}

func ParseProperties(s string) (PropertySet, error) {
	p := make(PropertySet)
	if s == "" {
		return p, nil
	}
	for _, ss := range ParseDottedArgs(s) {
		if ss == "" {
			continue
		}
		ks := strings.SplitN(ss, "=", 2)
		if len(ks) > 1 {
			p[ks[0]] = ks[1]
		} else {
			p[ks[0]] = "t"
		}
	}

	return p, nil
}

func (h PropertySet) HasKey(k string) bool {
	if nil == h {
		return false
	}
	_, ok := h[k]
	return ok
}

func (h PropertySet) GetInt(k string, defaults ...int) (int, bool) {
	d := 0
	if len(defaults) > 0 {
		d = defaults[0]
	}
	if nil == h {
		return d, false
	}
	if ns, ok := h[k]; ok {
		i, e := strconv.ParseInt(ns, 10, 32)
		if nil != e {
			return d, false
		} else {
			return int(i), true
		}
	}
	return d, false
}

func (h PropertySet) PopInt(k string, defaults ...int) (int, bool) {
	if v, ok := h.GetInt(k, defaults...); ok {
		delete(h, k)
		return v, ok
	} else {
		return v, ok
	}
}

func (h PropertySet) GetInt64(k string, defaults ...int64) (int64, bool) {
	var d int64 = 0
	if len(defaults) > 0 {
		d = defaults[0]
	}
	if nil == h {
		return d, false
	}
	if ns, ok := h[k]; ok {
		i, e := strconv.ParseInt(ns, 10, 32)
		if nil != e {
			return d, false
		} else {
			return i, true
		}
	}
	return d, false
}

func (h PropertySet) PopInt64(k string, defaults ...int64) (int64, bool) {
	if v, ok := h.GetInt64(k, defaults...); ok {
		delete(h, k)
		return v, ok
	} else {
		return v, ok
	}
}

func (h PropertySet) GetUInt(k string, defaults ...uint) (uint, bool) {
	var d uint = 0
	if len(defaults) > 0 {
		d = defaults[0]
	}
	if ns, ok := h[k]; ok {
		i, e := strconv.ParseUint(ns, 10, 32)
		if nil != e {
			return d, false
		} else {
			return uint(i), true
		}
	}
	return d, false
}

func (h PropertySet) PopUInt(k string, defaults ...uint) (uint, bool) {
	if v, ok := h.GetUInt(k, defaults...); ok {
		delete(h, k)
		return v, ok
	} else {
		return v, ok
	}
}

func (h PropertySet) GetString(k string, defaults ...string) (string, bool) {
	d := ""
	if len(defaults) > 0 {
		d = defaults[0]
	}
	if nil == h {
		return d, false
	}
	if ns, ok := h[k]; ok {
		return ns, true
	}
	return d, false
}

func (h PropertySet) PopString(k string, defaults ...string) (string, bool) {
	if v, ok := h.GetString(k, defaults...); ok {
		delete(h, k)
		return v, ok
	} else {
		return v, ok
	}
}

func (h PropertySet) GetBool(k string, defaults ...bool) (bool, bool) {
	d := false
	if len(defaults) > 0 {
		d = defaults[0]
	}
	if nil == h {
		return d, false
	}
	if ns, ok := h[k]; ok {
		switch strings.ToLower(ns) {
		case "t", "true", "yes", "ok", "y":
			return true, true
		case "f", "false", "no", "n":
			return false, true
		}
	}
	return d, false
}

func (h PropertySet) PopBool(k string, defaults ...bool) (bool, bool) {
	if v, ok := h.GetBool(k, defaults...); ok {
		delete(h, k)
		return v, ok
	} else {
		return v, ok
	}
}
