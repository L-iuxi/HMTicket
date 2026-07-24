// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"Ticket/app/user/internal/middleware"
	"Ticket/app/user/internal/svc"
	"Ticket/app/user/internal/types"
	db "Ticket/app/user/model"
	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileReq) (resp *types.CommonResp, err error) {
	//获取当前用户
	userID := middleware.GetUserID(l.ctx)

	if userID == 0 {
		return nil, xerr.NewErrCode(xerr.USER_NOT_LOGIN)
	}
	// 用户是否存在
	var user db.User

	err = l.svcCtx.DB.
		First(&user, userID).
		Error

	if err != nil {
		return nil, xerr.NewErrCode(xerr.USER_NOT_FOUND)
	}
	// 传入的修改邮箱地址是否重复
	if req.Email != "" && req.Email != user.Email {

		var count int64

		err = l.svcCtx.DB.
			Model(&db.User{}).
			Where("email = ?", req.Email).
			Count(&count).
			Error

		if err != nil {
			return nil, err
		}

		if count > 0 {
			return nil, xerr.NewErrCode(xerr.EMAIL_ALREADY_USED)
		}
	}
	// 手机号是否重复
	if req.Phone != "" && req.Phone != user.Phone {

		var count int64

		err = l.svcCtx.DB.
			Model(&db.User{}).
			Where("phone = ?", req.Phone).
			Count(&count).
			Error

		if err != nil {
			return nil, err
		}

		if count > 0 {
			return nil, xerr.NewErrCode(xerr.PHONE_ALREADY_USED)
		}
	}

	if req.Phone == user.Phone {
		return nil, xerr.NewErrCode(xerr.NO_DATA_TO_UPDATE)
	}

	if req.Username == user.Username {
		return nil, xerr.NewErrCode(xerr.NO_DATA_TO_UPDATE)
	}
	if req.Email == user.Email {
		return nil, xerr.NewErrCode(xerr.NO_DATA_TO_UPDATE)
	}
	// 更新字段
	updates := make(map[string]interface{})

	// 密码修改
	if req.NewPassword != "" {
		if req.OldPassword == "" {
			return nil, xerr.NewErrMsg("修改密码需要提供旧密码")
		}
		// 验证旧密码
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
		if err != nil {
			return nil, xerr.NewErrCode(xerr.OLD_PASSWORD_ERROR)
		}
		// 新密码 hash
		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(hash)
	}

	if req.Username != "" {
		updates["username"] = req.Username
	}

	if req.Email != "" {
		updates["email"] = req.Email
	}

	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	updates["gender"] = req.Gender
	// 执行更新，写入数据库
	err = l.svcCtx.DB.
		Model(&user).
		Updates(updates).
		Error

	if err != nil {
		return nil, err
	}
	return &types.CommonResp{
		Message: "更新成功",
	}, nil
}
