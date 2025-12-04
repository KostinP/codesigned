'use client'

import { useEffect, useRef } from 'react'
import { usePathname, useSearchParams } from 'next/navigation'
import { logEvent } from '@/lib/tracking'
import { useSessionStore } from '@/features/session/store/session'

interface PageStats {
  pageStartTime: number
  pingsSent: number
}

export const useKeepAlive = (): void => {
  const timeout = useSessionStore((s) => s.inactivityTimeout)
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const startTimeRef = useRef<number>(Date.now())
  const pageStatsRef = useRef<PageStats>({
    pageStartTime: Date.now(),
    pingsSent: 0
  })
  const intervalRef = useRef<NodeJS.Timeout | null>(null)

  useEffect(() => {
    const getPageInfo = () => {
      return {
        pathname: pathname,
        url: `${window.location.origin}${pathname}${searchParams?.toString() ? `?${searchParams.toString()}` : ''}`,
        hostname: window.location.hostname,
        search: searchParams?.toString() || '',
        referrer: document.referrer || ''
      }
    }

    const handleVisibilityChange = () => {
      if (document.hidden) {
        if (intervalRef.current) {
          clearInterval(intervalRef.current)
          intervalRef.current = null
        }
      } else {
        if (!intervalRef.current) {
          startPinging()
        }
      }
    }

    const startPinging = () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }

      intervalRef.current = setInterval(() => {
        const currentTime = Date.now()
        const sessionDuration = Math.round((currentTime - startTimeRef.current) / 1000)
        const pageViewDuration = Math.round((currentTime - pageStatsRef.current.pageStartTime) / 1000)
        pageStatsRef.current.pingsSent++
        
        const pageInfo = getPageInfo()
        
        logEvent({ 
          event_type: 'activity_ping',
          event_data: {
            session_duration_seconds: sessionDuration,
            page_view_duration_seconds: pageViewDuration,
            page_pings_sent: pageStatsRef.current.pingsSent,
            tab_visible: !document.hidden,
            timestamp: new Date().toISOString(),
            ...pageInfo
          }
        })
      }, 10000)
    }

    // Событие начала просмотра новой страницы
    const pageInfo = getPageInfo()
    logEvent({
      event_type: 'page_view_start',
      event_data: {
        timestamp: new Date().toISOString(),
        ...pageInfo
      }
    })

    startPinging()
    document.addEventListener('visibilitychange', handleVisibilityChange)

    // Событие ухода со страницы
    const handleBeforeUnload = () => {
      const pageViewDuration = Math.round((Date.now() - pageStatsRef.current.pageStartTime) / 1000)
      const finalPageInfo = getPageInfo()
      
      logEvent({
        event_type: 'page_view_end',
        event_data: {
          page_view_duration_seconds: pageViewDuration,
          page_pings_sent: pageStatsRef.current.pingsSent,
          timestamp: new Date().toISOString(),
          ...finalPageInfo
        }
      })
    }

    window.addEventListener('beforeunload', handleBeforeUnload)

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }
      document.removeEventListener('visibilitychange', handleVisibilityChange)
      window.removeEventListener('beforeunload', handleBeforeUnload)
      
      // Событие окончания просмотра страницы при размонтировании
      const pageViewDuration = Math.round((Date.now() - pageStatsRef.current.pageStartTime) / 1000)
      const finalPageInfo = getPageInfo()
      
      logEvent({
        event_type: 'page_view_end',
        event_data: {
          page_view_duration_seconds: pageViewDuration,
          page_pings_sent: pageStatsRef.current.pingsSent,
          timestamp: new Date().toISOString(),
          ...finalPageInfo
        }
      })
    }
  }, [timeout, pathname, searchParams])

  useEffect(() => {
    // Сброс статистики страницы при смене страницы
    const now = Date.now()
    
    // Отправляем событие окончания предыдущей страницы
    if (pageStatsRef.current.pageStartTime > 0) {
      const previousPageDuration = Math.round((now - pageStatsRef.current.pageStartTime) / 1000)
      logEvent({
        event_type: 'page_view_end',
        event_data: {
          page_view_duration_seconds: previousPageDuration,
          page_pings_sent: pageStatsRef.current.pingsSent,
          timestamp: new Date().toISOString(),
          pathname: pathname, // Новая страница уже, но для контекста
        }
      })
    }
    
    // Сбрасываем статистику для новой страницы
    pageStatsRef.current = {
      pageStartTime: now,
      pingsSent: 0
    }
  }, [pathname])

  useEffect(() => {
    startTimeRef.current = Date.now()
  }, [timeout])
}