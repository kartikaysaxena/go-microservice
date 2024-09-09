package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var (
	redisClient *redis.Client
	rateLimit   = 10
	timeWindow  = 1 * time.Minute
)

// Initialize Redis connection
func initRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return redisClient
}

func getIPAddr(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	return strings.Split(IPAddress, ":")[0]
}

// User represents a user entity
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory data store and sync lock
var (
	users     = make(map[int]User)
	nextID    = 1
	usersLock sync.Mutex
)

func RateLimiter(rateLimit float32, timeWindow float32) error {
	ctx := context.Background()


	if count >= rateLimit {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return nil
	}

	pipe := redisClient.TxPipeline()
	pipe.Incr(ctx, redisKey)

	if count == 0 {
		pipe.Expire(ctx, redisKey, timeWindow)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		realIP := getIPAddr(r)

		log.WithFields(log.Fields{
			"IP": realIP,
		}).Info("Request received")

		log.Printf("Started %s %s", r.Method, r.RequestURI)

		ctx := context.Background()

		redisKey := "request_count:" + realIP

		count, err := redisClient.Get(ctx, redisKey).Int()
		if err != nil && err != redis.Nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	
		log.WithFields(log.Fields{
			"count": count,
		}).Infof("Request count for IP %s", realIP)
	

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v", time.Since(start))
	})
}

// CreateUser creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	user.ID = nextID
	nextID++
	users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser retrieves a user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	usersLock.Lock()
	user, exists := users[id]
	usersLock.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates a user by ID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	_, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	updatedUser.ID = id
	users[id] = updatedUser
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser deletes a user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	_, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}

// ListUsers lists all users
func ListUsers(w http.ResponseWriter, r *http.Request) {
	usersLock.Lock()
	defer usersLock.Unlock()

	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}

	json.NewEncoder(w).Encode(userList)
}

func main() {
	redisClient = initRedis()
	http.Handle("/user/list", loggingMiddleware(http.HandlerFunc(ListUsers)))
	http.Handle("/user/create", loggingMiddleware(http.HandlerFunc(CreateUser)))
	http.Handle("/user/get", loggingMiddleware(http.HandlerFunc(ListUsers)))
	http.Handle("/user/update", loggingMiddleware(http.HandlerFunc(UpdateUser)))
	http.Handle("/user/delete", loggingMiddleware(http.HandlerFunc(DeleteUser)))

	log.WithFields(log.Fields{
		"port": 8080,
	}).Info("Starting server")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
