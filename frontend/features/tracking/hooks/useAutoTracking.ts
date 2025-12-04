'use client'

import { useEffect, useRef } from 'react'
import { logEvent } from '@/lib/tracking'

export const useAutoTracking = (): void => {
  const pageViewSent = useRef(false)
  const lastScrollY = useRef(0)
  const scrollTimeout = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    // Отправляем page_view только один раз при загрузке
    if (!pageViewSent.current) {
      setTimeout(() => {
        logEvent({ 
          event_type: 'page_view',
          event_data: { 
            pathname: window.location.pathname,
            url: window.location.href,
            referrer: document.referrer || 'direct',
            timestamp: new Date().toISOString()
          }
        })
        pageViewSent.current = true
      }, 500)
    }

    const handleClick = (e: MouseEvent): void => {
      const target = e.target as HTMLElement
      
      // Отслеживаем только значимые клики
      const isSignificantClick = 
        target.matches('button, a, [data-track], [role="button"]') ||
        target.closest('button, a, [data-track], [role="button"]')
      
      if (isSignificantClick) {
        const name = target.getAttribute('data-track') || 
                    target.getAttribute('aria-label') ||
                    target.textContent?.trim().substring(0, 50) ||
                    target.tagName
        
        logEvent({ 
          event_type: 'click',
          event_data: { 
            name,
            element: target.tagName,
            href: target.getAttribute('href'),
            timestamp: new Date().toISOString()
          }
        })
      }
    }

    const handleScroll = (): void => {
      const currentScrollY = window.scrollY
      
      if (scrollTimeout.current) {
        clearTimeout(scrollTimeout.current)
      }

      scrollTimeout.current = setTimeout(() => {
        if (Math.abs(currentScrollY - lastScrollY.current) > 100) {
          const scrollDepth = Math.round((currentScrollY / (document.documentElement.scrollHeight - window.innerHeight)) * 100)
          
          logEvent({ 
            event_type: 'scroll',
            event_data: { 
              position: currentScrollY,
              depth_percentage: scrollDepth,
              timestamp: new Date().toISOString()
            }
          })
          lastScrollY.current = currentScrollY
        }
      }, 300)
    }

    window.addEventListener('click', handleClick)
    window.addEventListener('scroll', handleScroll)

    return () => {
      window.removeEventListener('click', handleClick)
      window.removeEventListener('scroll', handleScroll)
      if (scrollTimeout.current) {
        clearTimeout(scrollTimeout.current)
      }
    }
  }, [])
}