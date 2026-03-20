<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL } from '../config/api'

const isDrawerOpen = ref(false)
const isInstallDrawerOpen = ref(false)
const installUrl = ref('')
const isInstalling = ref(false)

const skills = ref<any[]>([])

const fetchSkills = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/skills/`)
    const data = await response.json()
    skills.value = data.map((s: any) => ({
      ...s,
      desc: s.description || '',
      color: s.color || '#3B9BFF'
    }))
  } catch (error) {
    console.error('Failed to fetch skills:', error)
  }
}

onMounted(() => {
  fetchSkills()
})

const openCreateDrawer = () => isDrawerOpen.value = true
const openInstallDrawer = () => isInstallDrawerOpen.value = true

const closeDrawers = () => {
  isDrawerOpen.value = false
  isInstallDrawerOpen.value = false
  installUrl.value = ''
  isInstalling.value = false
}

const handleInstall = async () => {
  if (!installUrl.value) return
  isInstalling.value = true
  try {
    const response = await fetch(`${API_BASE_URL}/api/skills/install`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url: installUrl.value })
    })
    if (response.ok) {
      await fetchSkills()
      closeDrawers()
    }
  } catch (error) {
    console.error('Installation failed:', error)
  } finally {
    isInstalling.value = false
  }
}

const isDeleteModalOpen = ref(false)
const skillToDelete = ref<any>(null)

const confirmDelete = (skill: any) => {
  skillToDelete.value = skill
  isDeleteModalOpen.value = true
}

const closeDeleteModal = () => {
  isDeleteModalOpen.value = false
  skillToDelete.value = null
}

const executeDelete = async () => {
  if (!skillToDelete.value) return
  try {
    const response = await fetch(`${API_BASE_URL}/api/skills/${skillToDelete.value.id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      skills.value = skills.value.filter(s => s.id !== skillToDelete.value.id)
      closeDeleteModal()
    }
  } catch (error) {
    console.error('Failed to delete skill:', error)
  }
}

const handleStart = (skill: any) => {
  console.log('🚀 Starting skill:', skill.name);
  alert(`正在启动技能: ${skill.name}`);
}

const isDetailModalOpen = ref(false)
const selectedSkillForDetail = ref<any>(null)

const showDetail = (skill: any) => {
  selectedSkillForDetail.value = skill
  isDetailModalOpen.value = true
}

const closeDetailModal = () => {
  isDetailModalOpen.value = false
  selectedSkillForDetail.value = null
}
</script>

<template>
  <BaseLayout>
    <main class="overflow-x-hidden flex flex-col grow p-8 gap-y-6 relative">
      
      <!-- Header Area -->
      <div class="flex justify-between items-end">
        <div>
          <h1 class="text-2xl font-bold text-white/95 mb-1">技能库</h1>
          <p class="text-xs text-white/40">定义和管理平台可用的标准技能，可供 Agent 绑定使用</p>
        </div>
        <div class="flex items-center gap-3">
          <!-- Add New Button -->
          <button 
            @click="openCreateDrawer"
            class="bg-white/5 hover:bg-white/10 text-white/80 px-5 py-2.5 rounded-xl flex items-center gap-2 border border-white/10 transition-all active:scale-95"
          >
            <iconify-icon icon="lucide:plus" class="text-lg"></iconify-icon>
            <span class="text-sm font-semibold">开发新技能</span>
          </button>
          <!-- Install Button -->
          <button 
            @click="openInstallDrawer"
            class="bg-[#3B9BFF] hover:bg-[#2A7FDB] text-white px-5 py-2.5 rounded-xl flex items-center gap-2 shadow-[0_0_15px_rgba(59,155,255,0.3)] transition-all active:scale-95"
          >
            <iconify-icon icon="lucide:download-cloud" class="text-lg"></iconify-icon>
            <span class="text-sm font-semibold">安装新技能</span>
          </button>
        </div>
      </div>

      <!-- Skills Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        <div v-for="skill in skills" :key="skill.id" 
          class="group p-6 bg-[#1A2536]/40 border border-white/10 rounded-2xl hover:border-[#3B9BFF]/40 hover:bg-[#1A2536]/70 transition-all duration-300 flex flex-col gap-5 relative overflow-hidden"
        >
          <div class="absolute -right-12 -top-12 w-24 h-24 blur-[60px] opacity-10 group-hover:opacity-25 transition-opacity" :style="{backgroundColor: skill.color}"></div>
          <div class="flex justify-between items-start relative z-10">
            <div class="flex gap-4">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl shadow-inner shrink-0" :style="{backgroundColor: (skill.color || '#3B9BFF')+'15', color: skill.color || '#3B9BFF'}">
                <iconify-icon :icon="skill.icon || 'lucide:terminal'"></iconify-icon>
              </div>
              <div class="min-w-0">
                <h3 class="text-base font-bold text-white/90 truncate">{{ skill.name }}</h3>
                <div class="mt-1 flex items-center gap-2">
                  <span class="text-[10px] px-1.5 py-0.5 rounded bg-white/5 text-white/40 border border-white/5 uppercase tracking-tighter">{{ skill.type }}</span>
                </div>
              </div>
            </div>
            <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-all transform translate-x-2 group-hover:translate-x-0">
              <button @click.stop="handleStart(skill)" class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/40 hover:text-[#50C878] hover:bg-[#50C878]/10 transition-all" title="启动技能"><iconify-icon icon="lucide:play" class="text-sm ml-0.5"></iconify-icon></button>
              <button @click.stop="confirmDelete(skill)" class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/40 hover:text-red-400 hover:bg-red-500/10 transition-all"><iconify-icon icon="lucide:trash-2" class="text-sm"></iconify-icon></button>
            </div>
          </div>
          <p class="text-sm text-white/50 leading-relaxed line-clamp-3 h-15">{{ skill.desc }}</p>
          <button @click.stop="showDetail(skill)" class="text-[10px] text-[#3B9BFF] hover:text-[#5FB4FF] transition-colors flex items-center gap-1 w-fit mt-[-10px]"><iconify-icon icon="lucide:info" class="text-xs"></iconify-icon>查看详情</button>
          <div class="flex items-center justify-between pt-4 border-t border-white/5">
            <div class="flex items-center gap-2">
              <span class="text-[10px] text-white/20 uppercase font-bold tracking-widest">Path</span>
              <span class="text-[10px] text-white/50 font-mono">{{ skill.version }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Drawer 1: Add New Skill (Manual) -->
      <Transition 
        enter-active-class="transition duration-300 ease-out" enter-from-class="translate-x-full" enter-to-class="translate-x-0"
        leave-active-class="transition duration-200 ease-in" leave-from-class="translate-x-0" leave-to-class="translate-x-full"
      >
        <aside v-if="isDrawerOpen" class="fixed top-0 right-0 w-[560px] h-full z-[100] bg-[#1A2536]/98 backdrop-blur-3xl border-l border-[#3B9BFF]/30 shadow-[-20px_0_50px_rgba(0,0,0,0.5)] flex flex-col">
          <div class="p-6 border-b border-white/10 flex justify-between items-center">
            <h2 class="text-xl font-bold text-white/95">开发新技能</h2>
            <button @click="closeDrawers" class="w-10 h-10 rounded-full flex items-center justify-center hover:bg-white/5 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
          </div>
          <div class="flex-1 overflow-y-auto p-8 space-y-8 custom-scrollbar">
            <div class="grid grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="text-[10px] text-white/40 font-bold uppercase tracking-widest">技能内部名称</label>
                <input type="text" placeholder="my_new_skill" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 text-sm text-white outline-none focus:border-[#3B9BFF]/50 transition-all font-mono">
              </div>
              <div class="space-y-2">
                <label class="text-[10px] text-white/40 font-bold uppercase tracking-widest">显示版本</label>
                <input type="text" value="v1.0.0" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 text-sm text-white outline-none focus:border-[#3B9BFF]/50 transition-all font-mono">
              </div>
            </div>
            <div class="space-y-2">
              <label class="text-[10px] text-white/40 font-bold uppercase tracking-widest">技能描述</label>
              <textarea placeholder="简述该技能的功能和输入输出..." rows="3" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-sm text-white outline-none focus:border-[#3B9BFF]/50 transition-all resize-none"></textarea>
            </div>
            <div class="space-y-4 pt-2 border-t border-white/5">
              <div class="flex justify-between items-center">
                <label class="text-[10px] text-white/40 font-bold uppercase tracking-widest">核心逻辑 (Python)</label>
                <button class="text-[10px] text-[#3B9BFF] hover:underline">加载示例代码</button>
              </div>
              <div class="bg-[#0F1928] border border-white/10 rounded-xl overflow-hidden shadow-inner">
                <textarea class="w-full h-80 p-4 text-[11px] font-mono text-[#5FB4FF] bg-transparent outline-none resize-none" spellcheck="false">def handler(params, context):
    pass</textarea>
              </div>
            </div>
          </div>
          <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md flex gap-4">
            <button @click="closeDrawers" class="flex-1 py-3 rounded-xl border border-white/10 text-sm text-white/60 hover:bg-white/5">取消</button>
            <button class="flex-1 py-3 rounded-xl bg-[#3B9BFF] text-white text-sm font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)]">保存至技能库</button>
          </div>
        </aside>
      </Transition>

      <!-- Drawer 2: Install Skill (via URL) -->
      <Transition 
        enter-active-class="transition duration-300 ease-out" enter-from-class="translate-x-full" enter-to-class="translate-x-0"
        leave-active-class="transition duration-200 ease-in" leave-from-class="translate-x-0" leave-to-class="translate-x-full"
      >
        <aside v-if="isInstallDrawerOpen" class="fixed top-0 right-0 w-[500px] h-full z-[100] bg-[#1A2536]/98 backdrop-blur-3xl border-l border-[#3B9BFF]/30 shadow-[-20px_0_50px_rgba(0,0,0,0.5)] flex flex-col">
          <div class="p-6 border-b border-white/10 flex justify-between items-center">
            <div>
              <h2 class="text-xl font-bold text-white/95">安装远程技能</h2>
              <p class="text-[10px] text-white/30 uppercase tracking-widest mt-1">Import skill from external URL</p>
            </div>
            <button @click="closeDrawers" class="w-10 h-10 rounded-full flex items-center justify-center hover:bg-white/5 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
          </div>
          <div class="flex-1 p-8 space-y-8">
            <div class="space-y-4">
              <label class="text-[10px] text-white/40 font-bold uppercase tracking-widest">技能包 HTTP(S) 地址</label>
              <div class="relative">
                <div class="absolute inset-y-0 left-4 flex items-center pointer-events-none text-white/30"><iconify-icon icon="lucide:link"></iconify-icon></div>
                <input 
                  v-model="installUrl"
                  type="text" 
                  placeholder="https://api.skills.com/package.json" 
                  class="w-full bg-white/5 border border-white/10 rounded-xl pl-12 pr-4 py-4 text-sm text-white outline-none focus:border-[#3B9BFF]/50 transition-all font-mono"
                >
              </div>
              <p class="text-[10px] text-white/20 italic">平台将尝试访问该地址，识别技能清单、图标及执行逻辑并自动完成安装。</p>
            </div>

            <div v-if="isInstalling" class="flex flex-col items-center justify-center py-12 gap-4">
               <iconify-icon icon="lucide:loader-2" class="text-4xl text-[#3B9BFF] animate-spin"></iconify-icon>
               <span class="text-sm text-white/50 animate-pulse">正在拉取并解析技能包...</span>
            </div>
          </div>
          <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md flex gap-4">
            <button @click="closeDrawers" class="flex-1 py-3 rounded-xl border border-white/10 text-sm text-white/60 hover:bg-white/5">取消</button>
            <button 
              @click="handleInstall"
              :disabled="!installUrl || isInstalling"
              class="flex-1 py-3 rounded-xl bg-[#3B9BFF] text-white text-sm font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)] disabled:opacity-50 disabled:cursor-not-allowed transition-all"
            >
              {{ isInstalling ? '安装中...' : '开始安装' }}
            </button>
          </div>
        </aside>
      </Transition>
      
      <!-- Drawer 4: Delete Confirmation Modal -->
      <Transition 
        enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0 scale-95" enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95"
      >
        <div v-if="isDeleteModalOpen" class="fixed inset-0 z-[120] flex items-center justify-center p-6">
          <div @click="closeDeleteModal" class="absolute inset-0 bg-black/60 backdrop-blur-md"></div>
          <div class="relative w-full max-w-md bg-[#1A2536] border border-red-500/20 rounded-3xl shadow-2xl overflow-hidden flex flex-col">
            <div class="p-6 border-b border-white/10 flex justify-between items-center bg-red-500/5">
              <div class="flex items-center gap-4">
                <div class="w-10 h-10 rounded-xl bg-red-500/10 flex items-center justify-center text-xl text-red-500">
                  <iconify-icon icon="lucide:alert-triangle"></iconify-icon>
                </div>
                <h2 class="text-xl font-bold text-white/95">确认删除技能</h2>
              </div>
              <button @click="closeDeleteModal" class="w-8 h-8 rounded-full flex items-center justify-center hover:bg-white/10 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
            </div>
            <div class="p-8 space-y-4">
              <p class="text-sm text-white/70 leading-relaxed">
                您确定要从技能库中移除 <span class="text-white font-bold">"{{ skillToDelete?.name }}"</span> 吗？此操作不可撤销。
              </p>
              <div class="bg-black/20 rounded-xl p-3 border border-white/5">
                <p class="text-[10px] text-white/30 uppercase font-bold tracking-widest mb-1">物理路径</p>
                <code class="text-[10px] text-white/40 break-all font-mono">{{ skillToDelete?.version }}</code>
              </div>
            </div>
            <div class="p-6 border-t border-white/10 bg-white/5 flex gap-3">
              <button @click="closeDeleteModal" class="flex-1 py-3 rounded-xl border border-white/10 text-sm text-white/60 hover:bg-white/5 transition-all">取消</button>
              <button @click="executeDelete" class="flex-1 py-3 rounded-xl bg-red-500/80 hover:bg-red-500 text-white text-sm font-bold shadow-[0_0_20px_rgba(239,68,68,0.2)] transition-all">确认删除</button>
            </div>
          </div>
        </div>
      </Transition>


      <!-- Drawer 3: Skill Details Modal -->
      <Transition 
        enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0 scale-95" enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95"
      >
        <div v-if="isDetailModalOpen" class="fixed inset-0 z-[110] flex items-center justify-center p-6">
          <div @click="closeDetailModal" class="absolute inset-0 bg-black/60 backdrop-blur-md"></div>
          <div class="relative w-full max-w-2xl bg-[#1A2536] border border-white/10 rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[80vh]">
            <div class="p-6 border-b border-white/10 flex justify-between items-center bg-white/5">
              <div class="flex items-center gap-4">
                <div class="w-10 h-10 rounded-xl flex items-center justify-center text-xl" :style="{backgroundColor: selectedSkillForDetail?.color+'15', color: selectedSkillForDetail?.color}">
                  <iconify-icon :icon="selectedSkillForDetail?.icon"></iconify-icon>
                </div>
                <h2 class="text-xl font-bold text-white/95">{{ selectedSkillForDetail?.name }}</h2>
              </div>
              <button @click="closeDetailModal" class="w-8 h-8 rounded-full flex items-center justify-center hover:bg-white/10 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
            </div>
            <div class="flex-1 overflow-y-auto p-8 custom-scrollbar space-y-6">
              <div class="space-y-2">
                <label class="text-[10px] text-white/40 font-bold uppercase tracking-widest">完整描述</label>
                <p class="text-sm text-white/70 leading-relaxed whitespace-pre-wrap">{{ selectedSkillForDetail?.desc }}</p>
              </div>
              <div class="space-y-2 pt-4 border-t border-white/5">
                <label class="text-[10px] text-white/40 font-bold uppercase tracking-widest">物理路径</label>
                <div class="bg-black/20 rounded-xl p-4 border border-white/5">
                  <code class="text-[11px] text-[#5FB4FF] break-all font-mono">{{ selectedSkillForDetail?.version }}</code>
                </div>
              </div>
            </div>
            <div class="p-4 border-t border-white/10 bg-white/5 flex justify-end">
              <button @click="closeDetailModal" class="px-6 py-2 rounded-xl bg-white/5 border border-white/10 text-sm text-white/60 hover:bg-white/10 transition-all">关闭</button>
            </div>
          </div>
        </div>
      </Transition>


      <div v-if="isDrawerOpen || isInstallDrawerOpen || isDetailModalOpen || isDeleteModalOpen" @click="closeDeleteModal(); closeDetailModal(); closeDrawers();" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-[80]"></div>
    </main>
  </BaseLayout>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
</style>
