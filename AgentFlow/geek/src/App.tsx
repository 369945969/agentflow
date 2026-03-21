import { createEffect, createSignal, onCleanup, onMount, Show } from "solid-js"
import "./App.css"
import Header from "./components/Header"
import Sidebar from "./components/Sidebar"
import Home from "./components/Home"
import ProjectList from "./components/ProjectList"
import Workspace from "./components/Workspace"
import DocDetail from "./components/DocDetail"
import StyleComparison from "./components/StyleComparison"
import PreviewWindow from "./components/PreviewWindow"

// @ts-ignore
declare module "solid-js" {
  namespace JSX {
    interface IntrinsicElements {
      "iconify-icon": any
    }
  }
}

const getStorage = (key: string) => {
  try {
    const val = globalThis.localStorage?.getItem?.(key)
    if (val) return val
  } catch {}
  
  const source = document.cookie || ""
  const list = source.split(";").map((part) => part.trim())
  for (const item of list) {
    if (!item.startsWith(`${key}=`)) continue
    return decodeURIComponent(item.slice(key.length + 1))
  }
  return ""
}

const setStorage = (key: string, value: string) => {
  try {
    globalThis.localStorage?.setItem?.(key, value)
  } catch {}
}

function App() {
  const params = new URLSearchParams(window.location.search)
  if (params.get("view") === "preview") {
    return <PreviewWindow />
  }

  const [sidebarWidth, setSidebarWidth] = createSignal(window.innerWidth * 0.36)
  const [isSidebarCollapsed, setIsSidebarCollapsed] = createSignal(false)
  const [currentView, setCurrentView] = createSignal(getStorage("app_name") ? "workspace" : "home")
  const [selectedDoc, setSelectedDoc] = createSignal<any>(null)
  const [transitioning, setTransitioning] = createSignal(false)
  const [transitionTitle, setTransitionTitle] = createSignal("Analyzing Input...")
  const [projectTitle, setProjectTitle] = createSignal(getStorage("app_name") || "")
  const [theme, setTheme] = createSignal<"dark" | "light">("dark")
  let isResizing = false
  createEffect(() => {
    if (!transitioning()) return
    const timer = window.setTimeout(() => {
      if (!transitioning()) return
      setTransitioning(false)
      if (currentView() === "home") setCurrentView("workspace")
    }, 5000)
    onCleanup(() => window.clearTimeout(timer))
  })

  const handleChatStart = async (text: string): Promise<void> => {
    if (currentView() !== "home") return

    setTransitioning(true)
    
    // User requested: Use language-specific limits for title
    const isEnglish = /^[A-Za-z0-9\s.,!?'"()-]+$/.test(text)
    const limit = isEnglish ? 50 : 20
    const title = text.length > limit ? text.slice(0, limit) + "..." : text
    
    setTransitionTitle(`Project: ${title}`)
    setProjectTitle(title)
    setStorage("app_name", title)

    return new Promise((resolve) => {
      setTimeout(() => {
        setCurrentView("workspace")
        setTransitioning(false)
        resolve()
      }, 2000)
    })
  }

  const startResizing = () => {
    isResizing = true
    document.addEventListener("mousemove", handleMouseMove)
    document.addEventListener("mouseup", stopResizing)
    document.body.style.cursor = "col-resize"
    document.body.style.userSelect = "none"
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!isResizing) return
    const newWidth = e.clientX
    // Min 250px, Max 80% screen
    if (newWidth > 250 && newWidth < window.innerWidth * 0.8) {
      setSidebarWidth(newWidth)
    }
  }

  const stopResizing = () => {
    isResizing = false
    document.removeEventListener("mousemove", handleMouseMove)
    document.removeEventListener("mouseup", stopResizing)
    document.body.style.cursor = ""
    document.body.style.userSelect = ""
  }

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

  onCleanup(() => {
    document.removeEventListener("mousemove", handleMouseMove)
    document.removeEventListener("mouseup", stopResizing)
  })

  return (
    <div
      class={`font-[-apple-system,BlinkMacSystemFont,'Segoe UI'] pf-BODY-main w-full h-screen flex flex-col ${theme()}`}
      style="line-height: 1.5;"
    >
      <Header onNavigate={setCurrentView} currentView={currentView()} />
      <div
        id="14:30"
        style={
          theme() === "light"
            ? "background: radial-gradient(circle at 20% 20%, rgba(56, 189, 248, 0.16) 0, transparent 55%), radial-gradient(circle at 80% 80%, rgba(59, 130, 246, 0.12) 0, transparent 55%), #F3FAFF;"
            : "background: linear-gradient(135deg, rgba(10, 14, 26, 1) 0%, rgba(26, 19, 50, 1) 50%, rgba(13, 27, 42, 1) 100%), radial-gradient(circle at 20% 30%, rgba(0, 240, 255, 0.101961) 0%, transparent 50%), radial-gradient(circle at 80% 70%, rgba(255, 0, 110, 0.101961) 0%, transparent 50%);"
        }
        class="flex w-full grow overflow-hidden"
      >
        <Show when={currentView() !== "project-list"}>
          <Sidebar 
            width={isSidebarCollapsed() ? 60 : sidebarWidth()} 
            isCollapsed={isSidebarCollapsed()}
            onToggle={() => setIsSidebarCollapsed(!isSidebarCollapsed())}
            onChatStart={handleChatStart}
            title={projectTitle()}
          />
          <Show when={!isSidebarCollapsed()}>
            <div
              class={
                theme() === "light"
                  ? "w-1 cursor-col-resize transition-colors duration-150 flex-shrink-0 bg-slate-200 hover:bg-slate-400"
                  : "w-1 cursor-col-resize transition-colors duration-150 flex-shrink-0 hover:bg-[#00F0FF]"
              }
              style={theme() === "light" ? "" : "background-color: rgba(0, 240, 255, 0.05);"}
              onMouseDown={startResizing}
            />
          </Show>
        </Show>

        <div class="relative flex flex-col grow overflow-hidden">
          <Show when={currentView() === "home"}>
            <Home onNavigate={(page) => setCurrentView(page)} />
          </Show>

          <Show when={currentView() === "project-list"}>
            <ProjectList onNavigate={(page) => setCurrentView(page)} />
          </Show>

          <Show when={currentView() === "workspace"}>
            <Workspace
              title={projectTitle()}
              onOpenFile={(doc: any) => {
                if (doc.type === "style-comparison") {
                  setCurrentView("style-comparison")
                } else {
                  setSelectedDoc(doc)
                  setCurrentView("doc-detail")
                }
              }}
            />
          </Show>

          <Show when={currentView() === "doc-detail"}>
            <DocDetail doc={selectedDoc()} onBack={() => setCurrentView("workspace")} />
          </Show>

          <Show when={currentView() === "style-comparison"}>
            <StyleComparison onBack={() => setCurrentView("workspace")} />
          </Show>

          <Show when={transitioning()}>
            <div
              style="z-index: 9999;"
              class={
                theme() === "light"
                  ? "absolute inset-0 flex flex-col items-center justify-center bg-[#EDF5FF] text-[#1F2933]"
                  : "absolute inset-0 flex flex-col items-center justify-center bg-[#0B0E14] text-[#00F0FF]"
              }
            >
              <div class="relative w-64 h-64 flex items-center justify-center mb-8">
                <div class="absolute inset-0 border-4 border-[#00F0FF] rounded-full opacity-20 animate-[ping_3s_cubic-bezier(0,0,0.2,1)_infinite]"></div>
                <div class="absolute inset-0 border-t-4 border-[#00F0FF] rounded-full animate-spin"></div>
                <div class="text-4xl font-mono font-bold tracking-widest animate-pulse">GEEK</div>
              </div>
              <div
                class={
                  theme() === "light"
                    ? "text-2xl mb-2 font-mono text-[#1F2933]"
                    : "text-2xl mb-2 font-mono text-[#E8F0FF]"
                }
              >
                {transitionTitle()}
              </div>
              <div
                class={
                  theme() === "light"
                    ? "text-sm font-mono text-[#5C6876]"
                    : "text-sm font-mono text-[#5C6876]"
                }
              >
                Initializing Environment...
              </div>
            </div>
          </Show>
        </div>
      </div>
    </div>
  )
}

export default App
