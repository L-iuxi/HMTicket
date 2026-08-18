// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"

	"Ticket/app/user/internal/middleware"
	"Ticket/app/user/internal/svc"
	"Ticket/app/user/internal/types"
	db "Ticket/app/user/model"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 查询个人信息
func (l *ProfileLogic) Profile() (*types.UserProfileResp, error) {

	userID := middleware.GetUserID(l.ctx)

	logx.Infof("Profile userID=%d", userID)

	var user db.User
	err := l.svcCtx.DB.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.USER_NOT_FOUND)
		}
		return nil, err
	}

	logx.Infof("Profile err=%v", err)

	var gender string
	if user.Gender == 0 {
		gender = "男"
	} else {
		gender = "女"
	}
	return &types.UserProfileResp{
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Gender:   gender,
		UserId:   user.UserID,
		Role:     user.Role,
	}, nil
}
