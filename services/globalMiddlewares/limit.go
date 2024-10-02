package globalmiddlewares

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/santhosh3/ECOM/utils"
)

func RateLimitingMiddleware(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Using request context instead of a global context
			ctx := r.Context()

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("IP not found"))
				return
			}

			// Get the current count for the IP
			val, err := rdb.Get(ctx, ip).Result()
			if err == redis.Nil {
				// Key does not exist, set the count to 1 with 1 min expiration
				err = rdb.Set(ctx, ip, 1, time.Minute).Err()
				if err != nil {
					utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("server not reachable"))
					return
				}
			} else if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("server not reachable"))
				return
			} else {
				// Convert the value to an integer
				count, _ := strconv.Atoi(val)
				if count >= 10 {
					utils.WriteError(w, http.StatusTooManyRequests, fmt.Errorf("the API is at capacity, try again later"))
					return
				}
				// Increment the request count
				err = rdb.Incr(ctx, ip).Err()
				if err != nil {
					utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("server not reachable"))
					return
				}
			}

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}


























// func PerClientRateLimiter(handlerFunc http.HandlerFunc) http.HandlerFunc {
// 	type Client struct {
// 		limiter  *rate.Limiter
// 		lastSeen time.Time
// 	}

// 	var (
// 		mutex   sync.Mutex
// 		clients = make(map[string]*Client)
// 	)

// 	go func() {
// 		for {
// 			time.Sleep(time.Minute)
// 			mutex.Lock()
// 			for ip, client := range clients {
// 				fmt.Println(ip, client.limiter, client.lastSeen)
// 				if time.Since(client.lastSeen) > 3*time.Minute {
// 					delete(clients, ip)
// 				}
// 			}
// 			mutex.Unlock()
// 		}
// 	}()

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ip, _, err := net.SplitHostPort(r.RemoteAddr)
// 		if err != nil {
// 			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("IP Not Found"))
// 			return
// 		}
// 		mutex.Lock()
// 		if _, found := clients[ip]; !found {
// 			clients[ip] = &Client{limiter: rate.NewLimiter(1, 4)}
// 		}
// 		clients[ip].lastSeen = time.Now()
// 		if !clients[ip].limiter.Allow() {
// 			mutex.Unlock()
// 			message := types.RateLimitStruct{
// 				Status: "Request Failed",
// 				Body:   "The API is at capacity, try again later.",
// 			}
// 			utils.WriteJSON(w, http.StatusTooManyRequests, message)
// 			return
// 		}
// 		mutex.Unlock()
// 		handlerFunc(w, r)
// 	}
// }
