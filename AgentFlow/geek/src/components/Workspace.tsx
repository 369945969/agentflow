import { type Component, createSignal, onMount, onCleanup, Show, For, createEffect } from "solid-js"
import { WebviewWindow } from "@tauri-apps/api/webviewWindow"

// Mock Data for Product Requirements Document
const mockDocs = [
  {
    id: "1",
    title: "Product Charter — CodeFlow AI",
    filename: "product_charter.md",
    lastModified: "10 mins ago",
    content: `# Product Charter — CodeFlow AI\n\n## 1) Product Positioning\nCodeFlow AI is an integrated development environment (IDE) built for the AI-native era. It seamlessly combines code editing, real-time collaboration, and intelligent assistance into a single, cohesive platform.\n\n## 2) Brand Keywords\n- **Intelligent**: Leveraging state-of-the-art LLMs to understand code context.\n- **Collaborative**: Real-time multiplayer editing and debugging.\n- **Efficient**: Streamlined workflows to reduce context switching.\n\n## 3) Core Problem / JTBD\nDevelopers spend too much time configuring tools and switching between context (IDE, chat, docs). CodeFlow AI unifies these into one workspace.\n\n## 4) Goals & Mission\n- **Mission**: To empower every developer to build software faster and with more confidence.\n- **Vision**: A world where coding is as fluid as thought.`,
  },
  {
    id: "2",
    title: 'Professional Developer "Efficiency Seeker"',
    filename: "persona_developer.md",
    lastModified: "2 hours ago",
    content: `# Professional Developer "Efficiency Seeker"\n\n**Name**: Amara\n**Age**: 28\n**Role**: Senior Frontend Engineer\n\n## Bio\nAmara works at a fast-paced SaaS startup. She values speed and precision. She hates repetitive tasks and slow tooling.\n\n## Goals\n- Ship features quickly without bugs.\n- Maintain a clean codebase.\n- Automate boring parts of development.\n\n## Pain Points\n- Context switching between JIRA, GitHub, and VS Code.\n- Waiting for CI/CD pipelines.\n- outdated documentation.`,
  },
  {
    id: "3",
    title: 'Learning Developer "Exploration Tester"',
    filename: "persona_learner.md",
    lastModified: "1 day ago",
    content: `# Learning Developer "Exploration Tester"\n\n**Name**: Diego\n**Age**: 22\n**Role**: Junior Developer / Student\n\n## Bio\nDiego is just starting his career. He is eager to learn but often gets stuck on configuration and complex patterns.\n\n## Goals\n- Learn best practices.\n- Understand the codebase quickly.\n- Get help when stuck.\n\n## Pain Points\n- Overwhelmed by complex toolchains.\n- Fear of breaking things.\n- Lack of immediate feedback.`,
  },
  {
    id: "4",
    title: "Product Charter — DataFlow AI",
    filename: "data_agency.md",
    lastModified: "2 days ago",
    content: `# Product Charter — DataFlow AI\n\n## Overview\nDataFlow AI is a sub-module focusing on data processing pipelines.\n\n## Key Features\n- Visual pipeline builder.\n- Real-time data monitoring.\n- Automated anomaly detection.`,
  },
  {
    id: "5",
    title: 'Business Analyst "Data Observer"',
    filename: "persona_analyst.md",
    lastModified: "3 days ago",
    content: `# Business Analyst "Data Observer"\n\n**Name**: Sarah\n**Role**: Product Analyst\n\n## Goals\n- Understand user behavior through data.\n- Create actionable reports.\n\n## Needs\n- Easy access to production data (sanitized).\n- Visualization tools.`,
  },
  {
    id: "6",
    title: 'Operations Manager "Decision Driver"',
    filename: "persona_manager.md",
    lastModified: "1 week ago",
    content: `# Operations Manager "Decision Driver"\n\n**Name**: Marcus\n**Role**: VP of Engineering\n\n## Goals\n- Ensure team velocity and health.\n- Make strategic technology decisions.\n\n## Needs\n- High-level overview of project status.\n- Resource allocation metrics.`,
  },
]

// Mock Data for Interface Design (HTML)
const mockDesignFiles = [
  {
    id: "d1",
    title: "Login Page",
    filename: "login.html",
    type: "html",
    lastModified: "1 hour ago",
    thumbnail: "https://static.paraflowcontent.com/public/resource/image/7641c341-bc09-4527-a6a4-1075f23ad867.jpeg",
    content: `<!DOCTYPE html>
<html>
<head>
    <style>
        body { background: #0f172a; color: white; font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
        .card { background: #1e293b; padding: 2rem; rounded: 1rem; border: 1px solid #334155; width: 300px; box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1); }
        h2 { text-align: center; margin-top: 0; color: #e2e8f0; }
        input { width: 100%; padding: 0.75rem; margin: 0.5rem 0; background: #0f172a; border: 1px solid #334155; color: white; border-radius: 0.5rem; box-sizing: border-box; }
        input:focus { outline: none; border-color: #3b82f6; }
        button { width: 100%; padding: 0.75rem; margin-top: 1rem; background: #3b82f6; color: white; border: none; border-radius: 0.5rem; cursor: pointer; font-weight: bold; transition: background 0.2s; }
        button:hover { background: #2563eb; }
    </style>
</head>
<body>
    <div class="card">
        <h2>Welcome Back</h2>
        <input type="email" placeholder="Email Address" />
        <input type="password" placeholder="Password" />
        <button>Sign In</button>
    </div>
</body>
</html>`,
  },
  {
    id: "d2",
    title: "Dashboard Layout",
    filename: "dashboard.html",
    type: "html",
    lastModified: "3 hours ago",
    thumbnail: "https://static.paraflowcontent.com/public/resource/image/8f64b451-6891-4d35-bce3-60d02eccab0e.jpeg",
    content: `<!DOCTYPE html>
<html>
<head>
    <style>
        body { background: #0f172a; color: white; font-family: sans-serif; margin: 0; display: grid; grid-template-columns: 240px 1fr; height: 100vh; }
        aside { background: #1e293b; border-right: 1px solid #334155; padding: 1.5rem; }
        h3 { color: #3b82f6; margin-top: 0; font-size: 1.5rem; margin-bottom: 2rem; }
        ul { list-style: none; padding: 0; }
        li { padding: 0.75rem 1rem; color: #94a3b8; cursor: pointer; border-radius: 0.5rem; margin-bottom: 0.5rem; transition: all 0.2s; }
        li:hover, li.active { background: #334155; color: white; }
        main { padding: 2rem; overflow-y: auto; }
        header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
        .btn { background: #3b82f6; color: white; border: none; padding: 0.5rem 1rem; border-radius: 0.5rem; cursor: pointer; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1.5rem; }
        .stat-card { background: #1e293b; padding: 1.5rem; border-radius: 1rem; border: 1px solid #334155; }
        .stat-val { font-size: 2rem; font-weight: bold; margin: 0.5rem 0; }
        .stat-label { color: #94a3b8; font-size: 0.875rem; }
    </style>
</head>
<body>
    <aside>
        <h3>CodeFlow</h3>
        <ul>
            <li class="active">Dashboard</li>
            <li>Projects</li>
            <li>Team</li>
            <li>Settings</li>
        </ul>
    </aside>
    <main>
        <header>
            <h1>Overview</h1>
            <button class="btn">+ New Project</button>
        </header>
        <div class="grid">
            <div class="stat-card">
                <div class="stat-label">Total Projects</div>
                <div class="stat-val">12</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Active Users</div>
                <div class="stat-val">1,234</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Revenue</div>
                <div class="stat-val">$45.2k</div>
            </div>
        </div>
    </main>
</body>
</html>`,
  },
]

const mockWireframeFiles = [
  {
    id: "w1",
    title: "用户故事｜注册与登录",
    filename: "US_01_用户注册.md",
    type: "markdown",
    lastModified: "5 hours ago",
    content: `# 用户故事：注册与登录

## 背景
面向汽车维修店技师/店主，首次使用平台需要快速完成注册与登录。

## 用户故事
- 作为维修店技师，我希望能用手机号快速注册，以便快速查看可用车型钥匙。
- 作为店主，我希望登录后能看到最近下单记录，以便复购。

## 价值与收益
- 降低首次使用门槛
- 提升新用户转化率

## 验收标准
1. 支持手机号/邮箱注册与登录
2. 注册后自动进入首页
3. 登录失败提示明确原因`,
  },
  {
    id: "w2",
    title: "用户故事｜下单与支付",
    filename: "US_02_下单与支付.md",
    type: "markdown",
    lastModified: "1 day ago",
    content: `# 用户故事：下单与支付

## 用户故事
- 作为技师，我希望一键加入购物车并快速结算，以节省采购时间。
- 作为采购人员，我希望支持多种支付方式，以适配公司流程。

## 价值与收益
- 提升下单效率
- 降低支付流失

## 验收标准
1. 支持购物车批量下单
2. 支付成功后生成订单与通知
3. 支持失败重试与状态恢复`,
  },
  {
    id: "w3",
    title: "用户故事｜订单查询",
    filename: "US_03_订单查询.md",
    type: "markdown",
    lastModified: "2 days ago",
    content: `# 用户故事：订单查询

## 用户故事
- 作为店主，我希望能查询历史订单，以便对账与复购。
- 作为技师，我希望查看订单物流状态，以便安排工单。

## 价值与收益
- 提升复购率
- 降低售后沟通成本

## 验收标准
1. 订单支持按时间/状态筛选
2. 订单详情包含物流与发票信息
3. 支持导出订单记录`,
  },
]

interface WorkspaceProps {
  onOpenFile?: (doc: any) => void
  title?: string
}

const Workspace: Component<WorkspaceProps> = (props) => {
  const base = import.meta.env.VITE_PROXY_URL ?? "http://localhost:4097"
  const [scale] = createSignal(1)
  const [position, setPosition] = createSignal({ x: 0, y: 0 })
  const [isDragging, setIsDragging] = createSignal(false)
  const [dragStart, setDragStart] = createSignal({ x: 0, y: 0 })
  const [activeTool] = createSignal<"select" | "hand">("select")
  const [isFullScreen, setIsFullScreen] = createSignal(false)
  const [theme, setTheme] = createSignal<"dark" | "light">("dark")

  // New State for Folder/File View
  const [currentView, setCurrentView] = createSignal<"canvas" | "folder" | "file">("canvas")
  const [activeFolder, setActiveFolder] = createSignal<string | null>("界面设计")
  const [activeFile, setActiveFile] = createSignal<any | null>(null)

  const [globalDocs, setGlobalDocs] = createSignal<{ path: string; name: string; preview: string; kind: "md" | "html" }[]>([])
  const [featureDocs, setFeatureDocs] = createSignal<{ path: string; name: string; preview: string; kind: "md" | "html" }[]>([])
  const [styleDocs, setStyleDocs] = createSignal<{ path: string; name: string; preview: string; kind: "md" | "html" }[]>([])
  const [screenDocs, setScreenDocs] = createSignal<{ path: string; name: string; preview: string; kind: "md" | "html" }[]>([])

  const cookieValue = (key: string) => {
    const source = document.cookie || ""
    const list = source.split(";").map((part) => part.trim())
    for (const item of list) {
      if (!item.startsWith(`${key}=`)) continue
      return decodeURIComponent(item.slice(key.length + 1))
    }
    return ""
  }


  // Helper to format HTML
  const formatHtml = (html: string) => {
    let formatted = ''
    const pad = '  '
    let indent = 0
    
    const tokens = html.split(/(<[^>]+>)/g).filter(s => s.trim().length > 0)
    const voidTags = new Set(['area', 'base', 'br', 'col', 'embed', 'hr', 'img', 'input', 'link', 'meta', 'param', 'source', 'track', 'wbr'])
    
    tokens.forEach(token => {
      let change = 0
      const isTag = token.startsWith('<')
      const isClosing = token.startsWith('</')
      const isSelfClosing = token.endsWith('/>') || token.startsWith('<!')
      
      if (isTag) {
        if (isClosing) {
          indent = Math.max(0, indent - 1)
        } else if (!isSelfClosing) {
          const tagName = token.match(/^<([a-zA-Z0-9-]+)/)?.[1]?.toLowerCase()
          if (tagName && !voidTags.has(tagName)) {
            change = 1
          }
        }
      }
      
      formatted += pad.repeat(indent) + token.trim() + '\n'
      indent += change
    })
    
    return formatted.trim()
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

  const [previewFile, setPreviewFile] = createSignal<any | null>(null)
  const [editValue, setEditValue] = createSignal("")
  const [toast, setToast] = createSignal<{ message: string } | null>(null)
  const [selectedHtml, setSelectedHtml] = createSignal<string | null>(null)
  const [sessionId, setSessionId] = createSignal("")
  const [focusName, setFocusName] = createSignal<string | null>(null)

  onMount(() => {
    try {
      const stored = localStorage.getItem("workspace_theme")
      if (stored === "light" || stored === "dark") setTheme(stored)
    } catch {}
  })

  onMount(() => {
    const handler = (event: Event) => {
      const detail = (event as CustomEvent).detail as { name?: string } | undefined
      if (!detail?.name) return
      setFocusName(detail.name)
    }
    window.addEventListener("workspace_focus_file", handler as EventListener)
    onCleanup(() => window.removeEventListener("workspace_focus_file", handler as EventListener))
  })

  createEffect(() => {
    try {
      localStorage.setItem("workspace_theme", theme())
    } catch {}
  })

  onMount(() => {
    const handler = (event: Event) => {
      const detail = (event as CustomEvent).detail as "dark" | "light" | undefined
      if (detail !== "dark" && detail !== "light") return
      setTheme(detail)
    }
    window.addEventListener("workspace_theme_toggle", handler as EventListener)
    onCleanup(() => window.removeEventListener("workspace_theme_toggle", handler as EventListener))
  })

  onMount(() => {
    const handleMessage = (event: MessageEvent) => {
      if (event.data.type === "ELEMENT_SELECTED") {
        setSelectedHtml(formatHtml(event.data.html))
      }
    }
    window.addEventListener("message", handleMessage)
    onCleanup(() => window.removeEventListener("message", handleMessage))
  })

  onMount(() => {
    const handler = (event: Event) => {
      const detail = (event as CustomEvent).detail as { path?: string; area?: string; sessionId?: string } | undefined
      if (!detail?.path || !detail.area) return
      const sid = detail.sessionId || sessionId() || cookieValue("app_session")
      if (!sid) return
      if (detail.sessionId && detail.sessionId !== sessionId()) {
         setSessionId(detail.sessionId)
      }
      void loadWorkspaceDocs(sid)
    }
    window.addEventListener("workspace_file_created", handler as EventListener)
    onCleanup(() => window.removeEventListener("workspace_file_created", handler as EventListener))
  })

  const fetchDoc = async (pathValue: string) => {
    const sid = sessionId() || cookieValue("app_session")
    const url = `${base}/workspace/file?path=${encodeURIComponent(pathValue)}&session=${encodeURIComponent(sid)}`
    const res = await fetch(url).catch(() => null)
    if (!res || !res.ok) return ""
    const data = await res.json().catch(() => null)
    if (!data?.content || typeof data.content !== "string") return ""
    return data.content as string
  }

  const loadWorkspaceDocs = async (sid: string) => {
    const url = `${base}/workspace/tree?session=${encodeURIComponent(sid)}`
    const res = await fetch(url).catch(() => null)
    if (!res || !res.ok) return
    const data = await res.json().catch(() => null)
    const items = (data?.items ?? []) as any[]
    const flat: { path: string; name: string; content: string; kind: "md" | "html" }[] = []
    const walk = (list: any[]) => {
      list.forEach((item) => {
        if (!item || typeof item !== "object") return
        if (item.type === "file") {
          const pathValue = String(item.path || "")
          const name = String(item.name || "")
          const ext = name.toLowerCase().endsWith(".html") ? "html" : name.toLowerCase().endsWith(".md") ? "md" : ""
          if (!ext) return
          const kind = ext === "html" ? "html" : "md"
          const content = typeof item.content === "string" ? item.content : ""
          flat.push({ path: pathValue, name, content, kind })
          return
        }
        if (Array.isArray(item.children)) walk(item.children)
      })
    }
    walk(items)
    const makePreview = (text: string) => {
      if (!text) return ""
      const trimmed = text.trim()
      const lines = trimmed.split("\n").slice(0, 16)
      const snippet = lines.join("\n")
      return snippet.length > 800 ? `${snippet.slice(0, 800)}...` : snippet
    }
    const compareName = (a: { name: string }, b: { name: string }) =>
      a.name.localeCompare(b.name, undefined, { numeric: true, sensitivity: "base" })
    const globals = flat
      .filter(
        (item) =>
          item.path.includes("Global&Context") ||
          item.name === "FlowLogic.md",
      )
      .map((item) => ({
        path: item.path,
        name: item.name,
        preview: makePreview(item.content),
        kind: item.kind,
      }))
      .sort(compareName)
    const features = flat
      .filter((item) => item.path.includes("Feature&Plan"))
      .map((item) => ({
        path: item.path,
        name: item.name,
        preview: makePreview(item.content),
        kind: item.kind,
      }))
      .sort(compareName)
    const styles = flat
      .filter(
        (item) =>
          item.path.includes("Style&Guide") &&
          item.name !== "FlowLogic.md",
      )
      .map((item) => ({
        path: item.path,
        name: item.name,
        preview: item.kind === "html" ? item.content : makePreview(item.content),
        kind: item.kind,
      }))
      .sort(compareName)
    const screens = flat
      .filter((item) => item.path.includes("Screen&Prototype"))
      .map((item) => ({
        path: item.path,
        name: item.name,
        preview: item.kind === "html" ? item.content : makePreview(item.content),
        kind: item.kind,
      }))
      .sort(compareName)
    setGlobalDocs(globals)
    setFeatureDocs(features)
    setStyleDocs(styles)
    setScreenDocs(screens)
  }

  onMount(() => {
    const sid = cookieValue("app_session")
    if (!sid) return
    setSessionId(sid)
    void loadWorkspaceDocs(sid)
  })

  const openDocPreview = async (doc: { path: string; name: string; kind: "md" | "html" }) => {
    try {
      // @ts-ignore
      if (window.__TAURI_INTERNALS__) {
        const label = `preview-${Date.now()}`
        const webview = new WebviewWindow(label, {
          url: `index.html?view=preview&path=${encodeURIComponent(doc.path)}&name=${encodeURIComponent(doc.name)}&kind=${doc.kind}`,
          title: doc.name,
          width: 1024,
          height: 768,
        })
        webview.once('tauri://created', function () {
          // webview window successfully created
        })
        webview.once('tauri://error', function (e) {
          // an error happened creating the webview window
          console.error("Window create error", e)
        })
        return
      }
    } catch (e) {
      console.error("Failed to create window", e)
    }

    const text = await fetchDoc(doc.path)
    if (!text) return
    if (doc.kind === "html") {
      setPreviewFile({
        kind: "html",
        name: doc.name,
        path: doc.path,
        content: text,
      })
      return
    }
    setPreviewFile({
      kind: "md",
      name: doc.name,
      path: doc.path,
      content: text,
    })
    setEditValue(text)
  }

  const saveMarkdown = async () => {
    const file = previewFile()
    if (!file || file.kind !== "md") return
    const url = `${base}/workspace/file`
    const res = await fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ path: file.path, content: editValue() }),
    }).catch(() => null)
    if (!res || !res.ok) return
    setPreviewFile(null)
    setToast({ message: "Markdown 保存成功" })
    setTimeout(() => {
      setToast(null)
    }, 1500)
    const sid = sessionId()
    if (!sid) return
    void loadWorkspaceDocs(sid)
  }

  const getMermaidHtml = (content: string) => {
    return `
      <!DOCTYPE html>
      <html>
      <head>
        <style>
          body {
            background-color: #141829;
            color: #E8F0FF;
            font-family: sans-serif;
            margin: 0;
            height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            overflow: hidden;
          }
          .mermaid {
            width: 100%;
            height: 100%;
            display: flex;
            justify-content: center;
            align-items: center;
          }
          /* Cyberpunk overrides for Mermaid */
          :root {
            --mermaid-font-family: 'Courier New', monospace;
          }
          /* Custom Glow Effects via SVG Filters if possible, or just CSS */
          svg {
            filter: drop-shadow(0 0 5px rgba(0, 240, 255, 0.3));
          }
          .node rect, .node circle, .node polygon, .node path {
            fill: #1A1F3A !important;
            stroke: #00F0FF !important;
            stroke-width: 2px !important;
          }
          .node .label {
            color: #E8F0FF !important;
            font-weight: bold;
          }
          .edgePath .path {
            stroke: #B026FF !important;
            stroke-width: 2px !important;
          }
          .arrowheadPath {
            fill: #B026FF !important;
          }
          .cluster rect {
            fill: rgba(0, 240, 255, 0.05) !important;
            stroke: rgba(0, 240, 255, 0.2) !important;
          }
        </style>
      </head>
      <body>
        <div class="mermaid">
          ${content}
        </div>
        <script type="module">
          import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';
          mermaid.initialize({ 
            startOnLoad: true, 
            theme: 'base',
            themeVariables: {
              background: '#141829',
              primaryColor: '#1A1F3A',
              secondaryColor: '#141829',
              tertiaryColor: '#141829',
              primaryTextColor: '#E8F0FF',
              secondaryTextColor: '#E8F0FF',
              tertiaryTextColor: '#E8F0FF',
              lineColor: '#B026FF',
              fontFamily: 'monospace'
            }
          });
        </script>
      </body>
      </html>
    `
  }

  // Inject Inspector Script
  const getPreviewContent = (html: string) => {
    const style = `
      <style>
        html, body {
          margin: 0;
          padding: 10px;
          -ms-overflow-style: none;
          scrollbar-width: none;
        }
        body::-webkit-scrollbar {
          display: none;
        }
      </style>
    `
    const script = `
      <script>
        document.body.dataset.inspectMode = 'true';
        let hoveredElement = null;
        let originalOutline = '';
        let originalCursor = '';

        window.addEventListener('message', (event) => {
          if (event.data.type === 'TOGGLE_INSPECT') {
            document.body.dataset.inspectMode = event.data.active;
            if (!event.data.active && hoveredElement) {
               hoveredElement.style.outline = originalOutline;
               hoveredElement.style.cursor = originalCursor;
               hoveredElement = null;
            }
          }
        });

        document.addEventListener('mouseover', (e) => {
          if (document.body.dataset.inspectMode !== 'true') return;
          e.stopPropagation();
          
          if (hoveredElement && hoveredElement !== e.target) {
             hoveredElement.style.outline = originalOutline;
             hoveredElement.style.cursor = originalCursor;
          }
          
          hoveredElement = e.target;
          originalOutline = hoveredElement.style.outline;
          originalCursor = hoveredElement.style.cursor;
          
          hoveredElement.style.outline = '2px solid #00F0FF';
          hoveredElement.style.cursor = 'default';
        });

        document.addEventListener('mouseout', (e) => {
           if (document.body.dataset.inspectMode !== 'true') return;
           // Don't clear immediately on mouseout to avoid flickering, 
           // let mouseover of next element handle it or body mouseout
        });

        document.addEventListener('click', (e) => {
          if (document.body.dataset.inspectMode !== 'true') return;
          e.preventDefault();
          e.stopPropagation();
          window.parent.postMessage({
            type: 'ELEMENT_SELECTED',
            html: e.target.outerHTML
          }, '*');
        });
      </script>
    `
    let out = html
    if (out.includes("</head>")) {
      out = out.replace("</head>", `${style}</head>`)
    } else {
      out = out.replace("<body", `${style}<body`)
    }
    return out.replace("</body>", `${script}</body>`)
  }

  let canvasRef: HTMLDivElement | undefined

  const handleMouseDown = (e: MouseEvent) => {
    if (activeTool() === "hand") {
      setIsDragging(true)
      setDragStart({ x: e.clientX - position().x, y: e.clientY - position().y })
    }
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (isDragging()) {
      setPosition({
        x: e.clientX - dragStart().x,
        y: e.clientY - dragStart().y,
      })
    }
  }

  const handleMouseUp = () => {
    setIsDragging(false)
  }

  const toggleFullScreen = () => {
    if (!document.fullscreenElement) {
      canvasRef?.requestFullscreen()
    } else {
      document.exitFullscreen()
    }
  }

  const exportWorkspace = async () => {
    const sid =
      cookieValue("app_session") ||
      sessionId()
    if (!sid) return
    const url = `${base}/workspace/export?session=${encodeURIComponent(sid)}`
    const res = await fetch(url).catch(() => null)
    if (!res || !res.ok) return
    const blob = await res.blob().catch(() => null)
    if (!blob) return
    const href = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = href
    a.download = `${sid}.zip`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(href)
  }

  const updateFullScreenState = () => {
    setIsFullScreen(!!document.fullscreenElement)
  }

  onMount(() => {
    document.addEventListener("fullscreenchange", updateFullScreenState)
  })

  onCleanup(() => {
    document.removeEventListener("fullscreenchange", updateFullScreenState)
  })

  // Folder/File Handlers

  const openFile = (file: any) => {
    console.log("Opening file:", file.filename)
    setActiveFile(file)
    setCurrentView("file")
    props.onOpenFile?.(file)
  }

  const backToCanvas = () => {
    setCurrentView("canvas")
    setActiveFolder(null)
  }

  const backToFolder = () => {
    setCurrentView("folder")
    setActiveFile(null)
  }
  
  const isDesignView = () => currentView() === "folder" && activeFolder() === "design"
  const isCanvasOrDesign = () => currentView() === "canvas" || isDesignView()

  return (
    <main
      id="12:31"
      style="flex-basis: 0%; padding: 0.75rem 1rem;"
      class={`
        overflow-x-hidden overflow-y-hidden flex flex-col grow shrink relative
        ${theme() === "light" ? "bg-[#F5F5FA] text-slate-900" : ""}
      `}
    >
      <Show when={toast()}>
        <div class="pointer-events-none fixed top-4 inset-x-0 flex justify-center z-[1200]">
          <div class="pointer-events-auto px-4 py-2 rounded-lg bg-[#00F0FF]/10 border border-[#00F0FF]/40 text-xs text-[#E8F0FF] shadow-lg">
            {toast()!.message}
          </div>
        </div>
      </Show>
      <Show when={currentView() !== "canvas"}>
        <div id="12:32" class="flex justify-between items-center mb-3 h-10 relative">
          <div class="flex items-center gap-x-4">
            <Show when={currentView() !== "canvas"}>
              <button
                onClick={currentView() === "folder" ? backToCanvas : backToFolder}
                class="hover:bg-[#00F0FF]/15 hover:shadow-[0_0_12px_rgba(0,240,255,0.25)] flex justify-center items-center w-10 h-10 rounded-lg text-[#00F0FF]"
                style="background-color: color-mix( in oklab , #1A1F3A 90% , transparent );"
              >
                <iconify-icon icon="lucide:arrow-left" class="text-xl"></iconify-icon>
              </button>
              <Show when={currentView() === "folder" && activeFolder() === "wireframe"}>
                <h1 style="color: rgba(232, 240, 255, 1);" class="text-2xl font-bold ml-2">
                  用户故事
                </h1>
              </Show>
            </Show>
          </div>

          <div class="flex items-center gap-x-2"></div>
        </div>
      </Show>

      <div
        id="12:53"
        ref={canvasRef}
        style={{
          "background-color": "transparent",
          "border-color": "transparent",
          cursor:
            isCanvasOrDesign()
              ? activeTool() === "hand"
                ? isDragging()
                  ? "grabbing"
                  : "grab"
                : "default"
              : "default",
        }}
        class="overflow-hidden relative grow shrink"
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
      >
        <Show when={isCanvasOrDesign()}>
          <div
            style={{
              transform: `translate(${position().x}px, ${position().y}px) scale(${scale()})`,
              "transform-origin": "center center",
              width: "100%",
              height: "100%",
              position: "absolute",
              top: 0,
              left: 0,
            }}
          >
            <Show when={currentView() === "canvas"}>
              <div
                id="12:54"
                style="top: 0; right: 0; bottom: 0; left: 0;"
                class="absolute inset-0 overflow-hidden flex"
              >
                <div class="w-full h-full grid grid-cols-2 grid-rows-2 gap-4">
                  <div
                    class={
                      theme() === "light"
                        ? "h-full border border-slate-200 rounded-2xl bg-white p-4 flex flex-col overflow-hidden shadow-sm"
                        : "h-full border border-[#00F0FF]/20 rounded-2xl bg-[#0F1624]/80 p-4 flex flex-col overflow-hidden"
                    }
                  >
                    <div class="flex items-center justify-between mb-3">
                      <div>
                        <div
                          class={
                            theme() === "light"
                              ? "text-base font-semibold text-slate-900"
                              : "text-base font-semibold text-[#E8F0FF]"
                          }
                        >
                          PRD & Architecture
                        </div>
                      </div>
                    </div>
                    <div class="flex-1 overflow-auto geek-scroll px-2">
                      <Show
                        when={globalDocs().length}
                        fallback={
                          <div class="text-xs text-[#5C6876]">
                            等待 Global&Context 目录下 PRD 与 Architecture 文档生成...
                          </div>
                        }
                      >
                        <div class="inline-flex gap-3 flex-wrap min-w-max">
                          <For each={globalDocs()}>
                            {(doc) => {
                              let el: HTMLButtonElement | undefined
                              createEffect(() => {
                                if (focusName() === doc.name && el) {
                                  el.scrollIntoView({ behavior: "smooth", block: "nearest", inline: "nearest" })
                                }
                              })
                              return (
                                <button
                                  ref={(node) => {
                                    el = node
                                  }}
                                  type="button"
                                  title={doc.name}
                                  onMouseDown={(e) => e.stopPropagation()}
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    void openDocPreview(doc)
                                  }}
                                  class={
                                    theme() === "light"
                                      ? focusName() === doc.name
                                        ? "min-w-[160px] max-w-[200px] text-left rounded-xl border border-[#FB3B3B] bg-white ring-1 ring-[#FB3B3B]/40 transition-all px-3 py-2 flex flex-col gap-2 shadow-sm"
                                        : "min-w-[160px] max-w-[200px] text-left rounded-xl border border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-all px-3 py-2 flex flex-col gap-2 shadow-sm"
                                      : focusName() === doc.name
                                        ? "min-w-[160px] max-w-[200px] text-left rounded-xl border border-[#FB3B3B] bg-[#141829] ring-1 ring-[#FB3B3B]/40 transition-all px-3 py-2 flex flex-col gap-2"
                                        : "min-w-[160px] max-w-[200px] text-left rounded-xl border border-[#00F0FF]/10 bg-[#141829]/80 hover:border-[#00F0FF]/60 hover:bg-[#141829] transition-all px-3 py-2 flex flex-col gap-2"
                                  }
                                >
                                  <div
                                    class={
                                      theme() === "light"
                                        ? "text-xs font-semibold text-slate-900 truncate"
                                        : "text-xs font-semibold text-[#E8F0FF] truncate"
                                    }
                                  >
                                    {doc.name}
                                  </div>
                                  <Show
                                    when={doc.kind === "html"}
                                    fallback={
                                      <div
                                        class={
                                          theme() === "light"
                                            ? "flex-1 text-[10px] text-slate-600 whitespace-pre-wrap leading-relaxed overflow-hidden"
                                            : "flex-1 text-[10px] text-[#8A97AA] whitespace-pre-wrap leading-relaxed overflow-hidden"
                                        }
                                      >
                                        {doc.preview}
                                      </div>
                                    }
                                  >
                                    <div
                                      class={
                                        theme() === "light"
                                          ? "mt-1 flex-1 rounded-md overflow-hidden border border-slate-200 bg-white"
                                          : "mt-1 flex-1 rounded-md overflow-hidden border border-[#00F0FF]/20 bg-black"
                                      }
                                    >
                                      <iframe
                                        srcdoc={doc.preview}
                                        class="w-[300%] h-[300%] border-none pointer-events-none scale-[0.3333] origin-top-left bg-white"
                                        tabindex="-1"
                                      />
                                    </div>
                                  </Show>
                                </button>
                              )
                            }}
                          </For>
                        </div>
                      </Show>
                    </div>
                  </div>

                  <div
                    class={
                      theme() === "light"
                        ? "h-full border border-slate-200 rounded-2xl bg-white p-4 flex flex-col overflow-hidden shadow-sm"
                        : "h-full border border-[#00F0FF]/20 rounded-2xl bg-[#0F1624]/80 p-4 flex flex-col overflow-hidden"
                    }
                  >
                    <div class="flex items-center justify-between mb-3">
                      <div>
                        <div
                          class={
                            theme() === "light"
                              ? "text-base font-semibold text-slate-900"
                              : "text-base font-semibold text-[#E8F0FF]"
                          }
                        >
                          Style Guide
                        </div>
                      </div>
                    </div>
                    <div class="flex-1 overflow-auto geek-scroll pr-2">
                      <Show
                        when={styleDocs().length}
                        fallback={
                          <div class="text-xs text-[#5C6876]">
                            等待 Style&Guide 目录下样式与交互文档生成...
                          </div>
                        }
                      >
                        <div class="flex gap-3">
                          <For each={styleDocs()}>
                            {(doc) => {
                              let el: HTMLButtonElement | undefined
                              createEffect(() => {
                                if (focusName() === doc.name && el) {
                                  el.scrollIntoView({ behavior: "smooth", block: "nearest", inline: "nearest" })
                                }
                              })
                              return doc.kind === "html" ? (
                                <button
                                  ref={(node) => {
                                    el = node
                                  }}
                                  type="button"
                                  title={doc.name}
                                  onMouseDown={(e) => e.stopPropagation()}
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    void openDocPreview(doc)
                                  }}
                                  class={
                                    theme() === "light"
                                      ? focusName() === doc.name
                                        ? "group relative flex-none w-[260px] aspect-square rounded-xl border border-[#FB3B3B] bg-slate-50 ring-1 ring-[#FB3B3B]/40 transition-all overflow-hidden shadow-sm"
                                        : "group relative flex-none w-[260px] aspect-square rounded-xl border border-slate-200 bg-slate-50 hover:border-slate-300 hover:bg-white transition-all overflow-hidden shadow-sm"
                                      : focusName() === doc.name
                                        ? "group relative flex-none w-[260px] aspect-square rounded-xl border border-[#FB3B3B] bg-[#050816]/80 ring-1 ring-[#FB3B3B]/40 transition-all overflow-hidden"
                                        : "group relative flex-none w-[260px] aspect-square rounded-xl border border-[#00F0FF]/25 bg-[#050816]/80 hover:border-[#00F0FF]/60 hover:bg-[#050816] transition-all overflow-hidden"
                                  }
                                >
                                  <div class="absolute inset-0">
                                    <div
                                      class={
                                        theme() === "light"
                                          ? "w-full h-full bg-white rounded-lg overflow-hidden relative"
                                          : "w-full h-full bg-black rounded-lg overflow-hidden relative"
                                      }
                                    >
                                      <iframe
                                        srcdoc={doc.preview}
                                        scrolling="no"
                                        class="w-[300%] h-[300%] border-none pointer-events-none scale-[0.3333] origin-top-left bg-white"
                                        tabindex="-1"
                                      />
                                      <div
                                        class={
                                          theme() === "light"
                                            ? "absolute inset-0 bg-transparent group-hover:bg-slate-900/5 transition-colors"
                                            : "absolute inset-0 bg-transparent group-hover:bg-white/5 transition-colors"
                                        }
                                      ></div>
                                    </div>
                                  </div>
                                </button>
                              ) : (
                                <button
                                  ref={(node) => {
                                    el = node
                                  }}
                                  type="button"
                                  title={doc.name}
                                  onMouseDown={(e) => e.stopPropagation()}
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    void openDocPreview(doc)
                                  }}
                                  class={
                                    theme() === "light"
                                      ? focusName() === doc.name
                                        ? "min-w-[160px] max-w-[200px] text-left rounded-xl border border-[#FB3B3B] bg-white ring-1 ring-[#FB3B3B]/40 transition-all px-3 py-2 flex flex-col gap-2 shadow-sm"
                                        : "min-w-[160px] max-w-[200px] text-left rounded-xl border border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-all px-3 py-2 flex flex-col gap-2 shadow-sm"
                                      : focusName() === doc.name
                                        ? "min-w-[160px] max-w-[200px] text-left rounded-xl border border-[#FB3B3B] bg-[#141829] ring-1 ring-[#FB3B3B]/40 transition-all px-3 py-2 flex flex-col gap-2"
                                        : "min-w-[160px] max-w-[200px] text-left rounded-xl border border-[#00F0FF]/10 bg-[#141829]/80 hover:border-[#00F0FF]/60 hover:bg-[#141829] transition-all px-3 py-2 flex flex-col gap-2"
                                  }
                                >
                                  <div
                                    class={
                                      theme() === "light"
                                        ? "text-xs font-semibold text-slate-900 truncate"
                                        : "text-xs font-semibold text-[#E8F0FF] truncate"
                                    }
                                  >
                                    {doc.name}
                                  </div>
                                  <div
                                    class={
                                      theme() === "light"
                                        ? "flex-1 text-[10px] text-slate-600 whitespace-pre-wrap leading-relaxed overflow-hidden"
                                        : "flex-1 text-[10px] text-[#8A97AA] whitespace-pre-wrap leading-relaxed overflow-hidden"
                                    }
                                  >
                                    {doc.preview}
                                  </div>
                                </button>
                              )
                            }}
                          </For>
                        </div>
                      </Show>
                    </div>
                  </div>

                  <div
                    class={
                      theme() === "light"
                        ? "h-full border border-slate-200 rounded-2xl bg-white p-4 flex flex-col overflow-hidden shadow-sm"
                        : "h-full border border-[#00F0FF]/20 rounded-2xl bg-[#0F1624]/80 p-4 flex flex-col overflow-hidden"
                    }
                  >
                    <div class="flex items-center justify-between mb-3">
                      <div>
                        <div
                          class={
                            theme() === "light"
                              ? "text-base font-semibold text-slate-900"
                              : "text-base font-semibold text-[#E8F0FF]"
                          }
                        >
                          User Stories
                        </div>
                      </div>
                    </div>
                    <div class="flex-1 overflow-auto geek-scroll px-2">
                      <Show
                        when={featureDocs().length}
                        fallback={
                          <div class="text-xs text-[#5C6876]">
                            等待 Feature&Plan 目录下用户故事文档生成...
                          </div>
                        }
                      >
                        <div class="inline-flex gap-3 flex-wrap min-w-max">
                          <For each={featureDocs()}>
                            {(doc) => {
                              let el: HTMLButtonElement | undefined
                              createEffect(() => {
                                if (focusName() === doc.name && el) {
                                  el.scrollIntoView({ behavior: "smooth", block: "nearest", inline: "nearest" })
                                }
                              })
                              return (
                                <button
                                  ref={(node) => {
                                    el = node
                                  }}
                                  type="button"
                                  title={doc.name}
                                  onMouseDown={(e) => e.stopPropagation()}
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    void openDocPreview(doc)
                                  }}
                                  class={
                                    theme() === "light"
                                      ? focusName() === doc.name
                                        ? "min-w-[160px] max-w-[200px] text-left rounded-xl border border-[#FB3B3B] bg-white ring-1 ring-[#FB3B3B]/40 transition-all px-3 py-2 flex flex-col gap-2 shadow-sm"
                                        : "min-w-[160px] max-w-[200px] text-left rounded-xl border border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-all px-3 py-2 flex flex-col gap-2 shadow-sm"
                                      : focusName() === doc.name
                                        ? "min-w-[160px] max-w-[200px] text-left rounded-xl border border-[#FB3B3B] bg-[#141829] ring-1 ring-[#FB3B3B]/40 transition-all px-3 py-2 flex flex-col gap-2"
                                        : "min-w-[160px] max-w-[200px] text-left rounded-xl border border-[#00F0FF]/10 bg-[#141829]/80 hover:border-[#00F0FF]/60 hover:bg-[#141829] transition-all px-3 py-2 flex flex-col gap-2"
                                  }
                                >
                                  <div
                                    class={
                                      theme() === "light"
                                        ? "text-xs font-semibold text-slate-900 truncate"
                                        : "text-xs font-semibold text-[#E8F0FF] truncate"
                                    }
                                  >
                                    {doc.name}
                                  </div>
                                  <Show
                                    when={doc.kind === "html"}
                                    fallback={
                                      <div
                                        class={
                                          theme() === "light"
                                            ? "flex-1 text-[10px] text-slate-600 whitespace-pre-wrap leading-relaxed overflow-hidden"
                                            : "flex-1 text-[10px] text-[#8A97AA] whitespace-pre-wrap leading-relaxed overflow-hidden"
                                        }
                                      >
                                        {doc.preview}
                                      </div>
                                    }
                                  >
                                    <div
                                      class={
                                        theme() === "light"
                                          ? "mt-1 flex-1 rounded-md overflow-hidden border border-slate-200 bg-white"
                                          : "mt-1 flex-1 rounded-md overflow-hidden border border-[#00F0FF]/20 bg-black"
                                      }
                                    >
                                      <iframe
                                        srcdoc={doc.preview}
                                        class="w-[300%] h-[300%] border-none pointer-events-none scale-[0.3333] origin-top-left bg-white"
                                        tabindex="-1"
                                      />
                                    </div>
                                  </Show>
                                </button>
                              )
                            }}
                          </For>
                        </div>
                      </Show>
                    </div>
                  </div>

                  <div
                    class={
                      theme() === "light"
                        ? "h-full border border-slate-200 rounded-2xl bg-white p-4 flex flex-col overflow-hidden shadow-sm"
                        : "h-full border border-[#00F0FF]/20 rounded-2xl bg-[#0F1624]/80 p-4 flex flex-col overflow-hidden"
                    }
                  >
                    <div class="flex items-center justify-between mb-3">
                      <div>
                        <div
                          class={
                            theme() === "light"
                              ? "text-base font-semibold text-slate-900"
                              : "text-base font-semibold text-[#E8F0FF]"
                          }
                        >
                          HTML Screens
                        </div>
                      </div>
                    </div>
                    <div
                      class={
                        theme() === "light"
                          ? "flex-1 overflow-x-auto overflow-y-hidden rounded-xl border border-slate-200 bg-slate-50 geek-scroll"
                          : "flex-1 overflow-x-auto overflow-y-hidden rounded-xl bg-black/80 border border-[#00F0FF]/10 geek-scroll"
                      }
                    >
                      <Show
                        when={screenDocs().length}
                        fallback={
                          <div class="w-full h-full flex items-center justify-center text-xs text-[#5C6876]">
                            等待 Screen&Prototype 目录下 HTML 原型生成...
                          </div>
                        }
                      >
                        <div class="flex gap-4 h-full items-stretch px-3 py-3">
                          <For each={screenDocs()}>
                            {(doc) => {
                              let el: HTMLButtonElement | undefined
                              createEffect(() => {
                                if (focusName() === doc.name && el) {
                                  el.scrollIntoView({ behavior: "smooth", block: "nearest", inline: "nearest" })
                                }
                              })
                              return (
                                <button
                                  ref={(node) => {
                                    el = node
                                  }}
                                  type="button"
                                  onMouseDown={(e) => e.stopPropagation()}
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    void openDocPreview(doc)
                                  }}
                                  class={
                                    theme() === "light"
                                      ? focusName() === doc.name
                                        ? "group relative flex-none w-[260px] aspect-square rounded-xl border border-[#FB3B3B] bg-white ring-1 ring-[#FB3B3B]/40 transition-all overflow-hidden shadow-sm"
                                        : "group relative flex-none w-[260px] aspect-square rounded-xl border border-slate-200 bg-white hover:border-slate-400 hover:bg-slate-50 transition-all overflow-hidden shadow-sm"
                                      : focusName() === doc.name
                                        ? "group relative flex-none w-[260px] aspect-square rounded-xl border border-[#FB3B3B] bg-[#141829]/80 ring-1 ring-[#FB3B3B]/40 transition-all overflow-hidden"
                                        : "group relative flex-none w-[260px] aspect-square rounded-xl border border-[#00F0FF]/20 bg-[#141829]/80 hover:border-[#00F0FF]/60 hover:bg-[#141829] transition-all overflow-hidden"
                                  }
                                >
                                <div class="absolute inset-0">
                                  <div
                                    class={
                                      theme() === "light"
                                        ? "w-full h-full bg-white rounded-lg overflow-hidden relative"
                                        : "w-full h-full bg-black rounded-lg overflow-hidden relative"
                                    }
                                  >
                                    <Show when={doc.kind === "html"}>
                                      <iframe
                                        srcdoc={doc.preview}
                                        class="w-[300%] h-[300%] border-none pointer-events-none scale-[0.3333] origin-top-left bg-white"
                                        tabindex="-1"
                                      />
                                    </Show>
                                    <div
                                      class={
                                        theme() === "light"
                                          ? "absolute inset-0 bg-transparent group-hover:bg-slate-900/5 transition-colors"
                                          : "absolute inset-0 bg-transparent group-hover:bg-white/5 transition-colors"
                                      }
                                    ></div>
                                  </div>
                                </div>
                                <div class="absolute bottom-0 left-0 right-0 px-3 pb-3">
                                  <div
                                    class={
                                      theme() === "light"
                                        ? "inline-flex max-w-full rounded bg-black/60 px-2 py-1 text-xs font-semibold text-white truncate"
                                        : "inline-flex max-w-full rounded bg-black/60 px-2 py-1 text-xs font-semibold text-[#E8F0FF] truncate"
                                    }
                                  >
                                    {doc.name}
                                  </div>
                                </div>
                                </button>
                              )
                            }}
                          </For>
                        </div>
                      </Show>
                    </div>
                  </div>
                </div>
              </div>
            </Show>

            <Show when={isCanvasOrDesign()}>
              <div class="pointer-events-none absolute inset-0">
                <div class="absolute bottom-4 left-1/2 -translate-x-1/2 flex items-center gap-x-3 pointer-events-auto">
                  <button
                    id="12:fullscreen"
                    onClick={toggleFullScreen}
                    class={
                      theme() === "light"
                        ? "flex justify-center items-center w-10 h-10 rounded-lg border border-slate-200 bg-white hover:border-[#00F0FF] hover:bg-slate-50 hover:shadow-[0_0_12px_rgba(0,240,255,0.25)]"
                        : "flex justify-center items-center w-10 h-10 rounded-lg bg-[#1A1F3A]/90 hover:bg-[#00F0FF]/10 hover:shadow-[0_0_12px_rgba(0,240,255,0.25)]"
                    }
                  >
                    <div class="bg-transparent flex justify-center items-center w-5 h-5">
                      <iconify-icon
                        style="color: rgba(0, 240, 255, 1);"
                        icon={isFullScreen() ? "lucide:minimize" : "lucide:maximize"}
                        class="text-base"
                      ></iconify-icon>
                    </div>
                  </button>
                  <button
                    onClick={exportWorkspace}
                    class={
                      theme() === "light"
                        ? "flex justify-center items-center w-10 h-10 rounded-lg border border-slate-200 bg-white hover:border-[#00F0FF] hover:bg-slate-50 hover:shadow-[0_0_12px_rgba(0,240,255,0.25)]"
                        : "flex justify-center items-center w-10 h-10 rounded-lg bg-[#1A1F3A]/90 hover:bg-[#00F0FF]/10 hover:shadow-[0_0_12px_rgba(0,240,255,0.25)]"
                    }
                  >
                    <div class="bg-transparent flex justify-center items-center w-5 h-5">
                      <iconify-icon
                        style="color: rgba(0, 240, 255, 1);"
                        icon="lucide:download-cloud"
                        class="text-base"
                      ></iconify-icon>
                    </div>
                  </button>
                </div>
              </div>
            </Show>

            {/* Design Folder View - Inside Transform */}
            <Show when={isDesignView()}>
              <div class="absolute inset-0 p-12 overflow-y-auto pointer-events-auto">
                 <div class="grid grid-cols-3 gap-8 w-full max-w-7xl mx-auto pb-12">
                    <For each={mockDesignFiles}>
                      {(file) => (
                        <div
                          onClick={(e) => {
                             e.stopPropagation()
                             setPreviewFile(file)
                          }}
                          class="group cursor-pointer relative w-full aspect-[4/3] rounded-2xl border border-[#B026FF]/20 bg-[#1A1F3A] hover:border-[#B026FF]/60 hover:shadow-[0_0_30px_rgba(176,38,255,0.3)] hover:scale-105 transition-all duration-300"
                        >
                          <div class="w-full h-full rounded-2xl overflow-hidden relative bg-white">
                             <iframe 
                               srcdoc={file.content} 
                               class="w-[300%] h-[300%] border-none pointer-events-none scale-[0.3333] origin-top-left"
                               tabindex="-1"
                               scrolling="no"
                             />
                             <div class="absolute inset-0 bg-transparent hover:bg-black/10 transition-colors"></div>
                             
                             <div class="absolute bottom-0 left-0 right-0 p-6 bg-gradient-to-t from-[#141829] via-[#141829]/80 to-transparent pointer-events-none">
                                <h3 class="text-[#E8F0FF] text-xl font-bold mb-1">{file.title}</h3>
                                <p class="text-[#8A97AA] text-sm font-mono">{file.filename}</p>
                             </div>
                          </div>
                        </div>
                      )}
                    </For>
                 </div>
              </div>
            </Show>
          </div>
        </Show>

        <Show when={currentView() === "folder" && activeFolder() === "docs"}>
          <div class="w-full h-full p-8 overflow-hidden flex flex-col">
            <div class="mb-8 shrink-0">
              <h2 class="text-3xl font-bold text-[#E8F0FF] mb-2">产品需求文档</h2>
              <p class="text-[#8A97AA]">产品需求与规格说明</p>
            </div>
            <div class="flex-1 overflow-auto pb-4">
              <div class="flex gap-6 h-full min-w-max px-2">
                <For each={mockDocs}>
                  {(doc) => (
                    <button
                      type="button"
                      onMouseDown={(e) => e.stopPropagation()}
                      onClick={(e) => {
                        e.stopPropagation()
                        openFile(doc)
                      }}
                      class="text-left group cursor-pointer w-[400px] h-full flex flex-col rounded-xl border border-[#00F0FF]/10 bg-[#1A1F3A]/90 hover:border-[#00F0FF]/50 hover:bg-[#00F0FF]/5 transition-all relative overflow-hidden z-[100] pointer-events-auto"
                    >
                      <div class="p-4 border-b border-[#00F0FF]/10 bg-[#1A1F3A] flex items-center gap-3 shrink-0">
                        <div class="w-8 h-8 rounded bg-[#00F0FF]/10 flex items-center justify-center text-[#00F0FF]">
                          <iconify-icon icon="lucide:file-text" class="text-lg"></iconify-icon>
                        </div>
                        <div class="flex-1 min-w-0">
                          <h3 class="text-[#E8F0FF] font-medium text-sm truncate">{doc.title}</h3>
                          <p class="text-[#8A97AA] text-[10px] truncate">{doc.filename}</p>
                        </div>
                      </div>

                      <div class="flex-1 p-6 overflow-hidden bg-[#141829]/50 relative">
                        <div class="absolute inset-0 p-6 text-[10px] text-[#8A97AA] font-mono leading-relaxed whitespace-pre-wrap opacity-70 group-hover:opacity-100 transition-opacity">
                          {doc.content}
                        </div>
                        <div class="absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t from-[#1A1F3A] to-transparent pointer-events-none"></div>
                      </div>
                    </button>
                  )}
                </For>
              </div>
            </div>
          </div>
        </Show>

        {/* Wireframe Folder View */}
        <Show when={currentView() === "folder" && activeFolder() === "wireframe"}>
          <div class="flex grow min-w-0 h-full overflow-auto">
            <main class="flex flex-col grow shrink p-6 overflow-auto">
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mt-4 min-w-max">
                <For each={mockWireframeFiles}>
                  {(file) => (
                    <div
                      onClick={() => openFile(file)}
                      class="group cursor-pointer p-4 rounded-xl border border-[#00F0FF]/10 bg-[#1A1F3A]/90 hover:border-[#00F0FF]/50 hover:bg-[#00F0FF]/5 transition-all flex flex-col gap-4"
                    >
                      <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-lg bg-[#00F0FF]/10 flex items-center justify-center text-[#00F0FF]">
                          <iconify-icon icon="lucide:file-text" class="text-xl"></iconify-icon>
                        </div>
                        <div class="flex-1 min-w-0">
                          <h3 class="text-[#E8F0FF] font-medium truncate">{file.filename}</h3>
                          <p class="text-[#8A97AA] text-xs">{file.lastModified}</p>
                        </div>
                      </div>
                      <div class="h-32 bg-[#141829] rounded-lg p-3 overflow-hidden opacity-60 group-hover:opacity-100 transition-opacity">
                        <div class="w-full h-full overflow-hidden text-[10px] font-mono text-[#8A97AA] leading-tight whitespace-pre-wrap">
                          {file.content}
                        </div>
                      </div>
                    </div>
                  )}
                </For>
              </div>
            </main>
          </div>
        </Show>

        {/* Enhanced File Detail View */}
        <Show when={currentView() === "file" && activeFile()}>
          <div class="w-full h-full p-0 flex flex-col bg-[#141829]">
            <div class="border-b border-[#00F0FF]/10 p-6 flex items-center gap-4 bg-[#1A1F3A]/50">
              <div
                class={`w-12 h-12 rounded-xl flex items-center justify-center text-2xl ${
                  activeFile().type === "html"
                    ? "bg-[#B026FF]/10 text-[#B026FF]"
                    : activeFile().type === "flowchart"
                      ? "bg-[#FF006E]/10 text-[#FF006E]"
                      : activeFile().type === "markdown"
                        ? "bg-[#00F0FF]/10 text-[#00F0FF]"
                      : "bg-[#00F0FF]/10 text-[#00F0FF]"
                }`}
              >
                <iconify-icon
                  icon={
                    activeFile().type === "html"
                      ? "lucide:layout"
                      : activeFile().type === "flowchart"
                        ? "lucide:workflow"
                        : activeFile().type === "markdown"
                          ? "lucide:file-text"
                        : "lucide:file-text"
                  }
                ></iconify-icon>
              </div>
              <div>
                <h1 class="text-2xl font-bold text-[#E8F0FF]">{activeFile().title}</h1>
                <p class="text-[#8A97AA] text-sm font-mono">
                  {activeFile().filename} • {activeFile().lastModified}
                </p>
              </div>
            </div>

            <div class="flex-1 overflow-y-auto p-0 w-full relative">
              <Show when={activeFile().type === "html"}>
                <div class="w-full h-full bg-white">
                  <iframe srcdoc={activeFile().content} class="w-full h-full border-none" title="Preview" />
                </div>
              </Show>

              <Show when={activeFile().type === "markdown"}>
                <div class="p-8 max-w-4xl mx-auto w-full">
                  <div class="text-[#E8F0FF]" innerHTML={parseMarkdown(activeFile().content)}></div>
                </div>
              </Show>

              <Show when={activeFile().type === "flowchart"}>
                <div
                  style="background: linear-gradient(135deg, rgba(10, 14, 26, 1) 0%, rgba(26, 19, 50, 1) 50%, rgba(13, 27, 42, 1) 100%), radial-gradient(circle at 20% 30%, rgba(0, 240, 255, 0.101961) 0%, transparent 50%), radial-gradient(circle at 80% 70%, rgba(255, 0, 110, 0.101961) 0%, transparent 50%);"
                  class="flex grow min-w-0 h-full overflow-y-auto"
                >
                  <main
                    style="flex-basis: 0%; padding: 1.5rem 2rem;"
                    class="overflow-x-hidden flex flex-col grow shrink"
                  >
                    <div class="flex justify-between items-center mb-8">
                      <div class="flex items-center gap-y-4 gap-x-4">
                        <button 
                          onClick={() => setActiveFile(null)}
                          class="hover:bg-[#00F0FF]/10 p-2 rounded-lg transition-colors text-[#00F0FF]"
                        >
                          <iconify-icon icon="lucide:arrow-left" class="text-xl"></iconify-icon>
                        </button>
                        <h1 style="color: rgba(232, 240, 255, 1);" class="text-2xl font-bold">
                          用户故事
                        </h1>
                      </div>
                    </div>
                    <div class="w-full flex-1 bg-[#141829] relative overflow-hidden rounded-xl border border-[#00F0FF]/20 min-h-[500px]">
                      <iframe 
                        srcdoc={getMermaidHtml(activeFile().content)} 
                        class="w-full h-full border-none" 
                        title="Flowchart" 
                      />
                      <div class="absolute bottom-4 right-4 p-2 bg-[#1A1F3A]/80 backdrop-blur rounded-lg border border-[#00F0FF]/20 text-[#00F0FF] text-xs flex items-center gap-2">
                        <iconify-icon icon="simple-icons:mermaid"></iconify-icon>
                        <span>Mermaid Renderer</span>
                      </div>
                    </div>
                  </main>
                </div>
              </Show>

              <Show when={!activeFile().type}>
                <div class="p-8 max-w-4xl mx-auto w-full prose prose-invert prose-cyan max-w-none">
                  <div class="text-[#E8F0FF] whitespace-pre-wrap font-sans leading-7">{activeFile().content}</div>
                </div>
              </Show>
            </div>
          </div>
        </Show>

        <Show when={previewFile()}>
          <div
            class="fixed inset-0 z-[999] bg-black/80 flex flex-col items-center justify-center animate-fade-in"
            onClick={() => setPreviewFile(null)}
          >
            <div class="relative w-[70%] h-[80%] flex gap-4" onClick={(e) => e.stopPropagation()}>
              <Show when={!previewFile()?.kind || previewFile()?.kind === "html"}>
                <div class="flex-1 bg-white rounded-xl overflow-hidden shadow-2xl border border-[#00F0FF]/20">
                  <iframe
                    id="preview-iframe"
                    srcdoc={getPreviewContent(previewFile().content)}
                    class="w-full h-full border-none bg-white"
                  />
                </div>
              </Show>
              <Show when={previewFile()?.kind === "md"}>
                <div
                  class={
                    theme() === "light"
                      ? "flex-1 flex flex-col bg-white rounded-xl overflow-hidden shadow-2xl border border-slate-200"
                      : "flex-1 flex flex-col bg-[#141829] rounded-xl overflow-hidden shadow-2xl border border-[#00F0FF]/20"
                  }
                >
                  <div
                    class={
                      theme() === "light"
                        ? "px-6 py-4 border-b border-slate-200 flex items-center justify-between bg-slate-50"
                        : "px-6 py-4 border-b border-[#00F0FF]/10 flex items-center justify-between bg-[#050816]/80"
                    }
                  >
                    <div>
                      <div
                        class={
                          theme() === "light"
                            ? "text-sm font-semibold text-slate-900"
                            : "text-sm font-semibold text-[#E8F0FF]"
                        }
                      >
                        编辑 Markdown
                      </div>
                      <div
                        class={
                          theme() === "light"
                            ? "text-xs text-slate-500 mt-1 truncate max-w-md"
                            : "text-xs text-[#8A97AA] mt-1 truncate max-w-md"
                        }
                      >
                        {previewFile().name}
                      </div>
                    </div>
                  </div>
                  <div class="flex-1 p-3">
                    <textarea
                      class={
                        theme() === "light"
                          ? "w-full h-full bg-white border border-slate-200 focus:border-slate-400 outline-none rounded-lg px-3 py-2 text-sm font-mono text-slate-900 leading-relaxed resize-none overflow-auto"
                          : "w-full h-full bg-[#050816] border border-[#00F0FF]/15 focus:border-[#00F0FF]/60 outline-none rounded-lg px-3 py-2 text-sm font-mono text-[#E8F0FF] leading-relaxed resize-none overflow-auto"
                      }
                      value={editValue()}
                      onInput={(e) => setEditValue(e.currentTarget.value)}
                      spellcheck={false}
                    />
                  </div>
                  <div
                    class={
                      theme() === "light"
                        ? "px-6 py-4 border-t border-slate-200 bg-slate-50 flex justify-end gap-3"
                        : "px-6 py-4 border-t border-[#00F0FF]/10 bg-[#050816]/80 flex justify-end gap-3"
                    }
                  >
                    <button
                      type="button"
                      class={
                        theme() === "light"
                          ? "px-4 py-2 text-sm rounded-lg border border-transparent text-slate-500 hover:text-slate-900 hover:bg-slate-100 transition-colors"
                          : "px-4 py-2 text-sm rounded-lg border border-transparent text-[#8A97AA] hover:text-[#E8F0FF] hover:bg-white/5 transition-colors"
                      }
                      onClick={(e) => {
                        e.stopPropagation()
                        setPreviewFile(null)
                      }}
                    >
                      取消
                    </button>
                    <button
                      type="button"
                      class={
                        theme() === "light"
                          ? "px-4 py-2 text-sm rounded-lg bg-[#0F172A] text-white font-medium hover:bg-[#020617] transition-colors"
                          : "px-4 py-2 text-sm rounded-lg bg-[#00F0FF] text-[#050816] font-medium hover:bg-[#33F2FF] transition-colors"
                      }
                      onClick={(e) => {
                        e.stopPropagation()
                        void saveMarkdown()
                      }}
                    >
                      保存
                    </button>
                  </div>
                </div>
              </Show>
              <Show when={selectedHtml() && (!previewFile()?.kind || previewFile()?.kind === "html")}>
                <div class="w-[400px] flex flex-col bg-[#1A1F3A] rounded-xl overflow-hidden border border-[#00F0FF]/20 shadow-2xl animate-fade-in-right">
                  <div class="p-4 border-b border-[#00F0FF]/10 flex justify-between items-center bg-[#141829]/50">
                    <h3 class="text-[#E8F0FF] font-semibold text-sm">Selected Element</h3>
                    <button
                      onClick={() => setSelectedHtml(null)}
                      class="text-[#8A97AA] hover:text-white"
                    >
                      <iconify-icon icon="lucide:x" />
                    </button>
                  </div>
                  <div class="flex-1 overflow-auto p-4 font-mono text-xs text-[#E8F0FF] whitespace-pre-wrap">
                    {selectedHtml()}
                  </div>
                </div>
              </Show>
            </div>
          </div>
        </Show>
      </div>
    </main>
  )
}

export default Workspace
