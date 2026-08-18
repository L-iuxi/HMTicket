<template>
  <div class="page">
    <button class="glass-btn sm ghost" style="margin-bottom: 16px;" @click="$router.push('/admin/events')">← 返回活动列表</button>
    <h1 class="page-title">场次与票种</h1>

    <div v-if="loading" class="hint">加载中…</div>
    <div v-else-if="shows.length === 0" class="hint">暂无场次</div>
    <div v-else class="list">
      <div v-for="s in shows" :key="s.ShowID" class="glass-card block">
        <div class="block-head">
          <div>
            <div class="t">{{ s.Name }}</div>
            <div class="meta">{{ s.ShowTime }} · {{ s.Venue || '—' }}</div>
          </div>
          <button class="glass-btn primary sm" @click="openTicketType(s)">＋ 加票种</button>
        </div>
        <div v-if="!ticketTypes[s.ShowID]" class="meta" style="padding: 8px 0;">加载票种…</div>
        <div v-else-if="ticketTypes[s.ShowID].length === 0" class="meta" style="padding: 8px 0;">暂无票种</div>
        <div v-else class="tt-list">
          <div v-for="t in ticketTypes[s.ShowID]" :key="t.TicketTypeID" class="tt-row">
            <span style="font-weight: 700;">{{ t.Name }}</span>
            <span style="color: var(--yellow-600); font-weight: 700;">¥{{ t.Price }}</span>
            <span class="meta">库存 {{ t.Stock }} · 限购 {{ t.MaxPerUser }}</span>
          </div>
        </div>
      </div>
    </div>

    <button class="glass-btn primary" style="margin-top: 20px;" @click="openShow">＋ 新建场次</button>

    <!-- 新建场次 -->
    <el-dialog v-model="showDialog" title="新建场次" width="420px" align-center>
      <el-input v-model="showForm.name" placeholder="场次名称" style="margin-bottom: 12px;" />
      <el-input v-model="showForm.venue" placeholder="场地（选填）" style="margin-bottom: 12px;" />
      <div style="display: flex; gap: 12px; margin-bottom: 12px;">
        <el-date-picker v-model="showForm.showTime" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="开始时间" style="flex: 1;" />
        <el-date-picker v-model="showForm.endTime" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="结束时间" style="flex: 1;" />
      </div>
      <template #footer>
        <button class="glass-btn" @click="showDialog = false">取消</button>
        <button class="glass-btn primary" :disabled="saving" @click="onCreateShow">{{ saving ? '保存中…' : '保存' }}</button>
      </template>
    </el-dialog>

    <!-- 新建票种 -->
    <el-dialog v-model="ttDialog" title="新建票种" width="420px" align-center>
      <el-input v-model="ttForm.name" placeholder="票种名称（如 VIP票）" style="margin-bottom: 12px;" />
      <div style="display: flex; gap: 12px; margin-bottom: 12px;">
        <div style="flex: 1;">
          <div class="lbl">价格</div>
          <el-input-number v-model="ttForm.price" :min="0" :precision="2" style="width: 100%;" />
        </div>
        <div style="flex: 1;">
          <div class="lbl">库存</div>
          <el-input-number v-model="ttForm.stock" :min="0" style="width: 100%;" />
        </div>
        <div style="flex: 1;">
          <div class="lbl">限购</div>
          <el-input-number v-model="ttForm.maxPerUser" :min="1" style="width: 100%;" />
        </div>
      </div>
      <template #footer>
        <button class="glass-btn" @click="ttDialog = false">取消</button>
        <button class="glass-btn primary" :disabled="saving" @click="onCreateTicketType">{{ saving ? '保存中…' : '保存' }}</button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getShowList, getTicketTypeList, createShow, createTicketType } from '../../api'

const route = useRoute()
const eventId = route.params.id

const shows = ref([])
const ticketTypes = reactive({})
const loading = ref(true)
const saving = ref(false)
const showDialog = ref(false)
const ttDialog = ref(false)
const currentShow = ref(null)

const showForm = reactive({ name: '', venue: '', showTime: '', endTime: '' })
const ttForm = reactive({ name: '', price: 0, stock: 0, maxPerUser: 1 })

async function load() {
  loading.value = true
  try {
    shows.value = (await getShowList(eventId))?.Shows || []
    for (const s of shows.value) {
      ticketTypes[s.ShowID] = (await getTicketTypeList(s.ShowID))?.TicketTypes || []
    }
  } catch (e) {
  } finally {
    loading.value = false
  }
}

function openShow() {
  Object.assign(showForm, { name: '', venue: '', showTime: '', endTime: '' })
  showDialog.value = true
}
function openTicketType(s) {
  currentShow.value = s
  Object.assign(ttForm, { name: '', price: 0, stock: 0, maxPerUser: 1 })
  ttDialog.value = true
}

async function onCreateShow() {
  if (!showForm.name || !showForm.showTime || !showForm.endTime) {
    ElMessage.warning('请填写名称和起止时间')
    return
  }
  saving.value = true
  try {
    await createShow({ eventId: Number(eventId), name: showForm.name, showTime: showForm.showTime, endTime: showForm.endTime, venue: showForm.venue })
    ElMessage.success('场次已创建')
    showDialog.value = false
    await load()
  } catch (e) {
  } finally {
    saving.value = false
  }
}

async function onCreateTicketType() {
  if (!ttForm.name || ttForm.price <= 0 || ttForm.stock <= 0) {
    ElMessage.warning('请填写名称、价格和库存')
    return
  }
  saving.value = true
  try {
    await createTicketType({
      eventId: Number(eventId),
      showId: currentShow.value.ShowID,
      name: ttForm.name,
      price: ttForm.price,
      stock: ttForm.stock,
      maxPerUser: ttForm.maxPerUser,
      sortOrder: 0
    })
    ElMessage.success('票种已创建')
    ttDialog.value = false
    await load()
  } catch (e) {
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.hint { padding: 24px 0; color: var(--ink-soft); }
.list { display: flex; flex-direction: column; gap: 16px; }
.block { padding: 18px 22px; }
.block-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.t { font-weight: 800; font-size: 17px; margin-bottom: 4px; }
.meta { color: var(--ink-soft); font-size: 13px; }
.tt-list { display: flex; flex-direction: column; gap: 8px; }
.tt-row { display: flex; gap: 24px; align-items: center; padding: 8px 14px; border-radius: 12px; background: rgba(255, 184, 0, 0.06); }
.lbl { font-size: 12px; color: var(--ink-soft); margin-bottom: 4px; }
</style>
