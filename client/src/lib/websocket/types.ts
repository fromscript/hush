export interface WebSocketMessage {
  type: 'message' | 'system' | 'join'
  payload: string
  timestamp?: number
  sessionId?: string
}

export interface WebSocketState {
  status: 'disconnected' | 'connecting' | 'connected'
  messages: WebSocketMessage[]
  error: string | null
  connect: (url: string, token: string) => void
  sendMessage: (message: WebSocketMessage) => void
}
