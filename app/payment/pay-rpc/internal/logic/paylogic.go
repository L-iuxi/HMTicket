package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"Ticket/app/payment/pay-rpc/internal/svc"
	"Ticket/app/payment/pay-rpc/payment"
	"Ticket/common/xerr"
	"Ticket/internal/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayLogic {
	return &PayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 模拟支付
func (l *PayLogic) Pay(in *payment.PayReq) (*payment.PayResp, error) {
	//唯一建名
	lockKey := fmt.Sprintf("lock:payment:%s", in.OrderNo)
	//拿锁
	lockValue, err := l.svcCtx.Lock.Lock(l.ctx, lockKey, 30*time.Second)
	if err != nil {

		if errors.Is(err, redis.ErrLockFailed) {
			return nil,
				xerr.NewErrCode(xerr.PAY_DUPLICATE)
		}

		return nil, err
	}

	defer func() {
		err := l.svcCtx.Lock.Unlock(l.ctx, lockKey, lockValue)
		if err != nil {
			logx.Errorf(
				"释放支付锁失败:%v", err,
			)
		}

	}()
	// 校验支付金额
	if in.Amount != in.TotalPrice {
		return &payment.PayResp{
			Success: false,
			Message: "支付金额错误",
		}, nil
	}

	// 模拟支付成功
	return &payment.PayResp{
		Success: true,
		Message: "支付成功",
	}, nil
}
