'use client'

import { useAutoTracking } from '@/features/tracking/hooks/useAutoTracking'
import { useSectionTracking } from '@/features/tracking/hooks/useSectionTracking'
import { useKeepAlive } from '@/features/session/hooks/useKeepAlive'
import { useVisitorTracker } from '@/features/tracking/hooks/useVisitorTracker'
// import { usePageView } from '@/features/tracking/hooks/usePageView'

export default function ClientTracking() {
  // Все клиентские хуки вызываются здесь
  useAutoTracking()       // Основные события (page_view, click, scroll)
  useSectionTracking()    // Специально для отслеживания секций
  useVisitorTracker()     // Visitor ID
  useKeepAlive()          // Активность сессии
  // usePageView()        //

  // Этот компонент ничего не рендерит
  return null
}