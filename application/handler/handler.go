package handler

import "github.com/zhashkevych/goutalk/chat"

type Handler struct {
	chatter chat.UseCase
}

func NewHandler(uc chat.UseCase) *Handler {
	return &Handler{
		chatter: uc,
	}
}
