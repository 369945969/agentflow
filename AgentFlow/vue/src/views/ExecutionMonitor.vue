<script setup lang="ts">
import { ref } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'

// viewMode: 'list' or 'detail'
const viewMode = ref<'list' | 'detail'>('list')
const selectedRecord = ref<any>(null)

const historyRecords = ref([
  { id: 'EXE-20240317-001', workflow: '客户服务自动化流程', startTime: '2024-03-17 14:30:25', duration: '4.5s', status: 'Running', color: '#3B9BFF' },
  { id: 'EXE-20240317-002', workflow: '数据报表生成', startTime: '2024-03-17 12:15:10', duration: '12.8s', status: 'Completed', color: '#50C878' },
  { id: 'EXE-20240317-003', workflow: '情感分析批处理', startTime: '2024-03-17 10:05:00', duration: '2.1s', status: 'Failed', color: '#F87171' },
  { id: 'EXE-20240316-098', workflow: '智能翻译助手', startTime: '2024-03-16 22:45:12', duration: '1.2s', status: 'Completed', color: '#50C878' },
  { id: 'EXE-20240316-097', workflow: '客户服务自动化流程', startTime: '2024-03-16 20:30:00', duration: '5.4s', status: 'Completed', color: '#50C878' },
])

const enterDetail = (record: any) => {
  selectedRecord.value = record
  viewMode.value = 'detail'
}

const goBackToList = () => {
  viewMode.value = 'list'
  selectedRecord.value = null
}
</script>

<template>
  <BaseLayout>
    <main class="overflow-x-hidden flex flex-col grow shrink p-8 gap-y-6 relative">
      
      <!-- LEVEL 1: LIST VIEW -->
      <template v-if="viewMode === 'list'">
        <div class="flex justify-between items-end mb-2">
          <div>
            <h1 class="text-2xl font-bold text-white/95 mb-1">执行历史</h1>
            <p class="text-xs text-white/40">查看和追踪所有编排流程的执行记录</p>
          </div>
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-3 px-3 py-1.5 bg-white/5 border border-white/10 rounded-xl">
              <iconify-icon icon="lucide:search" class="text-white/30 text-sm"></iconify-icon>
              <input type="text" placeholder="搜索流水号/流程名..." class="bg-transparent border-none outline-none text-xs text-white/70 w-48">
            </div>
            <button class="bg-white/5 hover:bg-white/10 text-white/70 px-4 py-2 rounded-xl border border-white/10 text-xs flex items-center gap-2">
              <iconify-icon icon="lucide:filter"></iconify-icon>
              筛选
            </button>
          </div>
        </div>

        <!-- History Table -->
        <div class="flex-1 bg-[#1A2536]/50 border border-white/10 rounded-2xl overflow-hidden backdrop-blur-xl">
          <table class="w-full border-collapse">
            <thead>
              <tr class="border-b border-white/5 bg-white/5">
                <th class="text-left p-4 text-[10px] font-bold text-white/30 uppercase tracking-widest">执行流水号</th>
                <th class="text-left p-4 text-[10px] font-bold text-white/30 uppercase tracking-widest">编排流程</th>
                <th class="text-left p-4 text-[10px] font-bold text-white/30 uppercase tracking-widest">开始时间</th>
                <th class="text-left p-4 text-[10px] font-bold text-white/30 uppercase tracking-widest">耗时</th>
                <th class="text-left p-4 text-[10px] font-bold text-white/30 uppercase tracking-widest">状态</th>
                <th class="text-right p-4 text-[10px] font-bold text-white/30 uppercase tracking-widest">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="record in historyRecords" :key="record.id" 
                class="border-b border-white/5 hover:bg-white/5 transition-colors group cursor-pointer"
                @click="enterDetail(record)"
              >
                <td class="p-4 text-sm font-mono text-[#3B9BFF]">{{ record.id }}</td>
                <td class="p-4 text-sm text-white/80 font-medium">{{ record.workflow }}</td>
                <td class="p-4 text-xs text-white/40">{{ record.startTime }}</td>
                <td class="p-4 text-xs text-white/60 font-mono">{{ record.duration }}</td>
                <td class="p-4">
                  <div class="flex items-center gap-2">
                    <div class="w-1.5 h-1.5 rounded-full" :style="{ backgroundColor: record.color }"></div>
                    <span class="text-xs" :style="{ color: record.color }">{{ record.status }}</span>
                  </div>
                </td>
                <td class="p-4 text-right">
                  <button class="text-white/30 group-hover:text-[#3B9BFF] transition-colors"><iconify-icon icon="lucide:chevron-right" class="text-lg"></iconify-icon></button>
                </td>
              </tr>
            </tbody>
          </table>
          <div class="p-4 flex justify-center border-t border-white/5">
             <button class="text-[10px] text-white/20 hover:text-white/40 uppercase tracking-widest">加载更多记录</button>
          </div>
        </div>
      </template>

      <!-- LEVEL 2: DETAIL VIEW (The original content) -->
      <template v-else>
        <div class="flex justify-between items-center mb-2">
          <div class="flex items-center gap-4">
            <button @click="goBackToList" class="w-10 h-10 rounded-full bg-white/5 border border-white/10 flex items-center justify-center text-white/60 hover:text-white hover:bg-white/10">
              <iconify-icon icon="lucide:arrow-left" class="text-xl"></iconify-icon>
            </button>
            <div>
              <div class="flex items-center gap-3">
                <h1 class="text-xl font-bold text-white/95">{{ selectedRecord?.workflow }}</h1>
                <span class="text-[10px] font-mono text-white/30 bg-white/5 px-2 py-0.5 rounded border border-white/10">{{ selectedRecord?.id }}</span>
              </div>
              <p class="text-xs text-white/40 mt-0.5">执行详情与性能监控回放</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <div class="px-4 py-2 bg-white/5 border border-white/10 rounded-xl flex items-center gap-4">
               <div v-for="s in ['Waiting', 'Running', 'Completed', 'Failed']" :key="s" class="flex items-center gap-1.5">
                  <div class="w-1.5 h-1.5 rounded-full" :class="{'bg-[#3B9BFF]': s === 'Running', 'bg-[#50C878]': s === 'Completed', 'bg-[#F87171]': s === 'Failed', 'bg-white/20': s === 'Waiting'}"></div>
                  <span class="text-[10px] text-white/40">{{ s }}</span>
               </div>
            </div>
          </div>
        </div>

        <!-- Detail Layout (Reused original code with minor fixes) -->
        <div class="flex-1 bg-[#1A2536]/30 border border-white/10 rounded-2xl overflow-hidden relative backdrop-blur-md flex flex-col">
          <!-- Toolbar -->
          <div class="p-4 border-b border-white/5 bg-white/5 flex justify-between items-center">
             <div class="flex gap-2">
                <button class="px-3 py-1.5 rounded-lg bg-[#3B9BFF]/20 text-[#3B9BFF] border border-[#3B9BFF]/30 text-[10px] font-bold uppercase tracking-wider">执行拓扑图</button>
                <button class="px-3 py-1.5 rounded-lg bg-white/5 text-white/40 border border-white/5 text-[10px] font-bold uppercase tracking-wider hover:text-white/60 transition-all">执行日志</button>
             </div>
             <div class="flex items-center gap-2">
                <button class="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center text-white/40 hover:text-white"><iconify-icon icon="lucide:zoom-in"></iconify-icon></button>
                <button class="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center text-white/40 hover:text-white"><iconify-icon icon="lucide:zoom-out"></iconify-icon></button>
             </div>
          </div>

          <!-- Flow Content (Mockup of the sequence) -->
          <div class="flex-1 flex items-center justify-center overflow-auto p-20">
            <div class="flex items-center gap-12">
              
              <!-- Start Node -->
              <div class="flex flex-col items-center gap-3">
                <div class="w-16 h-16 rounded-2xl bg-[#50C878]/20 border-2 border-[#50C878] flex items-center justify-center text-[#50C878] shadow-[0_0_20px_rgba(80,200,120,0.3)]">
                  <iconify-icon icon="lucide:play" class="text-2xl"></iconify-icon>
                </div>
                <span class="text-xs font-bold text-white/80">开始</span>
              </div>

              <div class="w-12 h-0.5 bg-[#50C878]/40 relative">
                <div class="absolute -top-1 right-0 w-2 h-2 rounded-full bg-[#50C878]"></div>
              </div>

              <!-- Process Node -->
              <div class="flex flex-col items-center gap-3">
                <div class="w-16 h-16 rounded-2xl bg-[#3B9BFF]/20 border-2 border-[#3B9BFF] flex items-center justify-center text-[#3B9BFF] shadow-[0_0_20px_rgba(59,155,255,0.3)] relative">
                  <iconify-icon icon="lucide:brain" class="text-2xl"></iconify-icon>
                  <div class="absolute -top-2 -right-2 bg-[#3B9BFF] text-[8px] font-bold px-1.5 py-0.5 rounded-full text-white animate-pulse">RUNNING</div>
                </div>
                <span class="text-xs font-bold text-white/80">意图识别</span>
              </div>

              <div class="w-12 h-0.5 bg-white/10 relative">
                <div class="absolute -top-1 right-0 w-2 h-2 rounded-full bg-white/10"></div>
              </div>

              <!-- Next Node (Locked) -->
              <div class="flex flex-col items-center gap-3 opacity-30">
                <div class="w-16 h-16 rounded-2xl bg-white/5 border border-white/10 flex items-center justify-center text-white/40">
                  <iconify-icon icon="lucide:message-square" class="text-2xl"></iconify-icon>
                </div>
                <span class="text-xs font-bold text-white/40">响应生成</span>
              </div>

            </div>
          </div>
        </div>
      </template>

    </main>
  </BaseLayout>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
</style>
