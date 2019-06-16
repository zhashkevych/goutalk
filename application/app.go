package application

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/goutalk/chat/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
	"github.com/zhashkevych/goutalk/application/handler"
	"github.com/zhashkevych/goutalk/chat"
	repo "github.com/zhashkevych/goutalk/chat/repository/mongo"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	DBName = "goutalk"
)

type App struct {
	httpServer *http.Server

	mongoDB *mongo.Database

	chatUsecase chat.UseCase
	userRepo    chat.UserRepository
	roomRepo    chat.RoomRepository
}

func NewApp() *App {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	mongoDB := client.Database(DBName)

	userRepo := repo.NewUserRepository(mongoDB)
	roomRepo := repo.NewRoomsRepository(mongoDB)

	return &App{
		mongoDB:     mongoDB,
		userRepo:    userRepo,
		roomRepo:    roomRepo,
		chatUsecase: usecase.NewChatEngine(userRepo, roomRepo),
	}
}

func (a *App) Run(addr string) error {
	ctx := context.Background()

	h, err := a.getHandler()
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: h,
	}

	log.Printf("Starting HTTP application on port %s", addr)
	go func() {
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
		log.Print("Stopping HTTP application")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()

		err := a.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

func (a *App) getHandler() (http.Handler, error) {
	ginHandler := gin.New()
	ginHandler.Use(gin.Recovery(), gin.Logger())

	h := handler.NewHandler(a.chatUsecase)

	ginHandler.POST("/login", h.Login)

	ginHandler.GET("/users", handler.AuthMiddleware, h.GetUsers)
	ginHandler.GET("/users/:id", handler.AuthMiddleware, h.GetUserByID)

	ginHandler.POST("/rooms", handler.AuthMiddleware, h.CreateRoom)
	ginHandler.POST("/rooms/:id/join", handler.AuthMiddleware, h.JoinRoom)
	ginHandler.POST("/rooms/:id/leave", handler.AuthMiddleware, h.LeaveRoom)
	ginHandler.GET("/rooms", handler.AuthMiddleware, h.GetRooms)
	ginHandler.GET("/rooms/:id", handler.AuthMiddleware, h.GetRoomByID)
	ginHandler.DELETE("/rooms/:id", handler.AuthMiddleware, h.DeleteRoom)

	ginHandler.POST("/message")

	return ginHandler, nil
}
