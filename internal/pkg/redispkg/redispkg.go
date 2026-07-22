package redispkg

import (
	"github.com/redis/go-redis/v9"
)

// 依赖注入
type StockRepository struct {
	client *redis.Client
}

func NewStockRepository(Client *redis.Client) *StockRepository {
	return &StockRepository{
		client: Client,
	}
}

func (r *StockRepository) Client() *redis.Client {

	return r.client
}
