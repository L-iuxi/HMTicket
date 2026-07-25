package redis

import "fmt"

// Redis key 命名集中管理，避免硬编码散落各处。

const (
	StockPrefix     = "stock:ticket:"
	OrderLockPrefix = "lock:order:"
	IdempPrefix     = "idem:"
	UserLimitPrefix = "user:limit:"
)

func StockKey(ticketTypeID uint64) string {
	return fmt.Sprintf("%s%d", StockPrefix, ticketTypeID)
}

func OrderLockKey(key string) string {
	return OrderLockPrefix + key
}

func IdempotentKey(parts ...interface{}) string {
	s := IdempPrefix
	for i, p := range parts {
		if i > 0 {
			s += ":"
		}
		s += fmt.Sprintf("%v", p)
	}
	return s
}

func UserLimitKey(userID uint64) string {
	return fmt.Sprintf("%s%d", UserLimitPrefix, userID)
}
