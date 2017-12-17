package xql

import (
    "strings"
)

type PropertySet struct {
    FieldName    string
    PropertyName string

    extras map[string]string
}

func ParseProperties(s string) (*PropertySet, error) {
    p := &PropertySet{
        extras: make(map[string]string),
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

func (h PropertySet) GetInt(k string, defaults ... int) (int, bool) {
    d := 0
    if len(defaults) > 0 {
        d = defaults[0]
    }
    if nil == h {
        return d, false
    }
    if ns, ok := h[k]; ok {
        if n, t := ns.(int); t {
            return n, t
        }
        if n, t := ns.(int32); t {
            return int(n), t
        }
        if n, t := ns.(int16); t {
            return int(n), t
        }
        if n, t := ns.(int64); t {
            return int(n), t
        }
    }
    return d, false
}

func (h PropertySet) GetInt64(k string, defaults ... int64) (int64, bool) {
    var d int64 = 0
    if len(defaults) > 0 {
        d = defaults[0]
    }
    if nil == h {
        return d, false
    }
    if ns, ok := h[k]; ok {
        if n, t := ns.(int); t {
            return int64(n), t
        }
        if n, t := ns.(int32); t {
            return int64(n), t
        }
        if n, t := ns.(int16); t {
            return int64(n), t
        }
        if n, t := ns.(int64); t {
            return int64(n), t
        }
    }
    return d, false
}

func (h PropertySet) GetUInt(k string, defaults ... uint) (uint, bool) {
    var d uint = 0
    if len(defaults) > 0 {
        d = defaults[0]
    }
    if ns, ok := h[k]; ok {
        if n, t := ns.(uint); t {
            return n, t
        }
        if n, t := ns.(int); t {
            return uint(n), t
        }
        if n, t := ns.(int32); t {
            return uint(n), t
        }
        if n, t := ns.(int16); t {
            return uint(n), t
        }
        if n, t := ns.(int64); t {
            return uint(n), t
        }
    }
    return d, false
}

func (h PropertySet) GetString(k string, defaults ... string) (string, bool) {
    d := ""
    if len(defaults) > 0 {
        d = defaults[0]
    }
    if nil == h {
        return d, false
    }
    if ns, ok := h[k]; ok {
        if s, t := ns.(string); t {
            return s, t
        }
    }
    return d, false
}

func (h PropertySet) GetBool(k string, defaults ... bool) (bool, bool) {
    d := false
    if len(defaults) > 0 {
        d = defaults[0]
    }
    if nil == h {
        return d, false
    }
    if ns, ok := h[k]; ok {
        if b, t := ns.(bool); t {
            return b, t
        }
        if s, t := ns.(string); t {
            switch strings.ToLower(s) {
            case "t", "true", "yes", "ok":
                return true, true
            case "f", "false", "no":
                return false, true
            }
        }

    }
    return d, false
}
