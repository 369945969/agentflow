import { type Component, createSignal, onMount, onCleanup } from "solid-js"
import userAvatar from "../assets/img.jpeg"

interface HeaderProps {
  onNavigate?: (page: string) => void
  currentView?: string
}

const Header: Component<HeaderProps> = (props) => {
  const [theme, setTheme] = createSignal<"dark" | "light">("dark")

  onMount(() => {
    try {
      const stored = localStorage.getItem("workspace_theme")
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

  const toggleTheme = () => {
    const next = theme() === "dark" ? "light" : "dark"
    setTheme(next)
    try {
      localStorage.setItem("workspace_theme", next)
    } catch {}
    window.dispatchEvent(new CustomEvent("workspace_theme_toggle", { detail: next }))
  }

  return (
    <header
      id="14:2"
      class={
        theme() === "light"
          ? "w-full shrink-0 bg-white border-b border-slate-200"
          : "w-full shrink-0"
      }
      style={
        theme() === "light"
          ? ""
          : "background: linear-gradient(135deg, rgba(26, 11, 46, 1) 0%, rgba(22, 33, 62, 1) 100%);"
      }
    >
      <nav id="14:3" style="padding: 1rem 1.5rem;" class="flex justify-between items-center w-full">
        <div id="14:4" class="flex items-center gap-x-8">
          <div id="14:5" class="flex items-center gap-x-3">
            <button
              id="14:5-btn"
              onClick={(e) => {
                e.preventDefault()
                props.onNavigate?.("workspace")
              }}
              class="bg-transparent p-0 border-0 cursor-pointer flex items-center gap-x-3"
            >
              <div id="14:6" class="bg-transparent flex justify-center items-center w-8 h-8">
                <iconify-icon
                  id="14:7"
                  style="color: rgba(0, 240, 255, 1);"
                  icon="lucide:sparkles"
                  class="text-2xl"
                ></iconify-icon>
              </div>
              <span
                id="14:8"
                class="text-xl font-bold"
                style={theme() === "light" ? "color: rgb(15,23,42);" : "color: rgba(232, 240, 255, 1);"}
              >
                极客设计工坊
              </span>
            </button>
          </div>
          <div id="14:9" class="hidden"></div>
        </div>
        <div id="14:16" class="flex items-center gap-x-4">
          <div id="14:26" class="flex items-center gap-x-2 relative">
            <button
              class="bg-transparent p-0 border-0 cursor-pointer flex justify-center items-center w-6 h-6 hover:bg-[#00F0FF]/10 rounded-full transition-colors"
              onClick={toggleTheme}
            >
              <iconify-icon
                icon={theme() === "light" ? "lucide:moon-star" : "lucide:sun-medium"}
                class="text-sm"
                style="color: rgba(0, 240, 255, 1);"
              ></iconify-icon>
            </button>
            <div class="flex justify-center items-center w-9 h-9">
              <img
                id="14:27"
                style="border-color: color-mix( in oklab , #00F0FF 40% , transparent );"
                alt="User avatar with professional appearance"
                src={userAvatar}
                class="w-9 h-9 object-cover border-[2px] border-solid rounded-full"
              />
            </div>
          </div>
        </div>
      </nav>
    </header>
  )
}

export default Header
