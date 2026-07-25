package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "SUCCESS"

	// 全局
	message[SERVER_COMMON_ERROR] = "服务器开小差啦，稍后再来试一下"
	message[REQUEST_PARAM_ERROR] = "参数错误"
	message[TOKEN_EXPIRE_ERROR] = "token失效，请重新登录"
	message[TOKEN_GENERATE_ERROR] = "生成token失败"
	message[DB_ERROR] = "数据库繁忙，请稍后再试"
	message[DB_UPDATE_AFFECTED_ZERO_ERROR] = "更新数据影响行数为0"
	message[DB_ERROR_NOT_FOUND] = "数据未找到"
	message[DB_TRANSACTION_ERROR] = "数据库事务错误"
	message[REDIS_ERROR] = "Redis错误"

	// 用户模块
	message[USER_NOT_FOUND] = "用户不存在"
	message[WRONG_PASSWORD] = "密码错误"
	message[USER_ALREADY_EXISTS] = "用户名/邮箱/电话已被注册"
	message[EMPTY_USER_INFO] = "用户名，密码，电话号邮箱地址不能为空"
	message[EMAIL_ALREADY_USED] = "邮箱已被使用"
	message[PHONE_ALREADY_USED] = "手机号已被使用"
	message[OLD_PASSWORD_ERROR] = "旧密码错误"
	message[USER_NOT_LOGIN] = "未登录"

	// 活动/场次/票种模块
	message[EVENT_NOT_FOUND] = "活动不存在"
	message[EVENT_NAME_EMPTY] = "活动名称不能为空"
	message[EVENT_LOCATION_EMPTY] = "活动地点不能为空"
	message[EVENT_START_TIME_EMPTY] = "开始时间不能为空"
	message[EVENT_END_TIME_EMPTY] = "结束时间不能为空"
	message[EVENT_TIME_INVALID] = "活动结束时间必须晚于开始时间"
	message[EVENT_ALREADY_EXISTS] = "同名称同时间的活动已存在"
	message[EVENT_STATUS_INVALID] = "当前活动状态不允许此操作"
	message[SHOW_NOT_FOUND] = "场次不存在"
	message[TICKET_TYPE_NOT_FOUND] = "票种不存在"
	message[TICKET_TYPE_EVENT_NOT_FOUND] = "所属活动不存在"
	message[TICKET_TYPE_STATUS_INVALID] = "当前活动状态不允许修改票种"
	message[NO_DATA_TO_UPDATE] = "没有需要修改的数据"
	message[SHOW_NAME_EMPTY] = "场次名称不能为空"
	message[SHOW_TIME_INVALID] = "场次时间必须位于活动时间范围内"
	message[SHOW_EVENT_NOT_MATCH] = "场次不属于当前活动"

	// 票务模块
	message[TICKET_NOT_FOUND] = "票不存在"
	message[TICKET_VERIFY_FAIL] = "验票失败"
	message[TICKET_REFUND_FAIL] = "退票失败"
	message[TICKET_TRANSFER_FAIL] = "转赠失败"
	message[TICKET_CREATE_FAIL] = "出票失败"

	// 库存模块
	message[STOCK_NOT_ENOUGH] = "库存不足"
	message[STOCK_INIT_FAIL] = "库存初始化失败"
	message[STOCK_UPDATE_FAIL] = "库存更新失败"
	message[STOCK_DEDUCT_FAIL] = "库存扣减失败"
	message[STOCK_RELEASE_FAIL] = "库存释放失败"

	// 订单模块
	message[ORDER_NOT_FOUND] = "订单不存在"
	message[ORDER_CREATE_FAIL] = "订单创建失败"
	message[ORDER_CANCEL_FAIL] = "订单取消失败"
	message[ORDER_UPDATE_FAIL] = "订单更新失败"
	message[ORDER_STATUS_INVALID] = "当前订单不可修改"
	message[ORDER_DUPLICATE] = "订单正在创建，请勿重复提交"
	message[ORDER_TIMEOUT] = "订单超时"

	// 支付模块
	message[PAY_FAIL] = "支付失败"
	message[PAY_ORDER_NOT_FOUND] = "支付订单不存在"
	message[PAY_STATUS_INVALID] = "支付状态异常"
	message[PAY_DUPLICATE] = "支付处理中，请勿重复操作"
}

func MapErrMsg(errcode uint32) string {
	if msg, ok := message[errcode]; ok {
		return msg
	}
	return "服务器开小差啦，稍后再来试一试"
}

func IsCodeErr(errcode uint32) bool {
	_, ok := message[errcode]
	return ok
}
