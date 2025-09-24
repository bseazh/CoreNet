import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'CoreNet Drive',
  description: 'Google Drive-like interface for CoreNet',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body className="min-h-screen bg-gray-50 text-gray-900">
        {children}
      </body>
    </html>
  )
}

