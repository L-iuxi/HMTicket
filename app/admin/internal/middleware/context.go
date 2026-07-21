package middleware

import "context"

type ctxKey string

const userRoleKey ctxKey = "userRole"
const userIDKey ctxKey = "userID"

// 把userid放进context
func contextWithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// 从context里面获取userid
func GetUserID(ctx context.Context) uint {
	id, ok := ctx.Value(userIDKey).(uint)
	if !ok {
		return 0
	}
	return id
}

// 把userrole放进context
func contextWithUserRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, userRoleKey, role)
}

// 从上下文获取userrole
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(userRoleKey).(string)
	return role, ok
}
