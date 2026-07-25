package order

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"Ticket/app/event/event-rpc/eventclient"
	"Ticket/app/inventory/inventory-rpc/inventoryclient"
	"Ticket/app/order/order-api/internal/common"
	"Ticket/app/order/order-api/internal/svc"
	"Ticket/app/order/order-api/internal/types"
	"Ticket/app/order/order-rpc/orderclient"
	"Ticket/common/mq"
	"Ticket/common/xerr"
	"Ticket/internal/redis"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type BuyTicketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBuyTicketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BuyTicketLogic {
	return &BuyTicketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BuyTicketLogic) BuyTicket(req *types.BuyTicketReq) (*types.BuyTicketResp, error) {
	userID := common.GetUserID(l.ctx)
	if userID == 0 {
		return nil, xerr.NewErrCode(xerr.USER_NOT_LOGIN)
	}

	if req.RequestId == "" {
		return &types.BuyTicketResp{Status: "fail", Message: "缺少请求ID"}, nil
	}

	/*1. 限流
	*用户级限流，对于单个用户，只允许一秒内最多五个购买请求到达
	*令牌桶限流：初始50个token，每秒钟补10个，没有拿到令牌的请求不允许执行
	 */

	limitkey := fmt.Sprintf("limit:buy:user:%d", userID)
	ok, err := l.svcCtx.RateLimiter.Allow(l.ctx, limitkey, 5, time.Second)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, xerr.NewErrMsg("请求过于频繁请稍后重试")
	}

	ticketTypeId := req.TicketTypeID
	key := fmt.Sprintf("bucket:ticket:%d", ticketTypeId)
	ok, err = l.svcCtx.TokenBucket.Allow(l.ctx, key, 50, 10)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, xerr.NewErrMsg("请求过于频繁请稍后重试")
	}

	/*2. 分布式锁
	*用户id+活动id+场次id+票种id，保证第一个请求处理的时候另一个请求到达不会处理
	 */
	lockKey := fmt.Sprintf("lock:order:%d:%d:%d:%d", userID, req.EventID, req.ShowID, req.TicketTypeID)
	lockValue, err := l.svcCtx.Lock.Lock(l.ctx, lockKey, 30*time.Second)
	if err != nil {
		if err == redis.ErrLockFailed {
			return &types.BuyTicketResp{Status: "fail", Message: "订单正在创建，请勿重复提交"}, nil
		}
		return nil, err
	}
	defer func() {
		if unlockErr := l.svcCtx.Lock.Unlock(l.ctx, lockKey, lockValue); unlockErr != nil {
			logx.Errorf("释放锁失败: %v", unlockErr)
		}
	}()

	/*3. 幂等检查
	*检查请求是否已经存在，建相比较以上加了一个requestid，保证两次相同请求不会产生充分结果
	 */
	idemKey := redis.IdempotentKey(userID, req.EventID, req.ShowID, req.TicketTypeID, req.RequestId)
	val, err := l.svcCtx.Idempotent.Get(l.ctx, idemKey)
	if err != nil {
		return nil, fmt.Errorf("幂等检查失败: %w", err)
	}
	if strings.HasPrefix(val, "ok:") {
		return &types.BuyTicketResp{
			OrderNo: strings.TrimPrefix(val, "ok:"),
			Status:  "unpaid",
			Message: "订单已存在",
		}, nil
	}
	if val == "pending" {
		return &types.BuyTicketResp{Status: "fail", Message: "请求处理中，请稍后"}, nil
	}
	if val == "failed" {
		return &types.BuyTicketResp{Status: "fail", Message: "上次请求失败，请更换RequestId重试"}, nil
	}

	if err := l.svcCtx.Idempotent.Start(l.ctx, idemKey, 5*time.Minute); err != nil {
		if err == redis.ErrDuplicateRequest {
			return &types.BuyTicketResp{Status: "fail", Message: "请勿重复提交"}, nil
		}
		return nil, fmt.Errorf("幂等标记失败: %w", err)
	}

	/*4. 查票种 + 校验活动 + 限购

	 */
	ticketTypeResp, err := l.svcCtx.EventRpc.GetTicketType(l.ctx, &eventclient.GetTicketTypeReq{
		TicketTypeId: req.TicketTypeID,
	})
	if err != nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return nil, fmt.Errorf("查询票种失败: %w", err)
	}
	if ticketTypeResp == nil || ticketTypeResp.TicketType == nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return &types.BuyTicketResp{Status: "fail", Message: "票种不存在"}, nil
	}

	totalPrice := ticketTypeResp.TicketType.Price * float64(req.Quantity)

	eventResp, err := l.svcCtx.EventRpc.GetEvent(l.ctx, &eventclient.GetEventReq{EventId: req.EventID})
	if err != nil || eventResp.Event == nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return &types.BuyTicketResp{Status: "fail", Message: "活动不存在"}, nil
	}
	if eventResp.Event.Status != "selling" {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return &types.BuyTicketResp{Status: "fail", Message: "活动未开售"}, nil
	}

	// 限购检查
	orderList, err := l.svcCtx.OrderRpc.GetOrderList(l.ctx, &orderclient.GetOrderListReq{
		UserId: uint32(userID),
	})
	if err != nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return nil, fmt.Errorf("查询订单列表失败: %w", err)
	}
	bought := int32(0)
	for _, o := range orderList.Orders {
		if o.TicketTypeId == uint32(req.TicketTypeID) &&
			o.Status != "failed" && o.Status != "cancelled" {
			bought += o.Quantity
		}
	}
	if bought+req.Quantity > ticketTypeResp.TicketType.MaxPerUser {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return &types.BuyTicketResp{Status: "fail", Message: "超出限购数量"}, nil
	}

	// 5. 扣库存（Redis Lua 原子）

	deductResp, err := l.svcCtx.InventoryRpc.DeductStock(l.ctx, &inventoryclient.DeductStockReq{
		TicketTypeId: req.TicketTypeID,
		Quantity:     req.Quantity,
	})
	if err != nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return nil, fmt.Errorf("库存扣减失败: %w", err)
	}
	if !deductResp.Success {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		return &types.BuyTicketResp{Status: "fail", Message: deductResp.Message}, nil
	}

	// 6. 发 MQ 异步建订单

	//提前创建订单号供用户调用支付
	orderNo := uuid.NewString()

	msg := mq.CreateOrderMessage{
		OrderNo:      orderNo,
		UserID:       uint32(userID),
		EventID:      req.EventID,
		ShowID:       uint32(req.ShowID),
		TicketTypeID: uint32(req.TicketTypeID),
		Quantity:     req.Quantity,
		TotalPrice:   totalPrice,
		RequestID:    req.RequestId,
		IdemKey:      idemKey,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		l.releaseStock(req.TicketTypeID, req.Quantity)
		return nil, err
	}

	if err := l.svcCtx.MQ.SendMsg(l.ctx, body, mq.RoutingOrderCreate); err != nil {
		l.svcCtx.Idempotent.Failed(l.ctx, idemKey, time.Minute)
		l.releaseStock(req.TicketTypeID, req.Quantity)
		return nil, fmt.Errorf("发送MQ消息失败: %w", err)
	}

	return &types.BuyTicketResp{
		OrderNo:    orderNo,
		Status:     "pending",
		TotalPrice: totalPrice,
		Message:    "下单成功，订单正在创建，请稍后查询",
	}, nil
}

// releaseStock 回滚库存
func (l *BuyTicketLogic) releaseStock(ticketTypeID uint64, quantity int32) {
	_, err := l.svcCtx.InventoryRpc.ReleaseStock(l.ctx, &inventoryclient.ReleaseStockReq{
		TicketTypeId: ticketTypeID,
		Quantity:     quantity,
	})
	if err != nil {
		l.Errorf("回滚库存失败: ticketTypeId=%d, quantity=%d, err=%v", ticketTypeID, quantity, err)
	}
}
