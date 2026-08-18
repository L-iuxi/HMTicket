<template>
  <div>
    <NavBar />
    <div class="page">
      <h1 class="page-title">个人中心</h1>
      <div class="glass-card" style="max-width: 520px; padding: 28px 30px;">
        <div v-if="!loaded" class="hint">加载中…</div>
        <template v-else>
          <div class="field">
            <label>用户名</label>
            <el-input v-model="form.username" />
          </div>
          <div class="field">
            <label>邮箱</label>
            <el-input v-model="form.email" />
          </div>
          <div class="field">
            <label>手机号</label>
            <el-input v-model="form.phone" />
          </div>
          <div class="field">
            <label>性别</label>
            <el-select v-model="form.gender" style="width: 100%;">
              <el-option label="未知" :value="0" />
              <el-option label="男" :value="1" />
              <el-option label="女" :value="2" />
            </el-select>
          </div>
          <div class="field">
            <label>新密码（留空不改）</label>
            <el-input v-model="form.newPassword" type="password" show-password />
          </div>
          <div class="field">
            <label>旧密码（改密码时必填）</label>
            <el-input v-model="form.oldPassword" type="password" show-password />
          </div>
          <button class="glass-btn primary block" :disabled="saving" @click="onSave">
            {{ saving ? '保存中…' : '保存修改' }}
          </button>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import NavBar from '../components/NavBar.vue'
import { getProfile, updateProfile } from '../api'
import { useUserStore } from '../stores/user'

const store = useUserStore()
const loaded = ref(false)
const saving = ref(false)
const form = reactive({ username: '', email: '', phone: '', gender: 0, newPassword: '', oldPassword: '' })

onMounted(async () => {
  try {
    const p = await getProfile()
    form.username = p.username
    form.email = p.email
    form.phone = p.phone
    form.gender = p.gender || 0
  } catch (e) {
  } finally {
    loaded.value = true
  }
})

async function onSave() {
  saving.value = true
  try {
    const payload = { email: form.email, phone: form.phone, gender: form.gender, username: form.username }
    if (form.newPassword) payload.newPassword = form.newPassword
    if (form.oldPassword) payload.oldPassword = form.oldPassword
    await updateProfile(payload)
    ElMessage.success('保存成功')
    await store.fetchProfile()
  } catch (e) {
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.hint { color: var(--ink-soft); }
.field { margin-bottom: 16px; }
.field label { display: block; margin-bottom: 6px; font-size: 13px; color: var(--ink-soft); font-weight: 600; }
</style>
