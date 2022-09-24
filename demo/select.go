package orm

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/demo/internal/errs"
	"reflect"
	"strings"
)

// Selector 用于构造 SELECT 语句
type Selector[T any] struct {
	sb strings.Builder
	args []any
	table string
	where []Predicate
	model *Model

	db *DB
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}
	// 先看一下你返回了哪些列
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	rows.Next()
	vals := make([]any, 0, len(cols))
	eleVals := make([]reflect.Value, 0, len(cols))
	for _, col := range cols {
		fd := s.model.columnMap[col]
		// fd.Type 是 int，那么  reflect.New(fd.typ) 是 *int
		fdVal := reflect.New(fd.typ)
		eleVals = append(eleVals, fdVal.Elem())

		// 因为 Scan 要指针，所以我们在这里，不需要调用 Elem
		vals = append(vals, fdVal.Interface())
	}
	// 要把 cols 映射过去字段

	err = rows.Scan(vals...)
	if err != nil {
		return nil, err
	}
	// 咋办呢？我已经有 vals 了，接下来咋办？ vals= [123, "Ming", 18, "Deng"]

	// 反射放回去 t 里面

	t := new(T)
	tVal := reflect.ValueOf(t).Elem()
	for i, col := range cols {
		fd := s.model.columnMap[col]
		tVal.FieldByName(fd.goName).Set(eleVals[i])
	}

	return t, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	// var db *sql.DB
	// q, err := s.Build()
	// if err != nil {
	// 	return nil, err
	// }
	// rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	// if err != nil {
	// 	return nil, err
	// }
	// 想办法，把 rows 所有行转换为 []*T
	panic("implement me")
}

// From 指定表名，如果是空字符串，那么将会使用默认表名
func (s *Selector[T]) From(tbl string) *Selector[T] {
	s.table = tbl
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	t := new(T)
	var err error
	s.model, err = s.db.r.Get(t)
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT * FROM ")
	if s.table == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.tableName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.table)
	}

	// 构造 WHERE
	if len(s.where) > 0 {
		// 类似这种可有可无的部分，都要在前面加一个空格
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}
	s.sb.WriteString(";")
	return &Query{
		SQL: s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}
	switch exp := e.(type) {
	case Column:
		s.sb.WriteByte('`')
		fd, ok := s.model.fieldMap[exp.name]
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteByte('?')
		s.args = append(s.args, exp.val)
	case Predicate:
		_, lp := exp.left.(Predicate)
		if lp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			s.sb.WriteByte(')')
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')

		_, rp := exp.right.(Predicate)
		if rp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			s.sb.WriteByte(')')
		}
	default:
		return errs.NewErrUnsupportedExpressionType(exp)
	}
	return nil
}

// Where 用于构造 WHERE 查询条件。如果 ps 长度为 0，那么不会构造 WHERE 部分
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

// cols 是用于 WHERE 的列，难以解决 And Or 和 Not 等问题
// func (s *Selector[T]) Where(cols []string, args...any) *Selector[T] {
// 	s.whereCols = cols
// 	s.args = append(s.args, args...)
// }

// 最为灵活的设计
// func (s *Selector[T]) Where(where string, args...any) *Selector[T] {
// 	s.where = where
// 	s.args = append(s.args, args...)
// }

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}