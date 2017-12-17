package xql

import (
    "strings"
    "regexp"
)

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func Camel2Underscore(s string) string {
    var a []string
    for _, sub := range camel.FindAllStringSubmatch(s, -1) {
        if sub[1] != "" {
            a = append(a, sub[1])
        }
        if sub[2] != "" {
            a = append(a, sub[2])
        }
    }
    return strings.ToLower(strings.Join(a, "_"))
}


func inSlice(a string, ls []string) bool {
    for _, s := range ls {
        if a == s {
            return true
        }
    }
    return false
}

func getSkips(tags []string) (skips []string) {
    if nil == tags || len(tags) < 1 {
        return
    }
    for _, tag := range tags {
        if strings.HasPrefix(tag, "skips:") {
            s := strings.TrimLeft(tag, "skips:")
            for _, n := range strings.Split(s, ";") {
                if n != "" {
                    skips = append(skips, n)
                }
            }
            return
        }
    }
    return
}
