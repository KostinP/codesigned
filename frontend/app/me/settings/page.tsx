import { UserSessions } from '@/features/session/components/UserSessions'
import { InactivityTimeout } from '@/features/session/components/InactivityTimeout'

export default function SettingsPage() {
  return (
    <div className="max-w-2xl mx-auto py-10 px-4">
      <h1 className="text-2xl font-bold mb-6">Настройки аккаунта</h1>
      <UserSessions />
      <InactivityTimeout />
    </div>
  )
}
