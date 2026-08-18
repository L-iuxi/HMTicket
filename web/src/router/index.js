import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/login', component: () => import('../views/Login.vue') },
  { path: '/register', component: () => import('../views/Register.vue') },
  { path: '/', component: () => import('../views/Home.vue') },
  { path: '/event/:id', component: () => import('../views/EventDetail.vue') },
  { path: '/orders', component: () => import('../views/Orders.vue'), meta: { auth: true } },
  { path: '/tickets', component: () => import('../views/Tickets.vue'), meta: { auth: true } },
  { path: '/ticket/:id', component: () => import('../views/TicketDetail.vue'), meta: { auth: true } },
  { path: '/profile', component: () => import('../views/Profile.vue'), meta: { auth: true } },
  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    meta: { auth: true, admin: true },
    children: [
      { path: '', redirect: '/admin/events' },
      { path: 'events', component: () => import('../views/admin/Events.vue') },
      { path: 'event/:id', component: () => import('../views/admin/EventManage.vue') },
      { path: 'orders', component: () => import('../views/admin/Orders.vue') },
      { path: 'stock', component: () => import('../views/admin/Stock.vue') }
    ]
  }
]

const router = createRouter({ history: createWebHashHistory(), routes })

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  const role = localStorage.getItem('role')
  if (to.meta.auth && !token) return '/login'
  if (to.meta.admin && role !== 'admin') return '/'
  return true
})

export default router
