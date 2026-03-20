#!/bin/bash

export DEEPSEEK_API_KEY=sk-6696ebd5e41f4fd9a1a218b57a85ad6b
# OpenCode 多模型路由动态验证脚本 (从 API 动态查询 ID)
# 运行前请确保：1. Go 后端已启动(3000) 2. opencode server 已启动(3001)

set -euo pipefail

BASE_URL="http://localhost:3001"
BACKEND_API="http://localhost:3000/api/models/"

log_step() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🧪 $1"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

log_kv() {
    echo "   - $1: $2"
}

preview_json() {
    local label=$1
    local json=$2
    echo "   - $label(前 30 行):"
    echo "$json" | jq '.' 2>/dev/null | sed -n '1,30p' | sed 's/^/     /'
}

DEEPSEEK_API_KEY="${DEEPSEEK_API_KEY:-}"

mem0_enabled() {
    [ -n "${DEEPSEEK_API_KEY}" ]
}

echo "🔍 开始功能动态验证..."

# 1. 检查 jq 是否安装 (用于解析 JSON)
if ! command -v jq &> /dev/null; then
    echo "❌ 错误: 未安装 jq，无法解析 API 返回的 JSON。"
    echo "请运行 'brew install jq' (macOS) 或 'sudo apt install jq' (Linux)"
    exit 1
fi

# 2. 动态获取模型 ID (通过 Go 后端 API)
log_step "读取后端模型列表"
log_kv "GET" "$BACKEND_API"
MODELS_JSON=$(curl -s --compressed "$BACKEND_API")

if [ -z "$MODELS_JSON" ] || [ "$MODELS_JSON" == "[]" ]; then
    echo "❌ 错误: 后端返回的模型列表为空。请确保 Go 后端正在运行且已配置模型。"
    exit 1
fi
preview_json "后端 models 列表" "$MODELS_JSON"

DEFAULT_ROW=$(echo "$MODELS_JSON" | jq -c 'map(select(.is_default == true))[0] // .[0]')
DEFAULT_MODEL_ID=$(echo "$DEFAULT_ROW" | jq -r '.id')
DEFAULT_MODEL_NAME=$(echo "$DEFAULT_ROW" | jq -r '.name')
DEFAULT_MODEL_BASE_URL=$(echo "$DEFAULT_ROW" | jq -r '.base_url // empty')
DEFAULT_MODEL_MODEL_NAME=$(echo "$DEFAULT_ROW" | jq -r '.model_name // empty')
DEFAULT_MODEL_PROVIDER_LABEL=$(echo "$DEFAULT_ROW" | jq -r '.provider // empty')

OTHER_ROW=$(echo "$MODELS_JSON" | jq -c 'map(select(.is_default == false))[0] // empty')
OTHER_MODEL_ID=$(echo "$OTHER_ROW" | jq -r '.id // empty')
OTHER_MODEL_NAME=$(echo "$OTHER_ROW" | jq -r '.name // empty')
OTHER_MODEL_BASE_URL=$(echo "$OTHER_ROW" | jq -r '.base_url // empty')
OTHER_MODEL_MODEL_NAME=$(echo "$OTHER_ROW" | jq -r '.model_name // empty')
OTHER_MODEL_PROVIDER_LABEL=$(echo "$OTHER_ROW" | jq -r '.provider // empty')

log_step "解析后端默认模型 (DuckDB 记录)"
log_kv "duckdb_model_id" "$DEFAULT_MODEL_ID"
log_kv "name" "$DEFAULT_MODEL_NAME"
log_kv "provider" "$DEFAULT_MODEL_PROVIDER_LABEL"
log_kv "base_url" "$DEFAULT_MODEL_BASE_URL"
log_kv "model_name" "$DEFAULT_MODEL_MODEL_NAME"

echo "   ✅ 默认模型: $DEFAULT_MODEL_NAME ($DEFAULT_MODEL_ID)"
[ -n "$OTHER_MODEL_ID" ] && [ "$OTHER_MODEL_ID" != "null" ] && echo "   ✅ 备选模型: $OTHER_MODEL_NAME ($OTHER_MODEL_ID)"

# 3. 检查 OpenCode 服务状态
log_step "检查 OpenCode 服务健康"
log_kv "GET" "$BASE_URL/global/health"
echo -n "🚀 检查 OpenCode 服务 (Port 3001): "
HEALTH=$(curl -s --compressed -o /dev/null -w "%{http_code}" "$BASE_URL/global/health" || echo "000")
if [ "$HEALTH" == "200" ]; then
    echo "✅ 正常"
else
    echo "❌ 失败 (状态码: $HEALTH)"
    echo "请运行 'bash start_opencode.sh' 启动服务。"
    exit 1
fi

log_step "读取 OpenCode provider 列表"
log_kv "GET" "$BASE_URL/provider"
OPENCODE_PROVIDER_JSON=$(curl -s --compressed "$BASE_URL/provider")
preview_json "OpenCode /provider" "$OPENCODE_PROVIDER_JSON"

provider_for_baseurl() {
    local base_url=$1
    if [ -z "$base_url" ]; then
        echo ""
        return 0
    fi

    echo "$OPENCODE_PROVIDER_JSON" | jq -r --arg url "$base_url" '
        .all[]
        | select(
            (.models | to_entries | any(.value.api.url == $url))
        )
        | .id
    ' | head -n 1
}

provider_has_model() {
    local provider_id=$1
    local model_id=$2
    if [ -z "$provider_id" ] || [ -z "$model_id" ]; then
        return 1
    fi

    echo "$OPENCODE_PROVIDER_JSON" | jq -e --arg pid "$provider_id" --arg mid "$model_id" '
        any(.all[]; .id==$pid and (.models[$mid] != null))
    ' >/dev/null 2>&1
}

pick_model() {
    local provider_id=$1
    local model_id
    model_id=$(echo "$OPENCODE_PROVIDER_JSON" | jq -r --arg pid "$provider_id" '.default[$pid] // empty')
    if [ -n "$model_id" ]; then
        echo "$model_id"
        return 0
    fi

    model_id=$(echo "$OPENCODE_PROVIDER_JSON" | jq -r --arg pid "$provider_id" '.all[] | select(.id==$pid) | (.models | keys | .[0])' | head -n 1)
    echo "$model_id"
}

resolve_provider() {
    local provider_id=$1
    local resolved
    resolved=$(echo "$OPENCODE_PROVIDER_JSON" | jq -r --arg pid "$provider_id" '.connected[]? | select(.==$pid)' | head -n 1)
    if [ -n "$resolved" ]; then
        echo "$resolved"
        return 0
    fi

    resolved=$(echo "$OPENCODE_PROVIDER_JSON" | jq -r '.connected[0] // empty')
    echo "$resolved"
}

capture_event_for_session() {
    local session_id=$1
    local outfile=$2

    curl -N -s --compressed "$BASE_URL/event" | while IFS= read -r line; do
        case "$line" in
            data:*)
                json="${line#data: }"
                type=$(echo "$json" | jq -r '.type // empty' 2>/dev/null)

                if [ "$type" = "message.updated" ]; then
                    sid=$(echo "$json" | jq -r '.properties.info.sessionID // empty' 2>/dev/null)
                    role=$(echo "$json" | jq -r '.properties.info.role // empty' 2>/dev/null)
                    if [ "$sid" = "$session_id" ] && [ "$role" = "assistant" ]; then
                        echo "$json" > "$outfile"
                        break
                    fi
                fi

                if [ "$type" = "session.error" ]; then
                    sid=$(echo "$json" | jq -r '.properties.sessionID // empty' 2>/dev/null)
                    if [ "$sid" = "$session_id" ]; then
                        echo "$json" > "$outfile"
                        break
                    fi
                fi
                ;;
        esac
    done
}

capture_assistant_text_for_session() {
    local session_id=$1
    local outfile=$2
    local errfile=$3

    local assistant_msg_id=""

    curl -N -s --compressed "$BASE_URL/event" | while IFS= read -r line; do
        case "$line" in
            data:*)
                json="${line#data: }"
                type=$(echo "$json" | jq -r '.type // empty' 2>/dev/null)

                if [ "$type" = "message.updated" ]; then
                    sid=$(echo "$json" | jq -r '.properties.info.sessionID // empty' 2>/dev/null)
                    role=$(echo "$json" | jq -r '.properties.info.role // empty' 2>/dev/null)
                    if [ "$sid" = "$session_id" ] && [ "$role" = "assistant" ]; then
                        assistant_msg_id=$(echo "$json" | jq -r '.properties.info.id // empty' 2>/dev/null)
                    fi
                fi

                if [ "$type" = "message.part.updated" ]; then
                    sid=$(echo "$json" | jq -r '.properties.part.sessionID // empty' 2>/dev/null)
                    mid=$(echo "$json" | jq -r '.properties.part.messageID // empty' 2>/dev/null)
                    ptype=$(echo "$json" | jq -r '.properties.part.type // empty' 2>/dev/null)
                    if [ "$sid" = "$session_id" ] && [ -n "$assistant_msg_id" ] && [ "$mid" = "$assistant_msg_id" ] && [ "$ptype" = "text" ]; then
                        echo "$json" | jq -r '.properties.part.text // empty' > "$outfile" 2>/dev/null || true
                    fi
                fi

                if [ "$type" = "session.error" ]; then
                    sid=$(echo "$json" | jq -r '.properties.sessionID // empty' 2>/dev/null)
                    if [ "$sid" = "$session_id" ]; then
                        echo "$json" > "$errfile"
                        break
                    fi
                fi

                if [ "$type" = "session.idle" ]; then
                    sid=$(echo "$json" | jq -r '.properties.sessionID // empty' 2>/dev/null)
                    if [ "$sid" = "$session_id" ] && [ -n "$assistant_msg_id" ] && [ -s "$outfile" ]; then
                        break
                    fi
                fi
                ;;
        esac
    done
}

send_and_assert_model() {
    local label=$1
    local expected_provider=$2
    local expected_model=$3
    local user_id=$4
    local duckdb_model_id=$5
    local message_text=$6

    log_step "创建 Session"
    log_kv "POST" "$BASE_URL/session"
    SESSION_ID=$(curl -s --compressed -X POST "$BASE_URL/session" -H "Content-Type: application/json" -d '{}' | jq -r '.id')
    if [ -z "$SESSION_ID" ] || [ "$SESSION_ID" == "null" ]; then
        echo "❌ 错误: 创建会话失败。"
        exit 1
    fi
    log_kv "sessionID" "$SESSION_ID"

    TMP_EVENT_FILE="/tmp/opencode_event_${SESSION_ID}.json"
    rm -f "$TMP_EVENT_FILE" 2>/dev/null || true

    log_step "监听 SSE /event (等待 assistant 或 session.error)"
    capture_event_for_session "$SESSION_ID" "$TMP_EVENT_FILE" &
    LISTENER_PID=$!
    log_kv "listenerPID" "$LISTENER_PID"

    log_step "发送消息到 Session"
    log_kv "POST" "$BASE_URL/session/$SESSION_ID/message"
    log_kv "期望模型" "$expected_provider/$expected_model"
    log_kv "X-User-ID" "$user_id"
    log_kv "X-Model-ID(duckdb)" "$duckdb_model_id"
    log_kv "variant(user_id)" "$user_id"
    BODY=$(jq -n --arg pid "$expected_provider" --arg mid "$expected_model" --arg text "$message_text" --arg v "$user_id" '{model:{providerID:$pid,modelID:$mid},variant:$v,parts:[{type:"text",text:$text}]}')
    preview_json "请求 body" "$BODY"
    curl -s --compressed -X POST "$BASE_URL/session/$SESSION_ID/message" \
      -H "Content-Type: application/json" \
      -H "X-User-ID: $user_id" \
      -H "X-Model-ID: $duckdb_model_id" \
      -d "$BODY" >/dev/null

    for _ in $(seq 1 40); do
        if [ -s "$TMP_EVENT_FILE" ]; then
            break
        fi
        sleep 0.5
    done

    if kill -0 "$LISTENER_PID" 2>/dev/null; then
        pkill -P "$LISTENER_PID" 2>/dev/null || true
        kill "$LISTENER_PID" 2>/dev/null || true
    fi

    if [ ! -s "$TMP_EVENT_FILE" ]; then
        echo "   ❌ 验证失败: 未收到事件回包"
        exit 1
    fi
    preview_json "捕获事件" "$(cat "$TMP_EVENT_FILE")"

    EVT_TYPE=$(jq -r '.type // empty' "$TMP_EVENT_FILE")
    if [ "$EVT_TYPE" = "session.error" ]; then
        ERR_MSG=$(jq -r '.properties.error.data.message // .properties.error.data // .properties.error // empty' "$TMP_EVENT_FILE")
        echo "   ❌ 验证失败: session.error $ERR_MSG"
        exit 1
    fi

    GOT_PROVIDER=$(jq -r '.properties.info.providerID // empty' "$TMP_EVENT_FILE")
    GOT_MODEL=$(jq -r '.properties.info.modelID // empty' "$TMP_EVENT_FILE")

    if [ "$GOT_PROVIDER" = "$expected_provider" ] && [ "$GOT_MODEL" = "$expected_model" ]; then
        echo "   ✅ $label: $GOT_PROVIDER/$GOT_MODEL"
        return 0
    fi

    echo "   ❌ 验证失败: $label (got: $GOT_PROVIDER/$GOT_MODEL)"
    exit 1
}

send_and_get_assistant_text() {
    local label=$1
    local expected_provider=$2
    local expected_model=$3
    local user_id=$4
    local duckdb_model_id=$5
    local message_text=$6

    log_step "创建 Session"
    log_kv "POST" "$BASE_URL/session"
    SESSION_ID=$(curl -s --compressed -X POST "$BASE_URL/session" -H "Content-Type: application/json" -d '{}' | jq -r '.id')
    if [ -z "$SESSION_ID" ] || [ "$SESSION_ID" == "null" ]; then
        echo "❌ 错误: 创建会话失败。"
        exit 1
    fi
    log_kv "sessionID" "$SESSION_ID"

    TMP_TEXT_FILE="/tmp/opencode_assistant_text_${SESSION_ID}.txt"
    TMP_ERR_FILE="/tmp/opencode_assistant_err_${SESSION_ID}.json"
    rm -f "$TMP_TEXT_FILE" "$TMP_ERR_FILE" 2>/dev/null || true

    log_step "监听 SSE /event (捕获 assistant 文本)"
    capture_assistant_text_for_session "$SESSION_ID" "$TMP_TEXT_FILE" "$TMP_ERR_FILE" &
    LISTENER_PID=$!
    log_kv "listenerPID" "$LISTENER_PID"

    log_step "发送消息到 Session"
    log_kv "POST" "$BASE_URL/session/$SESSION_ID/message"
    log_kv "期望模型" "$expected_provider/$expected_model"
    log_kv "X-User-ID" "$user_id"
    log_kv "X-Model-ID(duckdb)" "$duckdb_model_id"
    log_kv "variant(user_id)" "$user_id"

    BODY=$(jq -n --arg pid "$expected_provider" --arg mid "$expected_model" --arg text "$message_text" --arg v "$user_id" '{model:{providerID:$pid,modelID:$mid},variant:$v,parts:[{type:"text",text:$text}]}')
    preview_json "请求 body" "$BODY"
    curl -s --compressed -X POST "$BASE_URL/session/$SESSION_ID/message" \
      -H "Content-Type: application/json" \
      -H "X-User-ID: $user_id" \
      -H "X-Model-ID: $duckdb_model_id" \
      -d "$BODY" >/dev/null

    for _ in $(seq 1 60); do
        if [ -s "$TMP_ERR_FILE" ]; then
            break
        fi
        if [ -s "$TMP_TEXT_FILE" ]; then
            break
        fi
        sleep 0.5
    done

    if kill -0 "$LISTENER_PID" 2>/dev/null; then
        pkill -P "$LISTENER_PID" 2>/dev/null || true
        kill "$LISTENER_PID" 2>/dev/null || true
    fi

    if [ -s "$TMP_ERR_FILE" ]; then
        preview_json "捕获事件(session.error)" "$(cat "$TMP_ERR_FILE")"
        ERR_MSG=$(jq -r '.properties.error.data.message // .properties.error.data // .properties.error // empty' "$TMP_ERR_FILE")
        echo "   ❌ $label: session.error $ERR_MSG"
        exit 1
    fi

    if [ ! -s "$TMP_TEXT_FILE" ]; then
        echo "   ❌ $label: 未捕获到 assistant 文本"
        exit 1
    fi

    ASSISTANT_TEXT=$(cat "$TMP_TEXT_FILE")
    log_kv "assistant_text" "$ASSISTANT_TEXT"

    echo "$ASSISTANT_TEXT"
}

DEFAULT_PROVIDER=$(resolve_provider "$DEFAULT_MODEL_NAME")
if [ -z "$DEFAULT_PROVIDER" ]; then
    echo "❌ 错误: OpenCode 未发现任何已连接 provider。"
    exit 1
fi
log_step "选择用于验证的 provider/model"
log_kv "后端默认模型 name" "$DEFAULT_MODEL_NAME"
log_kv "duckdb_model_id" "$DEFAULT_MODEL_ID"

OPENCODE_PROVIDER_ID=$(provider_for_baseurl "$DEFAULT_MODEL_BASE_URL")
if [ -z "$OPENCODE_PROVIDER_ID" ]; then
    OPENCODE_PROVIDER_ID="$DEFAULT_PROVIDER"
fi
log_kv "OpenCode providerID" "$OPENCODE_PROVIDER_ID"

DEFAULT_OPENCODE_MODEL="$DEFAULT_MODEL_MODEL_NAME"
if ! provider_has_model "$OPENCODE_PROVIDER_ID" "$DEFAULT_OPENCODE_MODEL"; then
    DEFAULT_OPENCODE_MODEL=$(pick_model "$OPENCODE_PROVIDER_ID")
fi
if [ -z "$DEFAULT_OPENCODE_MODEL" ] || [ "$DEFAULT_OPENCODE_MODEL" = "null" ]; then
    echo "❌ 错误: 未能解析 OpenCode 模型 (provider: $OPENCODE_PROVIDER_ID)"
    exit 1
fi
log_kv "OpenCode 选中的 model" "$DEFAULT_OPENCODE_MODEL"

echo "4. 测试: 默认模型是否在 OpenCode 生效"
send_and_assert_model "默认" "$OPENCODE_PROVIDER_ID" "$DEFAULT_OPENCODE_MODEL" "user_001" "$DEFAULT_MODEL_ID" "你好（验证默认模型）"

if [ -n "$OTHER_MODEL_ID" ] && [ "$OTHER_MODEL_ID" != "null" ]; then
    log_step "解析后端备选模型 (DuckDB 记录)"
    log_kv "duckdb_model_id" "$OTHER_MODEL_ID"
    log_kv "name" "$OTHER_MODEL_NAME"
    log_kv "provider" "$OTHER_MODEL_PROVIDER_LABEL"
    log_kv "base_url" "$OTHER_MODEL_BASE_URL"
    log_kv "model_name" "$OTHER_MODEL_MODEL_NAME"

    OTHER_PROVIDER_ID=$(provider_for_baseurl "$OTHER_MODEL_BASE_URL")
    if [ -z "$OTHER_PROVIDER_ID" ]; then
        echo "5. 测试: 备选模型是否在 OpenCode 生效"
        echo "   - 跳过原因: OpenCode provider 列表中找不到匹配 base_url 的 provider"
    else
        OTHER_OPENCODE_MODEL="$OTHER_MODEL_MODEL_NAME"
        if ! provider_has_model "$OTHER_PROVIDER_ID" "$OTHER_OPENCODE_MODEL"; then
            OTHER_OPENCODE_MODEL=$(pick_model "$OTHER_PROVIDER_ID")
        fi

        if [ -n "$OTHER_OPENCODE_MODEL" ] && [ "$OTHER_OPENCODE_MODEL" != "null" ]; then
            echo "5. 测试: 备选模型是否在 OpenCode 生效"
            send_and_assert_model "备选" "$OTHER_PROVIDER_ID" "$OTHER_OPENCODE_MODEL" "user_002" "$OTHER_MODEL_ID" "你好（验证备选模型）"
        fi
    fi
fi

echo ""
log_step "记忆跨 Session 验证 (同一 UserID)"
MEM_USER_ID="memory_user_001"
MEM_SENTINEL="ALPHA-42"

log_kv "user_id" "$MEM_USER_ID"
log_kv "sentinel" "$MEM_SENTINEL"

if ! mem0_enabled; then
    echo "⚠️  跳过记忆验证: 未设置 DEEPSEEK_API_KEY（opencode-mem 需要用于记忆提取）"
    echo ""
    echo "✨ 动态验证结束！"
    exit 0
fi

MEM1_TEXT=$(send_and_get_assistant_text "记忆写入" "$OPENCODE_PROVIDER_ID" "$DEFAULT_OPENCODE_MODEL" "$MEM_USER_ID" "$DEFAULT_MODEL_ID" "请调用 memory({mode:\"add\", content:\"我的代号是 ${MEM_SENTINEL}\"}) 保存记忆。然后只回复 OK。")
log_kv "session1_result" "$MEM1_TEXT"

MEM2_TEXT=$(send_and_get_assistant_text "记忆读取" "$OPENCODE_PROVIDER_ID" "$DEFAULT_OPENCODE_MODEL" "$MEM_USER_ID" "$DEFAULT_MODEL_ID" "请调用 memory({mode:\"search\", query:\"代号\"}) 检索，并回答我的代号是什么。")
log_kv "session2_result" "$MEM2_TEXT"

if [[ "$MEM2_TEXT" == *"${MEM_SENTINEL}"* ]]; then
    echo "✅ 记忆验证成功: 新 Session 仍能回忆 (${MEM_SENTINEL})"
else
    echo "❌ 记忆验证失败: 新 Session 未回忆出 (${MEM_SENTINEL})"
fi

echo ""
echo "✨ 动态验证结束！"
