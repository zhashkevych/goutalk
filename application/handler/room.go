package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/goutalk/chat"
	"net/http"
	"time"
)

type createRoomInput struct {
	Name string `json:"name" binding:"required"`
}

type manageUserInput struct {
	ID string `json:"user_id" binding:"required"`
}

type room struct {
	ID        string    `json:"room_id"`
	Name      string    `json:"room_name"`
	CreatorID string    `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
	Members   []*user   `json:"members"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var inp createRoomInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &Response{
			"wrong json body",
		})
		return
	}

	ctx := c.Request.Context()
	user := c.MustGet(ctxKeyUser).(*chat.User)

	r, err := h.chatter.CreateRoom(ctx, inp.Name, user.ID.Hex())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			"failed to create room",
		})
		return
	}

	c.JSON(http.StatusOK, h.toRoom(ctx, r))
}

func (h *Handler) JoinRoom(c *gin.Context) {
	var inp manageUserInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &Response{
			"wrong json body",
		})
		return
	}

	id := c.Param("id")

	if err := h.chatter.AddRoomMember(c.Request.Context(), id, inp.ID); err != nil {
		if _, ok := err.(*chat.ErrorNotFound); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, &Response{
				"no room with ID " + id + " found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			"failed to add member",
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		"user successfully joined the room",
	})
}

func (h *Handler) LeaveRoom(c *gin.Context) {
	var inp manageUserInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &Response{
			"wrong json body",
		})
		return
	}

	id := c.Param("id")

	if err := h.chatter.RemoveRoomMeber(c.Request.Context(), id, inp.ID); err != nil {
		if _, ok := err.(*chat.ErrorNotFound); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, &Response{
				"no room with ID " + id + " found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			"failed to remove member",
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		"user successfully left the room",
	})
}

func (h *Handler) GetRooms(c *gin.Context) {
	ctx := c.Request.Context()

	rooms, err := h.chatter.GetAllRooms(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			"failed to get rooms from db",
		})
		return
	}

	if rooms == nil || len(rooms) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, &Response{
			"no rooms found",
		})
		return
	}

	c.JSON(http.StatusOK, h.toRooms(ctx, rooms))
}

func (h *Handler) GetRoomByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	room, err := h.chatter.GetRoomByID(ctx, id)
	if err != nil {
		if _, ok := err.(*chat.ErrorNotFound); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, &Response{
				"no room with ID " + id + " found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			"failed to get room",
		})
		return
	}

	c.JSON(http.StatusOK, h.toRoom(ctx, room))
}

func (h *Handler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	user := c.MustGet(ctxKeyUser).(*chat.User)

	if err := h.chatter.DeleteRoom(c.Request.Context(), id, user); err != nil {
		if _, ok := err.(*chat.ErrorNotFound); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, &Response{
				"no room with ID " + id + " found",
			})
			return
		}

		if err == chat.ErrMissingAccessRights {
			c.AbortWithStatusJSON(http.StatusForbidden, &Response{
				"no room with ID " + id + " found",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			"failed to delete room",
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		"room removed successfully",
	})
}

func (h *Handler) toRooms(ctx context.Context, rooms []*chat.Room) []*room {
	out := make([]*room, len(rooms))
	for i := range rooms {
		out[i] = h.toRoom(ctx, rooms[i])
	}

	return out
}

func (h *Handler) toRoom(ctx context.Context, r *chat.Room) *room {
	members, _ := h.chatter.GetRoomMembers(ctx, r.ID.Hex())
	return &room{
		ID:        r.ID.Hex(),
		Name:      r.Name,
		CreatorID: r.CreatorID,
		CreatedAt: r.CreatedAt,
		Members:   toUsers(members),
	}
}
