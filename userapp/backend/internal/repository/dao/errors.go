package dao

import (
	"errors"
	"github.com/rexshen5913/geek-learn-go/geektime-go /orm"
)

var (
	ErrDuplicateEmail = errors.New("dao: 邮件已经被注册过")
	ErrNoRows         = orm.ErrNoRows
)
