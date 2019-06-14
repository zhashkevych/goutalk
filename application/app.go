package application

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/goutalk/chat"
	"github.com/zhashkevych/goutalk/application/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server

	chatUsecase chat.UseCase
	userRepo    chat.UserRepository
	roomRepo    chat.RoomRepository
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(addr string) error {
	ctx := context.Background()

	h, err := getHandler()
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    addr,
		Handler: h,
	}

	log.Printf("Starting HTTP application on port %s", addr)
	go func() {
		// service connections
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(ctx, 10*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func (a *App) Stop() {
	ctx := context.Background()

	if a.httpServer != nil {
		log.Print("Stopping http application")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()

		err := a.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

func getHandler() (http.Handler, error) {
	ginHandler := gin.New()
	ginHandler.Use(gin.Recovery())

	h := handler.NewHandler()

	ginHandler.POST("/login", h.Login)

	ginHandler.GET("/users", handler.AuthMiddleware, h.GetUsers)
	ginHandler.GET("/users/:id", handler.AuthMiddleware, h.GetUserByID)

	ginHandler.POST("/rooms", handler.AuthMiddleware)
	ginHandler.POST("/rooms/:id/join", handler.AuthMiddleware)
	ginHandler.POST("/rooms/:id/leave", handler.AuthMiddleware)
	ginHandler.GET("/rooms", handler.AuthMiddleware)
	ginHandler.GET("/rooms/:id", handler.AuthMiddleware)
	ginHandler.DELETE("/rooms/:id", handler.AuthMiddleware)

	ginHandler.POST("/message")

	return ginHandler, nil
}
