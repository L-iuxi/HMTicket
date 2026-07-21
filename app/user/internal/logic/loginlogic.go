// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"

	"Ticket/app/user/internal/svc"
	"Ticket/app/user/internal/types"
	"Ticket/internal/pkg/db"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 登陆账号
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {

	user, err := l.FindUserByAccount(req.Account)
	if err != nil {
		return nil, errors.New("账号不存在")
	}
	// 比较hash之后的密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	//生成token
	token, err := l.svcCtx.Jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Token: token,
	}, nil
}

// 自动判断输入的是用户名还是邮箱地址还是电话号
func (l *LoginLogic) FindUserByAccount(account string) (*db.User, error) {
	var user db.User

	err := l.svcCtx.DB.
		Where(
			"username = ? OR phone = ? OR email = ?",
			account,
			account,
			account,
		).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
