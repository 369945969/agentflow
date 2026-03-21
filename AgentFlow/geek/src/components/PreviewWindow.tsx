import { createSignal, onMount, onCleanup, Show } from "solid-js"

const PreviewWindow = () => {
  const [content, setContent] = createSignal("")
  const [kind, setKind] = createSignal<"md" | "html">("md")
  const [name, setName] = createSignal("")
  const [theme, setTheme] = createSignal<"dark" | "light">("dark")

  const base = import.meta.env.VITE_PROXY_URL ?? "http://localhost:4097"

  const cookieValue = (key: string) => {
    const source = document.cookie || ""
    const list = source.split(";").map((part) => part.trim())
    for (const item of list) {
      if (!item.startsWith(`${key}=`)) continue
      return decodeURIComponent(item.slice(key.length + 1))
    }
    return ""
  }

  const parseMarkdown = (text: string): string => {
    if (!text) return ""
    let html = text

    html = html.replace(
      /^### (.*$)/gim,
      '<h3 class="text-lg font-semibold mb-3 mt-6" style="color: rgba(232, 240, 255, 1);">$1</h3>',
    )
    html = html.replace(
      /^## (.*$)/gim,
      '<h2 class="text-xl font-bold mb-4 mt-8 pb-2 border-b" style="color: rgba(232, 240, 255, 1); border-color: color-mix(in oklab, #00F0FF 30%, transparent);">$1</h2>',
    )
    html = html.replace(
      /^# (.*$)/gim,
      '<h1 class="text-3xl font-bold mb-4" style="color: rgba(232, 240, 255, 1);">$1</h1>',
    )

    html = html.replace(/\*\*(.*)\*\*/gim, "<strong>$1</strong>")
    html = html.replace(/\*(.*)\*/gim, "<em>$1</em>")
    html = html.replace(
      /`(.*)`/gim,
      '<code class="px-2 py-1 rounded" style="background-color: color-mix(in oklab, #00F0FF 10%, transparent); color: rgba(0, 240, 255, 1); font-family: monospace;">$1</code>',
    )

    html = html.replace(
      /^- (.*$)/gim,
      '<li class="ml-4 list-none flex items-start mb-2" style="color: rgba(184, 197, 217, 1);"><span class="w-2 h-2 mt-2 mr-3 rounded-full shrink-0" style="background-color: rgba(0, 240, 255, 1);"></span>$1</li>',
    )
    html = html.replace(
      /^\d+\. (.*$)/gim,
      '<li class="ml-4 list-decimal mb-2" style="color: rgba(184, 197, 217, 1);">$1</li>',
    )

    html = html.replace(/\n\n/g, '</p><p class="mb-4 leading-relaxed" style="color: rgba(184, 197, 217, 1);">')
    html = html.replace(/\n/g, "<br>")

    return `<p class="mb-4 leading-relaxed" style="color: rgba(184, 197, 217, 1);">${html}</p>`
  }

  const fetchDoc = async (pathValue: string) => {
    const sid = cookieValue("app_session")
    const url = `${base}/workspace/file?path=${encodeURIComponent(pathValue)}&session=${encodeURIComponent(sid)}`
    const res = await fetch(url).catch(() => null)
    if (!res || !res.ok) return ""
    const data = await res.json().catch(() => null)
    if (!data?.content || typeof data.content !== "string") return ""
    return data.content as string
  }

  onMount(async () => {
    const params = new URLSearchParams(window.location.search)
    const pathValue = params.get("path")
    const kindValue = params.get("kind") as "md" | "html" | null
    const nameValue = params.get("name")
    
    if (nameValue) setName(nameValue)
    if (kindValue) setKind(kindValue)
    
    try {
      const stored = localStorage.getItem("workspace_theme")
      if (stored === "light" || stored === "dark") setTheme(stored)
    } catch {}

    if (pathValue) {
      const text = await fetchDoc(pathValue)
      setContent(text)
    }

    const storageHandler = (e: StorageEvent) => {
      if (e.key === "workspace_theme" && (e.newValue === "light" || e.newValue === "dark")) {
        setTheme(e.newValue)
      }
    }
    window.addEventListener("storage", storageHandler)
    onCleanup(() => window.removeEventListener("storage", storageHandler))
  })

  return (
    <div class={theme() === "light" ? "light w-full h-screen" : "dark w-full h-screen"}>
      <div class={`w-full h-full overflow-auto geek-scroll ${theme() === "light" ? "bg-white text-slate-900" : "bg-[#141829] text-[#E8F0FF]"}`}>
        <Show when={kind() === "html"}>
          <iframe
            srcdoc={content()}
            class="w-full h-full border-none bg-white"
            title={name()}
          />
        </Show>
        <Show when={kind() === "md"}>
          <div class="max-w-4xl mx-auto p-8">
             <div innerHTML={parseMarkdown(content())} />
          </div>
        </Show>
      </div>
    </div>
  )
}

export default PreviewWindow
