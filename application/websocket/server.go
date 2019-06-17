package websocket

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server   *http.Server
	upgrader *websocket.Upgrader
	hub      *Hub
}

func NewServer(broadcast chan []byte) *Server {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	return &Server{
		upgrader: upgrader,
		hub:      newHub(broadcast),
	}
}

func (s *Server) Run(addr string) {
	wsHandler := s.getWSHandler()
	s.server = &http.Server{
		Addr:    addr,
		Handler: wsHandler,
	}

	log.Printf("Starting WebSocket application on port %s", addr)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen: %s", err)
		}
	}()
}

func (s *Server) Stop() {
	ctx := context.Background()

	if s.server != nil {
		log.Print("Stopping WebSocket application")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()

		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

func (s *Server) Broadcast(message []byte) {
	s.hub.broadcast <- message
}

func (s *Server) getWSHandler() http.Handler {
	ginHandler := gin.New()

	ginHandler.GET("/", func(c *gin.Context) {
		s.serveWS(c.Writer, c.Request)
	})

	return ginHandler
}

func (s *Server) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: s.hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.handleConnection()
}
