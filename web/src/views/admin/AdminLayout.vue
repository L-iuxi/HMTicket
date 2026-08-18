<template>
  <div class="admin">
    <aside class="glass-card side">
      <router-link to="/" class="brand">🎫 TicketX</router-link>
      <div class="side-sub">管理后台</div>
      <router-link to="/admin/events" class="item" :class="{ active: route.path.startsWith('/admin/event') }">活动管理</router-link>
      <router-link to="/admin/stock" class="item" :class="{ active: route.path === '/admin/stock' }">库存管理</router-link>
      <router-link to="/admin/orders" class="item" :class="{ active: route.path === '/admin/orders' }">订单管理</router-link>
      <div class="spacer"></div>
      <router-link to="/" class="item">← 返回前台</router-link>
    </aside>
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
const route = useRoute()
</script>

<style scoped>
.admin { display: flex; min-height: 100vh; padding: 16px; gap: 16px; }
.side { width: 220px; padding: 22px 14px; display: flex; flex-direction: column; gap: 6px; position: sticky; top: 16px; height: calc(100vh - 32px); }
.brand { font-size: 20px; font-weight: 800; margin: 0 10px 4px; }
.side-sub { font-size: 12px; color: var(--ink-soft); margin: 0 10px 12px; }
.item { padding: 11px 14px; border-radius: 12px; font-weight: 600; color: var(--ink-soft); transition: all .15s ease; }
.item:hover { background: rgba(255, 184, 0, 0.14); color: var(--ink); }
.item.active { background: rgba(255, 184, 0, 0.28); color: var(--ink); }
.spacer { flex: 1; }
.content { flex: 1; min-width: 0; }
</style>
