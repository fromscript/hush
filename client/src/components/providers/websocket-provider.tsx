// components/providers/websocket-provider.tsx
'use client'

import { createContext, useContext, useEffect, useRef, useState } from 'react'
import type { WebSocketMessage, WebSocketState } from '@/lib/websocket/types'

const WebSocketContext = createContext<WebSocketState | null>(null)

export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const [isMounted, setIsMounted] = useState(false)
  const [status, setStatus] = useState<'disconnected' | 'connecting' | 'connected'>('disconnected')
  const [messages, setMessages] = useState<WebSocketMessage[]>([])
  const [error, setError] = useState<string | null>(null)
  const ws = useRef<WebSocket | null>(null)

  useEffect(() => {
    setIsMounted(true)
    return () => {
      ws.current?.close()
      setIsMounted(false)
    }
  }, [])

  const connect = (url: string, token: string) => {
    if (!isMounted) return

    setStatus('connecting')

    try {
      ws.current = new WebSocket(`${url}?token=${token}`)

      ws.current.onopen = () => {
        if (isMounted) {
          setStatus('connected')
          setError(null)
        }
      }

      ws.current.onmessage = (event) => {
        if (!isMounted) return
        try {
          const message = JSON.parse(event.data)
          setMessages(prev => [...prev, message])
        } catch (err) {
          console.error('Invalid message format:', err)
        }
      }

      ws.current.onclose = () => {
        if (isMounted) {
          setStatus('disconnected')
          ws.current = null
        }
      }

      ws.current.onerror = (error) => {
        if (isMounted) {
          setStatus('disconnected')
          setError('Connection failed. Please check your settings.')
          console.error('WebSocket error:', error)
        }
      }

    } catch (err) {
      if (isMounted) {
        setStatus('disconnected')
        setError('Invalid WebSocket URL')
      }
    }
  }

  const sendMessage = (message: WebSocketMessage) => {
    if (isMounted && ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(message))
    }
  }

  if (!isMounted) return null

  return (
    <WebSocketContext.Provider value={{
      status,
      messages,
      error,
      connect,
      sendMessage
    }}>
      {children}
    </WebSocketContext.Provider>
  )
}

export const useWebSocket = () => {
  const context = useContext(WebSocketContext)
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider')
  }
  return context
}
