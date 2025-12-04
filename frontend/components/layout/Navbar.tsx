'use client'
import { useState } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import TelegramLoginButton from '../../features/auth/components/TelegramLoginButton'
import { Bell, Menu, UserCircle2, ChevronDown } from 'lucide-react'

export function Navbar() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  return (
    <header className="w-full bg-background border-b border-border px-4 py-3 flex items-center justify-between gap-4">
      <Link href="/" className="flex items-center gap-2">
        <Image src="/file.svg" alt="Логотип" width={36} height={36} />
        <div className="leading-tight text-[var(--primary)] font-bold text-lg">
          codesign.tech
          <div className="text-xs font-normal text-muted-foreground">образование для взрослых</div>
        </div>
      </Link>

      <nav className="hidden md:flex gap-6 items-center">
        <Dropdown label="Курсы" items={['Все курсы', 'Направления', 'Избранное']} />
        <Dropdown label="Блог" items={['Все статьи', 'Советы', 'Истории студентов']} />
      </nav>

      <div className="hidden md:flex flex-1 max-w-md">
        <Input placeholder="Поиск..." className="w-full rounded-xl" />
      </div>

      <div className="flex items-center gap-3">
        {isAuthenticated ? (
  <>
    <Button variant="ghost" size="icon" aria-label="Уведомления">
      <Bell className="w-5 h-5 text-[var(--primary)]" />
    </Button>
    <UserCircle2 className="w-8 h-8 text-[var(--primary)]" />
  </>
) : (
  <div className="rounded-full overflow-hidden">
    <TelegramLoginButton />
  </div>
)}

        <Menu className="w-6 h-6 md:hidden text-muted-foreground" />
      </div>
    </header>
  )
}

function Dropdown({ label, items }: { label: string; items: string[] }) {
  return (
    <div className="relative group">
      <button className="flex items-center gap-1 text-sm font-medium text-foreground hover:text-[var(--primary)]">
        {label}
        <ChevronDown className="w-4 h-4" />
      </button>
      <div className="absolute hidden group-hover:block mt-2 w-44 bg-background border border-border rounded-md shadow z-50">
        <ul className="py-2">
          {items.map(item => (
            <li key={item}>
              <a href="#" className="block px-4 py-2 text-sm text-foreground hover:bg-muted">
                {item}
              </a>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}
