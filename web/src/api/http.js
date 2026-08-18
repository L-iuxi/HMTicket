import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const http = axios.create({ baseURL: '', timeout: 12000 })

// 请求：自动带 token
http.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = 'Bearer ' + token
  return config
})

// 响应：统一 {code, msg, data}，code===200 才成功
http.interceptors.response.use(
  (res) => {
    const body = res.data
    if (body && typeof body === 'object' && 'code' in body) {
      if (body.code === 200) return body.data
      ElMessage.error(body.msg || '请求失败')
      return Promise.reject(new Error(body.msg || 'error'))
    }
    return body
  },
  (err) => {
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      router.push('/login')
      ElMessage.error('登录已过期，请重新登录')
    } else {
      ElMessage.error(err.message || '网络错误')
    }
    return Promise.reject(err)
  }
)

export default http
