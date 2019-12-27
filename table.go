package xql

// Table ...
// Struct defined for a table object.
type Table struct {
    columns      []*Column
    constraints  []*Constraint
    indexes      []*Index
    primary_keys []*Column
    m_columns    map[string]*Column
    x_columns    map[string]*Column
    j_columns    map[string]*Column
    entity       TableIdentified
    schema       string
}

// TableIdentified which make sure struct have a method TableName()
// This is mandatory for a struct can use as Table structure.
type TableIdentified interface {
    TableName() string
}

// TableIgnored
// Which allow struct to define a method Ignore() to tell ignore elements for table columns
type TableIgnored interface {
    Ignore() []string
}

// TableConstrainted
// Which allow struct to define a method Constraints() to cunstomize table constraints
// type, fields, statement
type TableConstrainted interface {
    Constraints() [][3]string
}

// TableIndexed
// Which allow struct to define a method Indexes() to define table indexes
// type, fields, statement
type TableIndexed interface {
    Indexes() [][2]string
}

// TablePreInsert
// Which entity implemented will be called when xql create a new struct instance.
type TablePreInsert interface {
    PreInsert(*Table, *Session) error
}

type TablePreUpdate interface {
    PreUpdate(*Table, *Session) error
}

type TablePreDelete interface {
    PreDelete(*Table, *Session) error
}

type TablePostInsert interface {
    PostInsert(*Table, *Session) error
}

type TablePostUpdate interface {
    PostUpdate(*Table, *Session) error
}

type TablePostDelete interface {
    PostDelete(*Table, *Session) error
}

type TableCreatable interface {
    Creatable() bool
}

type TableUpdatable interface {
    Updatable() bool
}

type TableReadable interface {
    Readable() bool
}

type TableDeletable interface {
    Deletable() bool
}

func (t *Table) TableName() string {
    if t.schema != "" {
        return t.schema + "." + t.entity.TableName()
    }
    return t.entity.TableName()
}

func (t *Table) BaseTableName() string {
    return t.entity.TableName()
}

func (t *Table) GetColumns() []*Column {
    return t.columns
}

func (t *Table) GetConstraints() []*Constraint {
    return t.constraints
}

func (t *Table) GetIndexes() []*Index {
    return t.indexes
}

func (t *Table) GetPrimaryKeys() []*Column {
    return t.primary_keys
}

func (t *Table) SetSchema(s string) {
    t.schema = s
}

func (t *Table) Schema() string {
    return t.schema
}

func (t *Table) GetColumn(name string) (*Column, bool) {
    if c, ok := t.m_columns[name]; ok {
        return c, true
    }
    if c, ok := t.x_columns[name]; ok {
        return c, true
    }
    if c, ok := t.j_columns[name]; ok {
        return c, true
    }
    return nil, false
}
