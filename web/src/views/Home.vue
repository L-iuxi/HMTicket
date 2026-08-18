<template>
  <div>
    <NavBar />
    <div class="page">
      <h1 class="page-title">热门活动</h1>
      <div v-if="loading" class="hint">加载中…</div>
      <div v-else-if="events.length === 0" class="hint">暂无活动，去管理后台创建一个吧</div>
      <div v-else class="event-grid">
        <div
          v-for="e in events"
          :key="e.EventID"
          class="glass-card event-card"
          @click="$router.push(`/event/${e.EventID}`)"
        >
          <div class="cover" :class="{ 'is-default': !e.CoverImage }" :style="{ backgroundImage: `url(${e.CoverImage || '/default-cover.png'})` }">
            <span class="status-tag">{{ statusText(e.Status) }}</span>
          </div>
          <div class="body">
            <div class="title">{{ e.Title }}</div>
            <div class="meta">📍 {{ e.Location }}</div>
            <div class="meta">🕐 {{ fmt(e.StartTime) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import NavBar from '../components/NavBar.vue'
import { getEventList } from '../api'

const events = ref([])
const loading = ref(true)

const STATUS = { draft: '草稿', ready: '待售', selling: '售票中', closed: '已结束' }
function statusText(s) { return STATUS[s] || s }
function fmt(ts) {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}

onMounted(async () => {
  try {
    events.value = (await getEventList())?.Events?.filter((e) => e.Status !== 'draft') || []
  } catch (e) {
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.hint { padding: 40px 0; text-align: center; color: var(--ink-soft); }
</style>
