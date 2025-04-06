// app/layout.tsx
import { Inter } from 'next/font/google'
import './globals.css'
import { Providers } from "@/app/providers";

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-sans',
  display: 'swap'
})

export default function RootLayout({
                                     children,
                                   }: {
  children: React.ReactNode
}) {
  return (
    <html
      lang="en"
      className={`${inter.variable}`}
      suppressHydrationWarning
    >
    <body className="min-h-screen bg-background font-sans antialiased">
    <Providers>
      <main className="relative flex flex-col">
        {children}
      </main>
    </Providers>
    </body>
    </html>
  )
}
