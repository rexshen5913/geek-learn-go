//go:build v13

package orm

import "github.com/rexshen5913/geek-learn-go/geektime-go /orm/v13/internal/errs"

// 将内部的 sentinel error 暴露出去
var (
	// ErrNoRows 代表没有找到数据
	ErrNoRows = errs.ErrNoRows
)
