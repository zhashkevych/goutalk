package handler

import "github.com/zhashkevych/goutalk/chat"

type Handler struct {
	chatter chat.UseCase
}

type Response struct {
	Message string `json:"message"`
}

func NewHandler(uc chat.UseCase) *Handler {
	return &Handler{
		chatter: uc,
	}
}
