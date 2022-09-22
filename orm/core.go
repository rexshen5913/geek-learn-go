

package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/valuer"
	"gitee.com/geektime-geekbang/geektime-go/orm/model"
)

type core struct {
	r          model.Registry
	dialect    Dialect
	valCreator valuer.Creator
	ms         []Middleware
}
