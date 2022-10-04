package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/cache"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/domainobject/entity"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
)

//go:generate mockgen -source=user.go -destination=mocks/user_mock.gen.go -package=usmocks UserService
type UserService interface {
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	Login(ctx context.Context, user entity.User) (entity.User, error)
	FindById(ctx context.Context, id uint64)(entity.User, error)
	EditProfile(ctx context.Context, user entity.User) error
}

type userService struct {
	cache        cache.Cache
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) EditProfile(ctx context.Context, user entity.User) error {
	return u.repo.UpdateUser(ctx, user)
}

func (u *userService) FindById(ctx context.Context, id uint64)(entity.User, error) {
	return u.repo.GetUserById(ctx, id)
}

func (u *userService) Login(ctx context.Context, input entity.User) (entity.User, error) {
	usr, err := u.repo.GetUserByEmail(ctx, input.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return entity.User{}, ErrInvalidUserOrPassword
	}

	if err != nil {
		return entity.User{}, err
	}

	encryptedPwd := u.encryptPwdByPbkdf2(input.Password, usr.Salt)
	if encryptedPwd != usr.Password {
		return entity.User{}, ErrInvalidUserOrPassword
	}
	return usr, nil
}

// func (u *userService) match(ct uint64, reqPwd string, pwd string) bool {
// 	encode := encryptPwdByMd5(reqPwd, generateSalt(ct))
// 	return pwd == encode
// }
//
// func (u *userService) cacheUser(ctx context.Context, key string, value *dto.User) {
// 	if u.cache != nil {
// 		if err := u.cache.Put(ctx, key, value, time.Hour); err != nil {
// 			logs.Errorf("cache user error: %v", err)
// 		}
// 	}
// }
//
// func (u *userService) invalidCache(ctx context.Context, key string) {
// 	if u.cache != nil {
// 		if err := u.cache.Delete(ctx, key); err != nil {
// 			logs.Errorf("could not remove data from cache: %v", err)
// 		}
// 	}
// }

func(u *userService) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := user.Check(); err != nil {
		return entity.User{}, fmt.Errorf("%w, 原因 %v", ErrInvalidNewUser, err)
	}
	// 开始加密
	// 每一个用户都是一个单独的 salt，会更加安全
	// 你也可以考虑组合多个 HASH 加密算法来存储
	salt := uuid.New().String()
	user.Salt = salt
	user.Password = u.encryptPwdByPbkdf2(user.Password, salt)
	return u.repo.CreateUser(ctx, user)
}

func (u *userService) encryptPwdByPbkdf2(raw string, salt string) string {
	// pbkdf2 需要比较多的 CPU 的资源。不过考虑到注册用户整体上是非常非常低频的，那么你也不会介意使用这种复杂的加密算法
	return fmt.Sprintf("%X", pbkdf2.Key([]byte(raw), []byte(salt), 4096, 32, sha1.New))
}


// func (u *userService) GetUserByEmail(ctx context.Context, req *dto.GetUserByEmailReq) (*dto.GetUserByEmailReply, error) {
// 	if len(req.Email) == 0 {
// 		return nil, errors.New("invalid email")
// 	}
//
// 	um, err := dao.GetUserByEmail(ctx, req.Email)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &dto.GetUserByEmailReply{
// 		User: um.ToPB(),
// 	}, nil
// }

func (u *userService) ServiceName() string {
	return "user"
}

