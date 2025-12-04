'use client'

import { useEffect, useRef } from 'react'
import { logEvent } from '@/lib/tracking'

interface SectionInfo {
  id: string
  distance: number
  element: Element
}

export const useSectionTracking = (): void => {
  const lastActiveSection = useRef<string>('')
  const sectionChangeTimeout = useRef<NodeJS.Timeout | null>(null)

  useEffect(() => {
    const generateSectionId = (section: Element): string => {
      if (section.id) return section.id
      
      const heading = section.querySelector('h1, h2, h3')
      if (heading?.textContent) {
        return heading.textContent.toLowerCase().replace(/[^a-zа-я0-9]/g, '_')
      }
      
      return `section_${Math.round(section.getBoundingClientRect().top)}`
    }

    const getSectionName = (section: Element): string => {
      const heading = section.querySelector('h1, h2, h3')
      return heading?.textContent?.trim() || 'Без названия'
    }

    const findMostVisibleSection = (): SectionInfo | null => {
      const sections = document.querySelectorAll('section')
      const viewportHeight = window.innerHeight
      const viewportCenter = viewportHeight / 2
      
      let mostVisibleSection: SectionInfo | null = null
      
      sections.forEach((section) => {
        const rect = section.getBoundingClientRect()
        const sectionCenter = rect.top + (rect.height / 2)
        const distanceFromCenter = Math.abs(sectionCenter - viewportCenter)
        
        if (mostVisibleSection === null || distanceFromCenter < mostVisibleSection.distance) {
          const sectionId = section.id || generateSectionId(section)
          mostVisibleSection = {
            id: sectionId,
            distance: distanceFromCenter,
            element: section
          }
        }
      })
      
      return mostVisibleSection
    }

    const handleSectionView = (): void => {
      const mostVisibleSection = findMostVisibleSection()
      
      if (!mostVisibleSection) return
      
      if (mostVisibleSection.id !== lastActiveSection.current) {
        if (sectionChangeTimeout.current) {
          clearTimeout(sectionChangeTimeout.current)
        }

        sectionChangeTimeout.current = setTimeout(() => {
          const sectionName = getSectionName(mostVisibleSection.element)
          logEvent({ 
            event_type: 'section_view',
            event_data: { 
              section_id: mostVisibleSection.id,
              section_name: sectionName,
              timestamp: new Date().toISOString()
            }
          })
          lastActiveSection.current = mostVisibleSection.id
        }, 500)
      }
    }

    let scrollTimeout: NodeJS.Timeout | null = null
    const throttledScroll = () => {
      if (scrollTimeout) return
      
      scrollTimeout = setTimeout(() => {
        handleSectionView()
        scrollTimeout = null
      }, 100)
    }

    window.addEventListener('scroll', throttledScroll, { passive: true })
    setTimeout(handleSectionView, 1000)

    return () => {
      window.removeEventListener('scroll', throttledScroll)
      if (sectionChangeTimeout.current) {
        clearTimeout(sectionChangeTimeout.current)
      }
      if (scrollTimeout) {
        clearTimeout(scrollTimeout)
      }
    }
  }, [])
}