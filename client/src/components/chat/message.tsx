// components/chat/message.tsx
'use client'

import { cn } from '@/lib/utils'
import {WebSocketMessage} from "@/lib/websocket/types";

interface MessageProps {
  message: WebSocketMessage
}

export function Message({ message }: MessageProps) {
  const isSystem = message.type === 'system'
  const content = JSON.parse(message.payload)

  return (
    <div className={cn(
      'p-4 rounded-lg max-w-[80%]',
      isSystem ? 'bg-gray-100 mx-auto' : 'bg-blue-500 text-white ml-auto'
    )}>
      <div className="text-sm font-medium mb-1">
        {isSystem ? 'System' : 'You'}
      </div>
      <div className="text-sm">{content.text}</div>
      {message.timestamp && (
        <div className="text-xs mt-2 opacity-70">
          {new Date(message.timestamp).toLocaleTimeString()}
        </div>
      )}
    </div>
  )
}
