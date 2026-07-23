package logic

import "fmt"

// 打印票种
func stockKey(ticketTypeId uint64) string {
	return fmt.Sprintf("stock:ticket:%d", ticketTypeId)
}
