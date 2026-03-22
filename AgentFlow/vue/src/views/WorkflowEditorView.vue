<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL } from '../config/api'

const route = useRoute()
const router = useRouter()
const workflowId = route.params.id as string
const isNew = workflowId === 'new'

interface Node {
  id: string
  name: string
  type: 'agent' | 'start' | 'parallel' | 'merge'
  x: number
  y: number
  icon: string
  color: string
  metadata?: any 
}

interface Connection {
  id: string
  fromId: string
  toId: string
}

// State
const workflowName = ref('未命名编排')
const workflowDescription = ref('')
const isSaveModalOpen = ref(false)
const isSaving = ref(false)

const nodes = ref<Node[]>([
  { id: 'start', name: '开始', type: 'start', x: 50, y: 200, icon: 'lucide:play', color: '#50C878' },
])
const connections = ref<Connection[]>([])
const selectedNodeId = ref<string | null>(null)
const isDraggingNode = ref<string | null>(null)
const dragOffset = { x: 0, y: 0 }
const canvasRef = ref<HTMLElement | null>(null)
const contextMenu = reactive({ show: false, x: 0, y: 0, type: '' as 'node' | 'conn', targetId: '' })
const activeDrawing = ref<{ fromId: string; startX: number; startY: number } | null>(null)
const mousePos = reactive({ x: 0, y: 0 })

// Digital Employees state
const digitalEmployees = ref<any[]>([])
const selectedNodeData = computed(() => {
  if (!selectedNodeId.value) return null
  const node = nodes.value.find(n => n.id === selectedNodeId.value)
  if (!node) return null
  if (node.type === 'agent' && node.metadata?.agentId) {
    return digitalEmployees.value.find(a => a.id === node.metadata.agentId)
  }
  return null
})

const fetchDigitalEmployees = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/`)
    digitalEmployees.value = await response.json()
  } catch (error) { console.error(error) }
}

const fetchWorkflow = async () => {
  if (isNew) return
  try {
    const response = await fetch(`${API_BASE_URL}/api/workflows/${workflowId}`)
    if (response.ok) {
      const data = await response.json()
      workflowName.value = data.name
      workflowDescription.value = data.description
      const graph = typeof data.graph === 'string' ? JSON.parse(data.graph) : data.graph
      nodes.value = graph.nodes || []
      connections.value = graph.connections || []
    }
  } catch (error) { console.error(error) }
}

const executeSave = async () => {
  isSaving.value = true
  const payload = {
    name: workflowName.value,
    description: workflowDescription.value,
    graph: {
      nodes: nodes.value,
      connections: connections.value
    }
  }
  
  try {
    const method = isNew ? 'POST' : 'PUT'
    const url = isNew ? `${API_BASE_URL}/api/workflows/` : `${API_BASE_URL}/api/workflows/${workflowId}`
    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    
    if (response.ok) {
      if (isNew) {
        const result = await response.json()
        router.replace(`/orchestration-editor/${result.id}`)
      }
      isSaveModalOpen.value = false
    }
  } catch (error) {
    console.error(error)
    alert('保存失败')
  } finally {
    isSaving.value = false
  }
}

// Canvas Helpers
const NODE_WIDTH = 220
const NODE_HEIGHT = 80 

const onSidebarDragStart = (event: DragEvent, nodeData: any) => {
  if (event.dataTransfer) {
    event.dataTransfer.setData('application/json', JSON.stringify(nodeData))
    event.dataTransfer.effectAllowed = 'copy'
  }
}

const onCanvasDrop = (event: DragEvent) => {
  event.preventDefault()
  if (!canvasRef.value) return
  const rect = canvasRef.value.getBoundingClientRect()
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top
  const data = event.dataTransfer?.getData('application/json')
  if (data) {
    const info = JSON.parse(data)
    nodes.value.push({
      id: `${info.type}-${Date.now()}`,
      name: info.name,
      type: info.type,
      x: x - NODE_WIDTH / 2,
      y: y - NODE_HEIGHT / 2,
      icon: info.icon || 'lucide:box',
      color: info.color || '#3B9BFF',
      metadata: info.agentId ? { agentId: info.agentId } : {}
    })
  }
}

const getCurvePath = (x1: number, y1: number, x2: number, y2: number) => {
  const dx = Math.abs(x2 - x1)
  let cx = Math.max(dx / 2, 50)
  if (x2 < x1) cx = Math.max(dx, 150)
  return `M ${x1} ${y1} C ${x1 + cx} ${y1}, ${x2 - cx} ${y2}, ${x2} ${y2}`
}

const getConnectionPath = (conn: Connection) => {
  const from = nodes.value.find(n => n.id === conn.fromId)
  const to = nodes.value.find(n => n.id === conn.toId)
  if (!from || !to) return ''
  return getCurvePath(from.x + NODE_WIDTH, from.y + NODE_HEIGHT / 2, to.x, to.y + NODE_HEIGHT / 2)
}

const startDragNode = (event: MouseEvent, nodeId: string) => {
  if (event.button !== 0) return 
  event.stopPropagation()
  selectedNodeId.value = nodeId
  isDraggingNode.value = nodeId
  const node = nodes.value.find(n => n.id === nodeId)
  if (node) {
    dragOffset.x = event.clientX - node.x
    dragOffset.y = event.clientY - node.y
  }
  contextMenu.show = false
}

const startConnecting = (event: MouseEvent, nodeId: string) => {
  event.stopPropagation()
  const node = nodes.value.find(n => n.id === nodeId)
  if (node) {
    activeDrawing.value = { fromId: nodeId, startX: node.x + NODE_WIDTH, startY: node.y + NODE_HEIGHT / 2 }
  }
}

const finishConnecting = (toNodeId: string) => {
  if (activeDrawing.value && activeDrawing.value.fromId !== toNodeId) {
    connections.value.push({ id: `conn-${Date.now()}`, fromId: activeDrawing.value.fromId, toId: toNodeId })
  }
  activeDrawing.value = null;
}

const openContextMenu = (event: MouseEvent, type: 'node' | 'conn', id: string) => {
  event.preventDefault()
  event.stopPropagation()
  contextMenu.show = true
  contextMenu.x = event.clientX
  contextMenu.y = event.clientY
  contextMenu.type = type
  contextMenu.targetId = id
}

const handleDelete = () => {
  if (contextMenu.type === 'node') {
    if (contextMenu.targetId === selectedNodeId.value) selectedNodeId.value = null
    nodes.value = nodes.value.filter(n => n.id !== contextMenu.targetId)
    connections.value = connections.value.filter(c => c.fromId !== contextMenu.targetId && c.toId !== contextMenu.targetId)
  } else if (contextMenu.type === 'conn') {
    connections.value = connections.value.filter(c => c.id !== contextMenu.targetId)
  }
  contextMenu.show = false
}

const onGlobalMouseMove = (event: MouseEvent) => {
  if (isDraggingNode.value) {
    const node = nodes.value.find(n => n.id === isDraggingNode.value)
    if (node) { node.x = event.clientX - dragOffset.x; node.y = event.clientY - dragOffset.y }
  }
  if (activeDrawing.value && canvasRef.value) {
    const rect = canvasRef.value.getBoundingClientRect()
    mousePos.x = event.clientX - rect.left; mousePos.y = event.clientY - rect.top
  }
}

const onGlobalMouseUp = (event: MouseEvent) => {
  isDraggingNode.value = null;
  if (activeDrawing.value && !(event.target as HTMLElement).closest('.node-input-port')) {
      activeDrawing.value = null;
  }
}

onMounted(() => {
  fetchDigitalEmployees()
  fetchWorkflow()
  window.addEventListener('mousemove', onGlobalMouseMove)
  window.addEventListener('mouseup', onGlobalMouseUp)
  window.addEventListener('click', () => { if (!isDraggingNode.value) contextMenu.show = false })
})
</script>

<template>
  <BaseLayout>
    <div class="flex w-full h-full overflow-hidden relative select-none">
      <!-- Sidebar Components -->
      <aside class="shrink-0 border-r border-white/10 flex flex-col w-[280px]" style="background: rgba(20, 30, 50, 0.8);">
        <div class="p-6 border-b border-white/5 flex justify-between items-center">
          <h3 class="text-lg font-semibold text-white/95">组件库</h3>
          <button @click="router.push('/orchestration-list')" class="text-xs text-white/40 hover:text-white flex items-center gap-1 transition-colors">
            <iconify-icon icon="lucide:arrow-left" class="text-[10px]"></iconify-icon>列表
          </button>
        </div>
        <div class="p-6 space-y-8 overflow-y-auto custom-scrollbar flex-1">
          <div class="space-y-4">
            <h4 class="text-[10px] font-bold text-white/30 uppercase tracking-[0.2em]">数字员工</h4>
            <div v-for="agent in digitalEmployees" :key="agent.id" 
              draggable="true" 
              @dragstart="onSidebarDragStart($event, {name: agent.name, type: 'agent', color: '#3B9BFF', icon: 'lucide:user', agentId: agent.id})" 
              class="p-3 bg-white/5 border border-white/10 rounded-xl cursor-grab hover:bg-white/10 transition-colors group"
            >
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg flex items-center justify-center bg-[#3B9BFF]/20 text-[#3B9BFF]"><iconify-icon icon="lucide:user"></iconify-icon></div>
                <span class="text-sm text-white/80 truncate">{{ agent.name }}</span>
              </div>
            </div>
          </div>
          <div class="space-y-4">
            <h4 class="text-[10px] font-bold text-white/30 uppercase tracking-[0.2em]">逻辑控制</h4>
            <div v-for="item in [{name:'并行执行', icon:'lucide:shuffle', type:'parallel', color:'#50C878'}, {name:'合并分支', icon:'lucide:combine', type:'merge', color:'#A78BFA'}]" :key="item.name" draggable="true" @dragstart="onSidebarDragStart($event, item)" class="p-3 bg-white/5 border border-white/10 rounded-xl cursor-grab hover:bg-white/10 transition-colors">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg flex items-center justify-center shadow-inner" :style="{backgroundColor: item.color+'20', color: item.color}"><iconify-icon :icon="item.icon"></iconify-icon></div>
                <span class="text-sm text-white/80">{{ item.name }}</span>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <!-- Main Canvas -->
      <main class="flex-1 relative bg-[#0F1928] overflow-hidden flex flex-col">
        <!-- Top Toolbar -->
        <div class="h-14 border-b border-white/5 bg-white/5 flex items-center justify-between px-6 shrink-0 z-40 backdrop-blur-md">
          <div class="flex items-center gap-2">
            <span class="text-white/40 text-xs font-mono uppercase tracking-widest">{{ isNew ? 'NEW' : 'EDIT' }}</span>
            <div class="h-4 w-px bg-white/10 mx-2"></div>
            <h2 class="text-white/90 font-semibold text-sm">{{ workflowName }}</h2>
          </div>
          <button @click="isSaveModalOpen = true" class="px-4 py-1.5 bg-[#3B9BFF] text-white text-xs font-bold rounded-lg shadow-lg hover:bg-[#2A7FDB] transition-all flex items-center gap-2">
            <iconify-icon icon="lucide:save"></iconify-icon>保存编排
          </button>
        </div>

        <div ref="canvasRef" @dragover.prevent @drop="onCanvasDrop" class="flex-1 relative overflow-hidden" style="background-image: radial-gradient(circle, rgba(255, 255, 255, 0.03) 1px, transparent 1px); background-size: 30px 30px;" @mousedown="selectedNodeId = null">
          <!-- Nodes -->
          <div v-for="node in nodes" :key="node.id" 
            class="absolute p-4 bg-[#1A2536]/90 backdrop-blur-xl border rounded-xl shadow-2xl transition-all group" 
            :style="{ left: node.x+'px', top: node.y+'px', width: NODE_WIDTH+'px', height: NODE_HEIGHT+'px', borderColor: selectedNodeId === node.id ? '#3B9BFF' : 'rgba(255,255,255,0.1)', zIndex: 10 }" 
            @mousedown.stop="startDragNode($event, node.id)" 
            @contextmenu="openContextMenu($event, 'node', node.id)"
          >
            <div class="absolute -left-6 top-0 bottom-0 w-12 z-50 cursor-crosshair node-input-port flex items-center justify-center" @mouseup.stop="finishConnecting(node.id)">
               <div class="w-4 h-4 bg-[#0F1928] border-2 border-white/20 rounded-full group-hover:border-[#3B9BFF] transition-all"></div>
            </div>
            <div class="absolute -right-2 top-1/2 -translate-y-1/2 w-4 h-4 bg-[#3B9BFF] border-2 border-white rounded-full z-20 cursor-crosshair hover:scale-125 transition-all" @mousedown.stop="startConnecting($event, node.id)"></div>
            <div class="flex items-center gap-3 h-full">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center text-lg shrink-0 shadow-inner" :style="{backgroundColor: node.color+'20', color: node.color}"><iconify-icon :icon="node.icon"></iconify-icon></div>
              <div class="flex-1 overflow-hidden text-left">
                <div class="text-sm font-bold text-white/90 truncate">{{ node.name }}</div>
                <div class="text-[9px] text-white/30 uppercase tracking-tighter">{{ node.type === 'agent' ? '数字员工' : node.type }}</div>
              </div>
            </div>
          </div>

          <!-- Connections -->
          <svg class="absolute inset-0 w-full h-full pointer-events-none z-[30]">
            <defs><marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto"><polygon points="0 0, 10 3.5, 0 7" fill="#3B9BFF" /></marker><filter id="glow"><feGaussianBlur stdDeviation="2" result="coloredBlur"/><feMerge><feMergeNode in="coloredBlur"/><feMergeNode in="SourceGraphic"/></feMerge></filter></defs>
            <g v-for="conn in connections" :key="conn.id" class="pointer-events-auto cursor-pointer group" @contextmenu="openContextMenu($event, 'conn', conn.id)">
              <path :d="getConnectionPath(conn)" stroke="rgba(59, 155, 255, 0.01)" stroke-width="20" fill="none" />
              <path :d="getConnectionPath(conn)" stroke="#3B9BFF" stroke-width="2.5" fill="none" marker-end="url(#arrowhead)" filter="url(#glow)" class="group-hover:stroke-white transition-colors" />
            </g>
            <path v-if="activeDrawing" :d="getCurvePath(activeDrawing.startX, activeDrawing.startY, mousePos.x, mousePos.y)" stroke="#3B9BFF" stroke-width="2" stroke-dasharray="6,6" fill="none" />
          </svg>
        </div>

        <!-- Save Modal (Matched style) -->
        <Transition enter-active-class="transition duration-300" enter-from-class="opacity-0 scale-95" enter-to-class="opacity-100 scale-100" leave-active-class="transition duration-200" leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95">
          <div v-if="isSaveModalOpen" class="fixed inset-0 z-[120] flex items-center justify-center p-6 bg-black/60 backdrop-blur-sm">
            <div class="relative w-full max-w-lg bg-[#1A2536] border border-white/10 rounded-3xl shadow-2xl overflow-hidden flex flex-col">
              <div class="p-6 border-b border-white/10 flex justify-between items-center bg-white/5">
                <h2 class="text-xl font-bold text-white/95">保存编排流程</h2>
                <button @click="isSaveModalOpen = false" class="w-8 h-8 rounded-full flex items-center justify-center hover:bg-white/10 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
              </div>
              <div class="p-8 space-y-6">
                <div class="space-y-2">
                  <label class="text-[10px] text-white/40 font-bold uppercase tracking-widest">编排名称</label>
                  <input v-model="workflowName" placeholder="给这个编排起个名字" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 h-[50px] text-sm text-white outline-none focus:border-[#3B9BFF]/50 transition-all">
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] text-white/40 font-bold uppercase tracking-widest">描述信息</label>
                  <textarea v-model="workflowDescription" placeholder="简要描述这个流程的作用..." rows="3" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-sm text-white outline-none focus:border-[#3B9BFF]/50 transition-all resize-none"></textarea>
                </div>
              </div>
              <div class="p-6 border-t border-white/10 bg-white/5 flex gap-3">
                <button @click="isSaveModalOpen = false" class="flex-1 py-3 rounded-xl border border-white/10 text-sm text-white/60 hover:bg-white/5 transition-all">取消</button>
                <button @click="executeSave" :disabled="isSaving" class="flex-1 py-3 rounded-xl bg-[#3B9BFF] text-white text-sm font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)] disabled:opacity-50 transition-all">
                  {{ isSaving ? '保存中...' : '确认保存' }}
                </button>
              </div>
            </div>
          </div>
        </Transition>

        <div v-if="contextMenu.show" class="fixed z-[100] bg-[#1A2536] border border-white/10 rounded-lg shadow-2xl py-1 min-w-[120px]" :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }" @click.stop>
          <button @click="handleDelete" class="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-red-500/10 flex items-center gap-3"><iconify-icon icon="lucide:trash-2"></iconify-icon>删除</button>
        </div>
      </main>

      <!-- Right Config Panel -->
      <aside class="shrink-0 border-l border-white/10 flex flex-col w-[320px]" style="background: rgba(20, 30, 50, 0.8);">
        <div class="p-6 border-b border-white/5"><h3 class="text-lg font-semibold text-white/95">配置面板</h3></div>
        
        <div v-if="selectedNodeId" class="p-6 space-y-6 overflow-y-auto custom-scrollbar flex-1 text-left">
          <div v-if="nodes.find(n => n.id === selectedNodeId)?.type === 'start'">
            <div class="space-y-4">
              <div class="text-xs text-[#50C878] font-bold uppercase tracking-widest">流程入口：开始</div>
              <div class="bg-white/5 p-4 rounded-xl border border-white/10 text-xs text-white/60 leading-relaxed space-y-2">
                <p>🚀 <span class="text-white">运行逻辑：</span></p>
                <p>这是整个编排工作流的<span class="text-[#50C878]">唯一触发起点</span>。</p>
                <p>当外部调用该编排时，系统会首先激活此节点，并将初始输入数据传递给后续连接的节点。</p>
              </div>
            </div>
          </div>

          <div v-else-if="nodes.find(n => n.id === selectedNodeId)?.type === 'agent'">
            <div class="space-y-4">
              <div class="text-xs text-[#3B9BFF] font-bold uppercase tracking-widest">数字员工属性</div>
              <div v-if="selectedNodeData" class="space-y-4">
                <div class="bg-white/5 p-4 rounded-xl border border-white/10">
                  <div class="text-[10px] text-white/30 uppercase mb-1 font-bold">名称</div>
                  <div class="text-sm text-white/90 font-semibold">{{ selectedNodeData.name }}</div>
                </div>
                <div class="bg-white/5 p-4 rounded-xl border border-white/10">
                  <div class="text-[10px] text-white/30 uppercase mb-1 font-bold">职责描述</div>
                  <div class="text-xs text-white/60 leading-relaxed">{{ selectedNodeData.description || '暂无描述' }}</div>
                </div>
                <div class="bg-white/5 p-4 rounded-xl border border-white/10">
                  <div class="text-[10px] text-white/30 uppercase mb-1 font-bold">底层模型</div>
                  <div class="text-xs text-white/60 font-mono">{{ selectedNodeData.model }}</div>
                </div>
              </div>
              <div v-else class="text-xs text-white/30 italic">未找到关联的员工数据，请重新拖入。</div>
            </div>
          </div>

          <div v-else-if="nodes.find(n => n.id === selectedNodeId)?.type === 'parallel'">
            <div class="space-y-4">
              <div class="text-xs text-[#50C878] font-bold uppercase tracking-widest">逻辑：并行执行</div>
              <div class="bg-white/5 p-4 rounded-xl border border-white/10 text-xs text-white/60 leading-relaxed space-y-2">
                <p>💡 <span class="text-white">运行逻辑：</span></p>
                <p>当上一个节点产生多条结果（如列表或数组）时，此节点会将结果拆分。</p>
                <p>每条结果将<span class="text-[#50C878]">分别且并行</span>地发送到下游的所有连接节点中独立处理。</p>
              </div>
            </div>
          </div>

          <div v-else-if="nodes.find(n => n.id === selectedNodeId)?.type === 'merge'">
            <div class="space-y-4">
              <div class="text-xs text-[#A78BFA] font-bold uppercase tracking-widest">逻辑：合并分支</div>
              <div class="bg-white/5 p-4 rounded-xl border border-white/10 text-xs text-white/60 leading-relaxed space-y-2">
                <p>💡 <span class="text-white">运行逻辑：</span></p>
                <p>这是一个同步等待点。</p>
                <p>系统会监控所有指向此节点的连线。只有当<span class="text-[#A78BFA]">所有上游节点</span>都达到完成状态后，此节点才会触发并汇总数据传向下一步。</p>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="flex-1 p-6 flex flex-col items-center justify-center text-white/20 text-center">
          <iconify-icon icon="lucide:settings-2" class="text-5xl mb-4"></iconify-icon>
          <p class="text-sm px-4 leading-relaxed">点击画布上的节点进行详细配置</p>
        </div>
      </aside>
    </div>
  </BaseLayout>
</template>

<style scoped>
:deep(main) { padding: 0 !important; }
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
</style>
