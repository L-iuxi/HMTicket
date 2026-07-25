package common

import (
	"context"
	"encoding/json"
	"strconv"
)

// key 必须用 plain string（非命名类型），go-zero Authorize 中间件
// 用 jwt.MapClaims 的 key（string）注入 context。命名类型会导致
// context.Value 接口比较时动态类型不匹配，始终返回 nil。
const userRoleKey = "role"
const userIDKey = "user_id"

// 把userid放进context
func ContextWithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// 从context里面获取userid
// go-zero Authorize 中间件把 jwt.MapClaims 非标准 claims 注入 context，
// 数值类型为 json.Number（jwt.WithJSONNumber），非 uint。
func GetUserID(ctx context.Context) uint {
	v := ctx.Value(userIDKey)
	switch val := v.(type) {
	case uint:
		return val
	case float64:
		return uint(val)
	case json.Number:
		n, err := val.Int64()
		if err != nil {
			return 0
		}
		return uint(n)
	case int64:
		return uint(val)
	case string:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return 0
		}
		return uint(n)
	default:
		return 0
	}
}

// 把userrole放进context
func ContextWithUserRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, userRoleKey, role)
}

// 从上下文获取userrole
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(userRoleKey).(string)
	return role, ok
}
