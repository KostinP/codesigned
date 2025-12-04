import axios from 'axios'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
})

axiosInstance.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token')
    const visitorId = localStorage.getItem('visitor_id')
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    if (visitorId) {
      config.headers['X-Visitor-ID'] = visitorId
    }
  }
  return config
})

// Добавляем интерсептор для автоматического обновления токена
axiosInstance.interceptors.response.use(
  (response) => {
    const newToken = response.headers['x-refresh-token']
    if (newToken && typeof window !== 'undefined') {
      localStorage.setItem('token', newToken)
      axiosInstance.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
    }
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // Токен истек или невалиден
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

// Экспортируем как именованный экспорт
export { axiosInstance as axios }
// ИЛИ экспортируем по умолчанию
// export default axiosInstance