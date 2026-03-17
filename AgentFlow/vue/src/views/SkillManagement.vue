<script setup lang="ts">
import { ref } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'

const activeTab = ref('全部')
const isDrawerOpen = ref(false)
const isInstallDrawerOpen = ref(false)
const installUrl = ref('')
const isInstalling = ref(false)

const skills = ref([
  { id: 1, name: 'Google Search', type: 'HTTP API', desc: '集成 Google 搜索引擎，允许 Agent 获取最新的网页、新闻和学术内容。', version: 'v2.1.0', icon: 'lucide:search', color: '#3B9BFF' },
  { id: 2, name: 'Python Runner', type: 'Script', desc: '安全的沙箱环境，支持执行 Python 脚本进行数据处理、绘图或数学运算。', version: 'v1.4.2', icon: 'lucide:code-2', color: '#50C878' },
  { id: 3, name: 'MySQL Connector', type: 'Database', desc: '标准数据库连接器，支持执行 SQL 查询、更新和结构化数据导出。', version: 'v3.0.1', icon: 'lucide:database', color: '#FFB846' },
  { id: 4, name: 'Email Dispatcher', type: 'Plugin', desc: '自动化邮件发送工具，支持 SMTP 配置及多模板 HTML 邮件推送。', version: 'v1.0.5', icon: 'lucide:mail', color: '#A78BFA' },
  { id: 5, name: 'PDF Parser', type: 'Script', desc: '高性能文档解析引擎，能够识别 PDF 中的文本、表格并提取结构化信息。', version: 'v2.2.0', icon: 'lucide:file-text', color: '#F472B6' },
  { id: 6, name: 'Weather API', type: 'HTTP API', desc: '调用全球实时天气数据，包括温度、湿度、风速及未来 7 天预报。', version: 'v1.2.0', icon: 'lucide:cloud-sun', color: '#60A5FA' }
])

const openCreateDrawer = () => isDrawerOpen.value = true
const openInstallDrawer = () => isInstallDrawerOpen.value = true

const closeDrawers = () => {
  isDrawerOpen.value = false
  isInstallDrawerOpen.value = false
  installUrl.value = ''
  isInstalling.value = false
}

const handleInstall = () => {
  if (!installUrl.value) return
  isInstalling.value = true
  // Simulate API fetch and installation
  setTimeout(() => {
    skills.value.unshift({
      id: Date.now(),
      name: 'New Linked Skill',
      type: 'HTTP API',
      desc: `从远程地址 ${installUrl.value} 成功识别并安装的技能。`,
      version: 'v1.0.0',
      icon: 'lucide:link',
      color: '#3B9BFF'
    })
    closeDrawers()
  }, 2000)
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
            <span class="text-sm font-semibold">新增新技能</span>
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

      <!-- Filter Bar -->
      <div class="flex items-center justify-between bg-white/5 border border-white/10 p-1.5 rounded-2xl backdrop-blur-md">
        <div class="flex gap-1">
          <button v-for="t in ['全部', 'HTTP API', 'Script', 'Database', 'Plugin']" :key="t" 
            @click="activeTab = t"
            :class="['px-4 py-1.5 rounded-xl text-xs transition-all font-medium', activeTab === t ? 'bg-[#3B9BFF]/20 text-[#3B9BFF] border border-[#3B9BFF]/30' : 'text-white/40 hover:text-white/60']"
          >{{ t }}</button>
        </div>
        <div class="flex items-center gap-3 px-3 py-1.5 bg-white/5 border border-white/10 rounded-xl mr-1">
          <iconify-icon icon="lucide:search" class="text-white/30 text-sm"></iconify-icon>
          <input type="text" placeholder="搜索技能库..." class="bg-transparent border-none outline-none text-xs text-white/70 w-48">
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
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl shadow-inner shrink-0" :style="{backgroundColor: skill.color+'15', color: skill.color}">
                <iconify-icon :icon="skill.icon"></iconify-icon>
              </div>
              <div class="min-w-0">
                <h3 class="text-base font-bold text-white/90 truncate">{{ skill.name }}</h3>
                <div class="mt-1 flex items-center gap-2">
                  <span class="text-[10px] px-1.5 py-0.5 rounded bg-white/5 text-white/40 border border-white/5 uppercase tracking-tighter">{{ skill.type }}</span>
                </div>
              </div>
            </div>
            <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-all transform translate-x-2 group-hover:translate-x-0">
              <button class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/40 hover:text-[#3B9BFF] hover:bg-[#3B9BFF]/10 transition-all"><iconify-icon icon="lucide:edit-3" class="text-sm"></iconify-icon></button>
              <button class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/40 hover:text-red-400 hover:bg-red-500/10 transition-all"><iconify-icon icon="lucide:trash-2" class="text-sm"></iconify-icon></button>
            </div>
          </div>
          <p class="text-sm text-white/50 leading-relaxed line-clamp-3 h-15">{{ skill.desc }}</p>
          <div class="flex items-center justify-between pt-4 border-t border-white/5">
            <div class="flex items-center gap-2">
              <span class="text-[10px] text-white/20 uppercase font-bold tracking-widest">Version</span>
              <span class="text-[10px] text-white/50 font-mono">{{ skill.version }}</span>
            </div>
            <div class="flex items-center gap-1 px-2 py-0.5 rounded bg-[#50C878]/10 border border-[#50C878]/20">
              <div class="w-1 h-1 rounded-full bg-[#50C878]"></div>
              <span class="text-[9px] text-[#50C878] font-bold uppercase">Ready</span>
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
            <h2 class="text-xl font-bold text-white/95">新增新技能</h2>
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

      <div v-if="isDrawerOpen || isInstallDrawerOpen" @click="closeDrawers" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-[90]"></div>
    </main>
  </BaseLayout>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
</style>
