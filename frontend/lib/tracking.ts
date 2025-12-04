import { axios } from '@/lib/axios'

export interface TrackingEvent {
  event_type: string  // изменяем с type на event_type
  event_data?: Record<string, unknown>  // изменяем с data на event_data
}

export const logEvent = async (event: TrackingEvent): Promise<void> => {
  const visitorId = typeof window !== 'undefined' ? localStorage.getItem('visitor_id') : null
  if (!visitorId) {
    console.warn('No visitor ID found for event tracking')
    return
  }

  try {
    const payload = {
      event_type: event.event_type,
      event_data: event.event_data || {},
      visitor_id: visitorId  // добавляем visitor_id в тело запроса
    }

    console.log('Sending event:', payload) // для отладки

    const response = await axios.post('/api/visitor/events', payload)
    
    if (response.status !== 200) {
      console.warn('Event tracking returned non-200 status:', response.status)
    }
  } catch (error) {
    console.error('Event tracking error', error)
    // В development показываем детали ошибки
    if (process.env.NODE_ENV === 'development') {
      console.error('Error details:', error)
    }
  }
}