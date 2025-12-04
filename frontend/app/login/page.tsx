'use client'

import {TelegramLoginButton} from '@/features/auth/components/TelegramLoginButton'

export default function LoginPage() {
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1 className="text-3xl font-bold mb-4">Вход через Telegram</h1>
      <TelegramLoginButton />
    </div>
  )
}
