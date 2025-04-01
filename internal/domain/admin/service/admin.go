package service

import (
	"blog/internal/domain/admin/model"
	"blog/internal/domain/admin/repository"
	"blog/pkg/config"
	"blog/pkg/jwt"
	"blog/utils"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req *model.RegisterReq) error
	Auth(email, password string) (res *model.LoginRes, err error)
}

type service struct {
	repo repository.Repository
}

// NewAdminService 创建服务实例
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(req *model.RegisterReq) error {
	admin, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return errors.WithStack(err)
	}

	if admin != nil {
		return errors.New("邮箱已存在")
	}

	hashPassword, err := utils.EncryptPassword(req.Password)

	if err != nil {
		return errors.WithStack(err)
	}

	newAdmin := &model.Admin{
		Email:        req.Email,
		HashPassword: hashPassword,
	}

	return errors.WithStack(s.repo.Create(newAdmin))
}

func (s *service) Auth(email, password string) (res *model.LoginRes, err error) {
	admin, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if admin == nil {
		return nil, errors.New("身份验证失败")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.HashPassword), []byte(password))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jwtCfg := config.Cfg.Jwt

	JWTTokenParams := jwt.JWTTokenParams{
		Payload:  model.JwtPaylod{ID: admin.ID},
		Duration: time.Minute * time.Duration(jwtCfg.ExpireMinute),
		Secret:   []byte(jwtCfg.Secret),
	}

	token, err := jwt.GenToken[model.JwtPaylod](&JWTTokenParams)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	refreshtoken, err := s.repo.GenRefreshToken(email)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res = &model.LoginRes{
		ID:           admin.ID,
		Token:        token,
		RefreshToken: refreshtoken,
	}

	return res, nil
}
