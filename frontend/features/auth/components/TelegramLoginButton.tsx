// features/auth/components/TelegramLoginButton.tsx
'use client'

import { useEffect } from 'react'

interface TelegramUser {
  id: number
  first_name: string
  last_name?: string
  username?: string
  photo_url?: string
  auth_date: number
  hash: string
}

// Правильно расширяем Window interface
declare global {
  interface Window {
    handleTelegramAuth?: (user: TelegramUser) => void
  }
}

export default function TelegramLoginButton() {
  useEffect(() => {
    const handleTelegramAuth = (user: TelegramUser) => {
      console.log('Telegram auth:', user)
      // TODO: Добавить реальную логику авторизации
    }

    // Создаем кнопку Telegram
    const script = document.createElement('script')
    script.src = 'https://telegram.org/js/telegram-widget.js?22'
    script.async = true
    script.setAttribute('data-telegram-login', 'codesigned_dev_bot')
    script.setAttribute('data-size', 'large')
    script.setAttribute('data-radius', '10')
    script.setAttribute('data-request-access', 'write')
    script.setAttribute('data-userpic', 'true')
    script.setAttribute('data-onauth', 'handleTelegramAuth(user)')
    
    const container = document.getElementById('telegram-login-container')
    if (container) {
      container.appendChild(script)
    }

    // Добавляем глобальную функцию для обработки авторизации
    window.handleTelegramAuth = handleTelegramAuth

    return () => {
      if (container) {
        container.innerHTML = ''
      }
      // Используем присваивание undefined вместо delete
      window.handleTelegramAuth = undefined
    }
  }, [])

  return (
    <div id="telegram-login-container" className="telegram-login-wrapper">
      {/* Кнопка будет вставлена скриптом */}
    </div>
  )
}