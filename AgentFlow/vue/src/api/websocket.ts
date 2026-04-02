// WebSocket client for agentServer communication
import { WEBSOCKET_URL } from '../config/api'

// Protocol types matching backend/internal/protocol/messages.go
export interface Attachment {
  type?: string
  url?: string
}

export interface MessageContent {
  text: string
  attachments?: Attachment[]
  mentions?: string[]
}

export interface ReturnStrategy {
  type: string
  final_result_marker: string
  chunk_size: number
  debounce_ms: number
}

export interface UserInteractionConfig {
  stream_mode: string
  enable_thinking: boolean
  return_strategy: ReturnStrategy
  timeout_ms: number
  max_retries: number
}

export interface MessageMetadata {
  override_config?: UserInteractionConfig
  required_skills?: string[]
  target_agent_id?: string
  model_id?: string
}

export interface Message {
  message_id: string
  session_id: string
  group_id?: string
  user_id: string
  type: string
  content: MessageContent
  metadata: MessageMetadata
  timestamp: number
}

export interface StreamChunkPayload {
  content: string
  is_thinking: boolean
  is_final: boolean
  progress?: number
  metadata?: Record<string, any>
}

export interface StreamEndPayload {
  content: string
  is_final: boolean
  usage?: Record<string, number>
}

export interface ErrorPayload {
  code: string
  message: string
  retryable: boolean
}

export interface ServerMessage {
  type: string
  message_id?: string
  session_id?: string
  payload?: StreamChunkPayload | StreamEndPayload | ErrorPayload | any
  timestamp?: number
}

type WebSocketEvent = 'open' | 'close' | 'error' | 'message' | 'connection-change'

type EventListener = (data?: any) => void

export class AgentWebSocket {
  private ws: WebSocket | null = null
  private url: string
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000 // Start with 1 second
  private maxReconnectDelay = 30000 // Max 30 seconds
  private listeners: Map<WebSocketEvent, EventListener[]> = new Map()
  private isConnected = false
  private connectionId: string | null = null
  private sessionId: string | null = null
  private userId: string | null = null
  private pendingSendQueue: Message[] = []

  constructor(url?: string) {
    this.url = url || WEBSOCKET_URL
  }

  connect(): void {
    if (this.ws && (this.ws.readyState === WebSocket.CONNECTING || this.ws.readyState === WebSocket.OPEN)) {
      console.warn('WebSocket already connecting or connected')
      return
    }

    try {
      console.log('[AgentWebSocket] Connecting to WebSocket:', this.url)
      this.ws = new WebSocket(this.url)
      this.setupEventListeners()
    } catch (error) {
      console.error('[AgentWebSocket] Failed to create WebSocket connection:', error)
      this.scheduleReconnect()
    }
  }

  private setupEventListeners(): void {
    if (!this.ws) return

    this.ws.onopen = (event) => {
      console.log('[AgentWebSocket] WebSocket connected:', this.url)
      this.isConnected = true
      this.reconnectAttempts = 0
      this.reconnectDelay = 1000
      this.emit('open', event)
      this.emit('connection-change', true)
      this.flushPendingQueue()
    }

    this.ws.onclose = (event) => {
      console.log('[AgentWebSocket] WebSocket disconnected:', {
        url: this.url,
        code: event.code,
        reason: event.reason,
        wasClean: event.wasClean
      })
      this.isConnected = false
      this.emit('close', event)
      this.emit('connection-change', false)
      
      // Server-initiated close codes (1000=normal, 1001=going away, 1005=no status)
      if (event.code !== 1000 && event.code !== 1001 && event.code !== 1005) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = (error) => {
      console.error('[AgentWebSocket] WebSocket error:', error)
      this.emit('error', error)
    }

    this.ws.onmessage = (event) => {
      try {
        console.log('[AgentWebSocket] Raw message received:', event.data)
        const data = JSON.parse(event.data) as ServerMessage
        this.handleIncomingMessage(data)
        this.emit('message', data)
      } catch (error) {
        console.error('[AgentWebSocket] Failed to parse WebSocket message:', error, event.data)
      }
    }
  }

  private handleIncomingMessage(message: ServerMessage): void {
    // Update connection/session info if provided
    if (message.session_id) {
      this.sessionId = message.session_id
    }
    
    console.log('[AgentWebSocket] Parsed server message:', {
      type: message.type,
      messageId: message.message_id,
      sessionId: message.session_id
    })
  }

  send(message: Message): void {
    if (!this.isConnected || !this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn('[AgentWebSocket] WebSocket not connected; queueing message and connecting', {
        isConnected: this.isConnected,
        readyState: this.ws?.readyState
      })
      this.pendingSendQueue.push(message)
      this.connect()
      return
    }

    try {
      console.log('[AgentWebSocket] Sending message:', message)
      this.ws.send(JSON.stringify(message))
    } catch (error) {
      console.error('[AgentWebSocket] Failed to send WebSocket message:', error, message)
    }
  }

  private flushPendingQueue(): void {
    if (!this.isConnected || !this.ws || this.ws.readyState !== WebSocket.OPEN) return
    if (this.pendingSendQueue.length === 0) return

    const queue = [...this.pendingSendQueue]
    this.pendingSendQueue = []
    queue.forEach(msg => this.send(msg))
  }

  sendText(text: string, userId: string, sessionId?: string, metadata?: Partial<MessageMetadata>): string {
    const messageId = this.generateMessageId()
    const resolvedUserId = userId || this.userId || ''
    const message: Message = {
      message_id: messageId,
      session_id: sessionId || this.sessionId || this.generateMessageId(),
      user_id: resolvedUserId,
      type: 'user_message',
      content: {
        text,
        attachments: [],
        mentions: []
      },
      metadata: metadata || {},
      timestamp: Date.now()
    }
    
    this.send(message)
    return messageId
  }

  sendGroupText(text: string, userId: string, groupId: string, sessionId?: string, metadata?: Partial<MessageMetadata>): string {
    const messageId = this.generateMessageId()
    const resolvedUserId = userId || this.userId || ''
    const message: Message = {
      message_id: messageId,
      session_id: sessionId || this.sessionId || this.generateMessageId(),
      group_id: groupId,
      user_id: resolvedUserId,
      type: 'user_message',
      content: {
        text,
        attachments: [],
        mentions: []
      },
      metadata: metadata || {},
      timestamp: Date.now()
    }
    
    this.send(message)
    return messageId
  }

  setUser(userId: string): void {
    this.userId = userId
    console.log('[AgentWebSocket] Active user set:', userId)
  }

  setSession(sessionId: string): void {
    this.sessionId = sessionId
  }

  clearSession(): void {
    this.sessionId = null
    console.log('[AgentWebSocket] Session cleared')
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close(1000, 'Client disconnect')
      this.ws = null
    }
    this.isConnected = false
    this.emit('connection-change', false)
  }

  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('[AgentWebSocket] Max reconnection attempts reached')
      return
    }

    this.reconnectAttempts++
    const delay = Math.min(this.reconnectDelay * Math.pow(1.5, this.reconnectAttempts - 1), this.maxReconnectDelay)
    
    console.log(`[AgentWebSocket] Scheduling reconnect attempt ${this.reconnectAttempts} in ${delay}ms`)
    
    setTimeout(() => {
      this.connect()
    }, delay)
  }

  private generateMessageId(): string {
    return 'msg_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
  }

  // Event emitter methods
  on(event: WebSocketEvent, listener: EventListener): void {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, [])
    }
    this.listeners.get(event)!.push(listener)
  }

  off(event: WebSocketEvent, listener: EventListener): void {
    const listeners = this.listeners.get(event)
    if (!listeners) return
    
    const index = listeners.indexOf(listener)
    if (index !== -1) {
      listeners.splice(index, 1)
    }
  }

  private emit(event: WebSocketEvent, data?: any): void {
    const listeners = this.listeners.get(event)
    if (!listeners) return
    
    listeners.forEach(listener => {
      try {
        listener(data)
      } catch (error) {
        console.error(`Error in ${event} event listener:`, error)
      }
    })
  }

  getConnectionStatus(): boolean {
    return this.isConnected
  }

  getConnectionId(): string | null {
    return this.connectionId
  }

  getSessionId(): string | null {
    return this.sessionId
  }
}

// Singleton instance for global use
let globalWebSocketInstance: AgentWebSocket | null = null

export function getWebSocketInstance(url?: string): AgentWebSocket {
  if (!globalWebSocketInstance) {
    globalWebSocketInstance = new AgentWebSocket(url)
  }
  return globalWebSocketInstance
}
