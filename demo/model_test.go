package orm

import (
	"errors"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/demo/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseModel(t *testing.T) {
	tests := []struct {
		name string
		input any
		want *Model
		wantErr error
	}{
		{
			name: "ptr",
			input: &TestModel{},
			want: &Model{
				tableName: "test_model",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
		},
		{
			name: "struct",
			input: TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name: "map",
			input: map[string]string{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name: "nil",
			input: nil,
			wantErr: errors.New("orm: 不支持 nil"),
		},
		{
			name: "nil with type",
			input: (*TestModel)(nil),
			want: &Model{
				tableName: "test_model",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
		},


		{
			name: "column tag",
			input: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type ColumnTag struct {
					ID uint64 `orm:"column=id"`
				}
				return &ColumnTag{}
			}(),
			want: &Model{
				tableName: "column_tag",
				fieldMap: map[string]*field{
					// 默认是 i_d
					"ID": {
						colName: "id",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &registry{}
			m, err := r.Register(tt.input)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, m)
		})
	}
}

func TestSwitch(t *testing.T) {
	Switch(nil)
	Switch((*TestModel)(nil))
}

func Switch(val any) {
	switch v := val.(type) {
	case nil:
		fmt.Println("hello, nil")
	case *TestModel:
		fmt.Printf("hello, test Model %v \n", v)
		if v == nil {
			fmt.Printf("hello, test Model nil %v \n", v)
		}
	}
}
