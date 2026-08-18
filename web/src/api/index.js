import http from './http'

// ===== 用户 =====
export const register = (data) => http.post('/api/user/register', data)
export const login = (data) => http.post('/api/user/login', data)
export const getProfile = () => http.get('/api/user/profile')
export const updateProfile = (data) => http.put('/api/user/profile', data)

// ===== 活动/场次/票种 =====
export const getEventList = () => http.get('/event/list')
export const getEvent = (id) => http.get(`/event/${id}`)
export const getShowList = (eventId) => http.get('/show/list', { params: { eventId } })
export const getTicketTypeList = (showId) => http.get('/ticket-type/list', { params: { showId } })
export const getStock = (ticketTypeId) => http.get(`/api/v1/inventory/stock/${ticketTypeId}`)

// ===== 订单 =====
export const buyTicket = (data) => http.post('/api/v1/order/buy', data)
export const payOrder = (data) => http.post('/api/v1/order/pay', data)
export const getOrderList = () => http.get('/api/v1/order/list')
export const cancelOrder = (orderNo) => http.post('/api/v1/order/cancel', { orderNo })

// ===== 票 =====
export const listTicket = (params) => http.get('/api/ticket/list', { params })
export const getTicket = (id) => http.get(`/api/ticket/detail/${id}`)
export const refundTicket = (ticketId) => http.post('/api/ticket/refund', { ticketId })

// ===== 管理端 =====
export const createEvent = (data) => http.post('/admin/event', data)
export const updateEvent = (data) => http.put('/admin/event', data)
export const updateEventStatus = (data) => http.put('/admin/event/status', data)
export const createShow = (data) => http.post('/admin/show', data)
export const updateShow = (data) => http.put('/admin/show', data)
export const createTicketType = (data) => http.post('/admin/ticket-type', data)
export const updateTicketType = (data) => http.put('/admin/ticket-type', data)
export const updateStock = (data) => http.put('/api/v1/admin/inventory/stock', data)
export const updateOrder = (data) => http.put('/api/v1/admin/order/update', data)
export const deleteOrder = (orderNo) => http.delete(`/api/v1/admin/order/${orderNo}`)
