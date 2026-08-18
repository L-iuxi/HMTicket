// 生成幂等 requestId（下单/支付防重）
export function genId() {
  return Date.now().toString(36) + Math.random().toString(36).slice(2, 10)
}
