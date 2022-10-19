package orm

type TableReference interface {
	tableAlias() string
}

// 普通表
type Table struct {
	entity any
	alias string
}

// A JOIN B
func TableOf(entity any) Table{
	return Table{
		entity: entity,
	}
}

func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias: alias,
	}
}

// func (t *Table) As(alias string) {
//
// }

// A JOIN B
// A JOIN (B JOIN C)
// (A JOIN B) JOIN (C JOIN D)
func (t Table) Join(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left: t,
		typ: "JOIN",
		right: right,
	}
}

// func (t Join) Column(right TableReference) *JoinBuilder {
//
// }

func (t Table) LeftJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left: t,
		typ: "LEFT JOIN",
		right: right,
	}
}

func (t Table) tableAlias() string {
	return t.alias
}

// func (t Table) C(name string) Column {
//
// }


// JOIN 查询
type Join struct {
	left TableReference
	// JOIN, LEFT JOIN, LEFT OUTER JOIN
	typ string
	right TableReference
	on []Predicate
	using []string
}


func (t Join) tableAlias() string {
	return ""
}

type JoinBuilder struct {
	left TableReference
	// JOIN, LEFT JOIN, LEFT OUTER JOIN
	typ string
	right TableReference

}

func (jb *JoinBuilder) On(ps...Predicate) Join {
	return Join{
		left: jb.left,
		typ: jb.typ,
		right: jb.right,
		on: ps,
	}
}

func (jb *JoinBuilder) Using(cols...string) Join {
	return Join{
		left: jb.left,
		typ: jb.typ,
		right: jb.right,
		using: cols,
	}
}

// 子查询
type Subquery struct {

}

// func (t Subquery) C(name string) Column {
//
// }
