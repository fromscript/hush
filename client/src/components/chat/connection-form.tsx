// components/chat/connection-form.tsx
'use client'

import { useState } from 'react'
import { useWebSocket } from '@/components/providers/websocket-provider'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

export function ConnectionForm() {
  const { connect } = useWebSocket()
  const [formData, setFormData] = useState({
    url: process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws',
    token: ''
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    connect(formData.url, formData.token)
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white p-6 rounded-xl shadow-sm">
      <div className="space-y-4">
        <div>
          <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-1">
            WebSocket URL
          </label>
          <Input
            id="url"
            value={formData.url}
            onChange={(e) => setFormData(prev => ({ ...prev, url: e.target.value }))}
            required
          />
        </div>

        <div>
          <label htmlFor="token" className="block text-sm font-medium text-gray-700 mb-1">
            Authentication Token
          </label>
          <Input
            id="token"
            type="password"
            value={formData.token}
            onChange={(e) => setFormData(prev => ({ ...prev, token: e.target.value }))}
            required
          />
        </div>

        <Button type="submit" className="w-full">
          Connect
        </Button>
      </div>
    </form>
  )
}
