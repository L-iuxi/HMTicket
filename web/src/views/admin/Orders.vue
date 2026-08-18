<template>
  <div class="page">
    <h1 class="page-title">订单管理</h1>
    <div class="glass-card" style="max-width: 520px; padding: 24px 28px;">
      <p class="meta" style="margin-top: 0;">后端暂无「全部订单列表」接口，此处按订单号修改数量 / 删除。</p>
      <div class="field">
        <label>订单号</label>
        <el-input v-model="orderNo" placeholder="如 3740c823-ca42-4400-9af0-cc7971848b98" />
      </div>
      <div class="field">
        <label>新数量</label>
        <el-input-number v-model="quantity" :min="1" style="width: 100%;" />
      </div>
      <div style="display: flex; gap: 12px;">
        <button class="glass-btn primary" :disabled="busy" @click="onUpdate">修改数量</button>
        <button class="glass-btn danger" :disabled="busy" @click="onDelete">删除订单</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { updateOrder, deleteOrder } from '../../api'

const orderNo = ref('')
const quantity = ref(1)
const busy = ref(false)

async function onUpdate() {
  if (!orderNo.value) return ElMessage.warning('请输入订单号')
  busy.value = true
  try {
    await updateOrder({ orderNo: orderNo.value, quantity: quantity.value })
    ElMessage.success('已修改')
  } catch (e) {
  } finally {
    busy.value = false
  }
}

async function onDelete() {
  if (!orderNo.value) return ElMessage.warning('请输入订单号')
  try {
    await ElMessageBox.confirm(`确认删除订单 ${orderNo.value}？`, '删除订单', { type: 'warning' })
  } catch (e) {
    return
  }
  busy.value = true
  try {
    await deleteOrder(orderNo.value)
    ElMessage.success('已删除')
    orderNo.value = ''
  } catch (e) {
  } finally {
    busy.value = false
  }
}
</script>

<style scoped>
.meta { color: var(--ink-soft); font-size: 13px; }
.field { margin-bottom: 16px; }
.field label { display: block; margin-bottom: 6px; font-size: 13px; color: var(--ink-soft); font-weight: 600; }
</style>
