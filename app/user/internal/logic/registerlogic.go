package logic

import (
	"context"
	"strings"

	"Ticket/app/user/internal/svc"
	"Ticket/app/user/internal/types"
	db "Ticket/app/user/model"
	"Ticket/common/xerr"
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
		return nil, xerr.NewErrCode(xerr.EMPTY_USER_INFO)
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
			return nil, xerr.NewErrCode(xerr.USER_ALREADY_EXISTS)
		}
		return nil, err
	}

	return &types.RegisterResp{
		Message: "注册成功^^",
	}, nil
}
