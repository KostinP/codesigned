'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getInactivityTimeout, setInactivityTimeout } from '@/features/user/api/user'
import { Slider } from '@/components/ui/slider'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { useState, useEffect } from 'react'

export const InactivityTimeout = () => {
  const queryClient = useQueryClient()
  const [localTimeout, setLocalTimeout] = useState<number>(60)

  const { data, isLoading } = useQuery({
    queryKey: ['inactivity-timeout'],
    queryFn: getInactivityTimeout,
  })

  useEffect(() => {
    if (data?.timeout_seconds) {
      setLocalTimeout(data.timeout_seconds / 60)
    }
  }, [data])

  const mutationOptions = {
    mutationFn: (timeoutMinutes: number) => setInactivityTimeout(timeoutMinutes * 60),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inactivity-timeout'] })
    },
  }

  const { mutate, isPending } = useMutation(mutationOptions)

  const handleSave = () => {
    mutate(localTimeout)
  }

  if (isLoading) return <p className="text-muted-foreground">Загрузка...</p>

  return (
    <div className="space-y-4 mt-8 p-4 border rounded-lg">
      <div>
        <Label className="text-base font-medium">Таймаут неактивности</Label>
        <p className="text-sm text-muted-foreground mt-1">
          Автоматический выход из системы после указанного периода неактивности
        </p>
      </div>
      
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium">{localTimeout} минут</span>
        </div>
        
        <Slider
          value={[localTimeout]}
          min={5}
          max={1440}
          step={5}
          onValueChange={(value) => setLocalTimeout(value[0])}
          className="my-4"
        />
        
        <div className="flex gap-2 text-xs text-muted-foreground">
          <span>5 мин</span>
          <span className="flex-1 text-center">12 часов</span>
          <span>24 часа</span>
        </div>
      </div>

      <Button 
        onClick={handleSave} 
        disabled={isPending}
        className="w-full"
      >
        {isPending ? 'Сохранение...' : 'Сохранить настройки'}
      </Button>
    </div>
  )
}