// app/providers.tsx
'use client'

import { ThemeProvider } from 'next-themes'
import { WebSocketProvider } from '@/components/providers/websocket-provider'

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <WebSocketProvider>
        {children}
      </WebSocketProvider>
    </ThemeProvider>
  )
}
