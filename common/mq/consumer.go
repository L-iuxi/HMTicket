package mq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// Handler 消费消息的回调函数类型
type Handler func(ctx context.Context, body []byte) error

// Consume 启动消费者，处理队列消息
// handler 返回 error 时消息 NACK 并重回队列（利用 DLX 做重试）
func (mq *RabbitMQ) Consume(queueName string, handler Handler) error {
	msgs, err := mq.Ch.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack = false，手动 ACK
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			// 从 header 读重试次数，3 次后不重回队列（进死信）
			retryCount := getRetryCount(d.Headers)
			if err := handler(context.Background(), d.Body); err != nil {
				logx.Errorf("[MQ] consume failed, retry=%d, err=%v", retryCount, err)
				if retryCount >= 3 {
					// 超过 3 次重试，NACK 不进队 → 自动进死信
					d.Nack(false, false)
				} else {
					// 更新重试次数 header，重回队列
					setRetryCount(&d, retryCount+1)
					d.Nack(false, true)
				}
				continue
			}
			d.Ack(false)
		}
		log.Println("[MQ] consumer channel closed")
	}()

	return nil
}

func getRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}
	if v, ok := headers["x-retry-count"]; ok {
		switch val := v.(type) {
		case int:
			return val
		case int32:
			return int(val)
		case int64:
			return int(val)
		}
	}
	return 0
}

func setRetryCount(d *amqp.Delivery, count int) {
	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}
	d.Headers["x-retry-count"] = count
}
