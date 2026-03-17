<script setup lang="ts">
import BaseLayout from '../components/BaseLayout.vue'
</script>

<template>
  <BaseLayout>
    <main id="26:65" orderindex="81" style="flex-basis: 0%;" class="overflow-x-hidden flex flex-col grow shrink gap-y-8 gap-x-8 p-8">

<div id="26:66" orderindex="80" class="flex justify-between items-center">
<h1 id="26:67" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-2xl font-semibold">Model Management</h1>
<div id="26:68" orderindex="81" class="flex items-center gap-y-4 gap-x-4">
<button id="26:69" orderindex="80" class="hover:bg-[#2A7FDB] flex items-center gap-y-2 gap-x-2 rounded-xl" style="background-color: rgba(59, 155, 255, 1); box-shadow: 0 0 20px rgba(59, 155, 255, 0.4); padding: 0.5rem 1.5rem;">
<div id="26:70" orderindex="80" class="flex justify-center items-center w-4 h-4">
<iconify-icon id="26:71" orderindex="80" style="color: rgba(255, 255, 255, 1);" icon="lucide:plus" class="text-sm"></iconify-icon>
</div>
<span id="26:72" orderindex="81" style="color: rgba(255, 255, 255, 1);" class="text-sm whitespace-nowrap">添加新模型</span>
</button>
<button id="26:73" orderindex="81" class="hover:border-[#5FB4FF]/50 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:74" orderindex="80" class="flex justify-center items-center w-4 h-4">
<iconify-icon id="26:75" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:download" class="text-sm"></iconify-icon>
</div>
<span id="26:76" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm whitespace-nowrap">导入配置</span>
</button>
<button id="26:77" orderindex="82" class="hover:border-[#5FB4FF]/50 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:78" orderindex="80" class="flex justify-center items-center w-4 h-4">
<iconify-icon id="26:79" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:upload" class="text-sm"></iconify-icon>
</div>
<span id="26:80" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm whitespace-nowrap">导出配置</span>
</button>
<div id="26:81" orderindex="83" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl">
<div id="26:82" orderindex="80" class="flex justify-center items-center w-4 h-4">
<iconify-icon id="26:83" orderindex="80" style="color: rgba(95, 180, 255, 1);" icon="lucide:grid-3x3" class="text-sm"></iconify-icon>
</div>
<div id="26:84" orderindex="81" class="flex justify-center items-center w-4 h-4">
<iconify-icon id="26:85" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" icon="lucide:list" class="text-sm"></iconify-icon>
</div>
</div>
</div>
</div>

<div id="26:86" orderindex="81" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="p-6 border-[1px] border-solid rounded-2xl">
<div id="26:87" orderindex="80" class="flex items-center gap-y-6 gap-x-6 mb-4">
<div id="26:88" orderindex="80" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(16px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" class="flex items-center min-w-[300px] border-[1px] border-solid rounded-xl">
<div id="26:89" orderindex="80" class="flex justify-center items-center w-5 h-5 mr-3">
<iconify-icon id="26:90" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" icon="lucide:search" class="text-base"></iconify-icon>
</div>
<input id="26:91" orderindex="81" style="flex-basis: 0%; color: color-mix( in oklab , #fff 70% , transparent );" type="text" placeholder="按模型名称/提供商搜索..." class="bg-transparent grow shrink outline-none">
</div>
<div id="26:92" orderindex="81" class="flex items-center gap-y-4 gap-x-4">
<select id="26:93" orderindex="80" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(16px); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" class="outline-none border-[1px] border-solid rounded-xl">
<option id="26:94" orderindex="80">全部提供商</option>
<option id="26:95" orderindex="81">OpenAI</option>
<option id="26:96" orderindex="82">Anthropic</option>
<option id="26:97" orderindex="83">本地模型</option>
<option id="26:98" orderindex="84">其他</option>
</select>
<select id="26:99" orderindex="81" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(16px); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" class="outline-none border-[1px] border-solid rounded-xl">
<option id="26:100" orderindex="80">全部状态</option>
<option id="26:101" orderindex="81">已激活</option>
<option id="26:102" orderindex="82">配置中</option>
<option id="26:103" orderindex="83">不可用</option>
</select>
<select id="26:104" orderindex="82" style="background-color: color-mix( in oklab , #fff 3% , transparent ); backdrop-filter: blur(16px); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" class="outline-none border-[1px] border-solid rounded-xl">
<option id="26:105" orderindex="80">全部类型</option>
<option id="26:106" orderindex="81">对话</option>
<option id="26:107" orderindex="82">嵌入</option>
<option id="26:108" orderindex="83">图像生成</option>
</select>
</div>
</div>
<div id="26:109" orderindex="81" class="overflow-x-auto flex gap-y-2 gap-x-2">
<label id="26:110" orderindex="80">
<input id="26:111" orderindex="80" class="peer overflow-x-hidden overflow-y-hidden whitespace-nowrap absolute w-px h-px -m-px" style="clip-path: inset( 50% ); padding: 0; border-width: 0;" type="radio" checked="">
<div id="26:112" orderindex="81" class="hover:opacity-80 whitespace-nowrap border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #3B9BFF 20% , transparent ); backdrop-filter: blur(12px); color: color-mix( in oklab , #fff 95% , transparent ); padding: 0.25rem 0.75rem; border-color: color-mix( in oklab , #3B9BFF 50% , transparent );">
<span id="26:113" orderindex="80" class="text-sm">添加时间</span>
</div>
</label>
<label id="26:114" orderindex="81">
<input id="26:115" orderindex="80" class="peer overflow-x-hidden overflow-y-hidden whitespace-nowrap absolute w-px h-px -m-px" style="clip-path: inset( 50% ); padding: 0; border-width: 0;" type="radio">
<div id="26:116" orderindex="81" class="hover:opacity-80 whitespace-nowrap border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(12px); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.25rem 0.75rem; border-color: color-mix( in oklab , #fff 10% , transparent );">
<span id="26:117" orderindex="80" class="text-sm">使用频率</span>
</div>
</label>
<label id="26:118" orderindex="82">
<input id="26:119" orderindex="80" class="peer overflow-x-hidden overflow-y-hidden whitespace-nowrap absolute w-px h-px -m-px" style="clip-path: inset( 50% ); padding: 0; border-width: 0;" type="radio">
<div id="26:120" orderindex="81" class="hover:opacity-80 whitespace-nowrap border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(12px); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.25rem 0.75rem; border-color: color-mix( in oklab , #fff 10% , transparent );">
<span id="26:121" orderindex="80" class="text-sm">名称</span>
</div>
</label>
</div>
</div>

<div id="26:122" orderindex="82" class="grid grid-cols-3 gap-y-6 gap-x-6">

<div id="26:123" orderindex="80" class="hover:shadow-[0_12px_32px_rgba(0,0,0,0.5),0_0_20px_rgba(59,155,255,0.25)] p-6 border-[1px] border-solid rounded-2xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:124" orderindex="80" class="flex justify-between items-start mb-4">
<div id="26:125" orderindex="80" class="flex items-center gap-y-3 gap-x-3">
<div id="26:126" orderindex="80" style="background-color: color-mix( in oklab , #3B9BFF 20% , transparent );" class="flex justify-center items-center w-12 h-12 rounded-xl">
<iconify-icon id="26:127" orderindex="80" icon="logos:openai-icon" class="text-xl"></iconify-icon>
</div>
<div id="26:128" orderindex="81">
<h3 id="26:129" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-lg font-semibold">GPT-4</h3>
<div id="26:130" orderindex="81" class="flex items-center gap-y-2 gap-x-2 mt-1">
<span id="26:131" orderindex="80" style="background-color: rgba(80, 200, 120, 0.15); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.25rem 0.5rem;" class="text-sm rounded-full">可用</span>
<span id="26:132" orderindex="81" style="background-color: color-mix( in oklab , #fff 10% , transparent ); color: color-mix( in oklab , #fff 50% , transparent ); padding: 0.25rem 0.5rem;" class="text-xs rounded-full">对话</span>
</div>
</div>
</div>
<button id="26:133" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:134" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:more-horizontal" class="text-sm"></iconify-icon>
</button>
</div>
<div id="26:135" orderindex="81" class="mb-4">
<div id="26:136" orderindex="80" class="flex justify-between mb-2">
<span id="26:137" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">本月调用</span>
<span id="26:138" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">15,432</span>
</div>
<div id="26:139" orderindex="81" class="flex justify-between mb-2">
<span id="26:140" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">平均响应时间</span>
<span id="26:141" orderindex="81" style="color: rgba(95, 180, 255, 1);" class="text-xs font-semibold">1.2s</span>
</div>
<div id="26:142" orderindex="82" class="flex justify-between">
<span id="26:143" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">关联Agent</span>
<span id="26:144" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">8个</span>
</div>
</div>
<div id="26:145" orderindex="82" class="flex justify-between items-center">
<span id="26:146" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">最后更新: 2小时前</span>
<div id="26:147" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="26:148" orderindex="80" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:149" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:settings" class="text-sm"></iconify-icon>
</button>
<button id="26:150" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:151" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:play" class="text-sm"></iconify-icon>
</button>
</div>
</div>
</div>

<div id="26:152" orderindex="81" class="hover:shadow-[0_12px_32px_rgba(0,0,0,0.5),0_0_20px_rgba(59,155,255,0.25)] p-6 border-[1px] border-solid rounded-2xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:153" orderindex="80" class="flex justify-between items-start mb-4">
<div id="26:154" orderindex="80" class="flex items-center gap-y-3 gap-x-3">
<div id="26:155" orderindex="80" style="background-color: color-mix( in oklab , #87B4FF 20% , transparent );" class="flex justify-center items-center w-12 h-12 rounded-xl">
<iconify-icon id="26:156" orderindex="80" style="color: rgba(135, 180, 255, 1);" icon="simple-icons:anthropic" class="text-xl"></iconify-icon>
</div>
<div id="26:157" orderindex="81">
<h3 id="26:158" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-lg font-semibold">Claude 3 Sonnet</h3>
<div id="26:159" orderindex="81" class="flex items-center gap-y-2 gap-x-2 mt-1">
<span id="26:160" orderindex="80" style="background-color: rgba(80, 200, 120, 0.15); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.25rem 0.5rem;" class="text-sm rounded-full">可用</span>
<span id="26:161" orderindex="81" style="background-color: color-mix( in oklab , #fff 10% , transparent ); color: color-mix( in oklab , #fff 50% , transparent ); padding: 0.25rem 0.5rem;" class="text-xs rounded-full">对话</span>
</div>
</div>
</div>
<button id="26:162" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:163" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:more-horizontal" class="text-sm"></iconify-icon>
</button>
</div>
<div id="26:164" orderindex="81" class="mb-4">
<div id="26:165" orderindex="80" class="flex justify-between mb-2">
<span id="26:166" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">本月调用</span>
<span id="26:167" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">8,765</span>
</div>
<div id="26:168" orderindex="81" class="flex justify-between mb-2">
<span id="26:169" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">平均响应时间</span>
<span id="26:170" orderindex="81" style="color: rgba(135, 180, 255, 1);" class="text-xs font-semibold">2.1s</span>
</div>
<div id="26:171" orderindex="82" class="flex justify-between">
<span id="26:172" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">关联Agent</span>
<span id="26:173" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">5个</span>
</div>
</div>
<div id="26:174" orderindex="82" class="flex justify-between items-center">
<span id="26:175" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">最后更新: 1天前</span>
<div id="26:176" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="26:177" orderindex="80" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:178" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:settings" class="text-sm"></iconify-icon>
</button>
<button id="26:179" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:180" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:play" class="text-sm"></iconify-icon>
</button>
</div>
</div>
</div>

<div id="26:181" orderindex="82" class="hover:shadow-[0_12px_32px_rgba(0,0,0,0.5),0_0_20px_rgba(59,155,255,0.25)] p-6 border-[1px] border-solid rounded-2xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:182" orderindex="80" class="flex justify-between items-start mb-4">
<div id="26:183" orderindex="80" class="flex items-center gap-y-3 gap-x-3">
<div id="26:184" orderindex="80" style="background-color: color-mix( in oklab , #00E5FF 20% , transparent );" class="flex justify-center items-center w-12 h-12 rounded-xl">
<iconify-icon id="26:185" orderindex="80" style="color: rgba(0, 229, 255, 1);" icon="lucide:server" class="text-xl"></iconify-icon>
</div>
<div id="26:186" orderindex="81">
<h3 id="26:187" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-lg font-semibold">Llama-2-7B</h3>
<div id="26:188" orderindex="81" class="flex items-center gap-y-2 gap-x-2 mt-1">
<span id="26:189" orderindex="80" style="background-color: rgba(255, 180, 70, 0.15); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.25rem 0.5rem;" class="text-sm rounded-full">配置中</span>
<span id="26:190" orderindex="81" style="background-color: color-mix( in oklab , #fff 10% , transparent ); padding-top: <; padding-right: span; padding-bottom: class; padding-left: =; color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs rounded-full">本地</span>
</div>
</div>
</div>
<button id="26:191" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:192" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:more-horizontal" class="text-sm"></iconify-icon>
</button>
</div>
<div id="26:193" orderindex="81" class="mb-4">
<div id="26:194" orderindex="80" class="flex justify-between mb-2">
<span id="26:195" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">本月调用</span>
<span id="26:196" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">1,234</span>
</div>
<div id="26:197" orderindex="81" class="flex justify-between mb-2">
<span id="26:198" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">平均响应时间</span>
<span id="26:199" orderindex="81" style="color: rgba(0, 229, 255, 1);" class="text-xs font-semibold">3.5s</span>
</div>
<div id="26:200" orderindex="82" class="flex justify-between">
<span id="26:201" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">关联Agent</span>
<span id="26:202" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">2个</span>
</div>
</div>
<div id="26:203" orderindex="82" class="flex justify-between items-center">
<span id="26:204" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">最后更新: 3天前</span>
<div id="26:205" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="26:206" orderindex="80" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:207" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:settings" class="text-sm"></iconify-icon>
</button>
<button id="26:208" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:209" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:play" class="text-sm"></iconify-icon>
</button>
</div>
</div>
</div>

<div id="26:210" orderindex="83" class="hover:shadow-[0_12px_32px_rgba(0,0,0,0.5),0_0_20px_rgba(59,155,255,0.25)] p-6 border-[1px] border-solid rounded-2xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:211" orderindex="80" class="flex justify-between items-start mb-4">
<div id="26:212" orderindex="80" class="flex items-center gap-y-3 gap-x-3">
<div id="26:213" orderindex="80" style="background-color: color-mix( in oklab , #5FB4FF 20% , transparent );" class="flex justify-center items-center w-12 h-12 rounded-xl">
<iconify-icon id="26:214" orderindex="80" style="color: rgba(95, 180, 255, 1);" icon="lucide:image" class="text-xl"></iconify-icon>
</div>
<div id="26:215" orderindex="81">
<h3 id="26:216" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-lg font-semibold">DALL-E 3</h3>
<div id="26:217" orderindex="81" class="flex items-center gap-y-2 gap-x-2 mt-1">
<span id="26:218" orderindex="80" style="background-color: rgba(80, 200, 120, 0.15); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.25rem 0.5rem;" class="text-sm rounded-full">可用</span>
<span id="26:219" orderindex="81" style="background-color: color-mix( in oklab , #fff 10% , transparent ); color: color-mix( in oklab , #fff 50% , transparent ); padding: 0.25rem 0.5rem;" class="text-xs rounded-full">图像生成</span>
</div>
</div>
</div>
<button id="26:220" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:221" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:more-horizontal" class="text-sm"></iconify-icon>
</button>
</div>
<div id="26:222" orderindex="81" class="mb-4">
<div id="26:223" orderindex="80" class="flex justify-between mb-2">
<span id="26:224" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">本月调用</span>
<span id="26:225" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">2,456</span>
</div>
<div id="26:226" orderindex="81" class="flex justify-between mb-2">
<span id="26:227" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">平均响应时间</span>
<span id="26:228" orderindex="81" style="color: rgba(95, 180, 255, 1);" class="text-xs font-semibold">8.7s</span>
</div>
<div id="26:229" orderindex="82" class="flex justify-between">
<span id="26:230" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">关联Agent</span>
<span id="26:231" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">3个</span>
</div>
</div>
<div id="26:232" orderindex="82" class="flex justify-between items-center">
<span id="26:233" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">最后更新: 6小时前</span>
<div id="26:234" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="26:235" orderindex="80" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:236" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:settings" class="text-sm"></iconify-icon>
</button>
<button id="26:237" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:238" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:play" class="text-sm"></iconify-icon>
</button>
</div>
</div>
</div>

<div id="26:239" orderindex="84" class="hover:shadow-[0_12px_32px_rgba(0,0,0,0.5),0_0_20px_rgba(59,155,255,0.25)] p-6 border-[1px] border-solid rounded-2xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:240" orderindex="80" class="flex justify-between items-start mb-4">
<div id="26:241" orderindex="80" class="flex items-center gap-y-3 gap-x-3">
<div id="26:242" orderindex="80" style="background-color: color-mix( in oklab , #87B4FF 20% , transparent );" class="flex justify-center items-center w-12 h-12 rounded-xl">
<iconify-icon id="26:243" orderindex="80" style="color: rgba(135, 180, 255, 1);" icon="lucide:layers" class="text-xl"></iconify-icon>
</div>
<div id="26:244" orderindex="81">
<h3 id="26:245" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-lg font-semibold">text-embedding-ada-002</h3>
<div id="26:246" orderindex="81" class="flex items-center gap-y-2 gap-x-2 mt-1">
<span id="26:247" orderindex="80" style="background-color: rgba(80, 200, 120, 0.15); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.25rem 0.5rem;" class="text-sm rounded-full">可用</span>
<span id="26:248" orderindex="81" style="background-color: color-mix( in oklab , #fff 10% , transparent ); color: color-mix( in oklab , #fff 50% , transparent ); padding: 0.25rem 0.5rem;" class="text-xs rounded-full">嵌入</span>
</div>
</div>
</div>
<button id="26:249" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:250" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:more-horizontal" class="text-sm"></iconify-icon>
</button>
</div>
<div id="26:251" orderindex="81" class="mb-4">
<div id="26:252" orderindex="80" class="flex justify-between mb-2">
<span id="26:253" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">本月调用</span>
<span id="26:254" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">45,678</span>
</div>
<div id="26:255" orderindex="81" class="flex justify-between mb-2">
<span id="26:256" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">平均响应时间</span>
<span id="26:257" orderindex="81" style="color: rgba(135, 180, 255, 1);" class="text-xs font-semibold">0.3s</span>
</div>
<div id="26:258" orderindex="82" class="flex justify-between">
<span id="26:259" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">关联Agent</span>
<span id="26:260" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">12个</span>
</div>
</div>
<div id="26:261" orderindex="82" class="flex justify-between items-center">
<span id="26:262" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">最后更新: 1小时前</span>
<div id="26:263" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="26:264" orderindex="80" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:265" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:settings" class="text-sm"></iconify-icon>
</button>
<button id="26:266" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:267" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:play" class="text-sm"></iconify-icon>
</button>
</div>
</div>
</div>

<div id="26:268" orderindex="85" class="hover:shadow-[0_12px_32px_rgba(0,0,0,0.5),0_0_20px_rgba(59,155,255,0.25)] p-6 border-[1px] border-solid rounded-2xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 12px rgba(59, 155, 255, 0.15); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:269" orderindex="80" class="flex justify-between items-start mb-4">
<div id="26:270" orderindex="80" class="flex items-center gap-y-3 gap-x-3">
<div id="26:271" orderindex="80" style="background-color: color-mix( in oklab , #fff 10% , transparent );" class="flex justify-center items-center w-12 h-12 rounded-xl">
<iconify-icon id="26:272" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" icon="lucide:alert-triangle" class="text-xl"></iconify-icon>
</div>
<div id="26:273" orderindex="81">
<h3 id="26:274" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-lg font-semibold">GPT-3.5-Turbo</h3>
<div id="26:275" orderindex="81" class="flex items-center gap-y-2 gap-x-2 mt-1">
<span id="26:276" orderindex="80" style="background-color: rgba(255, 100, 100, 0.15); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.25rem 0.5rem;" class="text-sm rounded-full">不可用</span>
<span id="26:277" orderindex="81" style="background-color: color-mix( in oklab , #fff 10% , transparent ); color: color-mix( in oklab , #fff 50% , transparent ); padding: 0.25rem 0.5rem;" class="text-xs rounded-full">对话</span>
</div>
</div>
</div>
<button id="26:278" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:279" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:more-horizontal" class="text-sm"></iconify-icon>
</button>
</div>
<div id="26:280" orderindex="81" class="mb-4">
<div id="26:281" orderindex="80" class="flex justify-between mb-2">
<span id="26:282" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">本月调用</span>
<span id="26:283" orderindex="81" style="color: color-mix( in oklab , #fff 30% , transparent );" class="text-xs font-semibold">--</span>
</div>
<div id="26:284" orderindex="81" class="flex justify-between mb-2">
<span id="26:285" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">平均响应时间</span>
<span id="26:286" orderindex="81" style="color: color-mix( in oklab , #fff 30% , transparent );" class="text-xs font-semibold">--</span>
</div>
<div id="26:287" orderindex="82" class="flex justify-between">
<span id="26:288" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">关联Agent</span>
<span id="26:289" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-xs font-semibold">0个</span>
</div>
</div>
<div id="26:290" orderindex="82" class="flex justify-between items-center">
<span id="26:291" orderindex="80" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">最后更新: 1周前</span>
<div id="26:292" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<button id="26:293" orderindex="80" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:294" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:settings" class="text-sm"></iconify-icon>
</button>
<button id="26:295" orderindex="81" class="hover:bg-white/10 opacity-50 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );" disabled="">
<iconify-icon id="26:296" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:play" class="text-sm"></iconify-icon>
</button>
</div>
</div>
</div>
</div>
</main>
<aside id="26:297" orderindex="82" style="background-color: color-mix( in oklab , #fff 8% , transparent ); backdrop-filter: blur(40px); top: 0; right: 0; translate: 100% 0; box-shadow: 0 0 48px rgba(0, 0, 0, 0.6); border-left-style: solid; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="fixed w-[480px] h-full z-50 border-l-[1px]">
<div id="26:298" orderindex="80" class="flex flex-col h-full">

<div id="26:299" orderindex="80" style="border-bottom-style: solid; border-color: color-mix( in oklab , #fff 10% , transparent );" class="flex justify-between items-center border-b-[1px] p-6">
<h2 id="26:300" orderindex="80" style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-xl font-semibold">模型详情</h2>
<button id="26:301" orderindex="81" class="hover:bg-white/10 flex justify-center items-center w-8 h-8 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #fff 5% , transparent ); border-color: color-mix( in oklab , #3B9BFF 20% , transparent );">
<iconify-icon id="26:302" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:x" class="text-sm"></iconify-icon>
</button>
</div>

<div id="26:303" orderindex="81" style="border-bottom-style: solid; border-color: color-mix( in oklab , #fff 10% , transparent );" class="overflow-x-auto flex border-b-[1px]">
<button id="26:304" orderindex="80" style="background-color: color-mix( in oklab , #3B9BFF 10% , transparent ); color: rgba(95, 180, 255, 1); border-bottom-style: solid; padding: 0.75rem 1rem; border-color: rgba(95, 180, 255, 1);" class="text-sm shrink-0 border-b-[2px]">基本信息</button>
<button id="26:305" orderindex="81" class="hover:text-white/95 hover:bg-white/5 text-sm shrink-0" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">连接配置</button>
<button id="26:306" orderindex="82" class="hover:text-white/95 hover:bg-white/5 text-sm shrink-0" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">模型参数</button>
<button id="26:307" orderindex="83" class="hover:text-white/95 hover:bg-white/5 text-sm shrink-0" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">成本配置</button>
<button id="26:308" orderindex="84" class="hover:text-white/95 hover:bg-white/5 text-sm shrink-0" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">性能监控</button>
<button id="26:309" orderindex="85" class="hover:text-white/95 hover:bg-white/5 text-sm shrink-0" style="color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem;">使用情况</button>
</div>

<div id="26:310" orderindex="82" style="flex-basis: 0%;" class="overflow-y-auto grow shrink p-6">

<div id="26:311" orderindex="80">
<div id="26:312" orderindex="80" class="flex items-center gap-y-4 gap-x-4 mb-6">
<div id="26:313" orderindex="80" style="background-color: color-mix( in oklab , #3B9BFF 20% , transparent );" class="flex justify-center items-center w-16 h-16 rounded-xl">
<iconify-icon id="26:314" orderindex="80" icon="logos:openai-icon" class="text-2xl"></iconify-icon>
</div>
<div id="26:315" orderindex="81" style="flex-basis: 0%;" class="grow shrink">
<input id="26:316" orderindex="80" class="focus:border-[#3B9BFF]/50 text-lg w-full outline-none font-semibold border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(16px); color: color-mix( in oklab , #fff 95% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" value="GPT-4" type="text">
<div id="26:317" orderindex="81" class="flex items-center gap-y-2 gap-x-2 mt-2">
<span id="26:318" orderindex="80" style="background-color: rgba(80, 200, 120, 0.15); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.25rem 0.5rem;" class="text-sm rounded-full">可用</span>
<span id="26:319" orderindex="81" style="color: color-mix( in oklab , #fff 50% , transparent );" class="text-xs">v4.0</span>
</div>
</div>
</div>
<div id="26:320" orderindex="81" class="mb-6">
<label id="26:321" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm block mb-2 font-medium">模型描述</label>
<textarea id="26:322" orderindex="81" class="focus:border-[#3B9BFF]/50 text-sm w-full h-32 outline-none border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(16px); color: color-mix( in oklab , #fff 70% , transparent ); padding: 0.75rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );" placeholder="请输入模型的详细描述...">OpenAI最新的大型语言模型，具有强大的理解和生成能力，支持多轮对话、代码生成、文本分析等多种任务场景。</textarea>
</div>
<div id="26:323" orderindex="82" class="grid grid-cols-2 gap-y-4 gap-x-4 mb-6">
<div id="26:324" orderindex="80">
<label id="26:325" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm block mb-2 font-medium">提供商</label>
<div id="26:326" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<iconify-icon id="26:327" orderindex="80" icon="logos:openai-icon" class="text-base"></iconify-icon>
<span id="26:328" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm">OpenAI</span>
</div>
</div>
<div id="26:329" orderindex="81">
<label id="26:330" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm block mb-2 font-medium">模型类型</label>
<span id="26:331" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm">对话模型</span>
</div>
</div>
<div id="26:332" orderindex="83" class="grid grid-cols-2 gap-y-4 gap-x-4 mb-6">
<div id="26:333" orderindex="80">
<label id="26:334" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm block mb-2 font-medium">添加者</label>
<div id="26:335" orderindex="81" class="flex items-center gap-y-2 gap-x-2">
<img id="26:336" orderindex="80" style="border-color: color-mix( in oklab , #3B9BFF 40% , transparent );" alt="Model creator profile photo" src="/assets/img/0235382d-15f6-4b07-b006-ec8d97c40e68.jpeg" class="w-8 h-8 border-[1px] border-solid rounded-full">
<span id="26:337" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm">管理员</span>
</div>
</div>
<div id="26:338" orderindex="81">
<label id="26:339" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm block mb-2 font-medium">添加时间</label>
<span id="26:340" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm">2024-01-15 10:30</span>
</div>
</div>
<div id="26:341" orderindex="84" class="mb-6">
<div id="26:342" orderindex="80" class="flex justify-between items-center mb-2">
<label id="26:343" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm font-medium">状态</label>
<div id="26:344" orderindex="81" class="relative">
<input id="26:345" orderindex="80" class="peer overflow-x-hidden overflow-y-hidden whitespace-nowrap absolute w-px h-px -m-px" style="clip-path: inset( 50% ); padding: 0; border-width: 0;" type="checkbox" checked="">
<div id="26:346" orderindex="81" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(12px); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="w-14 h-7 border-[1px] border-solid rounded-full"></div>
<div id="26:347" orderindex="82" style="background-color: rgba(95, 180, 255, 1);" class="absolute top-1 left-1 w-5 h-5 rounded-full"></div>
</div>
</div>
</div>
</div>
</div>

<div id="26:348" orderindex="83" style="border-top-style: solid; border-color: color-mix( in oklab , #fff 10% , transparent );" class="flex justify-between items-center border-t-[1px] p-6">
<div id="26:349" orderindex="80" class="flex items-center gap-y-3 gap-x-3">
<button id="26:350" orderindex="80" class="hover:border-[#5FB4FF]/50 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:351" orderindex="80" class="flex justify-center items-center w-4 h-4">
<iconify-icon id="26:352" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:play" class="text-sm"></iconify-icon>
</div>
<span id="26:353" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm whitespace-nowrap">测试连接</span>
</button>
<button id="26:354" orderindex="81" class="hover:border-[#5FB4FF]/50 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );">
<div id="26:355" orderindex="80" class="flex justify-center items-center w-4 h-4">
<iconify-icon id="26:356" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" icon="lucide:book-open" class="text-sm"></iconify-icon>
</div>
<span id="26:357" orderindex="81" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm whitespace-nowrap">API文档</span>
</button>
</div>
<div id="26:358" orderindex="81" class="flex items-center gap-y-3 gap-x-3">
<button id="26:359" orderindex="80" class="hover:border-white/20 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #fff 10% , transparent );">
<span id="26:360" orderindex="80" style="color: color-mix( in oklab , #fff 70% , transparent );" class="text-sm whitespace-nowrap">取消</span>
</button>
<button id="26:361" orderindex="81" class="hover:bg-[#2A7FDB] flex items-center gap-y-2 gap-x-2 rounded-xl" style="background-color: rgba(59, 155, 255, 1); box-shadow: 0 0 20px rgba(59, 155, 255, 0.4); padding: 0.5rem 1.5rem;">
<span id="26:362" orderindex="80" style="color: rgba(255, 255, 255, 1);" class="text-sm whitespace-nowrap">保存配置</span>
</button>
</div>
</div>
</div>
</aside>
  </BaseLayout>
</template>

<style scoped>
/* Add any view-specific styles here */
</style>
