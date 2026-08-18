<template>
  <div>
    <NavBar />
    <div class="page">
      <h1 class="page-title">我的订单</h1>
      <div v-if="loading" class="hint">加载中…</div>
      <div v-else-if="orders.length === 0" class="hint">还没有订单，去首页逛逛吧</div>
      <div v-else class="order-list">
        <div v-for="o in orders" :key="o.orderNo" class="glass-card order-card">
          <div class="order-head">
            <span class="order-no">订单号 {{ o.orderNo }}</span>
            <span class="tag" :class="statusCls(o.status)">{{ statusText(o.status) }}</span>
          </div>
          <div class="order-body">
            <div class="meta">活动 ID {{ o.eventId }} · 场次 ID {{ o.showId }} · 票种 ID {{ o.ticketTypeId }}</div>
            <div class="order-price">¥{{ o.totalPrice }} <span class="meta">× {{ o.quantity }}</span></div>
          </div>
          <div class="order-foot">
            <button v-if="o.status === 'unpaid'" class="glass-btn" :disabled="busy.has(o.orderNo)" @click="onCancel(o)">
              取消
            </button>
            <button v-if="o.status === 'unpaid'" class="glass-btn primary" :disabled="busy.has(o.orderNo)" @click="onPay(o)">
              {{ busy.has(o.orderNo) ? '处理中…' : '去支付' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import NavBar from '../components/NavBar.vue'
import { getOrderList, payOrder, cancelOrder } from '../api'
import { genId } from '../utils/id'

const orders = ref([])
const loading = ref(true)
const busy = reactive(new Set())

const STATUS = {
  unpaid: { text: '待支付', cls: 'yellow' },
  pending: { text: '处理中', cls: 'yellow' },
  creating: { text: '处理中', cls: 'yellow' },
  paid: { text: '已支付', cls: 'green' },
  cancelled: { text: '已取消', cls: 'gray' },
  fail: { text: '失败', cls: 'red' },
  failed: { text: '失败', cls: 'red' }
}
function statusText(s) { return (STATUS[s] || { text: s }).text }
function statusCls(s) { return (STATUS[s] || { cls: 'gray' }).cls }

async function load() {
  try {
    orders.value = (await getOrderList())?.orders || []
  } catch (e) {
  } finally {
    loading.value = false
  }
}

async function onPay(o) {
  busy.add(o.orderNo)
  try {
    const res = await payOrder({ orderNo: o.orderNo, requestId: genId() })
    ElMessage.success(res.message || '支付成功，正在出票')
    await load()
  } catch (e) {
  } finally {
    busy.delete(o.orderNo)
  }
}

async function onCancel(o) {
  try {
    await ElMessageBox.confirm(`确认取消订单 ${o.orderNo}？`, '取消订单', { type: 'warning' })
  } catch (e) {
    return
  }
  busy.add(o.orderNo)
  try {
    const res = await cancelOrder(o.orderNo)
    ElMessage.success(res.message || '已取消')
    await load()
  } catch (e) {
  } finally {
    busy.delete(o.orderNo)
  }
}

onMounted(load)
</script>

<style scoped>
.hint { padding: 40px 0; text-align: center; color: var(--ink-soft); }
.order-list { display: flex; flex-direction: column; gap: 16px; }
.order-card { padding: 18px 22px; }
.order-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.order-no { font-weight: 700; font-size: 14px; word-break: break-all; }
.order-body { display: flex; justify-content: space-between; align-items: center; }
.order-price { font-size: 20px; font-weight: 800; color: var(--yellow-600); }
.order-foot { display: flex; justify-content: flex-end; gap: 10px; margin-top: 14px; }
.meta { color: var(--ink-soft); font-size: 13px; }
</style>
