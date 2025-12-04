import type { Metadata, Viewport } from "next"
import { Inter, Nunito_Sans } from 'next/font/google'
import { Navbar } from "@/components/layout/Navbar"
import "./globals.css"
import Footer from "@/components/layout/FooterSection"
import ClientLayout from "@/components/ClientLayout"

const inter = Inter({
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-inter',
})

const nunitoSans = Nunito_Sans({
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-nunito-sans',
})

export const metadata: Metadata = {
  title: "codesigned.ru – образование для взрослых",
  description: "Современное интерактивное образование для взрослых",
  icons: {
    icon: "/favicon.ico",
  },
}

export const viewport: Viewport = {
  themeColor: "#545BE8",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="ru">
      <body className={`${inter.variable} ${nunitoSans.variable} antialiased`}>
        <Navbar />
        <main>
          <ClientLayout>
            {children}
          </ClientLayout>
        </main>
        <Footer />
      </body>
    </html>
  )
}