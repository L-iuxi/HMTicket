package logic

import (
	"context"
	"errors"
	"strings"

	"Ticket/app/user/internal/svc"
	"Ticket/app/user/internal/types"
	"Ticket/internal/pkg/db"
	"Ticket/internal/pkg/util"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 注册账号
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	//用户名，密码，电话号邮箱地址不能为空
	if req.Username == "" || req.Password == "" || req.Phone == "" || req.Email == "" {
		return nil, errors.New("用户名，密码，电话号邮箱地址不能为空")
	}
	var count int64
	//检查用户名是否存在
	err = l.svcCtx.DB.Model(&db.User{}).Where("username = ?", req.Username).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已经存在")
	}

	//检查电话号是否存在
	count = 0
	err = l.svcCtx.DB.Model(&db.User{}).Where("phone = ?", req.Phone).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该电话号已被注册")
	}

	//检查邮箱地址是否存在
	count = 0
	err = l.svcCtx.DB.Model(&db.User{}).Where("email = ?", req.Email).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该邮箱已被注册")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &db.User{
		UserID:   util.GenerateUserID(), //随即生成的六位id
		Username: req.Username,
		Password: string(hash),
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Role:     "user",
	}
	err = l.svcCtx.DB.Create(user).Error

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return nil, errors.New("用户名/邮箱/电话已被注册")
		}
		return nil, err
	}

	return &types.RegisterResp{
		Message: "注册成功^^",
	}, nil
}
