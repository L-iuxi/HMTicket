<template>
  <div>
    <NavBar />
    <div class="page" v-if="event">
      <div class="glass-card hero">
        <div class="hero-cover" :class="{ 'is-default': !event.CoverImage }" :style="{ backgroundImage: `url(${event.CoverImage || '/default-cover.png'})` }"></div>
        <div class="hero-body">
          <div class="hero-title-row">
            <h1 class="hero-title">{{ event.Title }}</h1>
            <span class="tag yellow">{{ statusText(event.Status) }}</span>
          </div>
          <p class="meta">📍 {{ event.Location }} · 🕐 {{ fmt(event.StartTime) }} ~ {{ fmt(event.EndTime) }}</p>
          <p class="desc">{{ event.Desc }}</p>
        </div>
      </div>

      <h2 class="page-title" style="margin-top: 36px;">场次</h2>
      <div v-if="!showsLoaded" class="hint">加载场次中…</div>
      <div v-else-if="event.Shows.length === 0" class="hint">暂无场次</div>
      <div v-for="s in event.Shows" :key="s.ShowID" class="glass-card show-card">
        <div class="show-head">
          <div class="show-name">{{ s.Name }}</div>
          <div class="meta">{{ s.ShowTime }} · {{ s.Venue }}</div>
        </div>
        <div class="ticket-list">
          <div v-if="!ticketTypes[s.ShowID]" class="hint">加载票种中…</div>
          <div v-for="t in ticketTypes[s.ShowID] || []" :key="t.TicketTypeID" class="ticket-row">
            <div class="t-name">{{ t.Name }}</div>
            <div class="t-price">¥{{ t.Price }}</div>
            <div class="meta">库存 {{ t.Stock }} · 限购 {{ t.MaxPerUser }}</div>
            <button class="glass-btn primary sm" :disabled="t.Stock <= 0" @click="openBuy(s, t)">
              {{ t.Stock <= 0 ? '售罄' : '购票' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="buyVisible" title="确认购票" width="380px" align-center>
      <div v-if="buyTarget">
        <div style="margin-bottom: 14px; font-weight: 700;">
          {{ event.Title }} · {{ buyTarget.show.Name }} · {{ buyTarget.type.Name }}
        </div>
        <div class="meta" style="margin-bottom: 10px;">单价 ¥{{ buyTarget.type.Price }} · 限购 {{ buyTarget.type.MaxPerUser }} 张</div>
        <el-input-number v-model="qty" :min="1" :max="buyTarget.type.MaxPerUser" size="large" />
        <div style="margin-top: 18px; font-size: 20px; font-weight: 800;">
          合计 ¥{{ (buyTarget.type.Price * qty).toFixed(2) }}
        </div>
      </div>
      <template #footer>
        <button class="glass-btn" @click="buyVisible = false">取消</button>
        <button class="glass-btn primary" :disabled="buying" @click="onBuy">
          {{ buying ? '下单中…' : '立即下单' }}
        </button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import NavBar from '../components/NavBar.vue'
import { getEvent, getTicketTypeList, buyTicket } from '../api'
import { useUserStore } from '../stores/user'
import { genId } from '../utils/id'

const route = useRoute()
const router = useRouter()
const store = useUserStore()

const event = ref(null)
const ticketTypes = reactive({})
const showsLoaded = ref(false)
const buyVisible = ref(false)
const buyTarget = ref(null)
const qty = ref(1)
const buying = ref(false)

const STATUS = { draft: '草稿', ready: '待售', selling: '售票中', closed: '已结束' }
function statusText(s) { return STATUS[s] || s }
function fmt(ts) {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}

function openBuy(show, type) {
  if (!store.isLogin) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  buyTarget.value = { show, type }
  qty.value = 1
  buyVisible.value = true
}

async function onBuy() {
  buying.value = true
  try {
    const res = await buyTicket({
      eventId: event.value.EventID,
      showId: buyTarget.value.show.ShowID,
      ticketTypeId: buyTarget.value.type.TicketTypeID,
      quantity: qty.value,
      requestId: genId()
    })
    ElMessage.success('下单成功，订单号 ' + res.orderNo)
    buyVisible.value = false
    router.push('/orders')
  } catch (e) {
    /* 拦截器已提示 */
  } finally {
    buying.value = false
  }
}

onMounted(async () => {
  try {
    event.value = await getEvent(route.params.id)
    for (const s of event.value.Shows || []) {
      getTicketTypeList(s.ShowID).then((list) => {
        ticketTypes[s.ShowID] = list?.TicketTypes || []
      })
    }
  } catch (e) {
  } finally {
    showsLoaded.value = true
  }
})
</script>

<style scoped>
.hint { padding: 24px 0; color: var(--ink-soft); }
.hero { display: flex; overflow: hidden; }
.hero-cover { width: 40%; min-height: 240px; background-size: cover; background-position: center; }
.hero-cover.is-default { background-size: contain; background-repeat: no-repeat; }
.hero-body { flex: 1; padding: 28px 30px; }
.hero-title-row { display: flex; align-items: center; gap: 12px; }
.hero-title { margin: 0; font-size: 28px; font-weight: 800; }
.meta { color: var(--ink-soft); font-size: 14px; line-height: 1.8; }
.desc { margin-top: 14px; font-size: 15px; line-height: 1.7; }
.show-card { padding: 20px 24px; margin-bottom: 18px; }
.show-head { display: flex; align-items: baseline; gap: 16px; margin-bottom: 14px; }
.show-name { font-weight: 800; font-size: 18px; }
.ticket-list { display: flex; flex-direction: column; gap: 10px; }
.ticket-row { display: flex; align-items: center; gap: 20px; padding: 12px 16px; border-radius: 14px; background: rgba(255, 184, 0, 0.06); }
.t-name { font-weight: 700; min-width: 120px; }
.t-price { font-weight: 800; font-size: 18px; color: var(--yellow-600); min-width: 80px; }
.ticket-row .meta { flex: 1; }
@media (max-width: 640px) { .hero { flex-direction: column; } .hero-cover { width: 100%; } }
</style>
