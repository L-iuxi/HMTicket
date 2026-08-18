<template>
  <div>
    <NavBar />
    <div class="page">
      <div v-if="ticket" class="glass-card" style="max-width: 640px; margin: 0 auto; padding: 32px;">
        <h1 class="page-title" style="margin-top: 0;">票详情</h1>
        <div class="row"><span class="k">票号</span><span class="v">{{ ticket.ticketId }}</span></div>
        <div class="row"><span class="k">订单号</span><span class="v">{{ ticket.orderNo }}</span></div>
        <div class="row"><span class="k">状态</span><span class="v"><span class="tag" :class="statusCls(ticket.status)">{{ statusText(ticket.status) }}</span></span></div>
        <div class="row"><span class="k">活动 ID</span><span class="v">{{ ticket.eventId }}</span></div>
        <div class="row"><span class="k">场次 ID</span><span class="v">{{ ticket.showId }}</span></div>
        <div class="row"><span class="k">票种 ID</span><span class="v">{{ ticket.ticketTypeId }}</span></div>
        <div class="row"><span class="k">数量</span><span class="v">{{ ticket.quantity }}</span></div>
        <div class="row"><span class="k">金额</span><span class="v">¥{{ ticket.totalPrice }}</span></div>
        <div class="row" v-if="ticket.qrCode"><span class="k">核销码</span><span class="v mono">{{ ticket.qrCode }}</span></div>
        <div class="row" v-if="ticket.realName"><span class="k">实名</span><span class="v">{{ ticket.realName }}</span></div>
        <div class="row" v-if="ticket.phone"><span class="k">手机</span><span class="v">{{ ticket.phone }}</span></div>
        <div class="row"><span class="k">创建时间</span><span class="v">{{ ticket.createdAt }}</span></div>
      </div>
      <div v-else class="hint">加载中…</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import NavBar from '../components/NavBar.vue'
import { getTicket } from '../api'

const route = useRoute()
const ticket = ref(null)

const STATUS = {
  unused: { text: '未使用', cls: 'green' },
  used: { text: '已使用', cls: 'gray' },
  refunded: { text: '已退票', cls: 'red' }
}
function statusText(s) { return (STATUS[s] || { text: s }).text }
function statusCls(s) { return (STATUS[s] || { cls: 'gray' }).cls }

onMounted(async () => {
  try {
    ticket.value = await getTicket(route.params.id)
  } catch (e) {
  }
})
</script>

<style scoped>
.hint { padding: 40px 0; text-align: center; color: var(--ink-soft); }
.row { display: flex; gap: 20px; padding: 10px 0; border-bottom: 1px dashed rgba(43, 34, 0, 0.1); }
.k { min-width: 90px; color: var(--ink-soft); font-size: 14px; }
.v { font-weight: 600; word-break: break-all; }
.mono { font-family: monospace; }
</style>
