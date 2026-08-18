<template>
  <div class="auth-wrap">
    <div class="glass-card auth-card">
      <h2 class="auth-title">登录</h2>
      <p class="auth-sub">欢迎回到 TicketX</p>
      <el-input v-model="form.account" size="large" placeholder="账号 / 用户名" style="margin-bottom: 14px;" />
      <el-input
        v-model="form.password"
        type="password"
        size="large"
        placeholder="密码"
        show-password
        style="margin-bottom: 22px;"
        @keyup.enter="onLogin"
      />
      <button class="glass-btn primary block" :disabled="loading" @click="onLogin">
        {{ loading ? '登录中…' : '登录' }}
      </button>
      <p class="auth-foot">
        没有账号？<router-link to="/register" class="link">去注册</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'

const router = useRouter()
const store = useUserStore()
const form = reactive({ account: '', password: '' })
const loading = ref(false)

async function onLogin() {
  if (!form.account || !form.password) {
    ElMessage.warning('请输入账号和密码')
    return
  }
  loading.value = true
  try {
    await store.login({ account: form.account, password: form.password })
    ElMessage.success('登录成功')
    router.push('/')
  } catch (e) {
    /* 错误已由拦截器提示 */
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
