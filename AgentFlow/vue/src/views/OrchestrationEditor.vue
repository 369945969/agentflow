<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'

interface Node {
  id: string
  name: string
  type: 'agent' | 'start' | 'branch' | 'parallel' | 'merge'
  x: number
  y: number
  icon: string
  color: string
}

interface Connection {
  id: string
  fromId: string
  toId: string
}

const nodes = ref<Node[]>([
  { id: 'start', name: '开始', type: 'start', x: 50, y: 200, icon: 'lucide:play', color: '#50C878' },
])

const connections = ref<Connection[]>([])
const isDraggingNode = ref<string | null>(null)
const dragOffset = { x: 0, y: 0 }
const canvasRef = ref<HTMLElement | null>(null)

const contextMenu = reactive({ show: false, x: 0, y: 0, type: '' as 'node' | 'conn', targetId: '' })
const activeDrawing = ref<{ fromId: string; startX: number; startY: number } | null>(null)
const mousePos = reactive({ x: 0, y: 0 })

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
      icon: info.icon,
      color: info.color
    })
  }
}

const getCurvePath = (x1: number, y1: number, x2: number, y2: number, isReverse: boolean = false) => {
  const dx = Math.abs(x2 - x1)
  let cx = Math.max(dx / 2, 50)
  if (x2 < x1) cx = Math.max(dx, 150)
  const curveOffset = isReverse ? -40 : 0
  const cp1x = x1 + cx, cp1y = y1 + curveOffset
  const cp2x = x2 - cx, cp2y = y2 + curveOffset
  return `M ${x1} ${y1} C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${x2} ${y2}`
}

const getConnectionPath = (conn: Connection) => {
  const from = nodes.value.find(n => n.id === conn.fromId)
  const to = nodes.value.find(n => n.id === conn.toId)
  if (!from || !to) return ''
  const isReverse = connections.value.findIndex(c => c.fromId === conn.fromId && c.toId === conn.toId) > 
                    connections.value.findIndex(c => c.fromId === conn.toId && c.toId === conn.fromId) &&
                    connections.value.some(c => c.fromId === conn.toId && c.toId === conn.fromId)
  return getCurvePath(from.x + NODE_WIDTH, from.y + NODE_HEIGHT / 2, to.x, to.y + NODE_HEIGHT / 2, isReverse)
}

const startDragNode = (event: MouseEvent, nodeId: string) => {
  if (event.button !== 0) return 
  event.stopPropagation()
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

// @ts-ignore
const _startReconnecting = (event: MouseEvent, conn: Connection) => {
  event.stopPropagation()
  const fromNode = nodes.value.find(n => n.id === conn.fromId)
  if (fromNode) {
    activeDrawing.value = { fromId: conn.fromId, startX: fromNode.x + NODE_WIDTH, startY: fromNode.y + NODE_HEIGHT / 2 }
    connections.value = connections.value.filter(c => c.id !== conn.id)
  }
}

// 合并连接修复逻辑
const finishConnecting = (toNodeId: string) => {
  if (activeDrawing.value && activeDrawing.value.fromId !== toNodeId) {
    // 允许任何节点连向此目标节点，无需去重
    connections.value.push({ 
      id: `conn-${Date.now()}`, 
      fromId: activeDrawing.value.fromId, 
      toId: toNodeId 
    });
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
  // 仅当没有击中节点输入口时才重置连线状态
  if (activeDrawing.value && !(event.target as HTMLElement).closest('.node-input-port')) {
      activeDrawing.value = null;
  }
}

onMounted(() => {
  window.addEventListener('mousemove', onGlobalMouseMove)
  window.addEventListener('mouseup', onGlobalMouseUp)
  window.addEventListener('click', () => contextMenu.show = false)
})
</script>

<template>
  <BaseLayout>
    <div class="flex w-full h-full overflow-hidden relative select-none">
      <aside class="shrink-0 border-r border-white/10 flex flex-col w-[280px]" style="background: rgba(20, 30, 50, 0.8);">
        <div class="p-6 border-b border-white/5"><h3 class="text-lg font-semibold text-white/95">组件库</h3></div>
        <div class="p-6 space-y-8 overflow-y-auto custom-scrollbar flex-1">
          <div class="space-y-4">
            <h4 class="text-[10px] font-bold text-white/30 uppercase tracking-[0.2em]">智能体角色</h4>
            <div v-for="item in [{name:'客服助手', icon:'lucide:message-circle', type:'agent', color:'#3B9BFF'}, {name:'分析专家', icon:'lucide:database', type:'agent', color:'#3B9BFF'}]" :key="item.name" draggable="true" @dragstart="onSidebarDragStart($event, item)" class="p-3 bg-white/5 border border-white/10 rounded-xl cursor-grab hover:bg-white/10 transition-colors">
              <div class="flex items-center gap-3"><div class="w-8 h-8 rounded-lg flex items-center justify-center bg-[#3B9BFF]/20 text-[#3B9BFF]"><iconify-icon :icon="item.icon"></iconify-icon></div><span class="text-sm text-white/80">{{ item.name }}</span></div>
            </div>
          </div>
          <div class="space-y-4">
            <h4 class="text-[10px] font-bold text-white/30 uppercase tracking-[0.2em]">逻辑控制</h4>
            <div v-for="item in [{name:'条件分支', icon:'lucide:git-branch', type:'branch', color:'#FFB846'}, {name:'并行执行', icon:'lucide:shuffle', type:'parallel', color:'#50C878'}, {name:'合并节点', icon:'lucide:combine', type:'merge', color:'#A78BFA'}]" :key="item.name" draggable="true" @dragstart="onSidebarDragStart($event, item)" class="p-3 bg-white/5 border border-white/10 rounded-xl cursor-grab">
              <div class="flex items-center gap-3"><div class="w-8 h-8 rounded-lg flex items-center justify-center shadow-inner" :style="{backgroundColor: item.color+'20', color: item.color}"><iconify-icon :icon="item.icon"></iconify-icon></div><span class="text-sm text-white/80">{{ item.name }}</span></div>
            </div>
          </div>
        </div>
      </aside>

      <main ref="canvasRef" @dragover.prevent @drop="onCanvasDrop" class="flex-1 relative bg-[#0F1928] overflow-hidden" style="background-image: radial-gradient(circle, rgba(255, 255, 255, 0.03) 1px, transparent 1px); background-size: 30px 30px;">
        <!-- Nodes -->
        <div v-for="node in nodes" :key="node.id" class="absolute p-4 bg-[#1A2536]/90 backdrop-blur-xl border rounded-xl shadow-2xl transition-all group" :style="{ left: node.x+'px', top: node.y+'px', width: NODE_WIDTH+'px', height: NODE_HEIGHT+'px', borderColor: isDraggingNode === node.id ? node.color : 'rgba(255,255,255,0.1)', zIndex: 10 }" @mousedown="startDragNode($event, node.id)" @contextmenu="openContextMenu($event, 'node', node.id)">
          <!-- 核心：扩大左侧捕获区域 -->
          <div 
            class="absolute -left-6 top-0 bottom-0 w-12 z-50 cursor-crosshair node-input-port flex items-center justify-center"
            @mouseup.stop="finishConnecting(node.id)"
          >
             <div class="w-4 h-4 bg-[#0F1928] border-2 border-white/20 rounded-full group-hover:border-[#3B9BFF] transition-all"></div>
          </div>
          <div class="absolute -right-2 top-1/2 -translate-y-1/2 w-4 h-4 bg-[#3B9BFF] border-2 border-white rounded-full z-20 cursor-crosshair hover:scale-125 transition-all" @mousedown.stop="startConnecting($event, node.id)"></div>
          <div class="flex items-center gap-3 h-full"><div class="w-10 h-10 rounded-lg flex items-center justify-center text-lg shrink-0 shadow-inner" :style="{backgroundColor: node.color+'20', color: node.color}"><iconify-icon :icon="node.icon"></iconify-icon></div><div class="flex-1 overflow-hidden"><div class="text-sm font-bold text-white/90 truncate">{{ node.name }}</div><div class="text-[9px] text-white/30 uppercase">{{ node.type }}</div></div></div>
        </div>

        <svg class="absolute inset-0 w-full h-full pointer-events-none z-[30]">
          <defs><marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto"><polygon points="0 0, 10 3.5, 0 7" fill="#3B9BFF" /></marker><filter id="glow"><feGaussianBlur stdDeviation="2" result="coloredBlur"/><feMerge><feMergeNode in="coloredBlur"/><feMergeNode in="SourceGraphic"/></feMerge></filter></defs>
          <g v-for="conn in connections" :key="conn.id" class="pointer-events-auto cursor-pointer group" @contextmenu="openContextMenu($event, 'conn', conn.id)">
            <path :d="getConnectionPath(conn)" stroke="rgba(59, 155, 255, 0.01)" stroke-width="20" fill="none" />
            <path :d="getConnectionPath(conn)" stroke="#3B9BFF" stroke-width="2.5" fill="none" marker-end="url(#arrowhead)" filter="url(#glow)" class="group-hover:stroke-white transition-colors" />
          </g>
          <path v-if="activeDrawing" :d="getCurvePath(activeDrawing.startX, activeDrawing.startY, mousePos.x, mousePos.y)" stroke="#3B9BFF" stroke-width="2" stroke-dasharray="6,6" fill="none" />
        </svg>

        <div v-if="contextMenu.show" class="fixed z-[100] bg-[#1A2536] border border-white/10 rounded-lg shadow-2xl py-1 min-w-[120px]" :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }" @click.stop>
          <button @click="handleDelete" class="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-red-500/10 flex items-center gap-3"><iconify-icon icon="lucide:trash-2"></iconify-icon>删除节点/线条</button>
        </div>
      </main>

      <aside class="shrink-0 border-l border-white/10 flex flex-col w-[320px]" style="background: rgba(20, 30, 50, 0.8);">
        <div class="p-6 border-b border-white/5"><h3 class="text-lg font-semibold text-white/95">配置面板</h3></div>
        <div class="flex-1 p-6 flex flex-col items-center justify-center text-white/20 text-center"><iconify-icon icon="lucide:settings-2" class="text-5xl mb-4"></iconify-icon><p class="text-sm">选中节点进行配置</p></div>
      </aside>
    </div>
  </BaseLayout>
</template>

<style scoped>
:deep(main) { padding: 0 !important; }
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
</style>