package service

import (
	"blog/internal/domain/admin/model"
	"blog/internal/domain/admin/repository"
	"blog/internal/response"
	"blog/pkg/config"
	"blog/pkg/jwt"
	"blog/utils"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	IfInit() (bool, error)
	Init(req *model.InitReq) (res model.InitRes, err error)
	Auth(email, password string) (res model.LoginRes, err error)
	RefreshToken(payload *model.JwtPayload, refreshToken string) (res model.RefreshTokenRes, err error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) IfInit() (bool, error) {
	return s.repo.HaveOne()
}

func (s *service) Init(req *model.InitReq) (res model.InitRes, err error) {
	have, err := s.IfInit()
	if err != nil {
		res.Code = response.CodeDatabaseError
		return res, errors.WithStack(err)
	}

	if have {
		res.Code = response.CodeAdminExist
		return res, errors.New("管理员已初始化")
	}

	hashPassword, err := utils.EncryptPassword(req.Password)

	if err != nil {
		res.Code = response.CodeServerError
		return res, errors.WithStack(err)
	}

	newAdmin := &model.Admin{
		Email:        req.Email,
		HashPassword: hashPassword,
	}

	return res, errors.WithStack(s.repo.Create(newAdmin))
}

func (s *service) genToken(payload *model.JwtPayload) (string, error) {
	jwtCfg := config.Cfg.Auth.JWT

	JWTTokenParams := jwt.JWTTokenParams{
		Payload:  *payload,
		Duration: time.Minute * time.Duration(jwtCfg.ExpireMinute),
		Secret:   []byte(jwtCfg.Secret),
	}

	token, err := jwt.GenToken[model.JwtPayload](&JWTTokenParams)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return token, err
}

func (s *service) Auth(email, password string) (res model.LoginRes, err error) {
	admin, err := s.repo.FindByEmail(email)
	if err != nil {
		return res, errors.WithStack(err)
	}

	if admin == nil {
		return res, errors.New("身份验证失败")
	}

	err = bcrypt.CompareHashAndPassword(admin.HashPassword, []byte(password))
	if err != nil {
		return res, errors.WithStack(err)
	}

	payload := &model.JwtPayload{
		ID: admin.ID,
	}

	token, err := s.genToken(payload)
	if err != nil {
		return res, errors.WithStack(err)
	}

	refreshToken, err := s.repo.GenRefreshToken(payload)
	if err != nil {
		return res, errors.WithStack(err)
	}

	res = model.LoginRes{
		Payload:      *payload,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (s *service) RefreshToken(payload *model.JwtPayload, refreshToken string) (res model.RefreshTokenRes, err error) {
	if err := s.repo.ValidateRefreshToken(payload, refreshToken); err != nil {
		if errors.Is(err, redis.Nil) {
			res.Code = response.CodeRefreshInvalid
			return res, errors.WithStack(err)
		}
		return res, errors.WithStack(err)
	}

	newToken, err := s.genToken(payload)
	if err != nil {
		return res, errors.WithStack(err)
	}

	res.Token = newToken
	return
}
