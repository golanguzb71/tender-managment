package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
	repository "tender-managment/internal/db/repo"
)

type WebSocketManager struct {
	connections map[int]*websocket.Conn
	mutex       sync.Mutex
}

var wsManager = WebSocketManager{connections: make(map[int]*websocket.Conn)}

func WebSocketHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := jwtParser(tokenString)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Accept all origins for now
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	wsManager.mutex.Lock()
	wsManager.connections[userID] = conn
	wsManager.mutex.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	wsManager.mutex.Lock()
	delete(wsManager.connections, userID)
	wsManager.mutex.Unlock()
}

func SendNotification(repo repository.BidRepository, userID int, message string, relationID string, relationType string) {
	err := repo.CreateNotification(userID, message, relationID, relationType)
	if err != nil {
		log.Println("Error inserting notification:", err)
		return
	}

	wsManager.mutex.Lock()
	defer wsManager.mutex.Unlock()

	if conn, ok := wsManager.connections[userID]; ok {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error sending WebSocket notification:", err)
		}
	}
}
