import { axios } from '@/lib/axios' 

export interface TelegramUser {
  id: number
  first_name: string
  last_name?: string
  username?: string
  photo_url?: string
  auth_date: number
  hash: string
}

export interface TelegramAuthResponse {
  token: string
  user_id: string
  session_id: string
  username: string
  full_name: string
  role: string
}

export interface AuthUser {
  id: string
  username: string
  full_name: string
  role: string
}

export const telegramAuth = async (telegramData: TelegramUser): Promise<TelegramAuthResponse> => {
  const { data } = await axios.post('/api/telegram/auth', telegramData)

  localStorage.setItem('token', data.token)
  localStorage.setItem('user', JSON.stringify({
    id: data.user_id,
    username: data.username,
    full_name: data.full_name,
    role: data.role
  }))
  
  axios.defaults.headers.common['Authorization'] = `Bearer ${data.token}`

  return data
}

export const logout = (): void => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  delete axios.defaults.headers.common['Authorization']
}