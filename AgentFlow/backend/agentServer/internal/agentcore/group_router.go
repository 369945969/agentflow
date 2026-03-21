package agentcore

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"agentServer/internal/protocol"
)

type GroupRouter struct {
	core *AgentCore
}

func NewGroupRouter(core *AgentCore) *GroupRouter {
	return &GroupRouter{core: core}
}

func (gr *GroupRouter) Handle(ctx context.Context, execCtx ExecutionContext) {
	member, err := gr.determineHandler(execCtx)
	if err == nil && member != nil {
		execCtx = gr.buildRoleContext(member, execCtx)
	}

	agent, err := gr.selectAgent(execCtx, member)
	if err != nil {
		if gr.core.emitter != nil {
			gr.core.emitter.SendToConnection(execCtx.ConnectionID, protocol.ServerMessage{
				Type:      "error",
				MessageID: execCtx.Message.MessageID,
				SessionID: execCtx.Session.SessionID,
				Timestamp: time.Now().UnixMilli(),
				Payload: protocol.ErrorPayload{
					Code:      "no_available_agent",
					Message:   err.Error(),
					Retryable: false,
				},
			})
		}
		return
	}

	execCtx.TargetAgent = agent
	execCtx.Session.LastAgentID = agent.AgentID
	_ = gr.core.sessionMgr.UpdateSession(&execCtx.Session)
	gr.core.Stream(ctx, execCtx)
}

func (gr *GroupRouter) determineHandler(ctx ExecutionContext) (*GroupMember, error) {
	if ctx.Session.GroupInfo == nil {
		return nil, errors.New("not a group session")
	}

	mentions := ctx.Message.Content.Mentions
	if len(mentions) > 0 {
		for _, mentionID := range mentions {
			for i := range ctx.Session.GroupInfo.Members {
				m := ctx.Session.GroupInfo.Members[i]
				if m.UserID == mentionID && m.IsActive {
					return &m, nil
				}
			}
		}
	}

	required := ctx.Message.Metadata.RequiredSkills
	if len(required) == 0 {
		required = gr.analyzeRequiredSkills(ctx.Message.Content.Text)
	}

	candidates := gr.findQualifiedMembers(ctx.Session.GroupInfo.Members, required)
	if len(candidates) == 0 {
		return nil, nil
	}

	return gr.selectBestCandidate(candidates, ctx, required), nil
}

func (gr *GroupRouter) analyzeRequiredSkills(text string) []string {
	skillKeywords := map[string][]string{
		"code":     {"代码", "编程", "bug", "函数", "写个", "实现"},
		"analysis": {"分析", "统计", "数据", "图表", "计算"},
		"writing":  {"写作", "写文章", "文案", "总结", "报告"},
		"image":    {"图片", "图像", "生成图", "画个"},
	}

	var skills []string
	for skill, kws := range skillKeywords {
		for _, kw := range kws {
			if kw != "" && strings.Contains(text, kw) {
				skills = append(skills, skill)
				break
			}
		}
	}
	return skills
}

func (gr *GroupRouter) findQualifiedMembers(members []GroupMember, required []string) []GroupMember {
	if len(required) == 0 {
		out := make([]GroupMember, 0, len(members))
		for _, m := range members {
			if m.IsActive {
				out = append(out, m)
			}
		}
		return out
	}

	reqSet := make(map[string]struct{}, len(required))
	for _, s := range required {
		reqSet[s] = struct{}{}
	}

	var out []GroupMember
	for _, m := range members {
		if !m.IsActive {
			continue
		}
		for _, ms := range m.Skills {
			if _, ok := reqSet[ms.SkillID]; ok {
				out = append(out, m)
				break
			}
		}
	}
	return out
}

func (gr *GroupRouter) selectBestCandidate(candidates []GroupMember, ctx ExecutionContext, required []string) *GroupMember {
	type scored struct {
		member GroupMember
		score  float64
	}

	reqSet := make(map[string]struct{}, len(required))
	for _, s := range required {
		reqSet[s] = struct{}{}
	}

	var scoredList []scored
	for _, m := range candidates {
		score := 0.0

		skillScore := 0.0
		for _, ms := range m.Skills {
			if _, ok := reqSet[ms.SkillID]; ok {
				skillScore += float64(ms.Priority)
			}
		}
		score += skillScore * 0.4

		descWords := strings.Fields(m.Description)
		msgWords := strings.Fields(ctx.Message.Content.Text)
		overlap := intersection(descWords, msgWords)
		semanticScore := float64(len(overlap)) / float64(len(descWords)+1)
		score += semanticScore * 10 * 0.3

		loadScore := 1.0
		score += loadScore * 10 * 0.2

		priorityScore := float64(m.Priority) / 10.0
		score += priorityScore * 10 * 0.1

		scoredList = append(scoredList, scored{member: m, score: score})
	}

	sort.Slice(scoredList, func(i, j int) bool { return scoredList[i].score > scoredList[j].score })
	if len(scoredList) == 0 {
		return nil
	}
	m := scoredList[0].member
	return &m
}

func intersection(a []string, b []string) []string {
	set := make(map[string]struct{}, len(a))
	for _, x := range a {
		set[x] = struct{}{}
	}
	var out []string
	for _, x := range b {
		if _, ok := set[x]; ok {
			out = append(out, x)
		}
	}
	return out
}

func (gr *GroupRouter) buildRoleContext(member *GroupMember, ctx ExecutionContext) ExecutionContext {
	if member == nil {
		return ctx
	}

	newCtx := ctx
	newCtx.RoleIdentity = &RoleIdentity{
		UserID:          member.UserID,
		RoleDescription: member.Description,
		Skills:          member.Skills,
		Persona:         member.Persona,
	}

	var skillDescs []string
	for _, s := range member.Skills {
		skillDescs = append(skillDescs, fmt.Sprintf("- %s: %s", s.SkillName, s.Description))
	}

	newCtx.SystemPrompt = fmt.Sprintf(
		"\n# 角色设定\n你是群聊成员 \"%s\"\n角色描述：%s\n人设风格：%s\n\n# 你的专业技能\n%s\n\n# 当前任务\n群ID：%s\n触发消息：\"%s\"\n\n请基于你的角色设定，以第一人称处理上述消息。回复时要体现你的专业背景和角色特点。\n",
		member.UserID,
		member.Description,
		member.Persona,
		strings.Join(skillDescs, "\n"),
		ctx.Session.SessionID,
		ctx.Message.Content.Text,
	)

	return newCtx
}

func (gr *GroupRouter) selectAgent(ctx ExecutionContext, member *GroupMember) (SubAgentConfig, error) {
	if ctx.Message.Metadata.TargetAgentID != "" {
		a, ok := gr.core.pool.Get(ctx.Message.Metadata.TargetAgentID)
		if ok {
			return a, nil
		}
	}

	a, ok := gr.core.pool.Default()
	if !ok {
		return SubAgentConfig{}, errors.New("no agent configured")
	}
	return a, nil
}
