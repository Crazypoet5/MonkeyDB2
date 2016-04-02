package exe

const (
	INT = iota
	FLOAT
	STRING
	ARRAY
	OBJECT
)

type Relation struct {
	ColumnNames []string
	Rows        []Row
}

type Row []Value

type Value struct {
	Kind int
	Raw  []byte
}

func NewRelation() *Relation {
	return &Relation{}
}

func (r *Relation) AddRow(row Row) {
	r.Rows = append(r.Rows, row)
}

func (r *Relation) SetColumnNames(strs []string) {
	copy(r.ColumnNames, strs)
}

func (r *Relation) SetColumnNameAlias(n int, alias string) {
	if n < len(r.ColumnNames) {
		r.ColumnNames[n] = alias
	}
}

func (r *Relation) RowNum() int {
	return len(r.Rows)
}

func (r *Relation) ColumnNum() int {
	if r.ColumnNames != nil {
		return len(r.ColumnNames)
	}
	if len(r.Rows) == 0 {
		return 0
	}
	return len(r.Rows[0])
}

func NewRow(values []Value) Row {
	return Row(values)
}

func NewValue(kind int, raw []byte) Value {
	return Value{
		Kind: kind,
		Raw:  raw,
	}
}

func StringToType(s string) int {
	switch s {
	case "INT":
		return INT
	case "FLOAT":
		return FLOAT
	case "STRING":
		return STRING
	case "ARRAY":
		return ARRAY
	case "OBJECT":
		return OBJECT
	}
	return STRING
}
