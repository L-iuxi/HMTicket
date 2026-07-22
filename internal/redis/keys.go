package redis

import "fmt"

const (
	StockPrefix = "stock:ticket:"

	OrderLockPrefix = "lock:order:"

	OrderTimeoutPrefix = "order:timeout:"

	UserLimitPrefix = "user:limit:"
)

func StockKey(ticketTypeID uint64) string {
	return fmt.Sprintf("%s%d", StockPrefix, ticketTypeID)
}

func OrderLockKey(orderNo string) string {
	return OrderLockPrefix + orderNo
}

func UserLimitKey(userID uint64) string {
	return fmt.Sprintf("%s%d", UserLimitPrefix, userID)
}

func OrderTimeoutKey(orderNo string) string {
	return OrderTimeoutPrefix + orderNo
}
