# Dialect

`Dialect` defines the low level of driver implements.

## Dialect Interface

```go
package xql

type Dialect interface {
    PrepareSQL(QuerySet) string
    Exec(QuerySet, ...args)
    
}
```