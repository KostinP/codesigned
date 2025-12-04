import Link from 'next/link'
import Image from 'next/image'

export default function FooterSection() {
  return (
    <footer className="bg-[var(--gray-bg)] py-12 px-4 text-muted-foreground">
      <div className="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-3 gap-8">
        <div className="flex flex-col gap-2">
          <div className="flex items-center gap-2">
            <Image src="/file.svg" alt="Logo" width={36} height={36} />
            <span className="font-bold text-foreground text-lg">codesign.tech</span>
          </div>
          <p className="text-sm">образование для взрослых</p>
        </div>

        <div className="flex flex-col gap-2">
          <span className="font-semibold text-foreground">Навигация</span>
          <Link href="#" className="hover:underline">Курсы</Link>
          <Link href="#" className="hover:underline">Блог</Link>
          <Link href="#" className="hover:underline">Контакты</Link>
        </div>

        <div className="flex flex-col gap-2">
          <span className="font-semibold text-foreground">Мы в соцсетях</span>
          <Link href="#" className="hover:underline">Telegram</Link>
          <Link href="#" className="hover:underline">YouTube</Link>
        </div>
      </div>

      <div className="text-center text-xs mt-8 text-muted-foreground">
        © 2025 codesign.tech
      </div>
    </footer>
  )
}
