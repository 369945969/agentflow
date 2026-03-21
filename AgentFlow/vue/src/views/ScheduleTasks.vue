<script setup lang="ts">
import BaseLayout from '../components/BaseLayout.vue'
</script>

<template>
  <BaseLayout>
    <main id="28:64" orderindex="81" style="flex-basis: 0%; padding: 1.5rem 2rem;" class="overflow-x-hidden flex flex-col grow shrink gap-y-8 gap-x-8">

<div id="28:65" orderindex="80" class="flex justify-between items-center">
<h1 id="28:66" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-2xl font-semibold">Schedule Tasks</h1>
<div id="28:67" orderindex="81" class="flex items-center gap-y-4 gap-x-4">
<button id="28:68" orderindex="80" class="hover:border-[#5FB4FF]/50 hover:shadow-[0_0_20px_rgba(59,155,255,0.3)] flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #3B9BFF 80% , transparent ); backdrop-filter: blur(24px); padding: 0.75rem 1.5rem; border-color: color-mix( in oklab , #3B9BFF 50% , transparent );">
<div id="28:69" orderindex="80" class="flex justify-center items-center w-4 h-4">
<iconify-icon id="28:70" orderindex="80" style="color: rgba(255, 255, 255, 1);" icon="lucide:plus" class="text-sm"></iconify-icon>
</div>
<span id="28:71" orderindex="81" style="color: rgba(255, 255, 255, 1);" class="text-sm whitespace-nowrap">创建新任务</span>
</button>
<div id="28:72" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="28:73" orderindex="80" class="hover:bg-white/10 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<span id="28:74" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm">批量操作</span>
<div id="28:75" orderindex="81" class="flex justify-center items-center w-4 h-4">
<iconify-icon id="28:76" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:chevron-down" class="text-sm"></iconify-icon>
</div>
</button>
<div id="28:77" orderindex="81" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="overflow-x-hidden overflow-y-hidden flex border-[1px] border-solid rounded-xl">
<button id="28:78" orderindex="80" class="hover:bg-white/10 flex justify-center items-center w-10 h-10" style="background-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="28:79" orderindex="80" style="color: rgba(95, 180, 255, 1);" icon="lucide:list" class="text-base"></iconify-icon>
</button>
<button id="28:80" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-10 h-10">
<iconify-icon id="28:81" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:calendar-days" class="text-base"></iconify-icon>
</button>
</div>
</div>
</div>
</div>

<div id="28:82" orderindex="81" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="flex flex-col gap-y-4 gap-x-4 p-6 border-[1px] border-solid rounded-2xl">
<div id="28:83" orderindex="80" class="flex items-center gap-y-4 gap-x-4">
<div id="28:84" orderindex="80" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(12px); padding: 0.75rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" class="flex items-center min-w-[320px] border-[1px] border-solid rounded-xl">
<div id="28:85" orderindex="80" class="flex justify-center items-center w-5 h-5 mr-3">
<iconify-icon id="28:86" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" icon="lucide:search" class="text-base"></iconify-icon>
</div>
<input id="28:87" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" type="text" placeholder="搜索任务名称或目标..." class="bg-transparent grow shrink outline-none">
</div>
<div id="28:88" orderindex="81" class="flex items-center gap-y-3 gap-x-3">
<select id="28:89" orderindex="80" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(12px); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" class="flex items-center outline-none border-[1px] border-solid rounded-xl">
<option id="28:90" orderindex="80" value="all">全部状态</option>
<option id="28:91" orderindex="81" value="enabled">已启用</option>
<option id="28:92" orderindex="82" value="disabled">已禁用</option>
</select>
<select id="28:93" orderindex="81" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(12px); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" class="flex items-center outline-none border-[1px] border-solid rounded-xl">
<option id="28:94" orderindex="80" value="all">全部类型</option>
<option id="28:95" orderindex="81" value="orchestration">编排流程</option>
<option id="28:96" orderindex="82" value="skill">单个技能</option>
</select>
<select id="28:97" orderindex="82" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(12px); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" class="flex items-center outline-none border-[1px] border-solid rounded-xl">
<option id="28:98" orderindex="80" value="all">全部频率</option>
<option id="28:99" orderindex="81" value="hourly">每小时</option>
<option id="28:100" orderindex="82" value="daily">每日</option>
<option id="28:101" orderindex="83" value="weekly">每周</option>
<option id="28:102" orderindex="84" value="monthly">每月</option>
<option id="28:103" orderindex="85" value="custom">自定义</option>
</select>
</div>
</div>
</div>

<div id="28:104" orderindex="82" class="flex flex-col gap-y-6 gap-x-6">

<div id="28:105" orderindex="80" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="flex justify-between items-start p-6 border-[1px] border-solid rounded-2xl">
<div id="28:106" orderindex="80" style="flex-basis: 0%;" class="grow shrink">
<div id="28:107" orderindex="80" class="flex justify-between items-start mb-4">
<div id="28:108" orderindex="80" style="flex-basis: 0%;" class="grow shrink">
<div id="28:109" orderindex="80" class="flex items-center gap-y-3 gap-x-3 mb-2">
<h3 id="28:110" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-xl font-semibold">客户服务自动化流程</h3>
<div id="28:111" orderindex="81" style="background-color: color-mix( in oklab , #3B9BFF 20% , transparent ); padding: 0.25rem 0.75rem; border-color: color-mix( in oklab , #3B9BFF 40% , transparent );" class="border-[1px] border-solid rounded-full">
<span id="28:112" orderindex="80" style="color: rgba(95, 180, 255, 1);" class="text-xs">编排</span>
</div>
</div>
<p id="28:113" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm mb-3">自动处理客户咨询，进行情感分析和智能回复生成</p>
<div id="28:114" orderindex="82" class="flex items-center gap-y-4 gap-x-4 mb-3">
<span id="28:115" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm">目标: 客户服务编排流程</span>
<span id="28:116" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm">Cron: 0 */2 * * *</span>
</div>
</div>
<div id="28:117" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<div id="28:118" orderindex="80" style="background-color: color-mix( in oklab , oklch( 79.2% 0.209 151.711 ) 20% , transparent );" class="flex justify-center items-center w-8 h-8 rounded-full">
<div id="28:119" orderindex="80" style="background-color: oklch( 79.2% 0.209 151.711 ); box-shadow: 0 0 8px rgba(80, 200, 120, 0.6);" class="w-3 h-3 rounded-full"></div>
</div>
<span id="28:120" orderindex="81" style="color: oklch( 79.2% 0.209 151.711 );" class="text-sm font-medium">已启用</span>
</div>
</div>
<div id="28:121" orderindex="81" class="flex justify-between items-center">
<div id="28:122" orderindex="80" class="flex items-center gap-y-6 gap-x-6">
<div id="28:123" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm">
<span id="28:124" orderindex="80" class="font-medium">下次执行:</span>
<span id="28:125" orderindex="81" style="color: rgba(95, 180, 255, 1);">2小时32分钟后</span>
</div>
<div id="28:126" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm flex items-center gap-y-4 gap-x-4">
<span id="28:127" orderindex="80"><span id="28:128" orderindex="80">总执行:</span> <span id="28:129" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="font-medium">247次</span></span>
<span id="28:130" orderindex="81"><span id="28:131" orderindex="80">成功率:</span> <span id="28:132" orderindex="81" style="color: oklch( 79.2% 0.209 151.711 );" class="font-medium">94.7%</span></span>
<span id="28:133" orderindex="82"><span id="28:134" orderindex="80">最近状态:</span> <span id="28:135" orderindex="81" style="color: oklch( 79.2% 0.209 151.711 );" class="font-medium">成功</span></span>
</div>
</div>
<div id="28:136" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="28:137" orderindex="80" class="hover:border-[#5FB4FF]/50 hover:shadow-[0_0_12px_rgba(59,155,255,0.3)] flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(12px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<iconify-icon id="28:138" orderindex="80" style="color: rgba(95, 180, 255, 1);" icon="lucide:play" class="text-sm"></iconify-icon>
<span id="28:139" orderindex="81" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-sm">立即执行</span>
</button>
<button id="28:140" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="28:141" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:edit" class="text-sm"></iconify-icon>
</button>
<button id="28:142" orderindex="82" class="hover:bg-white/10 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="28:143" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:more-horizontal" class="text-sm"></iconify-icon>
</button>
</div>
</div>
</div>
</div>

<div id="28:144" orderindex="81" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="flex justify-between items-start p-6 border-[1px] border-solid rounded-2xl">
<div id="28:145" orderindex="80" style="flex-basis: 0%;" class="grow shrink">
<div id="28:146" orderindex="80" class="flex justify-between items-start mb-4">
<div id="28:147" orderindex="80" style="flex-basis: 0%;" class="grow shrink">
<div id="28:148" orderindex="80" class="flex items-center gap-y-3 gap-x-3 mb-2">
<h3 id="28:149" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-xl font-semibold">数据分析Agent技能执行</h3>
<div id="28:150" orderindex="81" style="background-color: color-mix( in oklab , oklch( 70.7% 0.165 254.624 ) 20% , transparent ); padding: 0.25rem 0.75rem; border-color: color-mix( in oklab , oklch( 70.7% 0.165 254.624 ) 40% , transparent );" class="border-[1px] border-solid rounded-full">
<span id="28:151" orderindex="80" style="color: oklch( 70.7% 0.165 254.624 );" class="text-xs">技能</span>
</div>
</div>
<p id="28:152" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm mb-3">定时执行数据分析Agent的报表生成技能</p>
<div id="28:153" orderindex="82" class="flex items-center gap-y-4 gap-x-4 mb-3">
<span id="28:154" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm">目标: DataAnalyst Agent - 生成报表技能</span>
<span id="28:155" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm">Cron: 0 8 * * 1</span>
</div>
</div>
<div id="28:156" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<div id="28:157" orderindex="80" style="background-color: color-mix( in oklab , oklch( 79.2% 0.209 151.711 ) 20% , transparent );" class="flex justify-center items-center w-8 h-8 rounded-full">
<div id="28:158" orderindex="80" style="background-color: oklch( 79.2% 0.209 151.711 ); box-shadow: 0 0 8px rgba(80, 200, 120, 0.6);" class="w-3 h-3 rounded-full"></div>
</div>
<span id="28:159" orderindex="81" style="color: oklch( 79.2% 0.209 151.711 );" class="text-sm font-medium">已启用</span>
</div>
</div>
<div id="28:160" orderindex="81" class="flex justify-between items-center">
<div id="28:161" orderindex="80" class="flex items-center gap-y-6 gap-x-6">
<div id="28:162" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm">
<span id="28:163" orderindex="80" class="font-medium">下次执行:</span>
<span id="28:164" orderindex="81" style="color: rgba(95, 180, 255, 1);">3天14小时后</span>
</div>
<div id="28:165" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm flex items-center gap-y-4 gap-x-4">
<span id="28:166" orderindex="80"><span id="28:167" orderindex="80">总执行:</span> <span id="28:168" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="font-medium">42次</span></span>
<span id="28:169" orderindex="81"><span id="28:170" orderindex="80">成功率:</span> <span id="28:171" orderindex="81" style="color: oklch( 79.2% 0.209 151.711 );" class="font-medium">97.6%</span></span>
<span id="28:172" orderindex="82"><span id="28:173" orderindex="80">最近状态:</span> <span id="28:174" orderindex="81" style="color: oklch( 79.2% 0.209 151.711 );" class="font-medium">成功</span></span>
</div>
</div>
<div id="28:175" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="28:176" orderindex="80" class="hover:border-[#5FB4FF]/50 hover:shadow-[0_0_12px_rgba(59,155,255,0.3)] flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(12px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<iconify-icon id="28:177" orderindex="80" style="color: rgba(95, 180, 255, 1);" icon="lucide:play" class="text-sm"></iconify-icon>
<span id="28:178" orderindex="81" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-sm">立即执行</span>
</button>
<button id="28:179" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="28:180" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:edit" class="text-sm"></iconify-icon>
</button>
<button id="28:181" orderindex="82" class="hover:bg-white/10 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="28:182" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:more-horizontal" class="text-sm"></iconify-icon>
</button>
</div>
</div>
</div>
</div>

<div id="28:183" orderindex="82" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="flex justify-between items-start p-6 border-[1px] border-solid rounded-2xl">
<div id="28:184" orderindex="80" style="flex-basis: 0%;" class="grow shrink">
<div id="28:185" orderindex="80" class="flex justify-between items-start mb-4">
<div id="28:186" orderindex="80" style="flex-basis: 0%;" class="grow shrink">
<div id="28:187" orderindex="80" class="flex items-center gap-y-3 gap-x-3 mb-2">
<h3 id="28:188" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-xl font-semibold">系统备份流程</h3>
<div id="28:189" orderindex="81" style="background-color: color-mix( in oklab , #3B9BFF 20% , transparent ); padding: 0.25rem 0.75rem; border-color: color-mix( in oklab , #3B9BFF 40% , transparent );" class="border-[1px] border-solid rounded-full">
<span id="28:190" orderindex="80" style="color: rgba(95, 180, 255, 1);" class="text-xs">编排</span>
</div>
</div>
<p id="28:191" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm mb-3">每日凌晨自动执行系统数据备份和清理任务</p>
<div id="28:192" orderindex="82" class="flex items-center gap-y-4 gap-x-4 mb-3">
<span id="28:193" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm">目标: 系统备份编排流程</span>
<span id="28:194" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm">Cron: 0 2 * * *</span>
</div>
</div>
<div id="28:195" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<div id="28:196" orderindex="80" style="background-color: color-mix( in oklab , #fff 10% , transparent );" class="flex justify-center items-center w-8 h-8 rounded-full">
<div id="28:197" orderindex="80" style="background-color: color-mix( in oklab , #fff 50% , transparent );" class="w-3 h-3 rounded-full"></div>
</div>
<span id="28:198" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm font-medium">已禁用</span>
</div>
</div>
<div id="28:199" orderindex="81" class="flex justify-between items-center">
<div id="28:200" orderindex="80" class="flex items-center gap-y-6 gap-x-6">
<div id="28:201" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm">
<span id="28:202" orderindex="80" class="font-medium">下次执行:</span>
<span id="28:203" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );">已暂停</span>
</div>
<div id="28:204" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm flex items-center gap-y-4 gap-x-4">
<span id="28:205" orderindex="80"><span id="28:206" orderindex="80">总执行:</span> <span id="28:207" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="font-medium">168次</span></span>
<span id="28:208" orderindex="81"><span id="28:209" orderindex="80">成功率:</span> <span id="28:210" orderindex="81" style="color: oklch( 70.4% 0.191 22.216 );" class="font-medium">89.3%</span></span>
<span id="28:211" orderindex="82"><span id="28:212" orderindex="80">最近状态:</span> <span id="28:213" orderindex="81" style="color: oklch( 70.4% 0.191 22.216 );" class="font-medium">失败</span></span>
</div>
</div>
<div id="28:214" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="28:215" orderindex="80" class="hover:border-[#5FB4FF]/50 hover:shadow-[0_0_12px_rgba(59,155,255,0.3)] opacity-50 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(12px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<iconify-icon id="28:216" orderindex="80" style="color: rgba(95, 180, 255, 1);" icon="lucide:play" class="text-sm"></iconify-icon>
<span id="28:217" orderindex="81" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-sm">立即执行</span>
</button>
<button id="28:218" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="28:219" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:edit" class="text-sm"></iconify-icon>
</button>
<button id="28:220" orderindex="82" class="hover:bg-white/10 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="28:221" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:more-horizontal" class="text-sm"></iconify-icon>
</button>
</div>
</div>
</div>
</div>
</div>

<div id="28:222" orderindex="83" style="backdrop-filter: blur(24px); top: 0; right: 0; transform: NaN; translate: 100% 0; background: linear-gradient(180deg, rgba(20, 30, 50, 0.95) 0%, rgba(25, 40, 65, 0.95) 100%);" class="fixed w-[480px] h-full z-50">
<div id="28:223" orderindex="80" class="flex flex-col h-full">

<div id="28:224" orderindex="80" style="border-bottom-style: solid; border-color: color-mix( in oklab , #fff 10% , transparent );" class="flex justify-between items-center border-b-[1px] p-6">
<h2 id="28:225" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-xl font-semibold">任务详情</h2>
<button id="28:226" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-10 h-10 rounded-xl">
<iconify-icon id="28:227" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:x" class="text-lg"></iconify-icon>
</button>
</div>

<div id="28:228" orderindex="81" style="border-bottom-style: solid; border-color: color-mix( in oklab , #fff 10% , transparent );" class="overflow-x-auto flex border-b-[1px]">
<button id="28:229" orderindex="80" style="color: rgba(95, 180, 255, 1); border-bottom-style: solid; padding: 0.75rem 1rem; border-color: rgba(95, 180, 255, 1);" class="text-sm whitespace-nowrap border-b-[2px] font-medium">基本信息</button>
<button id="28:230" orderindex="81" class="hover:bg-white/5 text-sm whitespace-nowrap font-medium" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">目标配置</button>
<button id="28:231" orderindex="82" class="hover:bg-white/5 text-sm whitespace-nowrap font-medium" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">调度配置</button>
<button id="28:232" orderindex="83" class="hover:bg-white/5 text-sm whitespace-nowrap font-medium" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">执行配置</button>
<button id="28:233" orderindex="84" class="hover:bg-white/5 text-sm whitespace-nowrap font-medium" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">通知配置</button>
<button id="28:234" orderindex="85" class="hover:bg-white/5 text-sm whitespace-nowrap font-medium" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">执行历史</button>
<button id="28:235" orderindex="86" class="hover:bg-white/5 text-sm whitespace-nowrap font-medium" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">统计分析</button>
</div>

<div id="28:236" orderindex="82" style="flex-basis: 0%;" class="overflow-y-auto grow shrink p-6">
<div id="28:237" orderindex="80">

<div id="28:238" orderindex="80">
<div id="28:239" orderindex="80" style="margin-top: 0;" class="mb-4">
<label id="28:240" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm block mb-2 font-medium">任务名称</label>
<input id="28:241" orderindex="81" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(12px); color: color-mix( in oklab , #fff 95% , transparent ); padding: 0.75rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" value="客户服务自动化流程" type="text" class="w-full outline-none border-[1px] border-solid rounded-xl">
</div>
<div id="28:242" orderindex="81" style="margin-top: 0;" class="mb-4">
<label id="28:243" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm block mb-2 font-medium">任务描述</label>
<textarea id="28:244" orderindex="81" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(12px); color: color-mix( in oklab , #fff 95% , transparent ); padding: 0.75rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" placeholder="描述任务的作用和目标..." rows="3" class="w-full outline-none border-[1px] border-solid rounded-xl">自动处理客户咨询，进行情感分析和智能回复生成</textarea>
</div>
<div id="28:245" orderindex="82" style="margin-top: 0;" class="grid grid-cols-2 gap-y-4 gap-x-4 mb-4">
<div id="28:246" orderindex="80">
<label id="28:247" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm block mb-2 font-medium">创建者</label>
<div id="28:248" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm">系统管理员</div>
</div>
<div id="28:249" orderindex="81">
<label id="28:250" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm block mb-2 font-medium">创建时间</label>
<div id="28:251" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-sm">2024-01-15 09:30</div>
</div>
</div>
<div id="28:252" orderindex="83">
<label id="28:253" orderindex="80" class="flex items-center gap-y-3 gap-x-3">
<div id="28:254" orderindex="80" class="relative">
<input id="28:255" orderindex="80" class="peer overflow-x-hidden overflow-y-hidden whitespace-nowrap absolute w-px h-px -m-px" style="clip-path: inset( 50% ); padding: 0; border-width: 0;" type="checkbox" :checked="true">
<div id="28:256" orderindex="81" style="background-color: color-mix( in oklab , #3B9BFF 25% , transparent ); backdrop-filter: blur(12px); border-color: color-mix( in oklab , #3B9BFF 60% , transparent );" class="w-14 h-7 border-[1px] border-solid rounded-full">
<div id="28:257" orderindex="80" style="background-color: rgba(95, 180, 255, 1); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 12px rgba(95, 180, 255, 0.6);" class="absolute top-1 left-8 w-5 h-5 rounded-full"></div>
</div>
</div>
<span id="28:258" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm font-medium">启用任务</span>
</label>
</div>
</div>
</div>
</div>

<div id="28:259" orderindex="83" style="border-top-style: solid; border-color: color-mix( in oklab , #fff 10% , transparent );" class="flex justify-between items-center border-t-[1px] p-6">
<div id="28:260" orderindex="80" class="flex gap-y-3 gap-x-3">
<button id="28:261" orderindex="80" class="hover:border-[#5FB4FF]/50 hover:shadow-[0_0_20px_rgba(59,155,255,0.3)] flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #3B9BFF 80% , transparent ); backdrop-filter: blur(24px); padding: 0.75rem 1.5rem; border-color: color-mix( in oklab , #3B9BFF 50% , transparent );">
<span id="28:262" orderindex="80" style="color: rgba(255, 255, 255, 1);" class="text-sm whitespace-nowrap">保存更改</span>
</button>
<button id="28:263" orderindex="81" class="hover:bg-white/10 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); padding: 0.75rem 1.5rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<span id="28:264" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm">取消</span>
</button>
</div>
<div id="28:265" orderindex="81" class="flex gap-y-2 gap-x-2">
<button id="28:266" orderindex="80" class="hover:border-[#5FB4FF]/50 hover:shadow-[0_0_12px_rgba(59,155,255,0.3)] flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(12px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<iconify-icon id="28:267" orderindex="80" style="color: rgba(95, 180, 255, 1);" icon="lucide:play" class="text-sm"></iconify-icon>
<span id="28:268" orderindex="81" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-sm">立即执行</span>
</button>
<button id="28:269" orderindex="81" class="hover:border-red-500/50 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , oklch( 63.7% 0.237 25.331 ) 10% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , oklch( 63.7% 0.237 25.331 ) 30% , transparent );">
<iconify-icon id="28:270" orderindex="80" style="color: oklch( 70.4% 0.191 22.216 );" icon="lucide:trash-2" class="text-sm"></iconify-icon>
<span id="28:271" orderindex="81" style="color: oklch( 70.4% 0.191 22.216 );" class="text-sm">删除</span>
</button>
</div>
</div>
</div>
</div>
</main>

  </BaseLayout>
</template>

<style scoped>
/* Add any view-specific styles here */
</style>
