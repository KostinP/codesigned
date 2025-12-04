import { axios } from '@/lib/axios'
import type { AxiosError } from 'axios'

export interface VisitorResponse {
  visitor_id: string
  message: string
}

export interface VisitorEvent {
  event_type: string
  event_data: Record<string, unknown>
  visitor_id?: string
}

export const getVisitorId = async (): Promise<VisitorResponse> => {
  const { data } = await axios.get('/api/visitor')
  return data
}

export const logVisitorEvent = async (
  eventType: string,
  eventData: Record<string, unknown> = {}
): Promise<void> => {
  const visitorId = getVisitorFromStorage()
  
  try {
    const payload: VisitorEvent = {
      event_type: eventType,
      event_data: eventData
    }

    // Добавляем visitor_id если он есть
    if (visitorId) {
      payload.visitor_id = visitorId
    }

    const response = await axios.post('/api/visitor/events', payload)
    
    if (response.status !== 200) {
      console.warn('Event tracking returned non-200 status:', response.status)
    }
  } catch (error) {
    const axiosError = error as AxiosError
    // Не логировать ошибки в продакшене, только в development
    if (process.env.NODE_ENV === 'development') {
      console.warn('Failed to log visitor event (non-critical):', axiosError.message)
      
      // В development можно посмотреть детали ошибки
      if (axiosError.response) {
        console.warn('Error response:', axiosError.response.data)
      }
    }
  }
}

// Вспомогательная функция
export const getVisitorFromStorage = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('visitor_id')
  }
  return null
}