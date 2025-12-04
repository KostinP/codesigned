'use client'

import { useEffect, useRef } from 'react'
import { usePathname, useSearchParams } from 'next/navigation'
import { logEvent } from '@/lib/tracking'

export const usePageView = (): void => {
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const initialLoad = useRef(true)

  useEffect(() => {
    if (!pathname) return
    
    const url = pathname + (searchParams.toString() ? `?${searchParams.toString()}` : '')
    
    // Отправляем page_view только при SPA навигации, а не при первоначальной загрузке
    if (!initialLoad.current) {
      logEvent({
        event_type: 'page_view',
        event_data: {
          page: url,
          referrer: document.referrer || 'direct',
          title: document.title,
          url: window.location.href,
          timestamp: new Date().toISOString(),
          navigation_type: 'spa'
        }
      })
    }
    
    initialLoad.current = false
  }, [pathname, searchParams])
}