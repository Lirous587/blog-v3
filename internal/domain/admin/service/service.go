package service

import (
	"blog/internal/domain/admin/model"
	"blog/internal/domain/admin/repository/cache"
	"blog/internal/domain/admin/repository/db"
	"blog/pkg/config"
	"blog/pkg/jwt"
	"blog/pkg/response"
	"blog/utils"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	IfInit() (bool, error)
	Init(req *model.InitReq) *response.AppError
	Auth(email, password string) (res *model.LoginRes, appErr *response.AppError)
	RefreshToken(payload *model.JwtPayload, refreshToken string) (res *model.RefreshTokenRes, appErr *response.AppError)
}

type service struct {
	db    db.DB
	cache cache.Cache
}

func NewService(db db.DB, cache cache.Cache) Service {
	return &service{db: db, cache: cache}
}

func (s *service) IfInit() (bool, error) {
	return s.db.HaveOne()
}

func (s *service) Init(req *model.InitReq) (appErr *response.AppError) {
	have, err := s.IfInit()
	if err != nil {
		return response.NewAppError(response.CodeDatabaseError, errors.WithStack(err))
	}

	if have {
		return response.NewAppError(response.CodeAdminExist, errors.New("管理员已初始化"))
	}

	hashPassword, err := utils.EncryptPassword(req.Password)

	if err != nil {
		return response.NewAppError(response.CodeServerError, errors.WithStack(err))
	}

	newAdmin := &model.Admin{
		Email:        req.Email,
		HashPassword: hashPassword,
	}

	if err = s.db.Create(newAdmin); err != nil {
		return response.NewAppError(response.CodeDatabaseError, errors.WithStack(err))
	}

	return
}

func (s *service) genToken(payload *model.JwtPayload) (token string, err error) {
	jwtCfg := config.Cfg.Auth.JWT

	JWTTokenParams := jwt.JWTTokenParams{
		Payload:  *payload,
		Duration: time.Minute * time.Duration(jwtCfg.ExpireMinute),
		Secret:   []byte(jwtCfg.Secret),
	}

	token, err = jwt.GenToken[model.JwtPayload](&JWTTokenParams)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return
}

func (s *service) Auth(email, password string) (res *model.LoginRes, appErr *response.AppError) {
	admin, err := s.db.FindByEmail(email)
	if err != nil {
		return nil, response.NewAppError(response.CodeDatabaseError, errors.WithStack(err))
	}

	if admin == nil {
		return nil, response.NewAppError(response.CodeAuthFailed, errors.WithStack(err))
	}

	err = bcrypt.CompareHashAndPassword(admin.HashPassword, []byte(password))
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, errors.WithStack(err))
	}

	payload := &model.JwtPayload{
		ID: admin.ID,
	}

	token, err := s.genToken(payload)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, errors.WithStack(err))
	}

	refreshToken, err := s.cache.GenRefreshToken(payload)
	if err != nil {
		return nil, response.NewAppError(response.CodeDatabaseError, errors.WithStack(err))
	}

	res = &model.LoginRes{
		Payload:      *payload,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return
}

func (s *service) RefreshToken(payload *model.JwtPayload, refreshToken string) (res *model.RefreshTokenRes, appErr *response.AppError) {
	if err := s.cache.ValidateRefreshToken(payload, refreshToken); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, response.NewAppError(response.CodeRefreshInvalid, errors.WithStack(err))
		}
		return nil, response.NewAppError(response.CodeDatabaseError, err)
	}

	newToken, err := s.genToken(payload)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	res = &model.RefreshTokenRes{
		Token: newToken,
	}
	return
}
