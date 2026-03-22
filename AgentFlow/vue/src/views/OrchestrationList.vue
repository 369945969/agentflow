<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL } from '../config/api'

const router = useRouter()
const workflows = ref<any[]>([])
const searchQuery = ref('')

// Delete modal state
const isDeleteModalOpen = ref(false)
const workflowToDelete = ref<any>(null)
const isDeleting = ref(false)

const fetchWorkflows = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/workflows/`)
    const data = await response.json()
    workflows.value = data || []
  } catch (error) {
    console.error('Failed to fetch workflows:', error)
    workflows.value = []
  }
}

const filteredWorkflows = computed(() => {
  if (!workflows.value) return []
  if (!searchQuery.value) return workflows.value
  const query = searchQuery.value.toLowerCase()
  return workflows.value.filter(w => 
    (w.name && w.name.toLowerCase().includes(query)) || 
    (w.description && w.description.toLowerCase().includes(query))
  )
})

const goToEditor = (id?: string) => {
  if (id) {
    router.push(`/orchestration-editor/${id}`)
  } else {
    router.push('/orchestration-editor/new')
  }
}

const confirmDelete = (workflow: any) => {
  workflowToDelete.value = workflow
  isDeleteModalOpen.value = true
}

const executeDelete = async () => {
  if (!workflowToDelete.value) return
  isDeleting.value = true
  try {
    const response = await fetch(`${API_BASE_URL}/api/workflows/${workflowToDelete.value.id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      await fetchWorkflows()
      isDeleteModalOpen.value = false
      workflowToDelete.value = null
    }
  } catch (error) {
    console.error('Failed to delete workflow:', error)
    alert('删除失败')
  } finally {
    isDeleting.value = false
  }
}

onMounted(fetchWorkflows)
</script>

<template>
  <BaseLayout>
    <main class="overflow-x-hidden flex flex-col grow shrink gap-y-8 p-8 relative">
      <!-- Header Area -->
      <div class="flex justify-between items-center shrink-0">
        <div>
          <h1 class="text-3xl font-bold text-white mb-2">编排列表</h1>
          <p class="text-white/50 text-sm">管理和配置您的工作流编排，实现复杂的业务逻辑自动化。</p>
        </div>
        <button @click="goToEditor()" class="px-6 py-3 bg-[#3B9BFF] hover:bg-[#2A7FDB] text-white rounded-xl font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)] transition-all flex items-center gap-2 h-fit">
          <iconify-icon icon="lucide:plus"></iconify-icon>
          新建编排
        </button>
      </div>

      <!-- Search Area -->
      <div style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="p-6 border-[1px] border-solid rounded-2xl shrink-0">
        <div class="bg-white/5 border border-white/10 rounded-xl px-4 py-3 flex items-center gap-3 backdrop-blur-md w-full">
          <iconify-icon icon="lucide:search" class="text-white/40"></iconify-icon>
          <input v-model="searchQuery" type="text" placeholder="按名称/描述搜索编排..." class="bg-transparent border-none outline-none text-sm text-white/70 w-full">
        </div>
      </div>

      <!-- Workflows Content -->
      <div class="flex-1 min-h-0 overflow-y-auto custom-scrollbar pr-2">
        <div v-if="filteredWorkflows.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div 
            v-for="wf in filteredWorkflows" :key="wf.id"
            @click="goToEditor(wf.id)"
            class="group p-6 rounded-2xl border border-white/10 bg-[#1A2536]/40 backdrop-blur-xl hover:shadow-2xl hover:border-[#3B9BFF]/30 transition-all duration-300 cursor-pointer flex flex-col"
          >
            <div class="flex items-center gap-4 mb-4">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0 bg-[#3B9BFF]/20 text-[#3B9BFF]">
                <iconify-icon icon="lucide:workflow" class="text-2xl"></iconify-icon>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="text-white/95 font-semibold text-lg truncate">{{ wf.name }}</h3>
                <div class="text-[10px] text-white/30 truncate uppercase tracking-tighter">ID: {{ wf.id.substring(0,8) }}...</div>
              </div>
            </div>
            
            <p class="text-sm text-white/50 mb-6 line-clamp-2 h-10">{{ wf.description || '暂无编排描述' }}</p>
            
            <div class="flex items-center justify-between pt-4 border-t border-white/5 mt-auto shrink-0">
              <span class="text-[10px] text-white/30">{{ new Date(wf.created_at).toLocaleDateString() }}</span>
              <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button @click.stop="goToEditor(wf.id)" class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/60 hover:text-white hover:bg-[#3B9BFF]/20" title="编辑"><iconify-icon icon="lucide:edit-3"></iconify-icon></button>
                <button @click.stop="confirmDelete(wf)" class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/40 hover:text-red-400 hover:bg-red-500/10 transition-all" title="删除"><iconify-icon icon="lucide:trash-2"></iconify-icon></button>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="py-32 flex flex-col items-center justify-center text-white/20 border-2 border-dashed border-white/5 rounded-3xl">
          <div class="w-20 h-20 bg-white/5 rounded-full flex items-center justify-center mb-6">
            <iconify-icon icon="lucide:layout-template" class="text-5xl mb-4"></iconify-icon>
          </div>
          <p class="text-xl font-medium mb-2">暂无编排工作流</p>
          <p class="text-sm mb-8 text-white/10">创建一个可视化工作流来自动执行复杂任务</p>
          <button @click="goToEditor()" class="px-6 py-2 border border-white/10 rounded-xl hover:bg-white/5 transition-all text-white/60">开始创建</button>
        </div>
      </div>

      <!-- Delete Confirmation Modal -->
      <Transition 
        enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0 scale-95" enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95"
      >
        <div v-if="isDeleteModalOpen" class="fixed inset-0 z-[120] flex items-center justify-center p-6 bg-black/60 backdrop-blur-sm">
          <div @click="isDeleteModalOpen = false" class="absolute inset-0"></div>
          <div class="relative w-full max-w-md bg-[#1A2536] border border-red-500/20 rounded-3xl shadow-2xl overflow-hidden flex flex-col">
            <div class="p-6 border-b border-white/10 flex justify-between items-center bg-red-500/5">
              <div class="flex items-center gap-4">
                <div class="w-10 h-10 rounded-xl bg-red-500/10 flex items-center justify-center text-xl text-red-500">
                  <iconify-icon icon="lucide:alert-triangle"></iconify-icon>
                </div>
                <h2 class="text-xl font-bold text-white/95">确认删除编排</h2>
              </div>
              <button @click="isDeleteModalOpen = false" class="w-8 h-8 rounded-full flex items-center justify-center hover:bg-white/10 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
            </div>
            <div class="p-8 space-y-4">
              <p class="text-sm text-white/70 leading-relaxed">
                您确定要移除编排流程 <span class="text-white font-bold">"{{ workflowToDelete?.name }}"</span> 吗？此操作将永久删除该工作流配置。
              </p>
            </div>
            <div class="p-6 border-t border-white/10 bg-white/5 flex gap-3">
              <button @click="isDeleteModalOpen = false" class="flex-1 py-3 rounded-xl border border-white/10 text-sm text-white/60 hover:bg-white/5 transition-all">取消</button>
              <button @click="executeDelete" :disabled="isDeleting" class="flex-1 py-3 rounded-xl bg-red-500/80 hover:bg-red-500 text-white text-sm font-bold shadow-[0_0_20px_rgba(239,68,68,0.2)] transition-all">
                {{ isDeleting ? '正在删除...' : '确认删除' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </main>
  </BaseLayout>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
</style>
