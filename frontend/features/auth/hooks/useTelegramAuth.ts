import { useMutation } from '@tanstack/react-query'
import { telegramAuth, TelegramAuthResponse } from '@/features/auth/api/auth'
import { linkVisitor } from '@/features/user/api/user'
import { getVisitorFromStorage } from '@/features/tracking/api/visitor'
import { useAuthStore } from '@/features/auth/store/useAuthStore'

interface TelegramUser {
  id: number
  first_name: string
  last_name?: string
  username?: string
  photo_url?: string
  auth_date: number
  hash: string
}

export const useTelegramAuth = () => {
  const setUser = useAuthStore((state) => state.setUser)
  
  return useMutation({
    mutationFn: async (telegramData: TelegramUser): Promise<TelegramAuthResponse> => {
      const authData = await telegramAuth(telegramData)
      
      // Сохраняем пользователя в store
      setUser({
        id: authData.user_id,
        username: authData.username,
        full_name: authData.full_name,
        role: authData.role
      })

      // Связываем visitor с пользователем
      const visitorId = getVisitorFromStorage()
      if (visitorId) {
        try {
          await linkVisitor(authData.user_id, visitorId)
        } catch (error) {
          console.error('Failed to link visitor:', error)
          // Не прерываем процесс авторизации из-за этой ошибки
        }
      }
      
      return authData
    },
  })
}