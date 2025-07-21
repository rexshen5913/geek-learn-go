package repository

import (
	"errors"
	"github.com/rexshen5913/geek-learn-go/geektime-go /userapp/backend/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = errors.New("未找到指定的用户")
)
