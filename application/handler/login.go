package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/zhashkevych/goutalk/auth"
	"net/http"
)

type loginInput struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

type loginOutput struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"user_name"`
	Credentials string    `json:"credentials"`
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
		log.Errorf(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateAuthToken(user.ID, user.Username, user.Password)
	if err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &loginOutput{
		ID:          user.ID,
		Username:    user.Username,
		Credentials: token,
	})
}
