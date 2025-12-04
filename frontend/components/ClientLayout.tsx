'use client'

import dynamic from "next/dynamic"

// Динамически импортируем клиентские компоненты с отключенным SSR
const ClientTracking = dynamic(() => import('@/components/ClientTracking'), {
  ssr: false
})

interface ClientLayoutProps {
  children: React.ReactNode
}

export default function ClientLayout({ children }: ClientLayoutProps) {
  return (
    <>
      {children}
      <ClientTracking />
    </>
  )
}