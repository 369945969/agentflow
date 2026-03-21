import { type Component, createSignal, onCleanup, onMount } from "solid-js"

interface ProjectListProps {
  onNavigate: (page: string) => void;
}

const ProjectList: Component<ProjectListProps> = (props) => {
  const [theme, setTheme] = createSignal<"dark" | "light">("dark")

  onMount(() => {
    try {
      const stored = globalThis.localStorage?.getItem?.("workspace_theme") || ""
      if (stored === "light" || stored === "dark") setTheme(stored)
    } catch {}
    const handler = (event: Event) => {
      const detail = (event as CustomEvent).detail as "dark" | "light" | undefined
      if (detail !== "dark" && detail !== "light") return
      setTheme(detail)
    }
    window.addEventListener("workspace_theme_toggle", handler as EventListener)
    onCleanup(() => window.removeEventListener("workspace_theme_toggle", handler as EventListener))
  })

  return (
    <main 
      id="13:31" 
      style="flex-basis: 0%; padding: 1.5rem 2rem;" 
      class="overflow-x-hidden overflow-y-auto flex flex-col grow shrink"
    >
      <div id="13:32" class="flex justify-between items-center mb-8 shrink-0">
        <div id="13:33" class="flex items-center gap-y-4 gap-x-4">
          <h1
            id="13:34"
            style={theme() === "light" ? "color: #0F172A;" : "color: rgba(232, 240, 255, 1);"}
            class="text-2xl font-bold"
          >
            Skills 列表
          </h1>
          <div
            id="13:35"
            class={
              theme() === "light"
                ? "text-xs rounded-lg bg-sky-50 text-sky-700 border border-sky-200 px-2 py-1"
                : "text-xs rounded-lg"
            }
            style={
              theme() === "light"
                ? ""
                : "background-color: color-mix( in oklab , #00F0FF 15% , transparent ); color: rgba(0, 240, 255, 1); padding: 0.25rem 0.5rem;"
            }
          >
            12个 Skills
          </div>
        </div>
        <div id="13:36" class="flex items-center gap-y-4 gap-x-4">
          
          <button id="13:37" class="hover:shadow-[0_0_20px_rgba(0,240,255,0.4)] flex items-center gap-y-2 gap-x-2 rounded-lg" style="background-color: rgba(0, 240, 255, 1); color: rgba(10, 14, 26, 1); padding: 0.75rem 1.5rem;">
            <div id="13:38" class="bg-transparent flex justify-center items-center w-5 h-5">
              <iconify-icon id="13:39" icon="lucide:plus" class="text-base"></iconify-icon>
            </div>
            <span id="13:40" class="whitespace-nowrap font-semibold">新建 Skills</span>
          </button>
          
          <div id="13:41" class="flex items-center gap-y-2 gap-x-2">
            <button
              id="13:42"
              class={
                theme() === "light"
                  ? "flex justify-center items-center w-10 h-10 rounded-lg border border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-colors"
                  : "hover:bg-[#00F0FF]/15 hover:shadow-[0_0_12px_rgba(0,240,255,0.25)] flex justify-center items-center w-10 h-10 rounded-lg"
              }
              style={
                theme() === "light"
                  ? ""
                  : "background-color: color-mix( in oklab , #1A1F3A 90% , transparent );"
              }
            >
              <div id="13:43" class="bg-transparent flex justify-center items-center w-5 h-5">
                <iconify-icon id="13:44" style="color: rgba(0, 240, 255, 1);" icon="lucide:layout-grid" class="text-base"></iconify-icon>
              </div>
            </button>
            <button
              id="13:45"
              class={
                theme() === "light"
                  ? "flex justify-center items-center w-10 h-10 rounded-lg border border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-colors"
                  : "hover:bg-[#00F0FF]/15 hover:shadow-[0_0_12px_rgba(0,240,255,0.25)] flex justify-center items-center w-10 h-10 rounded-lg"
              }
              style={
                theme() === "light"
                  ? ""
                  : "background-color: color-mix( in oklab , #1A1F3A 90% , transparent );"
              }
            >
              <div id="13:46" class="bg-transparent flex justify-center items-center w-5 h-5">
                <iconify-icon id="13:47" style="color: rgba(184, 197, 217, 1);" icon="lucide:list" class="text-base"></iconify-icon>
              </div>
            </button>
          </div>
        </div>
      </div>
      
      <div id="13:48" class="flex justify-between items-center mb-6 shrink-0">
        
        <div id="13:49" class="overflow-x-auto flex gap-y-2 gap-x-2">
          <label id="13:50" class="whitespace-nowrap">
            <input id="13:51" class="peer overflow-x-hidden overflow-y-hidden whitespace-nowrap absolute w-px h-px -m-px" style="clip-path: inset( 50% ); padding: 0; border-width: 0;" type="radio" checked />
            <div
              id="13:52"
              class={
                theme() === "light"
                  ? "border rounded-lg px-4 py-2 bg-white text-slate-700 border-slate-300"
                  : "hover:border-[#00F0FF]/40 border-[1px] border-solid rounded-lg"
              }
              style={
                theme() === "light"
                  ? ""
                  : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); color: rgba(184, 197, 217, 1); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
              }
            >
              <span id="13:53">全部</span>
            </div>
          </label>
          <label id="13:54" class="whitespace-nowrap">
            <input id="13:55" class="peer overflow-x-hidden overflow-y-hidden whitespace-nowrap absolute w-px h-px -m-px" style="clip-path: inset( 50% ); padding: 0; border-width: 0;" type="radio" />
            <div
              id="13:56"
              class={
                theme() === "light"
                  ? "border rounded-lg px-4 py-2 bg-white text-slate-600 border-slate-200 hover:border-slate-400 hover:bg-slate-50"
                  : "hover:border-[#00F0FF]/40 border-[1px] border-solid rounded-lg"
              }
              style={
                theme() === "light"
                  ? ""
                  : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); color: rgba(184, 197, 217, 1); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
              }
            >
              <span id="13:57">设计稿</span>
            </div>
          </label>
          <label id="13:58" class="whitespace-nowrap">
            <input id="13:59" class="peer overflow-x-hidden overflow-y-hidden whitespace-nowrap absolute w-px h-px -m-px" style="clip-path: inset( 50% ); padding: 0; border-width: 0;" type="radio" />
            <div
              id="13:60"
              class={
                theme() === "light"
                  ? "border rounded-lg px-4 py-2 bg-white text-slate-600 border-slate-200 hover:border-slate-400 hover:bg-slate-50"
                  : "hover:border-[#00F0FF]/40 border-[1px] border-solid rounded-lg"
              }
              style={
                theme() === "light"
                  ? ""
                  : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); color: rgba(184, 197, 217, 1); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
              }
            >
              <span id="13:61">文档</span>
            </div>
          </label>
          <label id="13:62" class="whitespace-nowrap">
            <input id="13:63" class="peer overflow-x-hidden overflow-y-hidden whitespace-nowrap absolute w-px h-px -m-px" style="clip-path: inset( 50% ); padding: 0; border-width: 0;" type="radio" />
            <div
              id="13:64"
              class={
                theme() === "light"
                  ? "border rounded-lg px-4 py-2 bg-white text-slate-600 border-slate-200 hover:border-slate-400 hover:bg-slate-50"
                  : "hover:border-[#00F0FF]/40 border-[1px] border-solid rounded-lg"
              }
              style={
                theme() === "light"
                  ? ""
                  : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); color: rgba(184, 197, 217, 1); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
              }
            >
              <span id="13:65">原型</span>
            </div>
          </label>
        </div>
        
        <div id="13:66" class="flex items-center gap-y-4 gap-x-4">
          <span id="13:67" style="color: rgba(138, 151, 170, 1);" class="text-sm">排序</span>
          <div
            id="13:68"
            class={
              theme() === "light"
                ? "flex items-center gap-y-2 gap-x-2 border rounded-lg px-4 py-2 bg-white border-slate-200 hover:border-slate-400 hover:bg-slate-50 transition-colors"
                : "hover:border-[#00F0FF]/50 flex items-center gap-y-2 gap-x-2 border-[1px] border-solid rounded-lg"
            }
            style={
              theme() === "light"
                ? ""
                : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); padding: 0.5rem 1rem; border-color: color-mix( in oklab , #00F0FF 30% , transparent );"
            }
          >
            <span
              id="13:69"
              style={theme() === "light" ? "color: #0F172A;" : "color: rgba(232, 240, 255, 1);"}
              class="text-sm"
            >
              最近更新
            </span>
            <div id="13:70" class="bg-transparent flex justify-center items-center w-4 h-4">
              <iconify-icon id="13:71" style="color: rgba(0, 240, 255, 1);" icon="lucide:chevron-down" class="text-sm"></iconify-icon>
            </div>
          </div>
        </div>
      </div>
      
      <div id="13:72" class="grid grid-cols-3 gap-y-6 gap-x-6">
        
        <div
          id="13:73"
          class={
            theme() === "light"
              ? "group flex flex-col rounded-2xl border border-slate-200 bg-white hover:shadow-[0_8px_24px_rgba(15,23,42,0.08)]"
              : "hover:shadow-[0_8px_32px_rgba(0,240,255,0.12)] group flex flex-col border-[1px] border-solid rounded-2xl"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #1A1F3A 90% , transparent ); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 4px 20px rgba(0, 240, 255, 0.08); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
          }
        >
          
          <div id="13:74" class="relative">
            <img id="13:75" alt="AI协作设计平台登录界面预览，赛博朋克风格，深色主题配霓虹蓝色高亮" src="https://static.paraflowcontent.com/public/resource/image/4763fd91-8bb1-4a84-bba6-e2e240810c35.jpeg" class="w-full h-48 object-cover rounded-tl-2xl rounded-tr-2xl" />
            
            <div id="13:76" class="group-hover:opacity-100 opacity-0 flex absolute justify-center items-center gap-y-3 gap-x-3 rounded-tl-2xl rounded-tr-2xl" style="background-color: color-mix( in oklab , #0A0E1A 60% , transparent ); top: 0; right: 0; bottom: 0; left: 0;">
              <button id="13:77" class="hover:shadow-[0_0_20px_rgba(0,240,255,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(0, 240, 255, 1);">
                <div id="13:78" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:79" style="color: rgba(10, 14, 26, 1);" icon="lucide:external-link" class="text-base"></iconify-icon>
                </div>
              </button>
              <button id="13:80" class="hover:shadow-[0_0_20px_rgba(255,0,110,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(255, 0, 110, 1);">
                <div id="13:81" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:82" style="color: rgba(255, 255, 255, 1);" icon="lucide:trash-2" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
          
          <div id="13:83" class="flex flex-col gap-y-3 gap-x-3 p-5">
            <div id="13:84" class="flex justify-between items-start">
              <div id="13:85" class="flex flex-col gap-y-2 gap-x-2">
                <h3 id="13:86" style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">登录界面设计</h3>
                <p id="13:87" style="color: rgba(138, 151, 170, 1);" class="text-sm font-normal">AI协作设计平台用户登录页面</p>
              </div>
            </div>
            
            <div id="13:88" class="flex justify-between items-center mt-2">
              <div id="13:89" class="flex items-center gap-y-2 gap-x-2">
                <div id="13:90" style="background-color: color-mix( in oklab , #B026FF 20% , transparent ); color: rgba(176, 38, 255, 1); padding: 0.25rem 0.5rem;" class="text-xs rounded-lg">设计稿</div>
              </div>
              <span id="13:91" style="color: rgba(92, 104, 118, 1);" class="text-xs">2天前</span>
            </div>
            
            <div id="13:92" style="border-top-style: solid; border-color: color-mix( in oklab , #00F0FF 15% , transparent );" class="flex items-center gap-y-2 gap-x-2 mt-3 pt-3 border-t-[1px]">
              <button onClick={() => props.onNavigate('workspace')} id="13:93" class="hover:shadow-[0_0_15px_rgba(0,240,255,0.3)] flex grow shrink items-center gap-y-2 gap-x-2 rounded-lg" style="background-color: rgba(0, 240, 255, 1); flex-basis: 0%; color: rgba(10, 14, 26, 1); padding: 0.5rem 1rem;">
                <span id="13:94" class="text-sm font-semibold">打开</span>
              </button>
              <button
                id="13:95"
                class={
                  theme() === "light"
                    ? "flex justify-center items-center w-10 h-10 border rounded-lg border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-colors"
                    : "hover:bg-[#FF006E]/20 hover:border-[#FF006E]/50 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg"
                }
                style={
                  theme() === "light"
                    ? ""
                    : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                }
              >
                <div id="13:96" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:97" style="color: rgba(138, 151, 170, 1);" icon="lucide:more-horizontal" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
        </div>
        
        <div
          id="13:98"
          class={
            theme() === "light"
              ? "group flex flex-col rounded-2xl border border-slate-200 bg-white hover:shadow-[0_8px_24px_rgba(15,23,42,0.08)]"
              : "hover:shadow-[0_8px_32px_rgba(0,240,255,0.12)] group flex flex-col border-[1px] border-solid rounded-2xl"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #1A1F3A 90% , transparent ); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 4px 20px rgba(0, 240, 255, 0.08); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
          }
        >
          
          <div id="13:99" class="relative">
            <img id="13:100" alt="产品需求文档界面，现代化文档编辑器界面，深色主题" src="https://static.paraflowcontent.com/public/resource/image/0d5c8dbb-8c9b-4c73-aeb0-ab5d9c182769.jpeg" class="w-full h-48 object-cover rounded-tl-2xl rounded-tr-2xl" />
            
            <div id="13:101" class="group-hover:opacity-100 opacity-0 flex absolute justify-center items-center gap-y-3 gap-x-3 rounded-tl-2xl rounded-tr-2xl" style="background-color: color-mix( in oklab , #0A0E1A 60% , transparent ); top: 0; right: 0; bottom: 0; left: 0;">
              <button id="13:102" class="hover:shadow-[0_0_20px_rgba(0,240,255,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(0, 240, 255, 1);">
                <div id="13:103" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:104" style="color: rgba(10, 14, 26, 1);" icon="lucide:external-link" class="text-base"></iconify-icon>
                </div>
              </button>
              <button id="13:105" class="hover:shadow-[0_0_20px_rgba(255,0,110,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(255, 0, 110, 1);">
                <div id="13:106" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:107" style="color: rgba(255, 255, 255, 1);" icon="lucide:trash-2" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
          
          <div id="13:108" class="flex flex-col gap-y-3 gap-x-3 p-5">
            <div id="13:109" class="flex justify-between items-start">
              <div id="13:110" class="flex flex-col gap-y-2 gap-x-2">
                <h3 id="13:111" style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">产品需求文档</h3>
                <p id="13:112" style="color: rgba(138, 151, 170, 1);" class="text-sm font-normal">AI协作设计平台核心功能规划</p>
              </div>
            </div>
            
            <div id="13:113" class="flex justify-between items-center mt-2">
              <div id="13:114" class="flex items-center gap-y-2 gap-x-2">
        <div
          id="13:115"
          class={
            theme() === "light"
              ? "text-xs rounded-lg bg-sky-50 text-sky-700 px-2 py-1"
              : "text-xs rounded-lg"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #00F0FF 15% , transparent ); color: rgba(0, 240, 255, 1); padding: 0.25rem 0.5rem;"
          }
        >
          文档
        </div>
              </div>
              <span id="13:116" style="color: rgba(92, 104, 118, 1);" class="text-xs">3天前</span>
            </div>
            
            <div id="13:117" style="border-top-style: solid; border-color: color-mix( in oklab , #00F0FF 15% , transparent );" class="flex items-center gap-y-2 gap-x-2 mt-3 pt-3 border-t-[1px]">
              <button onClick={() => props.onNavigate('workspace')} id="13:118" class="hover:shadow-[0_0_15px_rgba(0,240,255,0.3)] flex grow shrink items-center gap-y-2 gap-x-2 rounded-lg" style="background-color: rgba(0, 240, 255, 1); flex-basis: 0%; color: rgba(10, 14, 26, 1); padding: 0.5rem 1rem;">
                <span id="13:119" class="text-sm font-semibold">打开</span>
              </button>
              <button
                id="13:120"
                class={
                  theme() === "light"
                    ? "flex justify-center items-center w-10 h-10 border rounded-lg border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-colors"
                    : "hover:bg-[#FF006E]/20 hover:border-[#FF006E]/50 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg"
                }
                style={
                  theme() === "light"
                    ? ""
                    : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                }
              >
                <div id="13:121" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:122" style="color: rgba(138, 151, 170, 1);" icon="lucide:more-horizontal" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
        </div>
        
        <div
          id="13:123"
          class={
            theme() === "light"
              ? "group flex flex-col rounded-2xl border border-slate-200 bg-white hover:shadow-[0_8px_24px_rgba(15,23,42,0.08)]"
              : "hover:shadow-[0_8px_32px_rgba(0,240,255,0.12)] group flex flex-col border-[1px] border-solid rounded-2xl"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #1A1F3A 90% , transparent ); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 4px 20px rgba(0, 240, 255, 0.08); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
          }
        >
          
          <div id="13:124" class="relative">
            <img id="13:125" alt="主界面线框图，界面原型设计，画布和聊天面板布局" src="https://static.paraflowcontent.com/public/resource/image/3e5724b3-fc8a-4d92-82df-d7bce259a9a7.jpeg" class="w-full h-48 object-cover rounded-tl-2xl rounded-tr-2xl" />
            
            <div id="13:126" class="group-hover:opacity-100 opacity-0 flex absolute justify-center items-center gap-y-3 gap-x-3 rounded-tl-2xl rounded-tr-2xl" style="background-color: color-mix( in oklab , #0A0E1A 60% , transparent ); top: 0; right: 0; bottom: 0; left: 0;">
              <button id="13:127" class="hover:shadow-[0_0_20px_rgba(0,240,255,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(0, 240, 255, 1);">
                <div id="13:128" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:129" style="color: rgba(10, 14, 26, 1);" icon="lucide:external-link" class="text-base"></iconify-icon>
                </div>
              </button>
              <button id="13:130" class="hover:shadow-[0_0_20px_rgba(255,0,110,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(255, 0, 110, 1);">
                <div id="13:131" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:132" style="color: rgba(255, 255, 255, 1);" icon="lucide:trash-2" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
          
          <div id="13:133" class="flex flex-col gap-y-3 gap-x-3 p-5">
            <div id="13:134" class="flex justify-between items-start">
              <div id="13:135" class="flex flex-col gap-y-2 gap-x-2">
                <h3 id="13:136" style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">主界面线框图</h3>
                <p id="13:137" style="color: rgba(138, 151, 170, 1);" class="text-sm font-normal">画布与聊天面板的布局结构设计</p>
              </div>
            </div>
            
            <div id="13:138" class="flex justify-between items-center mt-2">
              <div id="13:139" class="flex items-center gap-y-2 gap-x-2">
        <div
          id="13:140"
          class={
            theme() === "light"
              ? "text-xs rounded-lg bg-rose-50 text-rose-700 px-2 py-1"
              : "text-xs rounded-lg"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #FF006E 20% , transparent ); color: rgba(255, 0, 110, 1); padding: 0.25rem 0.5rem;"
          }
        >
          原型
        </div>
              </div>
              <span id="13:141" style="color: rgba(92, 104, 118, 1);" class="text-xs">5天前</span>
            </div>
            
            <div id="13:142" style="border-top-style: solid; border-color: color-mix( in oklab , #00F0FF 15% , transparent );" class="flex items-center gap-y-2 gap-x-2 mt-3 pt-3 border-t-[1px]">
              <button onClick={() => props.onNavigate('workspace')} id="13:143" class="hover:shadow-[0_0_15px_rgba(0,240,255,0.3)] flex grow shrink items-center gap-y-2 gap-x-2 rounded-lg" style="background-color: rgba(0, 240, 255, 1); flex-basis: 0%; color: rgba(10, 14, 26, 1); padding: 0.5rem 1rem;">
                <span id="13:144" class="text-sm font-semibold">打开</span>
              </button>
              <button
                id="13:145"
                class={
                  theme() === "light"
                    ? "flex justify-center items-center w-10 h-10 border rounded-lg border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-colors"
                    : "hover:bg-[#FF006E]/20 hover:border-[#FF006E]/50 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg"
                }
                style={
                  theme() === "light"
                    ? ""
                    : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                }
              >
                <div id="13:146" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:147" style="color: rgba(138, 151, 170, 1);" icon="lucide:more-horizontal" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
        </div>
        
        <div
          id="13:148"
          class={
            theme() === "light"
              ? "group flex flex-col rounded-2xl border border-slate-200 bg-white hover:shadow-[0_8px_24px_rgba(15,23,42,0.08)]"
              : "hover:shadow-[0_8px_32px_rgba(0,240,255,0.12)] group flex flex-col border-[1px] border-solid rounded-2xl"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #1A1F3A 90% , transparent ); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 4px 20px rgba(0, 240, 255, 0.08); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
          }
        >
          
          <div id="13:149" class="relative">
            <img id="13:150" alt="用户注册页面设计，现代UI设计，暗色主题配霓虹色彩" src="https://static.paraflowcontent.com/public/resource/image/01236d12-30dc-4e14-9a0d-8c3753c456d5.jpeg" class="w-full h-48 object-cover rounded-tl-2xl rounded-tr-2xl" />
            
            <div id="13:151" class="group-hover:opacity-100 opacity-0 flex absolute justify-center items-center gap-y-3 gap-x-3 rounded-tl-2xl rounded-tr-2xl" style="background-color: color-mix( in oklab , #0A0E1A 60% , transparent ); top: 0; right: 0; bottom: 0; left: 0;">
              <button id="13:152" class="hover:shadow-[0_0_20px_rgba(0,240,255,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(0, 240, 255, 1);">
                <div id="13:153" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:154" style="color: rgba(10, 14, 26, 1);" icon="lucide:external-link" class="text-base"></iconify-icon>
                </div>
              </button>
              <button id="13:155" class="hover:shadow-[0_0_20px_rgba(255,0,110,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(255, 0, 110, 1);">
                <div id="13:156" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:157" style="color: rgba(255, 255, 255, 1);" icon="lucide:trash-2" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
          
          <div id="13:158" class="flex flex-col gap-y-3 gap-x-3 p-5">
            <div id="13:159" class="flex justify-between items-start">
              <div id="13:160" class="flex flex-col gap-y-2 gap-x-2">
                <h3 id="13:161" style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">用户注册页面</h3>
                <p id="13:162" style="color: rgba(138, 151, 170, 1);" class="text-sm font-normal">新用户注册界面设计</p>
              </div>
            </div>
            
            <div id="13:163" class="flex justify-between items-center mt-2">
              <div id="13:164" class="flex items-center gap-y-2 gap-x-2">
          <div
            id="13:165"
            class={
              theme() === "light"
                ? "text-xs rounded-lg bg-violet-50 text-violet-700 px-2 py-1"
                : "text-xs rounded-lg"
            }
            style={
              theme() === "light"
                ? ""
                : "background-color: color-mix( in oklab , #B026FF 20% , transparent ); color: rgba(176, 38, 255, 1); padding: 0.25rem 0.5rem;"
            }
          >
            设计稿
          </div>
              </div>
              <span id="13:166" style="color: rgba(92, 104, 118, 1);" class="text-xs">1周前</span>
            </div>
            
            <div id="13:167" style="border-top-style: solid; border-color: color-mix( in oklab , #00F0FF 15% , transparent );" class="flex items-center gap-y-2 gap-x-2 mt-3 pt-3 border-t-[1px]">
              <button onClick={() => props.onNavigate('workspace')} id="13:168" class="hover:shadow-[0_0_15px_rgba(0,240,255,0.3)] flex grow shrink items-center gap-y-2 gap-x-2 rounded-lg" style="background-color: rgba(0, 240, 255, 1); flex-basis: 0%; color: rgba(10, 14, 26, 1); padding: 0.5rem 1rem;">
                <span id="13:169" class="text-sm font-semibold">打开</span>
              </button>
              <button id="13:170" class="hover:bg-[#FF006E]/20 hover:border-[#FF006E]/50 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 20% , transparent );">
                <div id="13:171" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:172" style="color: rgba(138, 151, 170, 1);" icon="lucide:more-horizontal" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
        </div>
        
        <div
          id="13:173"
          class={
            theme() === "light"
              ? "group flex flex-col rounded-2xl border border-slate-200 bg-white hover:shadow-[0_8px_24px_rgba(15,23,42,0.08)]"
              : "hover:shadow-[0_8px_32px_rgba(0,240,255,0.12)] group flex flex-col border-[1px] border-solid rounded-2xl"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #1A1F3A 90% , transparent ); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 4px 20px rgba(0, 240, 255, 0.08); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
          }
        >
          
          <div id="13:174" class="relative">
            <img id="13:175" alt="仪表板界面设计，数据可视化界面，现代化暗色主题" src="https://static.paraflowcontent.com/public/resource/image/12c3585a-9400-4043-8f86-67950e8d1a43.jpeg" class="w-full h-48 object-cover rounded-tl-2xl rounded-tr-2xl" />
            
            <div id="13:176" class="group-hover:opacity-100 opacity-0 flex absolute justify-center items-center gap-y-3 gap-x-3 rounded-tl-2xl rounded-tr-2xl" style="background-color: color-mix( in oklab , #0A0E1A 60% , transparent ); top: 0; right: 0; bottom: 0; left: 0;">
              <button id="13:177" class="hover:shadow-[0_0_20px_rgba(0,240,255,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(0, 240, 255, 1);">
                <div id="13:178" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:179" style="color: rgba(10, 14, 26, 1);" icon="lucide:external-link" class="text-base"></iconify-icon>
                </div>
              </button>
              <button id="13:180" class="hover:shadow-[0_0_20px_rgba(255,0,110,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(255, 0, 110, 1);">
                <div id="13:181" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:182" style="color: rgba(255, 255, 255, 1);" icon="lucide:trash-2" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
          
          <div id="13:183" class="flex flex-col gap-y-3 gap-x-3 p-5">
            <div id="13:184" class="flex justify-between items-start">
              <div id="13:185" class="flex flex-col gap-y-2 gap-x-2">
                <h3 id="13:186" style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">数据仪表板</h3>
                <p id="13:187" style="color: rgba(138, 151, 170, 1);" class="text-sm font-normal">项目数据统计与分析界面</p>
              </div>
            </div>
            
            <div id="13:188" class="flex justify-between items-center mt-2">
              <div id="13:189" class="flex items-center gap-y-2 gap-x-2">
          <div
            id="13:190"
            class={
              theme() === "light"
                ? "text-xs rounded-lg bg-violet-50 text-violet-700 px-2 py-1"
                : "text-xs rounded-lg"
            }
            style={
              theme() === "light"
                ? ""
                : "background-color: color-mix( in oklab , #B026FF 20% , transparent ); color: rgba(176, 38, 255, 1); padding: 0.25rem 0.5rem;"
            }
          >
            设计稿
          </div>
              </div>
              <span id="13:191" style="color: rgba(92, 104, 118, 1);" class="text-xs">2周前</span>
            </div>
            
            <div id="13:192" style="border-top-style: solid; border-color: color-mix( in oklab , #00F0FF 15% , transparent );" class="flex items-center gap-y-2 gap-x-2 mt-3 pt-3 border-t-[1px]">
              <button onClick={() => props.onNavigate('workspace')} id="13:193" class="hover:shadow-[0_0_15px_rgba(0,240,255,0.3)] flex grow shrink items-center gap-y-2 gap-x-2 rounded-lg" style="background-color: rgba(0, 240, 255, 1); flex-basis: 0%; color: rgba(10, 14, 26, 1); padding: 0.5rem 1rem;">
                <span id="13:194" class="text-sm font-semibold">打开</span>
              </button>
              <button id="13:195" class="hover:bg-[#FF006E]/20 hover:border-[#FF006E]/50 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 20% , transparent );">
                <div id="13:196" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:197" style="color: rgba(138, 151, 170, 1);" icon="lucide:more-horizontal" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
        </div>
        
        <div
          id="13:198"
          class={
            theme() === "light"
              ? "group flex flex-col rounded-2xl border border-slate-200 bg-white hover:shadow-[0_8px_24px_rgba(15,23,42,0.08)]"
              : "hover:shadow-[0_8px_32px_rgba(0,240,255,0.12)] group flex flex-col border-[1px] border-solid rounded-2xl"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #1A1F3A 90% , transparent ); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 4px 20px rgba(0, 240, 255, 0.08); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
          }
        >
          
          <div id="13:199" class="relative">
            <img id="13:200" alt="移动应用界面设计，手机APP设计稿，现代化移动端界面" src="https://static.paraflowcontent.com/public/resource/image/aa88354a-31d6-49db-a129-8fc48c67d219.jpeg" class="w-full h-48 object-cover rounded-tl-2xl rounded-tr-2xl" />
            
            <div id="13:201" class="group-hover:opacity-100 opacity-0 flex absolute justify-center items-center gap-y-3 gap-x-3 rounded-tl-2xl rounded-tr-2xl" style="background-color: color-mix( in oklab , #0A0E1A 60% , transparent ); top: 0; right: 0; bottom: 0; left: 0;">
              <button id="13:202" class="hover:shadow-[0_0_20px_rgba(0,240,255,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(0, 240, 255, 1);">
                <div id="13:203" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:204" style="color: rgba(10, 14, 26, 1);" icon="lucide:external-link" class="text-base"></iconify-icon>
                </div>
              </button>
              <button id="13:205" class="hover:shadow-[0_0_20px_rgba(255,0,110,0.5)] flex justify-center items-center w-10 h-10 rounded-lg" style="background-color: rgba(255, 0, 110, 1);">
                <div id="13:206" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:207" style="color: rgba(255, 255, 255, 1);" icon="lucide:trash-2" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
          
          <div id="13:208" class="flex flex-col gap-y-3 gap-x-3 p-5">
            <div id="13:209" class="flex justify-between items-start">
              <div id="13:210" class="flex flex-col gap-y-2 gap-x-2">
                <h3 id="13:211" style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">移动端应用</h3>
                <p id="13:212" style="color: rgba(138, 151, 170, 1);" class="text-sm font-normal">AI协作设计移动端应用界面</p>
              </div>
            </div>
            
            <div id="13:213" class="flex justify-between items-center mt-2">
              <div id="13:214" class="flex items-center gap-y-2 gap-x-2">
          <div
            id="13:215"
            class={
              theme() === "light"
                ? "text-xs rounded-lg bg-violet-50 text-violet-700 px-2 py-1"
                : "text-xs rounded-lg"
            }
            style={
              theme() === "light"
                ? ""
                : "background-color: color-mix( in oklab , #B026FF 20% , transparent ); color: rgba(176, 38, 255, 1); padding: 0.25rem 0.5rem;"
            }
          >
            设计稿
          </div>
              </div>
              <span id="13:216" style="color: rgba(92, 104, 118, 1);" class="text-xs">3周前</span>
            </div>
            
            <div id="13:217" style="border-top-style: solid; border-color: color-mix( in oklab , #00F0FF 15% , transparent );" class="flex items-center gap-y-2 gap-x-2 mt-3 pt-3 border-t-[1px]">
              <button onClick={() => props.onNavigate('workspace')} id="13:218" class="hover:shadow-[0_0_15px_rgba(0,240,255,0.3)] flex grow shrink items-center gap-y-2 gap-x-2 rounded-lg" style="background-color: rgba(0, 240, 255, 1); flex-basis: 0%; color: rgba(10, 14, 26, 1); padding: 0.5rem 1rem;">
                <span id="13:219" class="text-sm font-semibold">打开</span>
              </button>
              <button id="13:220" class="hover:bg-[#FF006E]/20 hover:border-[#FF006E]/50 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg" style="background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 20% , transparent );">
                <div id="13:221" class="bg-transparent flex justify-center items-center w-5 h-5">
                  <iconify-icon id="13:222" style="color: rgba(138, 151, 170, 1);" icon="lucide:more-horizontal" class="text-base"></iconify-icon>
                </div>
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <div id="13:223" class="flex justify-center items-center gap-y-2 gap-x-2 mt-12 shrink-0">
        <button
          id="13:224"
          class={
            theme() === "light"
              ? "flex justify-center items-center w-10 h-10 rounded-lg border border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-colors"
              : "hover:border-[#00F0FF]/50 hover:bg-[#00F0FF]/10 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
          }
        >
          <div id="13:225" class="bg-transparent flex justify-center items-center w-5 h-5">
            <iconify-icon id="13:226" style="color: rgba(0, 240, 255, 1);" icon="lucide:chevron-left" class="text-base"></iconify-icon>
          </div>
        </button>
        <button
          id="13:227"
          class={
            theme() === "light"
              ? "flex justify-center items-center w-10 h-10 rounded-lg border border-slate-900 bg-slate-900 text-white text-sm font-semibold"
              : "flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #00F0FF 15% , transparent ); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 12px rgba(0, 240, 255, 0.2); color: rgba(0, 240, 255, 1); border-color: color-mix( in oklab , #00F0FF 50% , transparent );"
          }
        >
          <span id="13:228" class="text-sm font-semibold">
            1
          </span>
        </button>
        <button
          id="13:229"
          class={
            theme() === "light"
              ? "flex justify-center items-center w-10 h-10 rounded-lg border border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-colors"
              : "hover:border-[#00F0FF]/50 hover:bg-[#00F0FF]/10 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
          }
        >
          <span
            id="13:230"
            style={theme() === "light" ? "color: #64748B;" : "color: rgba(184, 197, 217, 1);"}
            class="text-sm"
          >
            2
          </span>
        </button>
        <button
          id="13:231"
          class={
            theme() === "light"
              ? "flex justify-center items-center w-10 h-10 rounded-lg border border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-colors"
              : "hover:border-[#00F0FF]/50 hover:bg-[#00F0FF]/10 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg"
          }
          style={
            theme() === "light"
              ? ""
              : "background-color: color-mix( in oklab , #1A1F3A 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
          }
        >
          <span
            id="13:232"
            style={theme() === "light" ? "color: #64748B;" : "color: rgba(184, 197, 217, 1);"}
            class="text-sm"
          >
            3
          </span>
        </button>
      </div>
    </main>
  );
};

export default ProjectList;
