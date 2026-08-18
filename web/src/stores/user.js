import { defineStore } from 'pinia'
import { login as loginApi, getProfile } from '../api'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    role: localStorage.getItem('role') || '',
    profile: null
  }),
  getters: {
    isLogin: (s) => !!s.token,
    isAdmin: (s) => s.role === 'admin'
  },
  actions: {
    async login(payload) {
      const data = await loginApi(payload)
      this.token = data.token
      localStorage.setItem('token', data.token)
      await this.fetchProfile()
    },
    async fetchProfile() {
      this.profile = await getProfile()
      this.role = this.profile.role || ''
      localStorage.setItem('role', this.role)
    },
    logout() {
      this.token = ''
      this.role = ''
      this.profile = null
      localStorage.removeItem('token')
      localStorage.removeItem('role')
    }
  }
})
