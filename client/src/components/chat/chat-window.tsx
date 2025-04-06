// components/chat/chat-window.tsx
'use client'

import { useWebSocket } from '@/components/providers/websocket-provider'
import { Message } from '@/components/chat/message'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {useState} from "react";

export function ChatWindow() {
  const { messages, sendMessage, status } = useWebSocket()
  const [input, setInput] = useState('')
  const [roomId, setRoomId] = useState('')

  const handleSend = () => {
    if (input.trim() && status === 'connected') {
      sendMessage({
        type: 'message',
        payload: JSON.stringify({ text: input }),
        timestamp: Date.now()
      })
      setInput('')
    }
  }

  const handleJoinRoom = () => {
    if (roomId.trim() && status === 'connected') {
      sendMessage({
        type: 'join',
        payload: JSON.stringify({ roomId }),
        timestamp: Date.now()
      })
    }
  }

  return (
    <div className="flex-1 flex flex-col bg-white rounded-xl shadow-sm overflow-hidden">
      <div className="p-4 border-b">
        <div className="flex gap-2">
          <Input
            value={roomId}
            onChange={(e) => setRoomId(e.target.value)}
            placeholder="Enter room ID"
            className="flex-1"
          />
          <Button onClick={handleJoinRoom} variant="secondary">
            Join Room
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message, i) => (
          <Message key={i} message={message} />
        ))}
      </div>

      <div className="p-4 border-t">
        <div className="flex gap-2">
          <Input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            placeholder="Type a message..."
          />
          <Button onClick={handleSend}>
            Send
          </Button>
        </div>
      </div>
    </div>
  )
}
