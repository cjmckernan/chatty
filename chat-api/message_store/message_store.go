package message_store

import (
	"chat-api/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

const (
	userKey     = "user"
	topicSetKey = "topics"
)

var (
	ctx  = context.Background()
	rdb  *redis.Client
	host = "redis:6379"
)

type User struct {
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Message struct {
	SessionID string `json:"sessionId"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}


func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
	validateRedis()
	initTopics()
}

func initTopics() {
	topics := []string{"General", "Sports", "Technology", "Movies", "Music"}

	for _, topic := range topics {
		rdb.SAdd(ctx, topicSetKey, topic)
	}
}

func validateRedis() {
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis: connected", pong)
}

func getUser(username string) (*User, error) {
	userKey := "user:" + username

	userData, err := rdb.HGetAll(ctx, userKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if len(userData) == 0 {
		return nil, fmt.Errorf("user %s not found", username)
	}
	user := &User{
		Username:     userData["username"],
		PasswordHash: userData["password_hash"],
	}

	return user, nil
}

func createUserSessionID(username string) (string, error) {

	existingSessionID, err := getSessionIDByUsername(username)
	if err != nil {
		return "", err
	}
	if existingSessionID != "" {
		return existingSessionID, nil
	}

	sessionID, err := utils.GenerateSessionId()
	err = storeSessionId(sessionID, username)
	if err != nil {
		return "", fmt.Errorf("error storing sessionId")
	}

	err = storeSessionId(sessionID, username)
	if err != nil {
		return "", fmt.Errorf("error storing sessionId: %w", err)
	}
	return sessionID, nil
}

func getSessionIDByUsername(username string) (string, error) {
	userSessionsKey := "user_sessions:" + username
	sessionIDs, err := rdb.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve sessions for user: %w", err)
	}
	if len(sessionIDs) > 0 {
		return sessionIDs[0], nil
	}
	return "", nil
}

func storeSessionId(sessionID, username string) error {
	sessionKey := "session:" + sessionID
	_, err := rdb.HMSet(ctx, sessionKey, map[string]interface{}{
		"username":   username,
		"created_at": time.Now().UTC().Format(time.RFC3339),
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to store session %w", err)
	}
	userSessionsKey := "user_sessions:" + username
	_, err = rdb.SAdd(ctx, userSessionsKey, sessionID).Result()
	if err != nil {
		return fmt.Errorf("failed to link session to user: %w", err)
	}

	return nil
}

func UserExists(username string) (bool, error) {
	userKey := "user:" + username
	exists, err := rdb.Exists(ctx, userKey).Result()
	fmt.Println("retrieving user")
	if err != nil {
		return false, fmt.Errorf("Error validating user: %w ", err)
	}
	return exists > 0, nil
}

func Auth(username string, password string) (string, error) {
	user, err := getUser(username)
	if err != nil {
		return "", fmt.Errorf("error retrieving user %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	sessionID, err := createUserSessionID(username)
	if err != nil {
		return "", fmt.Errorf("error creating sessionID: %w", err)
	}
	return sessionID, nil

}

func GetMessagesByTopic(topic string) ([]Message, error) {
	topicKey := fmt.Sprintf("topic:%s:messages", topic)
	messagesData, err := rdb.LRange(ctx, topicKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	var messages []Message
	for _, messageJSON := range messagesData {
		var message Message
		if err := json.Unmarshal([]byte(messageJSON), &message); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}



func CreateUser(user User) (string, error) {
	ctx := context.Background()
	userKey := "user:" + user.Username

	_, err := rdb.HMSet(ctx, userKey, map[string]interface{}{
		"username":      user.Username,
		"password_hash": user.PasswordHash,
	}).Result()

	if err != nil {
		return "", fmt.Errorf("Failed to save user %w", err)
	}

	sessionID, err := createUserSessionID(user.Username)
	if err != nil {
		return "", fmt.Errorf("Error creating sessionID %w", err)
	}

	return sessionID, nil
}

func GetUsernameBySessionID(sessionID string) (string, error) {
	sessionKey := "session:" + sessionID
	sessionData, err := rdb.HGetAll(ctx, sessionKey).Result()
	if err != nil || len(sessionData) == 0 {
		return "", fmt.Errorf("session not found")
	}
	username := sessionData["username"]
	fmt.Println(username)

	return username, nil
}

func GetTopics() ([]string, error) {
	topics, err := rdb.SMembers(ctx, topicSetKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve topics: %w", err)
	}
	return topics, nil
}

func StoreMessage(topic, sessionId, username, text string) error { 
	message := Message{
		SessionID: sessionId,
		Username:  username,
		Text:      text,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	 }
  
  messageJSON, err := json.Marshal(message)
  if err != nil {
    return fmt.Errorf("failed to get message")
  }

  topicKey := fmt.Sprintf("topic:%s:messages", topic)
  if _, err := rdb.RPush(ctx, topicKey, messageJSON).Result(); err != nil {
    return fmt.Errorf("failed to store %w", err)
  }
  return nil
}
