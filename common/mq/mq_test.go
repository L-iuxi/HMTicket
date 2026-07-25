package mq

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"
)

func TestRabbitMQIntegration(t *testing.T) {
	// 1. 连接
	mq, err := NewRabbitMQ("")
	if err != nil {
		t.Fatalf("连接 RabbitMQ 失败: %v", err)
	}
	defer mq.Close()

	// 2. 初始化 exchange + queue
	if err := mq.InitExchange(); err != nil {
		t.Fatalf("InitExchange: %v", err)
	}
	if err := mq.InitDeadExchange(); err != nil {
		t.Fatalf("InitDeadExchange: %v", err)
	}

	queueName := "test.order.create.queue"
	if err := mq.InitQueue(queueName); err != nil {
		t.Fatalf("InitQueue: %v", err)
	}
	if err := mq.InitDeadQueue(); err != nil {
		t.Fatalf("InitDeadQueue: %v", err)
	}
	if err := mq.BindQueue(queueName, RoutingOrderCreate); err != nil {
		t.Fatalf("BindQueue: %v", err)
	}
	if err := mq.BindDeadQueue(); err != nil {
		t.Fatalf("BindDeadQueue: %v", err)
	}

	// 3. 测试消息
	testMsg := CreateOrderMessage{
		UserID:       1,
		EventID:      100,
		ShowID:       10,
		TicketTypeID: 5,
		Quantity:     2,
		TotalPrice:   199.0,
		RequestID:    "test-req-001",
		IdemKey:      "idem:1:100:10:5:test-req-001",
	}

	body, _ := json.Marshal(testMsg)

	// 4. 启动消费者
	var receivedMsg *CreateOrderMessage
	var mu sync.Mutex
	done := make(chan struct{})

	handler := func(ctx context.Context, b []byte) error {
		var msg CreateOrderMessage
		if err := json.Unmarshal(b, &msg); err != nil {
			return err
		}
		mu.Lock()
		receivedMsg = &msg
		mu.Unlock()
		close(done)
		return nil
	}

	if err := mq.Consume(queueName, handler); err != nil {
		t.Fatalf("Consume: %v", err)
	}

	// 5. 发消息
	time.Sleep(100 * time.Millisecond) // 等 consumer 就绪
	if err := mq.SendMsg(context.Background(), body, RoutingOrderCreate); err != nil {
		t.Fatalf("SendMsg: %v", err)
	}

	// 6. 等消费完成
	select {
	case <-done:
		mu.Lock()
		msg := receivedMsg
		mu.Unlock()

		if msg == nil {
			t.Fatal("未收到消息")
		}
		if msg.UserID != 1 {
			t.Errorf("UserID = %d, want 1", msg.UserID)
		}
		if msg.TotalPrice != 199.0 {
			t.Errorf("TotalPrice = %f, want 199.0", msg.TotalPrice)
		}
		t.Logf("消息正确接收: userId=%d, ticketTypeId=%d, quantity=%d, totalPrice=%.2f",
			msg.UserID, msg.TicketTypeID, msg.Quantity, msg.TotalPrice)

	case <-time.After(5 * time.Second):
		t.Fatal("消费超时（5s）")
	}
}
