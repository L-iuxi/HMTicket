package xerr

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorInterceptor 将 *CodeError 转换为 gRPC status.Error
// 非 CodeError 直接透传
func ErrorInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return resp, GrpcStatus(err)
	}
	return resp, nil
}

// GrpcStatus 将 error 转为 gRPC status.Error
func GrpcStatus(err error) error {
	var ce *CodeError
	if !AsCodeError(err, &ce) {
		return err
	}
	return status.Error(toGrpcCode(ce.errCode), ce.Error())
}

// toGrpcCode 将业务错误码映射为 gRPC status code
func toGrpcCode(errCode uint32) codes.Code {
	switch {
	// 全局
	case errCode == OK:
		return codes.OK
	case errCode == REQUEST_PARAM_ERROR:
		return codes.InvalidArgument
	case errCode == TOKEN_EXPIRE_ERROR || errCode == TOKEN_GENERATE_ERROR:
		return codes.Unauthenticated
	case errCode == SERVER_COMMON_ERROR || errCode == DB_ERROR ||
		errCode == DB_UPDATE_AFFECTED_ZERO_ERROR || errCode == DB_ERROR_NOT_FOUND ||
		errCode == DB_TRANSACTION_ERROR || errCode == REDIS_ERROR:
		return codes.Internal

	// 用户
	case errCode == USER_NOT_FOUND:
		return codes.NotFound
	case errCode == USER_ALREADY_EXISTS || errCode == EMAIL_ALREADY_USED || errCode == PHONE_ALREADY_USED:
		return codes.AlreadyExists
	case errCode == WRONG_PASSWORD || errCode == OLD_PASSWORD_ERROR:
		return codes.Unauthenticated
	case errCode == EMPTY_USER_INFO:
		return codes.InvalidArgument
	case errCode == USER_NOT_LOGIN:
		return codes.Unauthenticated

	// 活动/场次/票种
	case errCode == EVENT_NOT_FOUND || errCode == SHOW_NOT_FOUND ||
		errCode == TICKET_TYPE_NOT_FOUND || errCode == TICKET_TYPE_EVENT_NOT_FOUND:
		return codes.NotFound
	case errCode == EVENT_ALREADY_EXISTS:
		return codes.AlreadyExists
	case errCode == EVENT_NAME_EMPTY || errCode == EVENT_LOCATION_EMPTY ||
		errCode == EVENT_START_TIME_EMPTY || errCode == EVENT_END_TIME_EMPTY ||
		errCode == EVENT_TIME_INVALID || errCode == SHOW_NAME_EMPTY ||
		errCode == SHOW_TIME_INVALID:
		return codes.InvalidArgument
	case errCode == EVENT_STATUS_INVALID || errCode == TICKET_TYPE_STATUS_INVALID ||
		errCode == NO_DATA_TO_UPDATE || errCode == SHOW_EVENT_NOT_MATCH:
		return codes.FailedPrecondition

	// 票务
	case errCode == TICKET_NOT_FOUND:
		return codes.NotFound
	case errCode == TICKET_VERIFY_FAIL || errCode == TICKET_REFUND_FAIL ||
		errCode == TICKET_TRANSFER_FAIL || errCode == TICKET_CREATE_FAIL:
		return codes.FailedPrecondition

	// 库存
	case errCode == STOCK_NOT_ENOUGH:
		return codes.ResourceExhausted
	case errCode == STOCK_INIT_FAIL || errCode == STOCK_UPDATE_FAIL ||
		errCode == STOCK_DEDUCT_FAIL || errCode == STOCK_RELEASE_FAIL:
		return codes.Internal

	// 订单
	case errCode == ORDER_NOT_FOUND:
		return codes.NotFound
	case errCode == ORDER_STATUS_INVALID:
		return codes.FailedPrecondition
	case errCode == ORDER_DUPLICATE:
		return codes.AlreadyExists
	case errCode == ORDER_TIMEOUT:
		return codes.DeadlineExceeded
	case errCode == ORDER_CREATE_FAIL || errCode == ORDER_CANCEL_FAIL ||
		errCode == ORDER_UPDATE_FAIL:
		return codes.Internal

	// 支付
	case errCode == PAY_ORDER_NOT_FOUND:
		return codes.NotFound
	case errCode == PAY_STATUS_INVALID:
		return codes.FailedPrecondition
	case errCode == PAY_DUPLICATE:
		return codes.AlreadyExists
	case errCode == PAY_FAIL:
		return codes.Internal
	}

	return codes.Internal
}

// AsCodeError 类型断言，类似 errors.As
func AsCodeError(err error, target **CodeError) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*CodeError); ok {
		*target = e
		return true
	}
	return false
}
