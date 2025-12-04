'use client'

import { useEffect } from 'react'
import { getVisitorFromStorage, getVisitorId } from '@/features/tracking/api/visitor'

export const useVisitorTracker = (): void => {
  useEffect(() => {
    const visitorId = getVisitorFromStorage()
    if (!visitorId) {
      getVisitorId().then((res) => {
        if (typeof window !== 'undefined') {
          localStorage.setItem('visitor_id', res.visitor_id)
        }
      }).catch((error) => {
        console.error('Failed to get visitor ID:', error)
      })
    }
  }, [])
}