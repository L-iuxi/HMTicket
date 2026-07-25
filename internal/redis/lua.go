package redis

// 扣库存脚本
const DeductStockLua = `
local stock = tonumber(redis.call('GET', KEYS[1]))
local quantity = tonumber(ARGV[1])

if not stock then
    return -1
end

if stock < quantity then
    return 0
end

redis.call("DECRBY", KEYS[1], quantity)

return 1
`

// Lua释放锁
const UnlockLua = `
	if redis.call("GET",KEYS[1]) == ARGV[1] then
		return redis.call("DEL",KEYS[1])
	else 
		return 0
	end
`

// 用户固定窗口限流
const RateLimitLua = `
local count = redis.call("INCR", KEYS[1])

-- 第一次请求设置过期时间
if count == 1 then
    redis.call("EXPIRE", KEYS[1], ARGV[1])
end

-- 超过限制
if count > tonumber(ARGV[2]) then
    return 0
end

return 1
`

// 令牌桶限流
const TokenBucketLua = `
local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])


local values = redis.call(
    "HMGET",
    key,
    "tokens",
    "last_time"
)


local tokens = tonumber(values[1])
local lastTime = tonumber(values[2])


-- 第一次访问
if tokens == nil then
    tokens = capacity
    lastTime = now
end


-- 计算补充数量

local delta = now - lastTime

local add = delta * rate


tokens = math.min(
    capacity,
    tokens + add
)


-- 没令牌

if tokens < 1 then

    redis.call(
        "HSET",
        key,
        "tokens",
        tokens,
        "last_time",
        now
    )

    redis.call(
        "EXPIRE",
        key,
        60
    )

    return 0

end



-- 消耗令牌

tokens = tokens - 1



redis.call(
    "HSET",
    key,
    "tokens",
    tokens,
    "last_time",
    now
)


redis.call(
    "EXPIRE",
    key,
    60
)


return 1
`
