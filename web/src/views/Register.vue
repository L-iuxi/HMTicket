<template>
  <div class="auth-wrap">
    <div class="glass-card auth-card">
      <h2 class="auth-title">注册</h2>
      <p class="auth-sub">创建你的 TicketX 账号</p>
      <el-input v-model="form.username" size="large" placeholder="用户名" style="margin-bottom: 14px;" />
      <el-input v-model="form.email" size="large" placeholder="邮箱" style="margin-bottom: 14px;" />
      <el-input v-model="form.phone" size="large" placeholder="手机号（选填）" style="margin-bottom: 14px;" />
      <el-input
        v-model="form.password"
        type="password"
        size="large"
        placeholder="密码"
        show-password
        style="margin-bottom: 22px;"
      />
      <button class="glass-btn primary block" :disabled="loading" @click="onRegister">
        {{ loading ? '注册中…' : '注册' }}
      </button>
      <p class="auth-foot">
        已有账号？<router-link to="/login" class="link">去登录</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { register } from '../api'

const router = useRouter()
const form = reactive({ username: '', email: '', phone: '', password: '' })
const loading = ref(false)

async function onRegister() {
  if (!form.username || !form.email || !form.password) {
    ElMessage.warning('请填写用户名、邮箱和密码')
    return
  }
  loading.value = true
  try {
    await register({ ...form })
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } catch (e) {
    /* 拦截器已提示 */
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-title { margin: 0 0 4px; font-size: 26px; font-weight: 800; }
.auth-sub { color: var(--ink-soft); margin: 0 0 24px; font-size: 14px; }
.auth-foot { text-align: center; margin: 18px 0 0; font-size: 13px; color: var(--ink-soft); }
.link { color: var(--yellow-600); font-weight: 700; }
</style>
