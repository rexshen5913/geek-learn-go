package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/demo/internal/errs"
	"reflect"
)

// 带校验的 option 模式
type ModelOpt func(m *Model) error

type Model struct {
	// 结果体对应的表名
	tableName string
	// 字段名到字段的元数据
	fieldMap map[string]*field

	// 列名到字段的映射
	columnMap map[string]*field
}

func ModelWithTableName (name string) ModelOpt {
	return func(m *Model) error {
		// if name == "" {
		// 	return errs.ErrEmptyTableName
		// }
		m.tableName = name
		return nil
	}
}

func ModelWithColumnName(field string, colName string) ModelOpt {
	return func(m *Model) error {
		fd, ok := m.fieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.colName = colName
		return nil
	}
}

func ModelWithColumn(field string, col *field) ModelOpt {
	return func(m *Model) error {
		m.fieldMap[field] = col
		return nil
	}
}

func ModelWithColumnAutoIncrement(field string) ModelOpt {
	return func(m *Model) error {
		m.fieldMap[field].autoIncrement = true
		return nil
	}
}

type field struct {
	// 字段名
	goName string
	// 字段对应的列名
	colName string

	typ reflect.Type

	autoIncrement bool
}

type TableName interface {
	TableName() string
}


