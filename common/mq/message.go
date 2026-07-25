package mq

import "encoding/json"

// CreateOrderMessage 扣库存成功后发送到 order.create 的消息体
type CreateOrderMessage struct {
	OrderNo      string  `json:"orderno"`
	UserID       uint32  `json:"userId"`
	EventID      uint64  `json:"eventId"`
	ShowID       uint32  `json:"showId"`
	TicketTypeID uint32  `json:"ticketTypeId"`
	Quantity     int32   `json:"quantity"`
	TotalPrice   float64 `json:"totalPrice"`
	RequestID    string  `json:"requestId"`
	IdemKey      string  `json:"idemKey"` // 幂等 key，建订单成功后标记成功
}

type TimeOutCancelMessage struct {
	OrderNo      string `json:"orderno"`
	Quantity     int32  `json:"quantity"`
	TicketTypeID uint32 `json:"ticketTypeId"`
}

func MarshalMessage(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
