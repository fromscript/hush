'use client'

import { useWebSocket } from '@/components/providers/websocket-provider'
import { ConnectionForm } from '@/components/chat/connection-form'
import { ChatWindow } from '@/components/chat/chat-window'
import { LoadingSkeleton } from '@/components/chat/loading-skeleton'

export default function ChatPage() {
  const { status, error } = useWebSocket()

  return (
    <main className="container mx-auto p-4 h-full flex flex-col">
      <div className="max-w-4xl mx-auto w-full flex-1 flex flex-col">
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-3xl font-bold text-gray-900">WebSocket Chat</h1>
          <div className="flex items-center gap-2">
            <div className={`h-3 w-3 rounded-full ${
              status === 'connected' ? 'bg-green-500' :
                status === 'connecting' ? 'bg-yellow-500' : 'bg-red-500'
            }`} />
            <span className="text-sm text-gray-600 capitalize">{status}</span>
          </div>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 text-red-700 rounded-lg">
            {error}
          </div>
        )}

        {status === 'disconnected' ? (
          <ConnectionForm />
        ) : status === 'connecting' ? (
          <LoadingSkeleton />
        ) : (
          <ChatWindow />
        )}
      </div>
    </main>
  )
}
