<template>
  <div>
    <NavBar />
    <div class="page">
      <h1 class="page-title">我的票</h1>
      <div v-if="loading" class="hint">加载中…</div>
      <div v-else-if="list.length === 0" class="hint">还没有票，支付订单后自动出票</div>
      <div v-else class="order-list">
        <div v-for="t in list" :key="t.ticketId" class="glass-card order-card">
          <div class="order-head">
            <span class="order-no">票号 {{ t.ticketId }}</span>
            <span class="tag" :class="statusCls(t.status)">{{ statusText(t.status) }}</span>
          </div>
          <div class="order-body">
            <div class="meta">
              订单 {{ t.orderNo }}<br />
              活动 {{ t.eventId }} · 场次 {{ t.showId }} · 票种 {{ t.ticketTypeId }} · × {{ t.quantity }}
            </div>
            <div class="order-price">¥{{ t.totalPrice }}</div>
          </div>
          <div class="order-foot">
            <button class="glass-btn sm" @click="$router.push(`/ticket/${t.ticketId}`)">查看详情</button>
            <button v-if="t.status === 'unused'" class="glass-btn sm danger" :disabled="busy.has(t.ticketId)" @click="onRefund(t)">
              退票
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
import { listTicket, refundTicket } from '../api'

const list = ref([])
const loading = ref(true)
const busy = reactive(new Set())

const STATUS = {
  unused: { text: '未使用', cls: 'green' },
  used: { text: '已使用', cls: 'gray' },
  refunded: { text: '已退票', cls: 'red' }
}
function statusText(s) { return (STATUS[s] || { text: s }).text }
function statusCls(s) { return (STATUS[s] || { cls: 'gray' }).cls }

async function load() {
  try {
    const data = await listTicket({ page: 1, pageSize: 50 })
    list.value = data.list || []
  } catch (e) {
  } finally {
    loading.value = false
  }
}

async function onRefund(t) {
  try {
    await ElMessageBox.confirm(`确认退票 票号 ${t.ticketId}？`, '退票', { type: 'warning' })
  } catch (e) {
    return
  }
  busy.add(t.ticketId)
  try {
    await refundTicket(t.ticketId)
    ElMessage.success('退票成功')
    await load()
  } catch (e) {
  } finally {
    busy.delete(t.ticketId)
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
.order-body { display: flex; justify-content: space-between; align-items: center; gap: 16px; }
.order-price { font-size: 20px; font-weight: 800; color: var(--yellow-600); }
.order-foot { display: flex; justify-content: flex-end; gap: 10px; margin-top: 14px; }
.meta { color: var(--ink-soft); font-size: 13px; line-height: 1.7; }
</style>
