<template>
  <div class="page">
    <h1 class="page-title">库存管理</h1>
    <div class="glass-card" style="max-width: 520px; padding: 24px 28px;">
      <div class="field">
        <label>票种 ID</label>
        <el-input v-model="ticketTypeId" placeholder="输入票种 ID" />
      </div>
      <button class="glass-btn sm" :disabled="querying" @click="onQuery">查询当前库存</button>
      <div v-if="currentStock !== null" class="stock-info">
        <div class="meta">当前库存：<b style="font-size: 22px; color: var(--yellow-600);">{{ currentStock }}</b></div>
        <div class="field" style="margin-top: 14px;">
          <label>新库存</label>
          <el-input-number v-model="newStock" :min="0" style="width: 100%;" />
        </div>
        <button class="glass-btn primary" :disabled="saving" @click="onUpdate">{{ saving ? '保存中…' : '更新库存' }}</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getStock, updateStock } from '../../api'

const ticketTypeId = ref('')
const currentStock = ref(null)
const newStock = ref(0)
const querying = ref(false)
const saving = ref(false)

async function onQuery() {
  if (!ticketTypeId.value) return ElMessage.warning('请输入票种 ID')
  querying.value = true
  try {
    const data = await getStock(ticketTypeId.value)
    currentStock.value = data.stock
    newStock.value = data.stock
  } catch (e) {
    currentStock.value = null
  } finally {
    querying.value = false
  }
}

async function onUpdate() {
  if (!ticketTypeId.value) return ElMessage.warning('请输入票种 ID')
  saving.value = true
  try {
    await updateStock({ ticketTypeId: Number(ticketTypeId.value), stock: newStock.value })
    ElMessage.success('库存已更新')
    currentStock.value = newStock.value
  } catch (e) {
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.meta { color: var(--ink-soft); font-size: 13px; }
.field { margin-bottom: 16px; }
.field label { display: block; margin-bottom: 6px; font-size: 13px; color: var(--ink-soft); font-weight: 600; }
.stock-info { margin-top: 18px; padding-top: 18px; border-top: 1px dashed rgba(43, 34, 0, 0.12); }
</style>
