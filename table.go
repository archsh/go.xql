package xql

// Table ...
// Struct defined for a table object.
type Table struct {
    fields       []*Column
    constraints  []*Constraint
    indexes      []*Index
    primary_keys []*Column
    foreign_keys []*Column
    m_fields     map[string]*Column
    x_fields     map[string]*Column
    j_fields     map[string]*Column
    entity       TableIdentified
    schema       string
}

// TableIdentified which make sure struct have a method TableName()
// This is mandatory for a struct can use as Table structure.
type TableIdentified interface {
    TableName() string
}

// TableIgnored
// Which allow struct to define a method Ignore() to tell ignore elements for table fields
type TableIgnored interface {
    Ignore() []string
}

// TablePrimaryKeyed
// Which allow struct to define a method PrimaryKey() to have multiple primary keys
type TablePrimaryKeyed interface {
    PrimaryKey() []string
}

// TableUniqued
// Which allow struct to define a method Unique() to define unique constraint with multiple fields
type TableUniqued interface {
    Unique() []string
}


// TableChecked
// Which allow struct to define a method Check() to define customized check constraint on table
type TableChecked interface {
    Check() []string
}

// TableExcluded
// Which allow struct to define a method Exclude() to define a customized exclude constraint on table
type TableExcluded interface {
    Exclude() []string
}

// TableInitRequired
// Which entity implemented will be called when xql create a new struct instance.
type TableInitRequired interface {
    Initialize()
}


func (t *Table) TableName() string {
    if t.schema != "" {
        return t.schema + "." + t.entity.TableName()
    }
    return t.entity.TableName()
}

func (t *Table) GetFields() []*Column {
    return t.fields
}

func (t *Table) GetConstraints() []*Constraint {
    return t.constraints
}

func (t *Table) GetIndexes() []*Index {
    return t.indexes
}

func (t *Table) GetField(name string) (*Column, bool) {
    if c, ok := t.m_fields[name]; ok {
        return c, true
    }
    if c, ok := t.x_fields[name]; ok {
        return c, true
    }
    if c, ok := t.j_fields[name]; ok {
        return c, true
    }
    return nil, false
}

// DeclareTable
// Which declare a new Table instance according to a given entity.
func DeclareTable(entity TableIdentified, schema ...string) *Table {
    t := &Table{
        entity: entity,
        fields: makeColumns(entity, false),
    }
    if len(schema) > 0 {
        t.schema = schema[0]
    }
    t.constraints = makeConstraints(t.fields...)
    t.indexes = makeIndexes(t.fields...)
    t.x_fields = make(map[string]*Column)
    t.j_fields = make(map[string]*Column)
    for _, f := range t.fields {
        t.x_fields[f.FieldName] = f
        t.m_fields[f.ElemName] = f
        t.j_fields[f.Jtag] = f
    }
    return t
}
