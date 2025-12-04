import { axios } from '@/lib/axios' 

export interface UserSession {
  id: string
  user_id: string
  token: string
  ip_address: string
  city: string
  country: string
  user_agent: string
  created_at: string
  expires_at: string
  last_active_at: string
}

export interface InactivityTimeout {
  timeout_seconds: number
}

export const linkVisitor = async (userId: string, visitorId: string) => {
  await axios.post(`/api/users/${userId}/link-visitor`, {}, {
    headers: {
      'X-Visitor-ID': visitorId,
    },
  })
}

export const getUserSessions = async (): Promise<UserSession[]> => {
  const { data } = await axios.get('/api/me/sessions')
  return data
}

export const deleteUserSession = async (sessionId: string) => {
  await axios.delete(`/api/me/sessions/${sessionId}`)
}

export const getInactivityTimeout = async (): Promise<InactivityTimeout> => {
  const { data } = await axios.get('/api/me/inactivity-timeout')
  return data
}

export const setInactivityTimeout = async (timeoutSeconds: number) => {
  await axios.post('/api/me/inactivity-timeout', {
    timeout_seconds: timeoutSeconds,
  })
}