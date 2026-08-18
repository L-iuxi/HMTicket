<template>
  <div class="page">
    <div class="head">
      <h1 class="page-title" style="margin: 0;">活动管理</h1>
      <button class="glass-btn primary" @click="openCreate">＋ 新建活动</button>
    </div>

    <div v-if="loading" class="hint">加载中…</div>
    <div v-else class="list">
      <div v-for="e in events" :key="e.EventID" class="glass-card row">
        <div class="row-info">
          <div class="t">{{ e.Title }}</div>
          <div class="meta">📍 {{ e.Location }} · 库存 {{ e.TotalStock }}</div>
        </div>
        <span class="tag" :class="statusCls(e.Status)">{{ statusText(e.Status) }}</span>
        <div class="actions">
          <button class="glass-btn sm" @click="$router.push(`/admin/event/${e.EventID}`)">场次/票种</button>
          <button class="glass-btn sm" @click="openEdit(e)">编辑</button>
          <el-select v-model="e.Status" size="small" style="width: 110px;" @change="(v) => onStatus(e, v)">
            <el-option v-for="s in STATUS_KEYS" :key="s" :label="STATUS[s]" :value="s" />
          </el-select>
        </div>
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="editing ? '编辑活动' : '新建活动'" width="460px" align-center>
      <el-input v-model="form.title" placeholder="标题" style="margin-bottom: 12px;" />
      <el-input v-model="form.location" placeholder="地点" style="margin-bottom: 12px;" />
      <el-input v-model="form.coverImage" placeholder="封面图 URL（选填）" style="margin-bottom: 12px;" />
      <el-input v-model="form.description" type="textarea" :rows="2" placeholder="描述（选填）" style="margin-bottom: 12px;" />
      <div style="display: flex; gap: 12px;">
        <el-date-picker v-model="form.startTime" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="开始时间" style="flex: 1;" />
        <el-date-picker v-model="form.endTime" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="结束时间" style="flex: 1;" />
      </div>
      <template #footer>
        <button class="glass-btn" @click="dialogVisible = false">取消</button>
        <button class="glass-btn primary" :disabled="saving" @click="onSave">{{ saving ? '保存中…' : '保存' }}</button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getEventList, createEvent, updateEvent, updateEventStatus } from '../../api'

const events = ref([])
const loading = ref(true)
const dialogVisible = ref(false)
const editing = ref(null)
const saving = ref(false)
const form = reactive({ title: '', location: '', coverImage: '', description: '', startTime: '', endTime: '' })

const STATUS = { draft: '草稿', ready: '待售', selling: '售票中', closed: '已结束' }
const STATUS_KEYS = ['draft', 'ready', 'selling', 'closed']
function statusText(s) { return STATUS[s] || s }
function statusCls(s) { return ({ draft: 'gray', ready: 'yellow', selling: 'green', closed: 'gray' }[s]) || 'gray' }

async function load() {
  try {
    events.value = (await getEventList())?.Events || []
  } catch (e) {
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  Object.assign(form, { title: '', location: '', coverImage: '', description: '', startTime: '', endTime: '' })
  dialogVisible.value = true
}
function openEdit(e) {
  editing.value = e
  Object.assign(form, {
    title: e.Title, location: e.Location, coverImage: e.CoverImage, description: e.Description,
    startTime: fmt(e.StartTime), endTime: fmt(e.EndTime)
  })
  dialogVisible.value = true
}

function fmt(ts) {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

async function onSave() {
  if (!form.title || !form.location || !form.startTime || !form.endTime) {
    ElMessage.warning('请填写标题、地点、开始和结束时间')
    return
  }
  saving.value = true
  try {
    if (editing.value) {
      await updateEvent({ eventId: editing.value.EventID, ...form })
      ElMessage.success('已更新')
    } else {
      await createEvent({ ...form })
      ElMessage.success('已创建')
    }
    dialogVisible.value = false
    await load()
  } catch (e) {
  } finally {
    saving.value = false
  }
}

async function onStatus(e, status) {
  try {
    await updateEventStatus({ eventId: e.EventID, status })
    ElMessage.success('状态已更新')
  } catch (err) {
    await load()
  }
}

onMounted(load)
</script>

<style scoped>
.hint { padding: 40px 0; text-align: center; color: var(--ink-soft); }
.head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.list { display: flex; flex-direction: column; gap: 14px; }
.row { display: flex; align-items: center; gap: 20px; padding: 16px 22px; }
.row-info { flex: 1; min-width: 0; }
.t { font-weight: 800; font-size: 17px; margin-bottom: 4px; }
.meta { color: var(--ink-soft); font-size: 13px; }
.actions { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
</style>
