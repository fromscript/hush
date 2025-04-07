export interface WebSocketMessage {
  type: 'message' | 'system' | 'join'
  payload: any
  timestamp?: number
  sessionId?: string
}

export type WebSocketState = {
  status: 'disconnected' | 'connecting' | 'connected';
  messages: WebSocketMessage[];
  error: string | null;
  connect: (url: string, token: string) => void;
  sendMessage: (message: WebSocketMessage) => void;
  joinRoom: (roomID: string) => void;
  sendChatMessage: (content: string) => void;
}
