package agentcore

import (
	"context"
)

type MessageRouter struct {
	single *SingleChatHandler
	group  *GroupRouter
}

func NewMessageRouter(single *SingleChatHandler, group *GroupRouter) *MessageRouter {
	return &MessageRouter{single: single, group: group}
}

func (mr *MessageRouter) Route(ctx context.Context, execCtx ExecutionContext) {
	if execCtx.Session.Type == "group" {
		mr.group.Handle(ctx, execCtx)
		return
	}
	mr.single.Handle(ctx, execCtx)
}
