import { type Component, createSignal } from "solid-js"

interface StyleComparisonProps {
  onBack?: () => void
}

const StyleComparison: Component<StyleComparisonProps> = (props) => {
  const [selectedStyle, setSelectedStyle] = createSignal("tech") // Default to 'tech' (Tech Future)

  return (
    <div
      id="16:30"
      style="background: linear-gradient(135deg, rgba(10, 14, 26, 1) 0%, rgba(26, 19, 50, 1) 50%, rgba(13, 27, 42, 1) 100%), radial-gradient(circle at 20% 30%, rgba(0, 240, 255, 0.101961) 0%, transparent 50%), radial-gradient(circle at 80% 70%, rgba(255, 0, 110, 0.101961) 0%, transparent 50%);"
      class="flex grow min-w-0 h-full overflow-y-auto"
    >
      <main
        id="16:31"
        style="flex-basis: 0%; padding: 1.5rem 2rem;"
        class="overflow-x-hidden flex flex-col grow shrink"
      >
        <div id="16:32" class="flex justify-between items-center mb-8">
          <div id="16:33" class="flex items-center gap-y-4 gap-x-4">
            <button 
              onClick={props.onBack}
              class="hover:bg-[#00F0FF]/10 p-2 rounded-lg transition-colors text-[#00F0FF]"
            >
              <iconify-icon icon="lucide:arrow-left" class="text-xl"></iconify-icon>
            </button>
            <h1 id="16:34" style="color: rgba(232, 240, 255, 1);" class="text-2xl font-bold">
              风格选择对比
            </h1>
          </div>
        </div>

        <div id="16:45" style="flex-basis: 0%;" class="grid grow shrink grid-cols-3 gap-y-6 gap-x-6">
          {/* Card 1: 现代极简 */}
          <div
            id="16:46"
            class="hover:shadow-[0_8px_32px_rgba(0,240,255,0.12)] flex flex-col border-[1px] border-solid rounded-2xl transition-all"
            style={{
              "background-color":
                selectedStyle() === "modern"
                  ? "color-mix( in oklab , #00F0FF 10% , transparent )"
                  : "color-mix( in oklab , #1A1F3A 90% , transparent )",
              "box-shadow":
                selectedStyle() === "modern"
                  ? "0 0 24px rgba(0, 240, 255, 0.15), 0 4px 20px rgba(0, 240, 255, 0.1)"
                  : "0 4px 20px rgba(0, 240, 255, 0.08)",
              "border-color":
                selectedStyle() === "modern"
                  ? "color-mix( in oklab , #00F0FF 40% , transparent )"
                  : "color-mix( in oklab , #00F0FF 10% , transparent )",
            }}
          >
            <div id="16:47" class="p-5">
              <div id="16:48" class="flex justify-between items-center mb-4">
                <h3 id="16:49" style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">
                  现代极简
                </h3>
                <div
                  id="16:50"
                  style="background-color: color-mix( in oklab , #00FF9F 20% , transparent ); color: rgba(0, 255, 159, 1); padding: 0.25rem 0.5rem;"
                  class="text-xs rounded-lg"
                >
                  推荐
                </div>
              </div>

              <div id="16:51" class="mb-4">
                <p id="16:52" style="color: rgba(138, 151, 170, 1);" class="text-sm mb-2">
                  主色调
                </p>
                <div id="16:53" class="flex gap-y-2 gap-x-2">
                  <div
                    id="16:54"
                    style="background-color: rgba(37, 99, 235, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                  <div
                    id="16:55"
                    style="background-color: rgba(248, 250, 252, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                  <div
                    id="16:56"
                    style="background-color: rgba(100, 116, 139, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                  <div
                    id="16:57"
                    style="background-color: rgba(15, 23, 42, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                </div>
              </div>

              <div id="16:58" class="mb-4">
                <p id="16:59" style="color: rgba(138, 151, 170, 1);" class="text-sm mb-2">
                  字体风格
                </p>
                <div
                  id="16:60"
                  style="background-color: color-mix( in oklab , #141829 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
                  class="p-3 border-[1px] border-solid rounded-lg"
                >
                  <p
                    id="16:61"
                    style="color: rgba(232, 240, 255, 1);"
                    class="text-base mb-1 font-semibold"
                  >
                    Inter Sans Serif
                  </p>
                  <p id="16:62" style="color: rgba(184, 197, 217, 1);" class="text-sm">
                    简洁现代的无衬线字体
                  </p>
                </div>
              </div>
            </div>

            <div id="16:63" style="flex-basis: 0%;" class="grow shrink pr-5 pb-5 pl-5">
              <p id="16:64" style="color: rgba(138, 151, 170, 1);" class="text-sm mb-3">
                界面预览
              </p>
              <div
                id="16:65"
                style="background-color: color-mix( in oklab , #141829 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
                class="h-48 p-4 border-[1px] border-solid rounded-lg"
              >
                <img
                  id="16:66"
                  style="filter: brightness(90%) contrast(90%);"
                  alt="Clean modern interface with blue primary color, white background, minimal design elements"
                  src="https://static.paraflowcontent.com/public/resource/image/aa23b1b8-591b-488e-a0e1-af8f6631a4c6.jpeg"
                  class="w-full h-full object-cover rounded-sm"
                />
              </div>
            </div>

            <div id="16:67" class="pr-5 pb-5 pl-5">
              <button
                id="16:68"
                onClick={() => setSelectedStyle("modern")}
                class="hover:shadow-[0_0_20px_rgba(0,240,255,0.4)] flex justify-center items-center gap-y-2 gap-x-2 w-full rounded-lg cursor-pointer transition-all"
                style={{
                  "background-color":
                    selectedStyle() === "modern"
                      ? "rgba(0, 240, 255, 1)"
                      : "transparent",
                  color:
                    selectedStyle() === "modern"
                      ? "rgba(10, 14, 26, 1)"
                      : "rgba(0, 240, 255, 1)",
                  padding: "0.75rem 1rem",
                  "border-color": "color-mix( in oklab , #00F0FF 60% , transparent )",
                  "border-width": selectedStyle() === "modern" ? "0" : "1px",
                  "border-style": "solid",
                }}
              >
                {selectedStyle() === "modern" && (
                  <div class="bg-transparent flex justify-center items-center w-4 h-4">
                    <iconify-icon icon="lucide:check" class="text-sm"></iconify-icon>
                  </div>
                )}
                <span id="16:69" class="whitespace-nowrap font-semibold">
                  {selectedStyle() === "modern" ? "已选择" : "选择此风格"}
                </span>
              </button>
            </div>
          </div>

          {/* Card 2: 温暖人文 */}
          <div
            id="16:70"
            class="hover:shadow-[0_8px_32px_rgba(0,240,255,0.12)] flex flex-col border-[1px] border-solid rounded-2xl transition-all"
            style={{
              "background-color":
                selectedStyle() === "warm"
                  ? "color-mix( in oklab , #00F0FF 10% , transparent )"
                  : "color-mix( in oklab , #1A1F3A 90% , transparent )",
              "box-shadow":
                selectedStyle() === "warm"
                  ? "0 0 24px rgba(0, 240, 255, 0.15), 0 4px 20px rgba(0, 240, 255, 0.1)"
                  : "0 4px 20px rgba(0, 240, 255, 0.08)",
              "border-color":
                selectedStyle() === "warm"
                  ? "color-mix( in oklab , #00F0FF 40% , transparent )"
                  : "color-mix( in oklab , #00F0FF 10% , transparent )",
            }}
          >
            <div id="16:71" class="p-5">
              <div id="16:72" class="flex justify-between items-center mb-4">
                <h3 id="16:73" style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">
                  温暖人文
                </h3>
                <div
                  id="16:74"
                  style="background-color: color-mix( in oklab , #B026FF 20% , transparent ); color: rgba(176, 38, 255, 1); padding: 0.25rem 0.5rem;"
                  class="text-xs rounded-lg"
                >
                  创意
                </div>
              </div>

              <div id="16:75" class="mb-4">
                <p id="16:76" style="color: rgba(138, 151, 170, 1);" class="text-sm mb-2">
                  主色调
                </p>
                <div id="16:77" class="flex gap-y-2 gap-x-2">
                  <div
                    id="16:78"
                    style="background-color: rgba(234, 88, 12, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                  <div
                    id="16:79"
                    style="background-color: rgba(254, 243, 199, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                  <div
                    id="16:80"
                    style="background-color: rgba(146, 64, 14, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                  <div
                    id="16:81"
                    style="background-color: rgba(69, 26, 3, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                </div>
              </div>

              <div id="16:82" class="mb-4">
                <p id="16:83" style="color: rgba(138, 151, 170, 1);" class="text-sm mb-2">
                  字体风格
                </p>
                <div
                  id="16:84"
                  style="background-color: color-mix( in oklab , #141829 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
                  class="p-3 border-[1px] border-solid rounded-lg"
                >
                  <p
                    id="16:85"
                    style="color: rgba(232, 240, 255, 1);"
                    class="text-base mb-1 font-semibold"
                  >
                    Playfair Display
                  </p>
                  <p id="16:86" style="color: rgba(184, 197, 217, 1);" class="text-sm">
                    优雅的衬线字体，富有人文气息
                  </p>
                </div>
              </div>
            </div>

            <div id="16:87" style="flex-basis: 0%;" class="grow shrink pr-5 pb-5 pl-5">
              <p id="16:88" style="color: rgba(138, 151, 170, 1);" class="text-sm mb-3">
                界面预览
              </p>
              <div
                id="16:89"
                style="background-color: color-mix( in oklab , #141829 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
                class="h-48 p-4 border-[1px] border-solid rounded-lg"
              >
                <img
                  id="16:90"
                  style="filter: brightness(90%) contrast(90%);"
                  alt="Warm humanistic interface with orange and cream colors, serif typography, organic shapes and friendly design"
                  src="https://static.paraflowcontent.com/public/resource/image/8f64b451-6891-4d35-bce3-60d02eccab0e.jpeg"
                  class="w-full h-full object-cover rounded-sm"
                />
              </div>
            </div>

            <div id="16:91" class="pr-5 pb-5 pl-5">
              <button
                id="16:92"
                onClick={() => setSelectedStyle("warm")}
                class="hover:bg-[#00F0FF]/10 hover:shadow-[0_0_15px_rgba(0,240,255,0.3)] bg-transparent flex justify-center items-center gap-y-2 gap-x-2 w-full border-[1px] border-solid rounded-lg cursor-pointer transition-all"
                style={{
                  "background-color":
                    selectedStyle() === "warm"
                      ? "rgba(0, 240, 255, 1)"
                      : "transparent",
                  color:
                    selectedStyle() === "warm"
                      ? "rgba(10, 14, 26, 1)"
                      : "rgba(0, 240, 255, 1)",
                  padding: "0.75rem 1rem",
                  "border-color": "color-mix( in oklab , #00F0FF 60% , transparent )",
                  "border-width": selectedStyle() === "warm" ? "0" : "1px",
                  "border-style": "solid",
                }}
              >
                {selectedStyle() === "warm" && (
                  <div class="bg-transparent flex justify-center items-center w-4 h-4">
                    <iconify-icon icon="lucide:check" class="text-sm"></iconify-icon>
                  </div>
                )}
                <span id="16:93" class="whitespace-nowrap font-semibold">
                  {selectedStyle() === "warm" ? "已选择" : "选择此风格"}
                </span>
              </button>
            </div>
          </div>

          {/* Card 3: 科技未来 */}
          <div
            id="16:94"
            class="hover:shadow-[0_8px_32px_rgba(0,240,255,0.12)] flex flex-col border-[1px] border-solid rounded-2xl transition-all"
            style={{
              "background-color":
                selectedStyle() === "tech"
                  ? "color-mix( in oklab , #00F0FF 10% , transparent )"
                  : "color-mix( in oklab , #1A1F3A 90% , transparent )",
              "box-shadow":
                selectedStyle() === "tech"
                  ? "0 0 24px rgba(0, 240, 255, 0.15), 0 4px 20px rgba(0, 240, 255, 0.1)"
                  : "0 4px 20px rgba(0, 240, 255, 0.08)",
              "border-color":
                selectedStyle() === "tech"
                  ? "color-mix( in oklab , #00F0FF 40% , transparent )"
                  : "color-mix( in oklab , #00F0FF 10% , transparent )",
            }}
          >
            <div id="16:95" class="p-5">
              <div id="16:96" class="flex justify-between items-center mb-4">
                <h3 id="16:97" style="color: rgba(232, 240, 255, 1);" class="text-lg font-semibold">
                  科技未来
                </h3>
                <div
                  id="16:98"
                  style="background-color: color-mix( in oklab , #00F0FF 30% , transparent ); color: rgba(0, 240, 255, 1); padding: 0.25rem 0.5rem;"
                  class="text-xs rounded-lg"
                >
                  当前
                </div>
              </div>

              <div id="16:99" class="mb-4">
                <p id="16:100" style="color: rgba(138, 151, 170, 1);" class="text-sm mb-2">
                  主色调
                </p>
                <div id="16:101" class="flex gap-y-2 gap-x-2">
                  <div
                    id="16:102"
                    style="background-color: rgba(0, 240, 255, 1); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 8px rgba(0, 240, 255, 0.4); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                  <div
                    id="16:103"
                    style="background-color: rgba(10, 14, 26, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                  <div
                    id="16:104"
                    style="background-color: rgba(255, 0, 110, 1); box-shadow: 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 rgba(0, 0, 0, 0), 0 0 8px rgba(255, 0, 110, 0.4); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                  <div
                    id="16:105"
                    style="background-color: rgba(26, 31, 58, 1); border-color: color-mix( in oklab , #00F0FF 20% , transparent );"
                    class="w-8 h-8 border-[1px] border-solid rounded-lg"
                  ></div>
                </div>
              </div>

              <div id="16:106" class="mb-4">
                <p id="16:107" style="color: rgba(138, 151, 170, 1);" class="text-sm mb-2">
                  字体风格
                </p>
                <div
                  id="16:108"
                  style="background-color: color-mix( in oklab , #141829 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
                  class="p-3 border-[1px] border-solid rounded-lg"
                >
                  <p
                    id="16:109"
                    style="color: rgba(232, 240, 255, 1);"
                    class="text-base mb-1 font-semibold"
                  >
                    Roboto Mono
                  </p>
                  <p id="16:110" style="color: rgba(184, 197, 217, 1);" class="text-sm">
                    未来科技感的等宽字体
                  </p>
                </div>
              </div>
            </div>

            <div id="16:111" style="flex-basis: 0%;" class="grow shrink pr-5 pb-5 pl-5">
              <p id="16:112" style="color: rgba(138, 151, 170, 1);" class="text-sm mb-3">
                界面预览
              </p>
              <div
                id="16:113"
                style="background-color: color-mix( in oklab , #141829 80% , transparent ); border-color: color-mix( in oklab , #00F0FF 10% , transparent );"
                class="h-48 p-4 border-[1px] border-solid rounded-lg"
              >
                <img
                  id="16:114"
                  style="filter: brightness(90%) contrast(90%);"
                  alt="Futuristic cyberpunk interface with neon cyan and magenta accents, dark gradient background, glowing elements"
                  src="https://static.paraflowcontent.com/public/resource/image/1f18cd3f-c98a-4303-b6db-786a4e87ed15.jpeg"
                  class="w-full h-full object-cover rounded-sm"
                />
              </div>
            </div>

            <div id="16:115" class="pr-5 pb-5 pl-5">
              <button
                id="16:116"
                onClick={() => setSelectedStyle("tech")}
                class="hover:shadow-[0_0_20px_rgba(0,240,255,0.4)] flex justify-center items-center gap-y-2 gap-x-2 w-full rounded-lg cursor-pointer transition-all"
                style={{
                  "background-color":
                    selectedStyle() === "tech"
                      ? "rgba(0, 240, 255, 1)"
                      : "transparent",
                  color:
                    selectedStyle() === "tech"
                      ? "rgba(10, 14, 26, 1)"
                      : "rgba(0, 240, 255, 1)",
                  padding: "0.75rem 1rem",
                  "border-color": "color-mix( in oklab , #00F0FF 60% , transparent )",
                  "border-width": selectedStyle() === "tech" ? "0" : "1px",
                  "border-style": "solid",
                }}
              >
                {selectedStyle() === "tech" && (
                  <div id="16:117" class="bg-transparent flex justify-center items-center w-4 h-4">
                    <iconify-icon id="16:118" icon="lucide:check" class="text-sm"></iconify-icon>
                  </div>
                )}
                <span id="16:119" class="whitespace-nowrap font-semibold">
                  {selectedStyle() === "tech" ? "已选择" : "选择此风格"}
                </span>
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

export default StyleComparison
