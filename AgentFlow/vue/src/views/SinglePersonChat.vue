<script setup lang="ts">
import { ref, computed, nextTick, onBeforeUnmount, onMounted, watch } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL, WEBSOCKET_URL } from '../config/api'
import { getWebSocketInstance, type ServerMessage } from '../api/websocket'
import 'highlight.js/styles/atom-one-dark.css'
import { unified } from 'unified'
import remarkParse from 'remark-parse'
import remarkGfm from 'remark-gfm'
import remarkBreaks from 'remark-breaks'
import remarkRehype from 'remark-rehype'
import rehypeHighlight from 'rehype-highlight'
import rehypeStringify from 'rehype-stringify'
import prettier from 'prettier/standalone'
import prettierMarkdown from 'prettier/plugins/markdown'
import mermaid from 'mermaid'

const mdProcessor = unified()
  .use(remarkParse)
  .use(remarkGfm)
  .use(remarkBreaks)
  .use(remarkRehype, { allowDangerousHtml: false })
  .use(rehypeHighlight, { detect: true, ignoreMissing: true })
  .use(rehypeStringify)

let mermaidInitialized = false
let mermaidTimer: number | null = null
const mermaidPreviewOpen = ref(false)
const mermaidPreviewHtml = ref('')
const mermaidZoom = ref(1)
const mermaidPreviewBase = ref<{ w: number, h: number } | null>(null)
const mermaidPreviewRef = ref<HTMLElement | null>(null)

const clampMermaidZoom = (v: number) => Math.min(3, Math.max(0.3, Number.isFinite(v) ? v : 1))

const applyMermaidZoom = () => {
  const root = mermaidPreviewRef.value
  if (!root) return
  const svg = root.querySelector('svg') as SVGSVGElement | null
  if (!svg) return

  if (!mermaidPreviewBase.value) {
    svg.style.width = ''
    svg.style.height = ''
    svg.style.maxWidth = 'none'
    const rect = svg.getBoundingClientRect()
    if (rect.width > 0 && rect.height > 0) {
      mermaidPreviewBase.value = { w: rect.width, h: rect.height }
    } else {
      mermaidPreviewBase.value = { w: 800, h: 600 }
    }
  }

  const base = mermaidPreviewBase.value
  svg.style.maxWidth = 'none'
  svg.style.width = `${Math.round(base.w * mermaidZoom.value)}px`
  svg.style.height = `${Math.round(base.h * mermaidZoom.value)}px`
}

const openMermaidPreview = (html: string) => {
  mermaidPreviewHtml.value = html || ''
  mermaidZoom.value = 1
  mermaidPreviewBase.value = null
  mermaidPreviewOpen.value = true
  nextTick(() => applyMermaidZoom())
}

const closeMermaidPreview = () => {
  mermaidPreviewOpen.value = false
  mermaidPreviewHtml.value = ''
  mermaidPreviewBase.value = null
}

const zoomMermaidIn = () => {
  mermaidZoom.value = clampMermaidZoom(Number((mermaidZoom.value + 0.1).toFixed(2)))
  nextTick(() => applyMermaidZoom())
}

const zoomMermaidOut = () => {
  mermaidZoom.value = clampMermaidZoom(Number((mermaidZoom.value - 0.1).toFixed(2)))
  nextTick(() => applyMermaidZoom())
}

const resetMermaidZoom = () => {
  mermaidZoom.value = 1
  nextTick(() => applyMermaidZoom())
}

const handleMermaidWheel = (e: WheelEvent) => {
  if (!e.ctrlKey && !e.metaKey) return
  e.preventDefault()
  const delta = e.deltaY > 0 ? -0.1 : 0.1
  mermaidZoom.value = clampMermaidZoom(Number((mermaidZoom.value + delta).toFixed(2)))
  nextTick(() => applyMermaidZoom())
}

const scheduleMermaidRender = () => {
  if (typeof window === 'undefined') return
  if (mermaidTimer) window.clearTimeout(mermaidTimer)
  mermaidTimer = window.setTimeout(() => {
    nextTick(() => {
      const root = messageListRef.value
      if (!root) return
      const codeNodes = root.querySelectorAll('pre code')
      codeNodes.forEach(code => {
        const className = (code.getAttribute('class') || '').toLowerCase()
        const isMermaidLang =
          className.includes('language-mermaid') ||
          className.includes('lang-mermaid') ||
          className.includes('language-mermaidgraph') ||
          className.includes('lang-mermaidgraph')
        if (!isMermaidLang) return
        const pre = code.parentElement
        if (!pre || pre.tagName !== 'PRE') return
        const anyPre = pre as HTMLElement
        if (anyPre.dataset.mermaidDone === '1') return
        anyPre.dataset.mermaidDone = '1'
        const container = document.createElement('div')
        container.className = 'mermaid'
        let graphText = code.textContent || ''
        const trimmed = graphText.trimStart()
        if (trimmed.toLowerCase().startsWith('mermaidgraph')) {
          graphText = trimmed.slice('mermaidgraph'.length).trimStart()
        } else {
          graphText = trimmed
        }
        if (/^(TD|LR|RL|BT)\b/.test(graphText)) {
          graphText = `graph ${graphText}`
        }
        container.textContent = graphText
        pre.replaceWith(container)
      })

      if (!mermaidInitialized) return
      const nodes = root.querySelectorAll('.mermaid')
      if (nodes.length === 0) return
      try {
        mermaid.run({ nodes: Array.from(nodes) as any })
      } catch {
      }

      nodes.forEach(node => {
        const el = node as HTMLElement
        if (el.dataset.previewBound === '1') return
        el.dataset.previewBound = '1'
        el.style.cursor = 'zoom-in'
        el.addEventListener('click', (e) => {
          e.preventDefault()
          e.stopPropagation()
          openMermaidPreview(el.innerHTML)
        })
      })
    })
  }, 60)
}

const normalizeMarkdown = (content: string) => {
  if (!content) return ''
  return String(content)
}

const sanitizeMarkdown = (content: string) => {
  let text = String(content || '')
  if (!text) return ''

  text = text.replace(/\r\n/g, '\n')
  const lines = text.split('\n')
  const out: string[] = []
  let inFence = false

  const splitInlineMarkers = (line: string) => {
    let s = line
    s = s.replace(/([^\n])(```)/g, '$1\n\n```')
    s = s.replace(/([^\n])(#{1,6})(?=[^\s#])/g, '$1\n\n$2')
    s = s.replace(/([^\n])-(?=[\u4E00-\u9FFFA-Za-z0-9])/g, '$1\n-')
    return s
  }

  for (const line of lines) {
    const isFenceLine = line.trimStart().startsWith('```')
    if (isFenceLine) {
      inFence = !inFence
      out.push(line)
      continue
    }
    if (inFence) {
      out.push(line)
      continue
    }

    const expanded = splitInlineMarkers(line)
    if (expanded.includes('\n')) {
      expanded.split('\n').forEach(l => {
        out.push(l)
      })
      continue
    }

    if (/^\s*\d+\./.test(line)) {
      const segments: string[] = []
      const re = /(\d+)\.(?=[^\d\s])/g
      let lastIndex = 0
      let match: RegExpExecArray | null
      let seenFirst = false
      while ((match = re.exec(line)) !== null) {
        const index = match.index
        if (!seenFirst) {
          seenFirst = true
          continue
        }
        segments.push(line.slice(lastIndex, index).trimEnd())
        lastIndex = index
      }
      if (seenFirst && segments.length > 0) {
        segments.push(line.slice(lastIndex).trim())
        segments.forEach((seg, idx) => {
          const fixed = seg.replace(/^(\s*)(\d+)\.(?=[^\d\s])/, '$1$2. ')
          out.push(idx === 0 ? fixed : fixed.trimStart())
        })
        continue
      }
    }

    out.push(line)
  }

  const normalized = out
    .map(l => l.replace(/^(\s*[-*+])(\S)/, '$1 $2').replace(/^(\s*\d+)\.(\S)/, '$1. $2').replace(/^(#{1,6})(\S)/, '$1 $2'))
    .join('\n')

  return normalized
}

const formatMarkdown = (content: string) => {
  const input = sanitizeMarkdown(normalizeMarkdown(content || ''))
  if (!input) return ''
  try {
    const out = prettier.format(input, {
      parser: 'markdown',
      plugins: [prettierMarkdown],
      proseWrap: 'preserve',
      printWidth: 120
    })
    return typeof out === 'string' ? out : input
  } catch {
    return input
  }
}

const renderMarkdown = (content: string, streaming = false) => {
  const input = streaming ? normalizeMarkdown(content || '') : formatMarkdown(content || '')
  try {
    return String(mdProcessor.processSync(input))
  } catch {
    const esc: Record<string, string> = {
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      '"': '&quot;',
      "'": '&#39;'
    }
    return `<pre class="hljs"><code>${input.replace(/[&<>"']/g, (s: string) => esc[s] || s)}</code></pre>`
  }
}

// ========== Models ==========
const models = ref<any[]>([])
const selectedModelId = ref('')
const thinkingEnabled = ref(true)
const simplifiedOutput = ref(false)

// ========== Chat ==========
const messages = ref<any[]>([])
const inputText = ref('')
const isConnected = ref(false)
const typingQueues = new Map<string, Promise<void>>()
const collapsedAssistant = ref<Record<string, boolean>>({})

const isAssistantCollapsed = (baseId: string) => Boolean(collapsedAssistant.value[baseId])

const toggleAssistantCollapse = (baseId: string) => {
  collapsedAssistant.value = { ...collapsedAssistant.value, [baseId]: !isAssistantCollapsed(baseId) }
}

const messageGroups = computed(() => {
  const groups: any[] = []
  const assistantGroups = new Map<string, any>()

  for (const message of messages.value) {
    const baseId = typeof message.streamKey === 'string'
      ? message.streamKey.split(':')[0]
      : message.id

    if (message.type === 'agent' && !message.isError) {
      let group = assistantGroups.get(baseId)
      if (!group) {
        group = {
          id: `assistant:${baseId}`,
          baseId,
          kind: 'assistant',
          thinking: null,
          reply: null
        }
        assistantGroups.set(baseId, group)
        groups.push(group)
      }

      if (message.isThinking) {
        group.thinking = message
      } else {
        group.reply = message
      }
      continue
    }

    groups.push({
      id: message.id,
      kind: 'single',
      message
    })
  }

  return groups
})

// ========== WebSocket ==========
const ws = getWebSocketInstance()
const wsInitialized = ref(false)
const messageListRef = ref<HTMLElement | null>(null)
const thinkingViewports = new Map<string, HTMLElement>()
const shouldStickToBottom = ref(true)
const seenStreamEventKeys = new Map<string, number>()

const rememberStreamEvent = (key: string) => {
  const now = Date.now()
  seenStreamEventKeys.set(key, now)

  if (seenStreamEventKeys.size < 400) return
  const cutoff = now - 60_000
  for (const [k, ts] of seenStreamEventKeys) {
    if (ts < cutoff) seenStreamEventKeys.delete(k)
    if (seenStreamEventKeys.size < 300) break
  }
}

const handleWsConnectionChange = (connected: boolean) => {
  isConnected.value = connected
  console.log('[SinglePersonChat] WebSocket connection changed:', connected)
}

const STORAGE_KEY = 'agentflow.single_person_chat.v1'

type StoredMessage = {
  id: string
  type: 'user' | 'agent'
  content: string
  timestamp: number
  isThinking?: boolean
  isFinal?: boolean
  isStreaming?: boolean
  isError?: boolean
  streamKey?: string
}

type StoredConversation = {
  sessionId?: string
  messages: StoredMessage[]
  updatedAt: number
}

type StoredState = {
  users: Record<string, StoredConversation>
}

const safeParseJSON = (raw: string | null) => {
  if (!raw) return null
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

const readState = (): StoredState => {
  if (typeof window === 'undefined') return { users: {} }
  const parsed = safeParseJSON(window.localStorage.getItem(STORAGE_KEY))
  if (!parsed || typeof parsed !== 'object') return { users: {} }
  const users = (parsed as any).users
  if (!users || typeof users !== 'object') return { users: {} }
  return { users }
}

const writeState = (next: StoredState) => {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(next))
}

const sanitizeStoredMessages = (raw: any): StoredMessage[] => {
  if (!Array.isArray(raw)) return []
  return raw
    .map((m: any) => {
      if (!m || typeof m !== 'object') return null
      const id = typeof m.id === 'string' ? m.id : ''
      const type = m.type === 'user' ? 'user' : (m.type === 'agent' ? 'agent' : null)
      const content = typeof m.content === 'string' ? m.content : ''
      const timestamp = typeof m.timestamp === 'number' ? m.timestamp : Date.now()
      if (!id || !type) return null
      return {
        id,
        type,
        content,
        timestamp,
        isThinking: Boolean(m.isThinking),
        isFinal: Boolean(m.isFinal),
        isStreaming: false,
        isError: Boolean(m.isError),
        streamKey: typeof m.streamKey === 'string' ? m.streamKey : undefined
      } satisfies StoredMessage
    })
    .filter(Boolean) as StoredMessage[]
}

const getConversation = (userId: string): StoredConversation | null => {
  const state = readState()
  const conv = state.users[userId]
  if (!conv || typeof conv !== 'object') return null
  return {
    sessionId: typeof (conv as any).sessionId === 'string' ? (conv as any).sessionId : undefined,
    messages: sanitizeStoredMessages((conv as any).messages),
    updatedAt: typeof (conv as any).updatedAt === 'number' ? (conv as any).updatedAt : Date.now()
  }
}

const setConversation = (userId: string, conv: Partial<StoredConversation>) => {
  const state = readState()
  const prev = state.users[userId] || { messages: [], updatedAt: Date.now() }
  const next: StoredConversation = {
    sessionId: conv.sessionId ?? (prev as any).sessionId,
    messages: conv.messages ?? sanitizeStoredMessages((prev as any).messages),
    updatedAt: Date.now()
  }
  state.users[userId] = next
  writeState(state)
}

const createNewSessionId = () => {
  return `sess_${Date.now()}_${Math.random().toString(36).slice(2, 10)}`
}

let persistTimer: number | null = null
const schedulePersist = () => {
  if (!selectedUserId.value) return
  if (persistTimer) window.clearTimeout(persistTimer)
  persistTimer = window.setTimeout(() => {
    const userId = selectedUserId.value
    if (!userId) return
    const conv = getConversation(userId) || { messages: [], updatedAt: Date.now() }
    const sanitized: StoredMessage[] = messages.value
      .filter((m: any) => m && typeof m === 'object')
      .map((m: any) => {
        if (m.type === 'user') {
          return {
            id: String(m.id || ''),
            type: 'user',
            content: String(m.content || ''),
            timestamp: Number(m.timestamp || Date.now())
          } satisfies StoredMessage
        }
        if (m.type === 'agent') {
          if (m.isPlaceholder) return null
          return {
            id: String(m.id || ''),
            type: 'agent',
            content: typeof m.finalContent === 'string' ? m.finalContent : String(m.content || ''),
            timestamp: Number(m.timestamp || Date.now()),
            isThinking: Boolean(m.isThinking),
            isFinal: Boolean(m.isFinal),
            isStreaming: false,
            isError: Boolean(m.isError),
            streamKey: typeof m.streamKey === 'string' ? m.streamKey : undefined
          } satisfies StoredMessage
        }
        return null
      })
      .filter(Boolean) as StoredMessage[]

    setConversation(userId, {
      sessionId: conv.sessionId,
      messages: sanitized.slice(-300)
    })
  }, 250)
}

const handleWsMessage = (message: ServerMessage) => {
  console.log('[SinglePersonChat] Received message:', message)
  if (message.type === 'stream_chunk' || message.type === 'stream_end') {
    const payload: any = message.payload
    const content = typeof payload?.content === 'string' ? payload.content : ''
    const isThinking = Boolean(payload?.is_thinking)
    const isFinal = Boolean(payload?.is_final)
    const eventKey = `${message.type}|${message.message_id || ''}|${isThinking ? 1 : 0}|${isFinal ? 1 : 0}|${content}`
    if (seenStreamEventKeys.has(eventKey)) {
      return
    }
    rememberStreamEvent(eventKey)
    upsertAgentMessage(message)
    return
  }
  if (message.type === 'error') {
    console.error('[SinglePersonChat] WebSocket error payload:', message.payload)
    const payload = message.payload
    const errorText = payload?.message || '服务暂时不可用，请稍后重试'
    messages.value.push({
      id: `error:${message.message_id || Date.now()}`,
      type: 'agent',
      content: `错误：${errorText}`,
      isThinking: false,
      isFinal: true,
      isError: true,
      timestamp: message.timestamp || Date.now()
    })
    if (shouldStickToBottom.value) {
      scrollMessagesToBottom(true)
    }

    const maybeSessionError = typeof payload?.message === 'string' ? payload.message : ''
    const code = typeof payload?.code === 'string' ? payload.code : ''
    const isSessionInvalid = code === 'session_error' || /\bNot Found\b/i.test(maybeSessionError)
    if (isSessionInvalid && selectedUserId.value) {
      const newSessionId = createNewSessionId()
      ws.setSession(newSessionId)
      setConversation(selectedUserId.value, { sessionId: newSessionId })
    }
    schedulePersist()
  }
}

const trimTrailingBlankLines = (text: string) => text.replace(/\n{3,}$/g, '\n').replace(/\s+$/g, '')
const trimLeadingBlankLines = (text: string) => text.replace(/^\s*\n+/g, '')

const stripMemoryLines = (text: string) => {
  if (!text) return ''
  return text
    .split('\n')
    .filter(line => {
      const trimmed = line.trimStart()
      if (trimmed.startsWith('[MEMORY]')) return false
      if (trimmed.startsWith('[MEMORY')) return false
      return true
    })
    .join('\n')
}

const scrollMessagesToBottom = (smooth = false) => {
  nextTick(() => {
    const el = messageListRef.value
    if (!el) return
    el.scrollTo({
      top: el.scrollHeight,
      behavior: smooth ? 'smooth' : 'auto'
    })
  })
}

const handleMessageListScroll = () => {
  const el = messageListRef.value
  if (!el) return
  const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight
  shouldStickToBottom.value = distanceFromBottom < 80
}

const hasMeaningfulThinking = (message: any) => {
  if (!message) return false
  const content = typeof message.content === 'string' ? message.content.trim() : ''
  if (!content) return false
  return !/^(\.{3}|…+)$/.test(content)
}

const scrollThinkingViewport = (streamKey: string) => {
  nextTick(() => {
    const el = thinkingViewports.get(streamKey)
    if (!el) return
    el.scrollTop = el.scrollHeight
  })
  if (shouldStickToBottom.value) {
    scrollMessagesToBottom()
  }
}

const ensureTypingPlaceholder = (messageId: string) => {
  const answerKey = `${messageId || 'unknown'}:answer`
  const existing = messages.value.find((item: any) => item.type === 'agent' && item.streamKey === answerKey)
  const label = thinkingEnabled.value ? '正在推理中...' : '正在输入中...'
  if (existing) {
    existing.timestamp = Date.now()
    existing.isStreaming = true
    existing.isFinal = false
    existing.isPlaceholder = true
    existing.content = label
    return
  }

  messages.value.push({
    id: answerKey,
    streamKey: answerKey,
    type: 'agent',
    content: label,
    isThinking: false,
    isFinal: false,
    isStreaming: true,
    isPlaceholder: true,
    timestamp: Date.now()
  })
  if (shouldStickToBottom.value) {
    scrollMessagesToBottom()
  }
}

const typeOut = (messageObj: any, finalText: string) => {
  const text = String(finalText || '')
  messageObj.isPlaceholder = false
  messageObj.isStreaming = false
  messageObj.isFinal = true
  messageObj.isAnimating = true
  messageObj.displayContent = ''
  messageObj.finalContent = text
  messageObj.content = ''

  const total = text.length
  if (total === 0) {
    messageObj.isAnimating = false
    messageObj.content = ''
    return Promise.resolve()
  }

  const durationMs = Math.min(1800, Math.max(450, total * 2))
  const start = performance.now()

  return new Promise<void>((resolve) => {
    const step = (now: number) => {
      const t = Math.min(1, (now - start) / durationMs)
      const count = Math.min(total, Math.max(0, Math.floor(total * t)))
      messageObj.displayContent = text.slice(0, count)
      messageObj.timestamp = Date.now()
      if (shouldStickToBottom.value) {
        scrollMessagesToBottom()
      }
      if (t >= 1) {
        messageObj.isAnimating = false
        messageObj.content = text
        delete messageObj.displayContent
        scheduleMermaidRender()
        schedulePersist()
        resolve()
        return
      }
      requestAnimationFrame(step)
    }
    requestAnimationFrame(step)
  })
}

const upsertAgentMessage = (message: ServerMessage) => {
  const payload = message.payload
  if (!payload || typeof payload !== 'object') return

  const isThinking = Boolean(payload.is_thinking)

  if (message.type === 'stream_chunk') {
    if (isThinking) return
    ensureTypingPlaceholder(message.message_id || 'unknown')
    return
  }

  const answerKey = `${message.message_id || 'unknown'}:answer`
  const thinkingKey = `${message.message_id || 'unknown'}:thinking`
  const baseId = (message.message_id || 'unknown').toString()

  const finalize = async () => {
    const thinkingContent = typeof (payload as any).thinking_content === 'string' ? (payload as any).thinking_content : ''
    const finalContent = typeof (payload as any).content === 'string' ? (payload as any).content : ''

    let thinkingMessage = messages.value.find((item: any) => item.type === 'agent' && item.streamKey === thinkingKey)
    if (thinkingContent) {
      if (!thinkingMessage) {
        thinkingMessage = {
          id: thinkingKey,
          streamKey: thinkingKey,
          type: 'agent',
          content: '',
          isThinking: true,
          isFinal: true,
          isStreaming: false,
          timestamp: message.timestamp || Date.now()
        }
        messages.value.push(thinkingMessage)
      }
      thinkingMessage.isThinking = true
      thinkingMessage.isFinal = true
      thinkingMessage.isStreaming = false
    }

    let answerMessage = messages.value.find((item: any) => item.type === 'agent' && item.streamKey === answerKey && !item.isError)
    if (!answerMessage) {
      answerMessage = {
        id: answerKey,
        streamKey: answerKey,
        type: 'agent',
        content: '',
        isThinking: false,
        isFinal: true,
        isStreaming: false,
        timestamp: message.timestamp || Date.now()
      }
      messages.value.push(answerMessage)
    }

    if (shouldStickToBottom.value) {
      scrollMessagesToBottom(true)
    }

    const cleanThinking = trimTrailingBlankLines(trimLeadingBlankLines(stripMemoryLines(thinkingContent || '')))
    const cleanAnswer = trimTrailingBlankLines(trimLeadingBlankLines(stripMemoryLines(finalContent || '')))

    const placeholder = messages.value.find((item: any) => item.type === 'agent' && item.streamKey === answerKey && item.isPlaceholder)
    if (placeholder && placeholder !== answerMessage) {
      messages.value = messages.value.filter((m: any) => !(m.type === 'agent' && m.streamKey === answerKey && m.isPlaceholder))
    }

    try {
      if (thinkingMessage && thinkingContent) {
        await typeOut(thinkingMessage, cleanThinking)
        scrollThinkingViewport(thinkingKey)
      }
      await typeOut(answerMessage, cleanAnswer)
    } catch {
      if (thinkingMessage && thinkingContent) {
        thinkingMessage.isAnimating = false
        thinkingMessage.isStreaming = false
        thinkingMessage.isFinal = true
        thinkingMessage.content = cleanThinking
      }
      answerMessage.isAnimating = false
      answerMessage.isStreaming = false
      answerMessage.isFinal = true
      answerMessage.content = cleanAnswer
      scheduleMermaidRender()
      schedulePersist()
    }
  }

  const chain = (typingQueues.get(baseId) || Promise.resolve()).then(finalize).finally(() => {
    if (typingQueues.get(baseId) === chain) typingQueues.delete(baseId)
  })
  typingQueues.set(baseId, chain)
}

// Update agent settings in database
const updateAgentSettings = async () => {
  const userId = selectedUserId.value
  if (!userId) return
  
  const user = selectedUser.value
  if (!user) return
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/${userId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...user,
        thinking_enabled: thinkingEnabled.value,
        simplified_output: simplifiedOutput.value
      })
    })
    
    if (response.ok) {
      // Update local user data
      const userIndex = users.value.findIndex(u => u.id === userId)
      if (userIndex !== -1) {
        users.value[userIndex].thinking_enabled = thinkingEnabled.value
        users.value[userIndex].simplified_output = simplifiedOutput.value
      }
      console.log('数字员工设置已更新')
    } else {
      const error = await response.json()
      console.error('更新失败:', error)
    }
  } catch (error) {
    console.error('更新数字员工设置时发生错误:', error)
  }
}

// ========== WebSocket Functions ==========
const setupWebSocket = () => {
  if (wsInitialized.value) return
  wsInitialized.value = true
  console.log('[SinglePersonChat] Initializing WebSocket listeners', { url: WEBSOCKET_URL })
  
  // Set up event listeners if not already set up
  ws.on('connection-change', handleWsConnectionChange)
  ws.on('message', handleWsMessage)
  
  // Connect if not already connected
  if (!ws.getConnectionStatus()) {
    ws.connect()
  }
}

const clearChatHistory = () => {
  const userId = selectedUserId.value
  if (!userId) return
  messages.value = []
  const conv = getConversation(userId)
  setConversation(userId, { sessionId: conv?.sessionId, messages: [] })
  shouldStickToBottom.value = true
}

const sendMessage = () => {
  const text = inputText.value.trim()
  if (!text || !selectedUserId.value) return
  const userId = selectedUserId.value
  const existing = getConversation(userId)
  const sessionId = existing?.sessionId || createNewSessionId()
  if (!existing?.sessionId) {
    setConversation(userId, { sessionId })
  }
  ws.setSession(sessionId)
  console.log('[SinglePersonChat] sendMessage called', {
    selectedUserId: userId,
    isConnected: ws.getConnectionStatus(),
    sessionId,
    textLength: text.length
  })
  
  const userMessage = {
    id: Date.now().toString(),
    type: 'user',
    content: text,
    timestamp: Date.now()
  }
  
  messages.value.push(userMessage)
  shouldStickToBottom.value = true
  scrollMessagesToBottom(true)
  
  // Send via WebSocket
  const messageId = ws.sendText(text, userId, sessionId)
  if (messageId) {
    ensureTypingPlaceholder(messageId)
  }
  console.log('[SinglePersonChat] Message dispatched with ID:', messageId)
  inputText.value = ''
  schedulePersist()
}

const handleTextareaKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

const fetchModels = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/models/`)
    const data = await response.json()
    models.value = data
    const defaultModel = data.find((m: any) => m.is_default)
    if (defaultModel) selectedModelId.value = defaultModel.id
    else if (data.length > 0) selectedModelId.value = data[0].id
  } catch (error) {
    console.error('Failed to fetch models:', error)
  }
}

// ========== Users ==========
const users = ref<any[]>([])
const searchQuery = ref('')


// Fetch users sorted by ID descending (newest first)
const fetchUsers = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/`)
    if (response.ok) {
      users.value = await response.json()
      
      // Auto-select first user as default conversation user
      if (!selectedUserId.value && users.value.length > 0) {
        selectedUserId.value = users.value[0].id
      }
    } else {
      console.warn('Failed to fetch users')
      users.value = []
    }
  } catch (error) {
    console.error('Error fetching users:', error)
    users.value = []
  }
}

// Filtered users based on search query
const filteredUsers = computed(() => {
  if (!searchQuery.value.trim()) return users.value

  const query = searchQuery.value.toLowerCase().trim()
  return users.value.filter(user =>
    user.name.toLowerCase().includes(query)
  )
})

// Selected user for this conversation
const selectedUserId = ref('')
const selectedUser = computed(() => {
  if (!selectedUserId.value) return null
  return users.value.find(u => u.id === selectedUserId.value)
})
const activeUserName = computed(() => selectedUser.value?.name || '--')

// Users visible in collapsed sidebar mode (max 5, always include selected user)
const visibleUsersInCollapsedMode = computed(() => {
  const filtered = filteredUsers.value
  if (filtered.length <= 5) return filtered
  
  // Always include selected user if exists in filtered results
  const selectedUserObj = filtered.find(u => u.id === selectedUserId.value)
  const otherUsers = filtered.filter(u => u.id !== selectedUserId.value)
  
  // If selected user exists, include it plus 4 others
  if (selectedUserObj) {
    return [selectedUserObj, ...otherUsers.slice(0, 4)]
  } else {
    return otherUsers.slice(0, 5)
  }
})

// Sidebar state
const leftSidebarCollapsed = ref(false) // 控制左侧边栏是否收起，默认展开
const leftSidebarWidth = ref(320) // 左侧边栏宽度，可拖动调整
const isDragging = ref(false) // 是否正在拖动分隔线
const rightSidebarCollapsed = ref(true) // 控制右侧边栏是否收起，默认收起

// Toggle left sidebar
const toggleLeftSidebar = () => {
  leftSidebarCollapsed.value = !leftSidebarCollapsed.value
}

// Select user from collapsed sidebar
const selectUserFromCollapsed = (userId: string) => {
  selectedUserId.value = userId
}

// Drag handling for resizing sidebar
const startDrag = (event: MouseEvent) => {
  isDragging.value = true
  document.addEventListener('mousemove', handleDrag)
  document.addEventListener('mouseup', stopDrag)
  // Prevent text selection during drag
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'col-resize'
  event.preventDefault()
}

const handleDrag = (event: MouseEvent) => {
  if (!isDragging.value) return
  
  // Calculate new width based on mouse position
  const container = document.querySelector('.flex.w-full.h-full.overflow-hidden')
  if (container) {
    const containerRect = container.getBoundingClientRect()
    const mouseX = event.clientX - containerRect.left
    // Set minimum and maximum width constraints
    const minWidth = 200
    const maxWidth = containerRect.width - 400 // Leave space for right sidebar
    const newWidth = Math.max(minWidth, Math.min(mouseX, maxWidth))
    
    // Update left sidebar width
    leftSidebarWidth.value = newWidth
  }
}

const stopDrag = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', handleDrag)
  document.removeEventListener('mouseup', stopDrag)
  // Restore normal cursor and text selection
  document.body.style.userSelect = ''
  document.body.style.cursor = ''
}

// Watch for selected user changes to update settings
watch(selectedUserId, () => {
  const user = selectedUser.value
  if (user) {
    thinkingEnabled.value = user.thinking_enabled ?? true
    simplifiedOutput.value = user.simplified_output ?? false
    
    // Update WebSocket user
    ws.setUser(user.id)
    const conv = getConversation(user.id)
    if (conv?.sessionId) {
      ws.setSession(conv.sessionId)
    } else {
      ws.clearSession()
    }
    console.log('[SinglePersonChat] Selected user changed', {
      userId: user.id,
      userName: user.name
    })
    
    seenStreamEventKeys.clear()

    const restored = conv?.messages || []
    messages.value = restored.map((m: StoredMessage) => {
      if (m.type === 'user') {
        return {
          id: m.id,
          type: 'user',
          content: m.content,
          timestamp: m.timestamp
        }
      }
      return {
        id: m.id,
        streamKey: m.streamKey || m.id,
        type: 'agent',
        content: m.content,
        isThinking: Boolean(m.isThinking),
        isFinal: Boolean(m.isFinal),
        isStreaming: false,
        isError: Boolean(m.isError),
        timestamp: m.timestamp
      }
    })
    shouldStickToBottom.value = true
    nextTick(() => scrollMessagesToBottom())
    scheduleMermaidRender()
    schedulePersist()
  }
  
  // Ensure WebSocket is connected and set up
  setupWebSocket()
})

// Watch for settings changes to save to database
watch([thinkingEnabled, simplifiedOutput], () => {
  updateAgentSettings()
})

onMounted(() => {
  console.log('[SinglePersonChat] mounted', { websocketUrl: WEBSOCKET_URL })
  if (!mermaidInitialized) {
    mermaid.initialize({
      startOnLoad: false,
      theme: 'dark',
      securityLevel: 'strict'
    })
    mermaidInitialized = true
  }
  fetchModels()
  fetchUsers()
  setupWebSocket()
  nextTick(() => scrollMessagesToBottom())
  scheduleMermaidRender()
})

onBeforeUnmount(() => {
  thinkingViewports.clear()
  if (wsInitialized.value) {
    ws.off('connection-change', handleWsConnectionChange)
    ws.off('message', handleWsMessage)
  }
})
</script>

<template>
  <BaseLayout>
    <div class="flex w-full flex-1 min-h-0 overflow-hidden bg-[#0F1928]">
      <!-- Left Sidebar: Agent List -->
       <!-- Collapsed state: small sidebar with user avatars -->
       <div v-if="leftSidebarCollapsed" class="w-12 shrink-0 border-r border-white/10 bg-[#1A2536]/80 flex flex-col">
         <div class="p-4 border-b border-white/5 flex items-center justify-center">
           <button 
             @click="toggleLeftSidebar"
             class="p-2 rounded-lg hover:bg-white/10 text-white/60 hover:text-white transition-colors"
              title="展开 数字员工"
           >
             <iconify-icon icon="lucide:menu" class="text-xl"></iconify-icon>
           </button>
         </div>
         <div class="flex-1 overflow-y-auto py-2 custom-scrollbar">
           <div class="flex flex-col items-center gap-2">
              <button 
                v-for="user in visibleUsersInCollapsedMode" 
                :key="user.id"
                @click="selectUserFromCollapsed(user.id)"
                :title="user.name"
                :class="['w-8 h-8 rounded-full flex items-center justify-center text-xs border transition-all hover:scale-110', selectedUserId === user.id ? 'bg-[#3B9BFF] text-white border-[#3B9BFF]' : 'bg-white/10 text-white/70 border-transparent hover:bg-white/20']"
              >
                {{ user.name.charAt(0) }}
              </button>
           </div>
            <div v-if="filteredUsers.length > 5" class="text-center mt-2">
              <div class="text-[10px] text-white/30">+{{ filteredUsers.length - 5 }}</div>
           </div>
         </div>
       </div>
      
      <!-- Expanded state: full sidebar -->
      <div v-else :style="{ width: leftSidebarWidth + 'px' }" class="shrink-0 border-r border-white/10 bg-[#1A2536]/50 flex flex-col">
        <div class="p-6 border-b border-white/5 flex flex-col gap-4">
          <div class="flex justify-between items-center">
             <h2 class="text-lg font-semibold text-white/95">数字员工</h2>
            <button 
              @click="toggleLeftSidebar"
              class="p-2 rounded-full hover:bg-white/10 text-white/40 hover:text-white transition-colors"
               title="收起 数字员工"
            >
              <iconify-icon icon="lucide:chevron-left" class="text-xl"></iconify-icon>
            </button>
          </div>
          <div class="bg-black/20 border border-white/10 rounded-xl px-4 py-2.5 flex items-center gap-3">
            <iconify-icon icon="lucide:search" class="text-white/40 text-sm"></iconify-icon>
             <input v-model="searchQuery" type="text" placeholder="搜索数字员工名称..." class="bg-transparent border-none outline-none text-sm text-white/70 w-full">
          </div>
        </div>
        <div class="p-4 overflow-y-auto flex-1 custom-scrollbar">
          <div v-for="user in filteredUsers" :key="user.id" 
               @click="selectedUserId = user.id"
               :class="['p-4 rounded-xl cursor-pointer mb-2 transition-all', selectedUserId === user.id ? 'bg-[#3B9BFF]/20 border border-[#3B9BFF]/40' : 'bg-white/5 hover:bg-white/10']">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-white/10 flex items-center justify-center">
                <span class="text-sm text-white/70">{{ user.name[0] }}</span>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="text-sm font-medium text-white truncate">{{ user.name }}</h3>
                <p class="text-xs text-white/40 truncate">{{ user.description || 'No description' }}</p>
              </div>
            </div>
          </div>
          <div v-if="filteredUsers.length === 0" class="text-center py-8">
            <iconify-icon icon="lucide:users" class="w-8 h-8 mx-auto mb-2 text-white/20" />
            <div class="text-white/30 text-sm">{{ searchQuery ? '未找到匹配Agent' : '暂无Agent' }}</div>
          </div>
        </div>
        <!-- Draggable separator between left sidebar and chat area -->
        <div 
          class="w-1 cursor-col-resize hover:bg-[#3B9BFF]/50 bg-white/10 shrink-0"
          @mousedown="startDrag"
          :class="{ 'bg-[#3B9BFF]': isDragging }"
        ></div>
      </div>

      <!-- Main Chat Area -->
      <main class="flex-1 flex flex-col min-w-0 min-h-0 overflow-hidden">
        <div class="border-b border-white/10 bg-white/5 backdrop-blur-md shrink-0">
          <div class="h-16 flex items-center justify-between px-6">
            <div class="flex items-center gap-3">
              <div class="w-2 h-2 rounded-full" :class="isConnected ? 'bg-green-500' : 'bg-red-500'" :title="isConnected ? '已连接到服务器' : '未连接'"></div>
              <h2 class="text-lg font-semibold text-white/95">{{ activeUserName }}</h2>
            </div>
            <button
              class="inline-flex items-center gap-2 rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-xs text-white/70 hover:bg-white/10"
              @click="clearChatHistory"
            >
              <iconify-icon icon="lucide:trash-2" class="text-base"></iconify-icon>
              清空记录
            </button>
          </div>
        </div>
        <div
          ref="messageListRef"
          class="message-scroll flex-1 min-h-0 overflow-y-scroll p-6 space-y-6 custom-scrollbar overscroll-contain"
          @scroll="handleMessageListScroll"
        >
          <div v-if="messageGroups.length === 0" class="flex flex-col items-center justify-center h-full text-white/30 text-sm">
            <iconify-icon icon="lucide:user" class="text-4xl mb-4"></iconify-icon>
            开始你的单人对话...
          </div>
          <div v-else class="space-y-6">
            <div
              v-for="group in messageGroups"
              :key="group.id"
              class="flex"
              :class="group.kind === 'single' && group.message.type === 'user' ? 'justify-end' : 'justify-start'"
            >
                <div v-if="group.kind === 'assistant'" class="w-full max-w-[760px] space-y-3 text-[13px]">
                <!-- Original Message Container -->
                <div class="relative max-w-[760px] rounded-2xl p-4 border bg-white/10 text-white/90 border-white/10">
                  <div class="mb-3 flex items-center justify-between gap-2">
                    <div class="flex items-center gap-2">
                      <div class="w-8 h-8 rounded-lg bg-[#3B9BFF]/20 flex items-center justify-center text-[#3B9BFF] font-bold text-xs">
                        {{ selectedUser ? selectedUser.name.charAt(0) : 'AI' }}
                      </div>
                      <span class="text-xs font-semibold text-white/80">
                        {{ selectedUser ? selectedUser.name : 'AI 助手' }}
                      </span>
                    </div>
                    <button class="h-8 w-8 rounded-lg hover:bg-white/10" @click="toggleAssistantCollapse(group.baseId)">
                      <iconify-icon :icon="isAssistantCollapsed(group.baseId) ? 'lucide:chevrons-down' : 'lucide:chevrons-up'" class="text-base"></iconify-icon>
                    </button>
                  </div>
                  
                  <!-- Merged Thinking and Reply -->
                  <div v-if="isAssistantCollapsed(group.baseId)" class="rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-xs text-white/70">
                    内容已收拢
                  </div>
                  <template v-else>
                    <div v-if="group.thinking && (group.thinking.isAnimating || hasMeaningfulThinking(group.thinking))" class="mb-3 p-3 bg-amber-500/10 border border-amber-500/20 rounded-lg text-amber-100/70">
                      <div class="font-bold mb-1 flex items-center gap-2 text-xs"><iconify-icon icon="lucide:brain-circuit"></iconify-icon>推理过程</div>
                      <pre v-if="group.thinking.isAnimating" class="markdown-typing">{{ group.thinking.displayContent }}</pre>
                      <div v-else class="markdown-body break-words text-[13px]" v-html="renderMarkdown(group.thinking.content, false)"></div>
                    </div>
                    
                    <div>
                      <div v-if="group.reply?.isPlaceholder" class="mt-2 inline-flex items-center gap-2 text-xs opacity-70">
                        <span class="typing-dot"></span>
                        <span class="typing-dot"></span>
                        <span class="typing-dot"></span>
                        <span class="ml-1">{{ group.reply.content }}</span>
                      </div>
                      <pre v-else-if="group.reply?.isAnimating" class="markdown-typing">{{ group.reply.displayContent }}</pre>
                      <div v-else class="markdown-body break-words text-[13px]" v-html="renderMarkdown(group.reply ? group.reply.content : '', false)"></div>
                    </div>
                  </template>
                  <div class="text-xs mt-2 opacity-60">{{ new Date((group.reply || group.thinking).timestamp).toLocaleTimeString() }}</div>
                </div>
              </div>

              <div
                v-else-if="group.message.type === 'user'"
                class="flex flex-row-reverse items-start gap-4 w-full"
              >
                <!-- Avatar & Name -->
                <div class="flex flex-col items-center gap-1 shrink-0">
                  <div class="w-10 h-10 rounded-full bg-[#50C878]/20 flex items-center justify-center text-[#50C878]">
                    <iconify-icon icon="lucide:user" class="text-xl"></iconify-icon>
                  </div>
                  <span class="text-[10px] text-white/40">Human</span>
                </div>
                
                <!-- Original Message Container -->
                <div class="max-w-[760px] rounded-2xl p-4 border bg-[#3B9BFF] text-white border-[#3B9BFF] text-[13px]">
                  <div class="whitespace-pre-wrap break-words">{{ group.message.content }}</div>
                  <div class="text-xs mt-2 opacity-60 text-right">{{ new Date(group.message.timestamp).toLocaleTimeString() }}</div>
                </div>
              </div>

              <div
                v-else
                class="max-w-[760px] rounded-2xl p-4 border bg-red-500/10 text-red-50 border-red-400/30 text-[13px]"
              >
                <div class="mb-2 inline-flex items-center gap-2 rounded-full px-2.5 py-1 text-[11px] uppercase tracking-[0.18em] bg-red-400/15 text-red-200">
                  Error
                </div>
                <div class="whitespace-pre-wrap break-words">{{ group.message.content }}</div>
                <div class="text-xs mt-2 opacity-60">{{ new Date(group.message.timestamp).toLocaleTimeString() }}</div>
              </div>
            </div>
          </div>
        </div>
        <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md shrink-0">
          <div class="flex gap-4">
            <textarea v-model="inputText" @keydown="handleTextareaKeydown" placeholder="输入消息..." rows="1" class="chat-input flex-1 bg-white/5 border border-white/10 rounded-xl p-3 text-sm text-white/90 outline-none focus:border-[#3B9BFF]/50 resize-none overflow-y-auto max-h-32"></textarea>
            <button @click="sendMessage" class="w-12 h-12 bg-[#3B9BFF] rounded-xl flex items-center justify-center text-white shadow-[0_0_20px_rgba(59,155,255,0.3)]"><iconify-icon icon="lucide:send" class="text-xl"></iconify-icon></button>
          </div>
        </div>
      </main>

       <!-- Right Sidebar: Controls -->
       <!-- Collapsed state: small expand button -->
       <div v-if="rightSidebarCollapsed" class="w-12 shrink-0 border-l border-white/10 bg-[#1A2536]/80 flex flex-col items-center py-4">
         <button 
           @click="rightSidebarCollapsed = false"
           class="p-2 rounded-lg hover:bg-white/10 text-white/60 hover:text-white transition-colors"
           title="展开设置"
         >
           <iconify-icon icon="lucide:settings" class="text-xl"></iconify-icon>
         </button>
       </div>
       
       <!-- Expanded state: full sidebar -->
       <aside v-else class="w-[380px] shrink-0 border-l border-white/10 bg-[#1A2536] flex flex-col">
         <div class="p-6 border-b border-white/5 flex justify-between items-center">
           <h3 class="text-lg font-semibold text-white/95">会话设置</h3>
           <button 
             @click="rightSidebarCollapsed = true"
             class="p-2 rounded-full hover:bg-white/10 text-white/40 hover:text-white transition-colors"
             title="收起设置"
           >
             <iconify-icon icon="lucide:chevron-right" class="text-xl"></iconify-icon>
           </button>
         </div>
         
         <div class="p-6 space-y-8 overflow-y-auto custom-scrollbar flex-1">


           <!-- Toggles -->
           <div class="space-y-4">
             <div class="flex items-center justify-between p-3 bg-white/5 rounded-xl border border-white/10 cursor-pointer" @click="thinkingEnabled = !thinkingEnabled">
               <span class="text-sm text-white/90">开启 Thinking 模式</span>
               <div class="w-10 h-6 rounded-full p-1 transition-all" :class="thinkingEnabled ? 'bg-[#3B9BFF]' : 'bg-white/10'">
                 <div class="w-4 h-4 bg-white rounded-full transition-all" :class="thinkingEnabled ? 'translate-x-4' : 'translate-x-0'"></div>
               </div>
             </div>
             <div class="flex items-center justify-between p-3 bg-white/5 rounded-xl border border-white/10 cursor-pointer" @click="simplifiedOutput = !simplifiedOutput">
               <span class="text-sm text-white/90">简化输出模式</span>
               <div class="w-10 h-6 rounded-full p-1 transition-all" :class="simplifiedOutput ? 'bg-[#3B9BFF]' : 'bg-white/10'">
                 <div class="w-4 h-4 bg-white rounded-full transition-all" :class="simplifiedOutput ? 'translate-x-4' : 'translate-x-0'"></div>
               </div>
             </div>
           </div>
         </div>
       </aside>
    </div>
    <div
      v-if="mermaidPreviewOpen"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-6"
      @click="closeMermaidPreview"
    >
      <div ref="mermaidPreviewRef" class="relative w-full max-w-6xl max-h-[85vh] overflow-auto rounded-2xl border border-white/10 bg-[#0F1928] p-4" @click.stop @wheel="handleMermaidWheel">
        <button
          class="absolute right-3 top-3 inline-flex h-9 w-9 items-center justify-center rounded-full border border-white/10 bg-white/5 text-white/70 hover:bg-white/10"
          @click="closeMermaidPreview"
        >
          <iconify-icon icon="lucide:x" class="text-xl"></iconify-icon>
        </button>
        <div class="absolute left-3 top-3 inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-2 py-1.5 text-xs text-white/70">
          <button class="h-8 w-8 rounded-full hover:bg-white/10" @click="zoomMermaidOut">-</button>
          <button class="h-8 w-8 rounded-full hover:bg-white/10" @click="zoomMermaidIn">+</button>
          <button class="h-8 rounded-full px-2 hover:bg-white/10" @click="resetMermaidZoom">{{ Math.round(mermaidZoom * 100) }}%</button>
        </div>
        <div class="markdown-body">
          <div class="mermaid" v-html="mermaidPreviewHtml"></div>
        </div>
      </div>
    </div>
  </BaseLayout>
</template>

<style scoped>
:deep(main) {
  padding: 0 !important;
}
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
.message-scroll {
  overscroll-behavior: contain;
  scrollbar-gutter: stable;
  scrollbar-width: thin;
  scrollbar-color: rgba(59, 155, 255, 0.35) rgba(255, 255, 255, 0.06);
}
.message-scroll::-webkit-scrollbar {
  width: 8px;
}
.message-scroll::-webkit-scrollbar-thumb {
  background: rgba(59, 155, 255, 0.28);
  border-radius: 9999px;
}
.message-scroll::-webkit-scrollbar-thumb:hover {
  background: rgba(95, 180, 255, 0.45);
}
.message-scroll::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 9999px;
}
.chat-input {
  scrollbar-width: thin;
  scrollbar-color: rgba(59, 155, 255, 0.35) rgba(255, 255, 255, 0.06);
}
.chat-input::-webkit-scrollbar {
  width: 6px;
}
.chat-input::-webkit-scrollbar-thumb {
  background: rgba(59, 155, 255, 0.35);
  border-radius: 9999px;
}
.chat-input::-webkit-scrollbar-thumb:hover {
  background: rgba(95, 180, 255, 0.5);
}
.chat-input::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 9999px;
}
.thinking-viewport {
  overflow-y: auto;
  padding-right: 0.25rem;
  transition: max-height 180ms ease;
}
.thinking-viewport-collapsed {
  max-height: 6.5em;
  overflow-y: hidden;
  mask-image: linear-gradient(to bottom, black 62%, transparent 100%);
}
.thinking-viewport-expanded {
  max-height: 18em;
  mask-image: none;
}
.thinking-viewport::-webkit-scrollbar {
  width: 4px;
}
.thinking-viewport::-webkit-scrollbar-thumb {
  background: rgba(251, 191, 36, 0.28);
  border-radius: 9999px;
}

:deep(.markdown-body) {
  line-height: 1.8;
  word-break: break-word;
  font-size: inherit;
}

:deep(.markdown-body p) {
  margin: 0.6em 0;
}

:deep(.markdown-body h1),
:deep(.markdown-body h2),
:deep(.markdown-body h3),
:deep(.markdown-body h4),
:deep(.markdown-body h5),
:deep(.markdown-body h6) {
  margin: 0.9em 0 0.45em;
  font-weight: 650;
}

:deep(.markdown-body h1) { font-size: 1.35rem; }
:deep(.markdown-body h2) { font-size: 1.2rem; }
:deep(.markdown-body h3) { font-size: 1.1rem; }
:deep(.markdown-body h4) { font-size: 1.02rem; }
:deep(.markdown-body h5) { font-size: 0.98rem; }
:deep(.markdown-body h6) { font-size: 0.95rem; }

:deep(.markdown-body ul),
:deep(.markdown-body ol) {
  margin: 0.5em 0 0.5em 1.1em;
  padding-left: 1.1em;
  list-style-position: outside;
}

:deep(.markdown-body ul) { list-style-type: disc; }
:deep(.markdown-body ol) { list-style-type: decimal; }

:deep(.markdown-body li) {
  margin: 0.25em 0;
}

:deep(.markdown-body blockquote) {
  margin: 0.8em 0;
  padding: 0.5em 0.9em;
  border-left: 3px solid rgba(59, 155, 255, 0.55);
  background: rgba(255, 255, 255, 0.06);
  border-radius: 0.75rem;
}

:deep(.markdown-body code) {
  padding: 0.15em 0.35em;
  border-radius: 0.4rem;
  background: rgba(255, 255, 255, 0.08);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.9em;
}

:deep(.markdown-body pre) {
  margin: 0.8em 0;
  padding: 0.9em 1em;
  border-radius: 0.9rem;
  background: rgba(0, 0, 0, 0.35);
  overflow: auto;
}

:deep(.markdown-body .mermaid) {
  margin: 0.8em 0;
  padding: 0.9em 1em;
  border-radius: 0.9rem;
  background: rgba(0, 0, 0, 0.25);
  overflow: auto;
}

:deep(.markdown-body .mermaid svg) {
  max-width: 100%;
  height: auto;
}

:deep(.markdown-body pre code) {
  padding: 0;
  background: transparent;
  font-size: 0.9em;
}

.markdown-typing {
  margin: 0.25em 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.92em;
  line-height: 1.7;
  color: rgba(255, 255, 255, 0.9);
}

.typing-dot {
  width: 6px;
  height: 6px;
  border-radius: 9999px;
  background: rgba(255, 255, 255, 0.6);
  display: inline-block;
  animation: typingBounce 1s infinite ease-in-out;
}

.typing-dot:nth-child(1) { animation-delay: 0ms; }
.typing-dot:nth-child(2) { animation-delay: 120ms; }
.typing-dot:nth-child(3) { animation-delay: 240ms; }

@keyframes typingBounce {
  0%, 80%, 100% { transform: translateY(0); opacity: 0.4; }
  40% { transform: translateY(-3px); opacity: 1; }
}

:deep(.markdown-body a) {
  color: rgba(110, 200, 255, 1);
  text-decoration: underline;
  text-underline-offset: 2px;
}

:deep(.markdown-body hr) {
  border: none;
  border-top: 1px solid rgba(255, 255, 255, 0.12);
  margin: 1em 0;
}

</style>
