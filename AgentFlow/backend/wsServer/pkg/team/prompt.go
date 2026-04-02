package team

import (
	"fmt"
	"strings"
)

func BuildRouterPrompt(rules GroupRules, profiles []AgentProfile) string {
	var sb strings.Builder

	sb.WriteString("你是这个群组的管理员 Agent Team Leader。你的任务是分析用户的输入，并根据群组规则和成员能力，决定由哪个专家来回答。\n\n")
	
	// 1. 注入群组规则
	sb.WriteString("【群组规则】\n")
	switch rules.Mode {
	case "expert":
		sb.WriteString("- 专家模式：必须严格匹配用户问题领域与专家技能，只选择最合适的一位专家。\n")
	case "custom":
		sb.WriteString("- 自定义规则：" + rules.CustomRule + "\n")
	default: // free
		sb.WriteString("- 自由模式：选择最相关的专家即可，可以是多位（但在本系统中暂只支持单选）。\n")
	}
	sb.WriteString("\n")

	// 2. 注入成员画像
	sb.WriteString("【专家成员列表】\n")
	for _, p := range profiles {
		sb.WriteString(fmt.Sprintf("- ID: %s\n", p.ID))
		sb.WriteString(fmt.Sprintf("  名称: %s\n", p.Name))
		sb.WriteString(fmt.Sprintf("  描述: %s\n", p.Description))
		if len(p.Skills) > 0 {
			sb.WriteString(fmt.Sprintf("  技能: %s\n", strings.Join(p.Skills, ", ")))
		}
		sb.WriteString("\n")
	}

	// 3. 注入输出要求
	sb.WriteString(`【输出格式要求】
请仅输出一个 JSON 对象，不要包含任何 Markdown 格式或额外文本。
格式如下：
{
  "selected_agent_id": "成员ID",
  "reason": "选择该成员的理由",
  "instruction": "给该成员的具体执行指令（可选）"
}
如果用户的问题不需要特定专家回答（如闲聊），或没有合适的专家，请选择你自己（ID: "main_agent"）并直接回复。
`)

	return sb.String()
}

func BuildWorkerPrompt(profile AgentProfile, otherProfiles []AgentProfile, userInstruction string) string {
	var sb strings.Builder
	
	// 1. 角色设定
	sb.WriteString(fmt.Sprintf("你现在的身份是：%s\n", profile.Name))
	sb.WriteString(profile.Description + "\n\n")
	if profile.SystemPrompt != "" {
		sb.WriteString("核心指令：\n" + profile.SystemPrompt + "\n\n")
	}

	// 2. 团队上下文
	sb.WriteString("【团队上下文】\n")
	sb.WriteString("你处在一个协作群组中，其他成员包括：\n")
	for _, p := range otherProfiles {
		if p.ID == profile.ID {
			continue
		}
		sb.WriteString(fmt.Sprintf("- %s: %s\n", p.Name, p.Description))
	}
	sb.WriteString("\n")

	// 3. 当前任务
	if userInstruction != "" {
		sb.WriteString("【本次任务特别指令】\n")
		sb.WriteString(userInstruction + "\n")
	}

	// 4. 接力机制 (Chain of Thought)
	sb.WriteString(`
【回复要求】
1. 请直接回答用户的问题。
2. 如果你认为你的回答还不完整，或者需要其他专家的后续处理，请在回答的最后一行添加标记：
   @NextAgent: [建议的下一个专家名称或ID] 理由: [原因]
   
   例如：
   ...以上是数据分析结果。
   @NextAgent: writer_01 理由: 需要根据数据生成最终报告。
   
3. 如果你认为任务已结束，不需要添加 @NextAgent。
`)

	return sb.String()
}
