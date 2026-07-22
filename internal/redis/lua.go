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
