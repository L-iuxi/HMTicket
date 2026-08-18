<template>
  <header class="glass-nav">
    <router-link to="/" class="brand">🎫 TicketX</router-link>
    <router-link to="/" class="nav-link" :class="{ active: route.path === '/' }">首页</router-link>
    <template v-if="store.isLogin">
      <router-link to="/orders" class="nav-link" :class="{ active: route.path === '/orders' }">我的订单</router-link>
      <router-link to="/tickets" class="nav-link" :class="{ active: route.path === '/tickets' }">我的票</router-link>
    </template>
    <div class="spacer"></div>
    <template v-if="!store.isLogin">
      <router-link to="/login"><button class="glass-btn sm ghost">登录</button></router-link>
      <router-link to="/register"><button class="glass-btn sm primary">注册</button></router-link>
    </template>
    <template v-else>
      <router-link v-if="store.isAdmin" to="/admin" class="nav-link">管理后台</router-link>
      <router-link to="/profile" class="nav-link">{{ store.profile?.username || '我的' }}</router-link>
      <button class="glass-btn sm ghost" @click="logout">退出</button>
    </template>
  </header>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const route = useRoute()
const router = useRouter()
const store = useUserStore()

function logout() {
  store.logout()
  router.push('/')
}
</script>
