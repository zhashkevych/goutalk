package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhashkevych/goutalk/auth"
	"github.com/zhashkevych/goutalk/chat"
	"net/http"
)

// TODO: responses with message

type loginInput struct {
	Username string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginOutput struct {
	ID          string `json:"id"`
	Username    string `json:"user_name"`
	Credentials string `json:"credentials"`
}

func (h *Handler) Login(c *gin.Context) {
	var inp loginInput
	if err := c.BindJSON(&inp); err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := h.chatter.LoginUser(c.Request.Context(), inp.Username, inp.Password)
	if err != nil {
		if err == chat.ErrWrongPassword {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &Response{
				"wrong password for user " + inp.Username,
			})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID := user.ID.Hex()

	token, err := auth.GenerateAuthToken(userID, user.Username, user.Password)
	if err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &loginOutput{
		ID:          userID,
		Username:    user.Username,
		Credentials: token,
	})
}
