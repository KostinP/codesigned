'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getUserSessions, deleteUserSession, UserSession } from '@/features/user/api/user'
import { Button } from '@/components/ui/button'
import { formatDistanceToNow } from 'date-fns'
import { ru } from 'date-fns/locale'

export const UserSessions = () => {
  const queryClient = useQueryClient()

  const { data: sessions, isLoading, isError } = useQuery<UserSession[]>({
    queryKey: ['sessions'],
    queryFn: getUserSessions,
  })

  const { mutate: removeSession, isPending } = useMutation({
    mutationFn: deleteUserSession,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sessions'] })
    },
  })

  if (isLoading) return <p className="text-muted-foreground">Загрузка сессий...</p>
  if (isError) return <p className="text-destructive">Ошибка загрузки сессий</p>

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold">Активные сессии</h2>
      {sessions?.length === 0 ? (
        <p className="text-muted-foreground">Нет активных сессий</p>
      ) : (
        sessions?.map((session) => (
          <div key={session.id} className="border rounded-lg p-4 flex justify-between items-center">
            <div className="flex-1">
              <p className="font-medium">
                {session.city || 'Неизвестно'}, {session.country || 'Неизвестно'} — {session.ip_address}
              </p>
              <p className="text-sm text-muted-foreground mt-1">
                Устройство: {session.user_agent?.substring(0, 50)}...
              </p>
              <p className="text-sm text-muted-foreground">
                Последняя активность: {formatDistanceToNow(new Date(session.last_active_at), { addSuffix: true, locale: ru })}
              </p>
              <p className="text-sm text-muted-foreground">
                Создана: {formatDistanceToNow(new Date(session.created_at), { addSuffix: true, locale: ru })}
              </p>
            </div>
            <Button
              variant="destructive"
              size="sm"
              disabled={isPending}
              onClick={() => removeSession(session.id)}
              className="ml-4"
            >
              Выйти
            </Button>
          </div>
        ))
      )}
    </div>
  )
}