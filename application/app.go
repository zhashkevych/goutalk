package application

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zhashkevych/goutalk/application/ws"
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

	upgrader *websocket.Upgrader
	wsServer *http.Server
	hub      *ws.Hub

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

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	hub := ws.NewHub()

	return &App{
		mongoDB:     mongoDB,
		userRepo:    userRepo,
		roomRepo:    roomRepo,
		chatUsecase: usecase.NewChatEngine(userRepo, roomRepo, hub),
		upgrader:    upgrader,
		hub:         hub,
	}
}

func (a *App) Run(addr string) error {
	ctx := context.Background()

	h := a.getHandler()
	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: h,
	}

	log.Printf("Starting HTTP server on port %s", addr)
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen: %s", err)
		}
	}()

	wsHandler := a.getWSHandler()
	a.wsServer = &http.Server{
		Addr:    ":1030",
		Handler: wsHandler,
	}

	log.Printf("Starting WebSocket server on port 1030")
	go func() {
		if err := a.wsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen: %s", err)
		}
	}()

	a.hub.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(ctx, 10*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func (a *App) Stop() {
	ctx := context.Background()

	// shutting down HTTP server
	if a.httpServer != nil {
		log.Print("Stopping HTTP application")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()

		err := a.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	// shutting down WS server
	if a.wsServer != nil {
		log.Print("Stopping WebSocket application")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()

		err := a.wsServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

// get router for HTTP server
func (a *App) getHandler() http.Handler {
	ginHandler := gin.New()
	ginHandler.Use(gin.Recovery(), gin.Logger())

	h := handler.NewHandler(a.chatUsecase)

	ginHandler.POST("/login", h.Login)

	users := ginHandler.Group("/users", h.Authorize)
	{
		users.GET("/", h.GetUsers)
		users.GET("/:id", h.GetUserByID)
	}

	rooms := ginHandler.Group("/rooms", h.Authorize)
	{
		rooms.POST("/", h.CreateRoom)
		rooms.POST("/:id/join", h.JoinRoom)
		rooms.POST("/:id/leave", h.LeaveRoom)

		rooms.GET("/", h.GetRooms)
		rooms.GET("/:id", h.GetRoomByID)

		rooms.DELETE("/:id", h.DeleteRoom)
	}

	ginHandler.POST("/message", h.Authorize, h.SendMessage)

	return ginHandler
}

// get router for WS server
func (a *App) getWSHandler() http.Handler {
	ginHandler := gin.New()

	ginHandler.GET("/", func(c *gin.Context) {
		a.serveWS(c.Writer, c.Request)
	})

	return ginHandler
}

// On connect, create new client and append him to hub's connection pool
func (a *App) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := ws.NewClient(a.hub, conn)
	client.Run()
}
