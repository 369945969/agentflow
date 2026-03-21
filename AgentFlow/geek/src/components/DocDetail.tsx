import { type Component, createSignal, Show } from "solid-js"

interface DocDetailProps {
  onBack: () => void
  doc: any
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

const DocDetail: Component<DocDetailProps> = (props) => {
  const [isEditing, setIsEditing] = createSignal(false)
  const [editContent, setEditContent] = createSignal("")

  const startEdit = () => {
    setEditContent(props.doc?.content || "")
    setIsEditing(true)
  }

  const saveEdit = () => {
    if (props.doc) {
      props.doc.content = editContent()
    }
    setIsEditing(false)
  }

  const cancelEdit = () => {
    setIsEditing(false)
  }

  return (
    <div class="flex flex-col grow min-w-0 h-full overflow-hidden">
      {/* Top Header */}
      <div
        class="shrink-0 flex justify-between items-center px-6 py-4 border-b border-[#00F0FF]/15"
        style="background: linear-gradient(90deg, rgba(22, 33, 62, 0.95) 0%, rgba(15, 22, 36, 0.95) 100%); backdrop-filter: blur(10px);"
      >
        <div class="flex items-center gap-x-4">
          <button
            onClick={props.onBack}
            class="hover:bg-[#00F0FF]/15 hover:shadow-[0_0_12px_rgba(0,240,255,0.25)] flex items-center justify-center rounded-lg cursor-pointer transition-all duration-300"
            style="background-color: color-mix( in oklab , #1A1F3A 90% , transparent ); padding: 0.5rem;"
          >
            <iconify-icon
              style="color: rgba(0, 240, 255, 1);"
              icon="lucide:arrow-left"
              class="text-xl"
            ></iconify-icon>
          </button>
          <h1 style="color: rgba(232, 240, 255, 1);" class="text-xl font-bold">
            {props.doc?.title || "产品需求文档"}
          </h1>
        </div>
        <div class="flex items-center gap-x-2">
          <Show
            when={!isEditing()}
            fallback={
              <button
                onClick={saveEdit}
                class="hover:shadow-[0_0_20px_rgba(0,240,255,0.4)] flex items-center gap-x-2 rounded-lg cursor-pointer transition-all duration-300"
                style="background-color: rgba(0, 240, 255, 1); color: rgba(10, 14, 26, 1); padding: 0.5rem 1rem;"
              >
                <iconify-icon
                  style="color: rgba(10, 14, 26, 1);"
                  icon="lucide:save"
                  class="text-sm"
                ></iconify-icon>
                <span style="color: rgba(10, 14, 26, 1);" class="text-sm font-semibold">
                  保存
                </span>
              </button>
            }
          >
            <button
              onClick={startEdit}
              class="hover:bg-[#00F0FF]/15 hover:shadow-[0_0_12px_rgba(0,240,255,0.25)] flex items-center gap-x-2 rounded-lg cursor-pointer transition-all duration-300"
              style="background-color: color-mix( in oklab , #1A1F3A 90% , transparent ); padding: 0.5rem 0.75rem;"
            >
              <iconify-icon
                style="color: rgba(0, 240, 255, 1);"
                icon="lucide:edit"
                class="text-sm"
              ></iconify-icon>
              <span style="color: rgba(232, 240, 255, 1);" class="text-sm">
                编辑
              </span>
            </button>
          </Show>
        </div>
      </div>

      <div class="flex grow overflow-hidden">
        {/* Document List - Moved inside */}
        <aside class="shrink-0 min-w-fit h-full">
          <div
            style="background: linear-gradient(135deg, rgba(22, 33, 62, 0.5) 0%, rgba(15, 22, 36, 0.5) 100%);"
            class="flex flex-col w-64 h-full border-r border-[#00F0FF]/15"
          >
            <div
              style="border-bottom-style: solid; padding: 1rem 1.5rem; border-color: color-mix( in oklab , #00F0FF 15% , transparent );"
              class="flex items-center border-b-[1px]"
            >
              <div class="flex items-center gap-x-3">
                <div class="bg-transparent flex justify-center items-center w-6 h-6">
                  <iconify-icon
                    style="color: rgba(0, 240, 255, 1);"
                    icon="lucide:folder"
                    class="text-xl"
                  ></iconify-icon>
                </div>
                <h2 style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">
                  文件列表
                </h2>
              </div>
            </div>

            <div style="flex-basis: 0%; padding: 1rem 1.5rem;" class="overflow-y-auto grow shrink">
              <div class="flex flex-col gap-y-2">
                <div
                  style="background-color: color-mix( in oklab , #00F0FF 15% , transparent ); border-color: color-mix( in oklab , #00F0FF 30% , transparent );"
                  class="flex items-center gap-x-3 p-3 rounded-lg cursor-pointer border-[1px] border-solid"
                >
                  <div
                    style="background-color: color-mix( in oklab , #00F0FF 20% , transparent ); border-color: color-mix( in oklab , #00F0FF 40% , transparent );"
                    class="flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg shrink-0"
                  >
                    <iconify-icon
                      style="color: rgba(0, 240, 255, 1);"
                      icon="lucide:file-text"
                      class="text-base"
                    ></iconify-icon>
                  </div>
                  <div class="flex flex-col overflow-hidden">
                    <span
                      style="color: rgba(232, 240, 255, 1);"
                      class="text-sm font-semibold truncate"
                    >
                      {props.doc?.title || "产品需求文档"}
                    </span>
                    <span style="color: rgba(138, 151, 170, 1);" class="text-xs">
                      文档
                    </span>
                  </div>
                </div>

                <button
                  class="hover:bg-[#00F0FF]/10 flex items-center gap-x-3 p-3 rounded-lg cursor-pointer transition-colors w-full"
                >
                  <div
                    style="background-color: color-mix( in oklab , #B026FF 10% , transparent ); border-color: color-mix( in oklab , #B026FF 20% , transparent );"
                    class="flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg shrink-0"
                  >
                    <iconify-icon
                      style="color: rgba(176, 38, 255, 1);"
                      icon="lucide:image"
                      class="text-base"
                    ></iconify-icon>
                  </div>
                  <div class="text-left flex flex-col overflow-hidden">
                    <span style="color: rgba(232, 240, 255, 1);" class="text-sm font-semibold truncate">
                      登录界面设计
                    </span>
                    <span style="color: rgba(138, 151, 170, 1);" class="text-xs">
                      设计稿
                    </span>
                  </div>
                </button>
                <button
                  class="hover:bg-[#00F0FF]/10 flex items-center gap-x-3 p-3 rounded-lg cursor-pointer transition-colors w-full"
                >
                  <div
                    style="background-color: color-mix( in oklab , #FF006E 10% , transparent ); border-color: color-mix( in oklab , #FF006E 20% , transparent );"
                    class="flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-lg shrink-0"
                  >
                    <iconify-icon
                      style="color: rgba(255, 0, 110, 1);"
                      icon="lucide:layout"
                      class="text-base"
                    ></iconify-icon>
                  </div>
                  <div class="text-left flex flex-col overflow-hidden">
                    <span style="color: rgba(232, 240, 255, 1);" class="text-sm font-semibold truncate">
                      主界面线框图
                    </span>
                    <span style="color: rgba(138, 151, 170, 1);" class="text-xs">
                      线框图
                    </span>
                  </div>
                </button>
              </div>
            </div>
          </div>
        </aside>

        <main
          style="flex-basis: 0%; padding: 0;"
          class="overflow-x-hidden overflow-y-auto flex flex-col grow shrink bg-[#0A0E1A]/50"
        >
          <div
            style="background-color: color-mix( in oklab , #141829 80% , transparent ); flex-basis: 0%;"
            class="overflow-y-auto grow shrink"
          >
            <div class="p-8">
            <Show
              when={isEditing()}
              fallback={
                <div id="15:91" class="prose prose-invert" style="max-width: none;">
                  <div id="15:92" class="mb-8">
                    <div
                      id="15:93"
                      style="background-color: color-mix( in oklab , #00F0FF 15% , transparent ); color: rgba(0, 240, 255, 1); padding: 0.25rem 0.5rem;"
                      class="text-xs inline-block mb-4 rounded-lg"
                    >
                      文档
                    </div>
                    <h1 id="15:94" style="color: rgba(232, 240, 255, 1);" class="text-3xl mb-2 font-bold">
                      {props.doc?.title || "AI 协作设计平台产品需求文档"}
                    </h1>
                    <div id="15:95" class="flex items-center gap-y-4 gap-x-4 mb-6">
                      <span id="15:96" style="color: rgba(138, 151, 170, 1);" class="text-sm">
                        最后更新: {props.doc?.lastModified || "2024年12月21日"}
                      </span>
                      <span id="15:97" style="color: rgba(138, 151, 170, 1);" class="text-sm">
                        版本: v1.0
                      </span>
                    </div>
                  </div>

                  <div id="15:98">
                    <Show
                      when={props.doc?.content}
                      fallback={
                        <div class="mb-8">
                          <h2
                            class="text-xl font-bold mb-4 mt-8 pb-2 border-b"
                            style="color: rgba(232, 240, 255, 1); border-color: color-mix(in oklab, #00F0FF 30%, transparent);"
                          >
                            1. 产品概述
                          </h2>
                          <p class="mb-4 leading-relaxed" style="color: rgba(184, 197, 217, 1);">
                            AI
                            协作设计平台是一个创新的设计工具，结合人工智能技术与协作设计理念，为设计师和团队提供智能化的设计创作环境。平台采用无限画布的交互模式，让用户能够自由地在画布上创建、编辑和组织各类设计内容。
                          </p>
                          <p class="mb-4 leading-relaxed" style="color: rgba(184, 197, 217, 1);">
                            核心特色在于集成的 AI
                            助手功能，能够实时理解用户需求，自动生成设计稿、文档和风格指南，大幅提升设计效率和创作质量。
                          </p>

                          <h2
                            class="text-xl font-bold mb-4 mt-8 pb-2 border-b"
                            style="color: rgba(232, 240, 255, 1); border-color: color-mix(in oklab, #00F0FF 30%, transparent);"
                          >
                            2. 核心功能
                          </h2>

                          <h3 class="text-lg font-semibold mb-3 mt-6" style="color: rgba(232, 240, 255, 1);">
                            2.1 无限画布
                          </h3>
                          <ul class="ml-4 list-none">
                            <li class="text-base flex items-start mb-2" style="color: rgba(184, 197, 217, 1);">
                              <span
                                class="w-2 h-2 mt-2 mr-3 rounded-full shrink-0"
                                style="background-color: rgba(0, 240, 255, 1);"
                              ></span>
                              支持无限缩放和平移的画布空间，让设计思路不受边界限制
                            </li>
                            <li class="text-base flex items-start mb-2" style="color: rgba(184, 197, 217, 1);">
                              <span
                                class="w-2 h-2 mt-2 mr-3 rounded-full shrink-0"
                                style="background-color: rgba(0, 240, 255, 1);"
                              ></span>
                              支持多种内容类型：文档、图片、线框图、设计稿等
                            </li>
                            <li class="text-base flex items-start" style="color: rgba(184, 197, 217, 1);">
                              <span
                                class="w-2 h-2 mt-2 mr-3 rounded-full shrink-0"
                                style="background-color: rgba(0, 240, 255, 1);"
                              ></span>
                              灵活的内容卡片系统，支持拖拽、缩放和连接操作
                            </li>
                          </ul>
                        </div>
                      }
                    >
                      <div innerHTML={parseMarkdown(props.doc.content)} />
                    </Show>
                  </div>
                </div>
              }
            >
              <div id="15:91-edit" class="h-full flex flex-col">
                <div
                  id="15:93"
                  style="background-color: color-mix( in oklab , #00F0FF 15% , transparent ); color: rgba(0, 240, 255, 1); padding: 0.25rem 0.5rem;"
                  class="text-xs inline-block mb-4 rounded-lg"
                >
                  编辑 Markdown 源码
                </div>
                <textarea
                  id="15:editor"
                  value={editContent()}
                  onInput={(e) => setEditContent(e.currentTarget.value)}
                  class="flex-grow w-full p-4 rounded-lg outline-none resize-none"
                  style="background-color: color-mix(in oklab, #1A1F3A 90%, transparent); border: 1px solid color-mix(in oklab, #00F0FF 30%, transparent); color: rgba(184, 197, 217, 1); font-family: monospace; min-height: 400px;"
                  placeholder="在此输入 Markdown 源码..."
                ></textarea>
                <div class="flex justify-end gap-3 mt-4">
                  <button
                    id="15:cancel"
                    onClick={cancelEdit}
                    class="hover:bg-[#00F0FF]/10 rounded-lg px-4 py-2"
                    style="color: rgba(138, 151, 170, 1); border: 1px solid color-mix(in oklab, #00F0FF 30%, transparent);"
                  >
                    取消
                  </button>
                </div>
              </div>
            </Show>
          </div>
        </div>
      </main>
    </div>
  </div>
  )
}

export default DocDetail
