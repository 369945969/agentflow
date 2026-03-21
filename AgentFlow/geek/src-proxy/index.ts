import path from "node:path"
import fs from "node:fs/promises"
import type { Dirent } from "node:fs"
import { Database } from "bun:sqlite"

export type ProxyConfig = {
  baseUrl?: string
  host?: string
  port?: number
  username?: string
  password?: string
  pmRoot?: string
  pmDepth?: number
  pmMaxBytes?: number
  workspace?: string
}

type TreeFile = {
  path: string
  name: string
  type: "file"
  size: number
  updatedAt: number
  content?: string
}

type TreeDir = {
  path: string
  name: string
  type: "dir"
  children: TreeItem[]
}

type TreeItem = TreeFile | TreeDir

type AppItem = {
  uuid: string
  sessionId: string
  appName: string
  path: string
}

type UserRow = {
  id: string
  root: string
  current_uuid: string | null
}

const defaults = {
  baseUrl: process.env.OPENCODE_BASE_URL ?? "http://127.0.0.1:2345",
  host: process.env.PROXY_HOST ?? "0.0.0.0",
  port: Number(process.env.PROXY_PORT ?? "4097"),
  username: process.env.OPENCODE_SERVER_USERNAME ?? "opencode",
  password: process.env.OPENCODE_SERVER_PASSWORD ?? "123",
  pmRoot: process.env.PM_ROOT ?? process.cwd(),
  pmDepth: Number(process.env.PM_DEPTH ?? "3"),
  pmMaxBytes: Number(process.env.PM_MAX_BYTES ?? "1048576"),
  workspace: process.env.OPENCODE_WORKSPACE ?? path.join(process.cwd(), "workspace"),
}
const defaultProviderID = process.env.OPENCODE_DEFAULT_PROVIDER_ID ?? "geek_code"
const defaultModelID = process.env.OPENCODE_DEFAULT_MODEL_ID ?? "agent"

const text = new TextEncoder()
const dbPath = path.join(process.cwd(), "src-proxy", "store.sqlite")
const appsBase = path.join(process.cwd(), "apps")
const db = new Database(dbPath, { create: true })
const webLogPath = path.join(process.cwd(), "web.log")
const appLogPath = path.join(process.cwd(), "app.log")

db.exec(
  "create table if not exists users (id text primary key, root text not null, current_uuid text)",
)
db.exec(
  "create table if not exists apps (uuid text primary key, user_id text not null, session_id text not null, app_name text not null, path text not null, created_at integer not null)",
)
db.exec(
  "create table if not exists session_skills (session_id text primary key, skill text, injected_at integer)",
)

const readTree = async (root: string, depth: number, maxDepth: number, maxBytes: number, relativeTo?: string) => {
  const entries = await fs.readdir(root, { withFileTypes: true }).catch(() => [])
  const base = relativeTo ?? process.cwd()
  const list = await Promise.all(
    entries.map(async (entry: Dirent) => {
      const full = path.join(root, entry.name)
      const rel = path.relative(base, full)
      const stat = await fs.stat(full).catch(() => null)
      if (!stat) return null
      if (entry.isDirectory()) {
        if (depth >= maxDepth) {
          const item: TreeDir = { path: rel, name: entry.name, type: "dir", children: [] }
          return item
        }
        const children = await readTree(full, depth + 1, maxDepth, maxBytes, base)
        const item: TreeDir = { path: rel, name: entry.name, type: "dir", children }
        return item
      }
      const file = Bun.file(full)
      const size = file.size
      const ext = path.extname(entry.name).toLowerCase()
      const isText = [".md", ".txt", ".json", ".yaml", ".yml", ".js", ".ts", ".tsx", ".css", ".html"].includes(ext)
      const content = isText && size <= maxBytes ? await file.text() : undefined
      const item: TreeFile = { path: rel, name: entry.name, type: "file", size, updatedAt: stat.mtimeMs, content }
      return item
    }),
  )
  return list.filter((item): item is TreeItem => Boolean(item))
}

const normalizeRoot = (root: string) => {
  if (path.isAbsolute(root)) return root
  return path.join(process.cwd(), root)
}

const ensureUser = async (user: string) => {
  const row = db
    .query("select id, root, current_uuid from users where id = ?")
    .get(user) as UserRow | undefined
  if (row) {
    const root = normalizeRoot(row.root)
    if (root !== row.root) {
      db.query("update users set root = ? where id = ?").run(root, user)
    }
    await fs.mkdir(root, { recursive: true })
    return { id: row.id, root, current: row.current_uuid ?? "" }
  }
  const root = path.join(appsBase, user)
  db.query("insert into users (id, root, current_uuid) values (?, ?, '')").run(user, root)
  await fs.mkdir(root, { recursive: true })
  return { id: user, root, current: "" }
}

const parseCookies = (value: string) => {
  const map: Record<string, string> = {}
  if (!value) return map
  value.split(";").forEach((part) => {
    const index = part.indexOf("=")
    if (index < 0) return
    const key = part.slice(0, index).trim()
    const val = part.slice(index + 1).trim()
    if (!key) return
    map[key] = decodeURIComponent(val)
  })
  return map
}

const setCookie = (headers: Headers, key: string, value: string) => {
  const line = `${key}=${encodeURIComponent(value)}; Path=/; SameSite=Lax`
  headers.append("Set-Cookie", line)
}

const createAuth = (username: string, password: string) => {
  if (!username || !password) return ""
  return `Basic ${Buffer.from(`${username}:${password}`).toString("base64")}`
}

const stamp = () => new Date().toISOString()

const logLine = async (filePath: string, line: string) => {
  await fs.appendFile(filePath, line + "\n").catch(() => {})
}

const logWeb = async (line: string) => logLine(webLogPath, line)
const logApp = async (line: string) => logLine(appLogPath, line)

const createLoggedStream = (res: Response, label: string) => {
  if (!res.body) {
    void logApp(`${stamp()} ${label} EMPTY`)
    return null
  }
  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  return new ReadableStream<Uint8Array>({
    async pull(controller) {
      const result = await reader.read().catch(() => null)
      if (!result || result.done) {
        controller.close()
        void logApp(`${stamp()} ${label} DONE`)
        return
      }
      const chunkText = decoder.decode(result.value, { stream: true })
      if (chunkText) void logApp(`${stamp()} ${label} ${chunkText}`)
      controller.enqueue(result.value)
    },
    cancel() {
      reader.cancel().catch(() => {})
    },
  })
}

const baseHeaders = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
  "Access-Control-Allow-Headers": "*",
}

export const createProxy = (input: ProxyConfig = {}) => {
  const config = { ...defaults, ...input }
  const clients = new Set<ReadableStreamDefaultController<Uint8Array>>()
  const auth = createAuth(config.username, config.password)
  let backendConnected = false

  const broadcast = (data: string) => {
    const payload = text.encode(data)
    clients.forEach((client) => {
      try {
        client.enqueue(payload)
      } catch {
        clients.delete(client)
      }
    })
  }

  const sendStatus = () => {
    broadcast(
      `data: ${JSON.stringify({ type: "proxy.status", connected: backendConnected, ts: Date.now() })}\n\n`,
    )
  }

  const setBackendConnected = (value: boolean) => {
    if (backendConnected === value) return
    backendConnected = value
    sendStatus()
  }

  const connectEvents = async () => {
    const headers = new Headers({ Accept: "text/event-stream" })
    if (auth) headers.set("Authorization", auth)
    try {
      await logApp(`${stamp()} CONNECT ${config.baseUrl}/event`)
      const res = await fetch(`${config.baseUrl}/event`, { headers })
      if (!res.ok || !res.body) {
        await logApp(`${stamp()} CONNECT_ERROR ${config.baseUrl}/event ${res.status}`)
        setBackendConnected(false)
        setTimeout(connectEvents, 2000)
        return
      }
      await logApp(`${stamp()} CONNECTED ${config.baseUrl}/event`)
      setBackendConnected(true)
      const reader = res.body.getReader()
      for (;;) {
        const result = await reader.read().catch(() => null)
        if (!result || result.done) {
          await logApp(`${stamp()} DISCONNECTED ${config.baseUrl}/event`)
          setBackendConnected(false)
          setTimeout(connectEvents, 2000)
          return
        }
        const chunk = new TextDecoder().decode(result.value)
        const snippet = chunk.length > 2000 ? chunk.slice(0, 2000) + "..." : chunk
        await logApp(`${stamp()} EVENT_IN ${snippet}`)
        broadcast(chunk)
      }
    } catch {
      await logApp(`${stamp()} CONNECT_ERROR ${config.baseUrl}/event`)
      setBackendConnected(false)
      setTimeout(connectEvents, 2000)
    }
  }

  const sseResponse = () => {
    const stream = new ReadableStream<Uint8Array>({
      start(controller) {
        clients.add(controller)
        logWeb(`${stamp()} SSE_CONNECT /events`)
        try {
          controller.enqueue(text.encode(`data: ${JSON.stringify({ type: "connected" })}\n\n`))
          controller.enqueue(
            text.encode(
              `data: ${JSON.stringify({ type: "proxy.status", connected: backendConnected, ts: Date.now() })}\n\n`,
            ),
          )
        } catch {
          clients.delete(controller)
          return
        }
        const timer = setInterval(() => {
          try {
            controller.enqueue(text.encode(":\n\n"))
          } catch {
            clearInterval(timer)
            clients.delete(controller)
          }
        }, 15000)
        return () => {
          clearInterval(timer)
          clients.delete(controller)
          logWeb(`${stamp()} SSE_DISCONNECT /events`)
        }
      },
      cancel() {},
    })
    return new Response(stream, {
      headers: {
        "Content-Type": "text/event-stream",
        "Cache-Control": "no-cache",
        Connection: "keep-alive",
        ...baseHeaders,
      },
    })
  }

  const pmResponse = async () => {
    const roots = [
      "01-Global-Context",
      "02-Feature-Plan",
      "03-Style-Guide",
      "04-Screen-Prototype",
    ]
    const data = await Promise.all(
      roots.map(async (name) => {
        const full = path.join(config.pmRoot, name)
        const exists = await fs.stat(full).then(() => true).catch(() => false)
        if (!exists) return { name, path: full, exists, items: [] as TreeItem[] }
        const items = await readTree(full, 0, config.pmDepth, config.pmMaxBytes)
        return { name, path: full, exists, items }
      }),
    )
    return Response.json({ root: config.pmRoot, folders: data }, { headers: baseHeaders })
  }

  const appsList = async (req: Request) => {
    const cookies = parseCookies(req.headers.get("cookie") ?? "")
    const user = cookies.user ?? "default"
    const userRow = await ensureUser(user)
    const rows = db
      .query(
        "select uuid, session_id as sessionId, app_name as appName, path from apps where user_id = ? order by created_at desc",
      )
      .all(user) as AppItem[]
    const cookieUuid = cookies.app_uuid ?? ""
    const cookieApp = rows.find((item) => item.uuid === cookieUuid)
    const current = cookieApp?.uuid ?? userRow.current ?? rows[0]?.uuid ?? ""
    if (current !== userRow.current) {
      db.query("update users set current_uuid = ? where id = ?").run(current, user)
    }
    const headers = new Headers(baseHeaders)
    if (current) {
      const active = rows.find((item) => item.uuid === current)
      if (active) {
        setCookie(headers, "app_uuid", active.uuid)
        setCookie(headers, "app_session", active.sessionId)
        setCookie(headers, "app_name", active.appName)
      }
    }
    return Response.json(
      {
        user,
        root: userRow.root,
        apps: rows,
        current,
      },
      { headers },
    )
  }

  const appsCreate = async (req: Request) => {
    const body = await req.json().catch(() => null)
    const appName =
      typeof body?.appName === "string"
        ? body.appName
        : typeof body?.name === "string"
          ? body.name
          : ""
    if (!appName) {
      return Response.json({ error: "appName required" }, { status: 400, headers: baseHeaders })
    }
    await logApp(`${stamp()} APP_CREATE_REQ ${JSON.stringify({ appName })}`)
    const cookies = parseCookies(req.headers.get("cookie") ?? "")
    const user = cookies.user ?? "default"
    await ensureUser(user)
    const uuid = crypto.randomUUID()
    const seed = `ses_${uuid}`
    const perm = [
      { permission: "edit", pattern: "*", action: "allow" },
      { permission: "bash", pattern: "*", action: "allow" },
      { permission: "read", pattern: "*", action: "allow" },
    ]
    const headers = new Headers({ "Content-Type": "application/json" })
    if (auth) headers.set("Authorization", auth)
    const sessionBody = { title: `${appName} (${seed})`, permission: perm }
    await logApp(`${stamp()} SESSION_CREATE_REQ ${JSON.stringify(sessionBody)}`)
    const res = await fetch(`${config.baseUrl}/session`, {
      method: "POST",
      headers,
      body: JSON.stringify(sessionBody),
    })
    const dataText = await res.text().catch(() => "")
    await logApp(`${stamp()} SESSION_CREATE_RES ${res.status} ${dataText}`)
    const data = dataText ? JSON.parse(dataText) : null
    if (!res.ok || !data?.id) {
      return Response.json({ error: "session create failed" }, { status: 500, headers: baseHeaders })
    }
    // User requested: Use config.workspace + sessionId
    const appPath = path.join(config.workspace, data.id)
    await fs.mkdir(appPath, { recursive: true })
    const app: AppItem = {
      uuid,
      sessionId: data.id,
      appName,
      path: appPath,
    }
    await logApp(`${stamp()} APP_CREATE_SESSION ${JSON.stringify({ uuid: app.uuid, sessionId: app.sessionId })}`)
    db.query(
      "insert into apps (uuid, user_id, session_id, app_name, path, created_at) values (?, ?, ?, ?, ?, ?)",
    ).run(app.uuid, user, app.sessionId, app.appName, app.path, Date.now())
    db.query("update users set current_uuid = ? where id = ?").run(app.uuid, user)
    const out = new Headers(baseHeaders)
    setCookie(out, "app_uuid", app.uuid)
    setCookie(out, "app_session", app.sessionId)
    setCookie(out, "app_name", app.appName)
    return Response.json(app, { headers: out })
  }

  const generateTitle = async (req: Request) => {
    const body = await req.json().catch(() => null)
    const textValue = body?.text ?? ""
    if (!textValue) return Response.json({ title: "New Project" }, { headers: baseHeaders })
    
    const headers = new Headers({ "Content-Type": "application/json" })
    if (auth) headers.set("Authorization", auth)
    
    // Create temp session
    const sessionRes = await fetch(`${config.baseUrl}/session`, {
        method: "POST",
        headers,
        body: JSON.stringify({ title: "TitleGen" })
    })
    const sessionData = await sessionRes.json()
    const sessionId = sessionData?.id
    
    if (!sessionId) return Response.json({ title: "Project" }, { headers: baseHeaders })
    
    // Add instruction to avoid repeating user input
    const prompt = `请分析以下用户输入，并用最简短的语言（5-10个字）概括其核心意图，作为项目标题。注意：绝对不要直接重复用户的输入，必须进行概括和重写。只返回标题文本。输入内容：${textValue}`
    
    
    let resultText = ""
    try {
        const promptRes = await fetch(`${config.baseUrl}/session/${sessionId}/prompt_async`, {
            method: "POST",
            headers,
            body: JSON.stringify({
                model: { providerID: defaultProviderID, modelID: defaultModelID },
                parts: [{ type: "text", text: prompt }],
            })
        })
    
    const reader = promptRes.body?.getReader()
    if (reader) {
        const decoder = new TextDecoder()
        while (true) {
            const { done, value } = await reader.read()
            if (done) break
            const chunk = decoder.decode(value)
            const lines = chunk.split("\n")
            for (const line of lines) {
                if (line.startsWith("data: ")) {
                    try {
                        const json = JSON.parse(line.slice(6))
                        console.log("Chunk:", JSON.stringify(json)) // Debug log
                        
                        // Handle text content
                        if (json.properties?.delta && json.properties.part?.type === "text") {
                            resultText += json.properties.delta
                        }
                    } catch {}
                }
            }
        }
    }
    } catch (e) {
        console.error("Generate title error:", e)
    }
    
    let title = resultText.trim().replace(/^["']|["']$/g, "").replace(/\n/g, " ")
    if (title.length > 20) title = title.slice(0, 20)
    
    // Cleanup
    try {
        await fetch(`${config.baseUrl}/session/${sessionId}`, { method: "DELETE", headers })
    } catch {}

    return Response.json({ title }, { headers: baseHeaders })
  }

  const ensureSessionDir = async (req: Request) => {
    const url = new URL(req.url)
    const parts = url.pathname.split("/")
    const sessionIdx = parts.indexOf("session")
    const sessionId = parts[sessionIdx + 1]
    if (!sessionId) return Response.json({ error: "No session ID" }, { status: 400, headers: baseHeaders })

    const appPath = path.join(config.workspace, sessionId)
    await fs.mkdir(appPath, { recursive: true })
    return Response.json({ status: "ok", path: appPath }, { headers: baseHeaders })
  }

  const readWorkspaceFile = async (req: Request) => {
    const url = new URL(req.url)
    const rawPath = url.searchParams.get("path") ?? ""
    const sessionId = url.searchParams.get("session") ?? ""
    
    if (!rawPath) {
      return Response.json({ error: "path required" }, { status: 400, headers: baseHeaders })
    }

    const trimmed = rawPath.startsWith("/") ? rawPath.slice(1) : rawPath
    const workspaceRoot = path.resolve(config.workspace)
    
    let full: string
    if (sessionId) {
      // If session is provided, path is relative to session folder
      full = path.join(workspaceRoot, sessionId, trimmed)
    } else {
      // Legacy behavior or absolute-ish path
      const workspacePrefix = "workspace/"
      const relative =
        trimmed.startsWith(workspacePrefix) ? trimmed.slice(workspacePrefix.length) : trimmed
      full = path.join(workspaceRoot, relative)
    }

    const resolvedFull = path.resolve(full)
    if (!resolvedFull.startsWith(workspaceRoot)) {
      return Response.json({ error: "forbidden" }, { status: 403, headers: baseHeaders })
    }
    try {
      const file = Bun.file(resolvedFull)
      const size = file.size
      const ext = path.extname(resolvedFull).toLowerCase()
      const isText =
        [".md", ".txt", ".json", ".yaml", ".yml", ".js", ".ts", ".tsx", ".css", ".html"].includes(
          ext,
        )
      const content = isText ? await file.text() : ""
      return Response.json({ path: rawPath, ext, size, content }, { headers: baseHeaders })
    } catch {
      return Response.json({ error: "not found" }, { status: 404, headers: baseHeaders })
    }
  }

  const writeWorkspaceFile = async (req: Request) => {
    const body = await req.json().catch(() => null) as { path?: string; content?: string } | null
    if (!body?.path || typeof body.path !== "string") {
      return Response.json({ error: "path required" }, { status: 400, headers: baseHeaders })
    }
    const content = typeof body.content === "string" ? body.content : ""
    const trimmed = body.path.startsWith("/") ? body.path.slice(1) : body.path
    const workspaceRoot = path.resolve(config.workspace)
    const workspacePrefix = "workspace/"
    const relative =
      trimmed.startsWith(workspacePrefix) ? trimmed.slice(workspacePrefix.length) : trimmed
    const full = path.join(workspaceRoot, relative)
    const resolvedFull = path.resolve(full)
    if (!resolvedFull.startsWith(workspaceRoot)) {
      return Response.json({ error: "forbidden" }, { status: 403, headers: baseHeaders })
    }
    try {
      const dir = path.dirname(full)
      await fs.mkdir(dir, { recursive: true })
      await fs.writeFile(full, content, "utf-8")
      return Response.json({ path: body.path, saved: true }, { headers: baseHeaders })
    } catch {
      return Response.json({ error: "write failed" }, { status: 500, headers: baseHeaders })
    }
  }

  const workspaceTree = async (req: Request) => {
    const url = new URL(req.url)
    const sessionId = url.searchParams.get("session") ?? ""
    if (!sessionId) {
      return Response.json({ error: "session required" }, { status: 400, headers: baseHeaders })
    }
    const workspaceRoot = path.resolve(config.workspace)
    const root = path.join(workspaceRoot, sessionId)
    const exists = await fs.stat(root).then(() => true).catch(() => false)
    if (!exists) {
      return Response.json({ root, items: [] as TreeItem[] }, { headers: baseHeaders })
    }
    const items = await readTree(root, 0, config.pmDepth, config.pmMaxBytes, root)
    return Response.json({ root, items }, { headers: baseHeaders })
  }

  const workspaceExport = async (req: Request) => {
    const url = new URL(req.url)
    const sessionId = url.searchParams.get("session") ?? ""
    if (!sessionId) {
      return Response.json({ error: "session required" }, { status: 400, headers: baseHeaders })
    }
    const workspaceRoot = path.resolve(config.workspace)
    const root = path.join(workspaceRoot, sessionId)
    const exists = await fs.stat(root).then(() => true).catch(() => false)
    if (!exists) {
      return Response.json({ error: "not found" }, { status: 404, headers: baseHeaders })
    }
    const proc = Bun.spawn(["zip", "-r", "-", "."], {
      cwd: root,
      stderr: "ignore",
      stdout: "pipe",
    })
    const stream = proc.stdout
    if (!stream) {
      return Response.json({ error: "export failed" }, { status: 500, headers: baseHeaders })
    }
    const headers = {
      ...baseHeaders,
      "Content-Type": "application/zip",
      "Content-Disposition": `attachment; filename="${sessionId}.zip"`,
    }
    return new Response(stream, { headers })
  }

  const proxyFetch = async (req: Request) => {
    const url = new URL(req.url)
    const target = new URL(url.pathname + url.search, config.baseUrl)
    const headers = new Headers(req.headers)
    headers.delete("host")
    headers.delete("connection")
    headers.delete("content-length")
    if (auth && !headers.get("authorization")) headers.set("authorization", auth)
    let bodyText = ""
    if (req.method !== "GET" && req.method !== "HEAD") {
      bodyText = await req.clone().text().catch(() => "")
    }
    const contentType = headers.get("content-type") ?? ""
    if (bodyText && contentType.includes("application/json")) {
      try {
        const parsed = JSON.parse(bodyText)
        const isMessage =
          req.method === "POST" &&
          target.pathname.startsWith("/session/") &&
          (target.pathname.endsWith("/message") || target.pathname.endsWith("/prompt_async"))
        if (isMessage) {
          if (!parsed?.model) {
            parsed.model = { providerID: defaultProviderID, modelID: defaultModelID }
          }
          // User requested: Inject system rules for directory restriction
          try {
            const parts = target.pathname.split("/")
            const sessionIdx = parts.indexOf("session")
            if (sessionIdx >= 0 && parts.length > sessionIdx + 1) {
              const sid = parts[sessionIdx + 1]
              const row = db.query("select path from apps where session_id = ?").get(sid) as { path: string } | undefined
                  if (row?.path) {
                    let rule = `[SYSTEM RULE]: You are strictly limited to operating within the directory: ${row.path}. You must not create, modify, or delete any files outside of this directory.`
                    
                    const skillInjected = db.query("select 1 from session_skills where session_id = ?").get(sid)
                    
                    if (!skillInjected) {
                      const skillPath = path.join(process.cwd(), "src-proxy", "product-manager", "SKILL.md")
                      const skillContent = await fs.readFile(skillPath, "utf-8").catch(() => "")
                      if (skillContent) {
                        rule += `\n\n[ROLE DEFINITION]: You are the Product Manager Agent. You MUST follow these rules and pipeline:\n${skillContent}\n\n[SYSTEM NOTE]: This is a role definition. Do NOT attempt to invoke a tool named "skill" or "product-manager". Use your standard tools (like write_file) to produce the outputs.`
                        
                        try {
                           db.query("insert or ignore into session_skills (session_id, skill, injected_at) values (?, ?, ?)").run(sid, "product-manager", Date.now())
                        } catch (e) {
                           console.error("DB Insert Error", e)
                        }
                      }
                    }

                    if (Array.isArray(parsed.parts)) {
                  const textPart = parsed.parts.find((p: any) => p.type === "text")
                  if (textPart) {
                    textPart.text = rule + "\n\n" + textPart.text
                  } else {
                    parsed.parts.unshift({ type: "text", text: rule })
                  }
                }
              }
            }
          } catch (e) {
            console.error("Inject system rule error:", e)
          }
        }
        bodyText = JSON.stringify(parsed)
      } catch {}
    }
    const bodySnippet =
      bodyText.length > 2000 ? bodyText.slice(0, 2000) + "..." : bodyText
    await logApp(`${stamp()} REQ ${req.method} ${target.pathname}${target.search} ${bodySnippet}`)
    try {
      const res = await fetch(target, {
        method: req.method,
        headers,
        body: bodyText || undefined,
        redirect: "manual",
      })
      await logApp(`${stamp()} RES ${req.method} ${target.pathname}${target.search} ${res.status}`)
      const out = new Headers(res.headers)
      out.set("Access-Control-Allow-Origin", "*")
      out.set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
      out.set("Access-Control-Allow-Headers", "*")
      const stream = createLoggedStream(res, `RES_STREAM ${req.method} ${target.pathname}${target.search}`)
      return new Response(stream ?? res.body, { status: res.status, headers: out })
    } catch (error) {
      const message = error instanceof Error ? error.message : "unknown error"
      await logApp(
        `${stamp()} RES_ERROR ${req.method} ${target.pathname}${target.search} ${message}`,
      )
      return Response.json({ error: message }, { status: 502, headers: baseHeaders })
    }
  }

  const fetchHandler = async (req: Request) => {
    if (req.method === "OPTIONS") return new Response(null, { status: 200, headers: baseHeaders })
    const url = new URL(req.url)
    await logWeb(`${stamp()} REQ ${req.method} ${url.pathname}${url.search}`)
    if (url.pathname === "/proxy/status" && req.method === "GET") {
      return Response.json(
        { connected: backendConnected, ts: Date.now() },
        { headers: baseHeaders },
      )
    }
    if (url.pathname === "/events") return sseResponse()
    if (url.pathname === "/pm/dirs") return pmResponse()
    if (url.pathname === "/apps" && req.method === "GET") return appsList(req)
    if (url.pathname === "/apps" && req.method === "POST") return appsCreate(req)
    if (url.pathname === "/generate-title" && req.method === "POST") return generateTitle(req)
    if (url.pathname.endsWith("/ensure") && req.method === "POST") return ensureSessionDir(req)
    if (url.pathname === "/workspace/file" && req.method === "GET") return readWorkspaceFile(req)
    if (url.pathname === "/workspace/file" && req.method === "POST") return writeWorkspaceFile(req)
    if (url.pathname === "/workspace/tree" && req.method === "GET") return workspaceTree(req)
    if (url.pathname === "/workspace/export" && req.method === "GET") return workspaceExport(req)
    return proxyFetch(req)
  }

  const start = () => {
    connectEvents()
    setInterval(() => sendStatus(), 2000)
    return Bun.serve({
      port: config.port,
      hostname: config.host,
      idleTimeout: 0,
      fetch: fetchHandler,
    })
  }

  const request = (pathValue: string, init: RequestInit = {}) => {
    const headers = new Headers(init.headers)
    if (auth && !headers.get("authorization")) headers.set("authorization", auth)
    const target = new URL(pathValue, config.baseUrl)
    return fetch(target, { ...init, headers })
  }

  const eventStream = () => {
    const headers = new Headers({ Accept: "text/event-stream" })
    if (auth) headers.set("Authorization", auth)
    return fetch(`${config.baseUrl}/event`, { headers })
  }

  return { start, fetch: fetchHandler, request, eventStream, config }
}

if (import.meta.main) {
  createProxy().start()
}
