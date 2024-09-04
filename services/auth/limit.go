package auth

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/santhosh3/ECOM/types"
	"github.com/santhosh3/ECOM/utils"
	"golang.org/x/time/rate"
)

func PerClientRateLimiter(handlerFunc http.HandlerFunc) http.HandlerFunc {
	type Client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mutex   sync.Mutex
		clients = make(map[string]*Client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mutex.Lock()
			for ip, client := range clients {
				fmt.Println(ip,client.limiter,client.lastSeen)
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mutex.Unlock()
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("IP Not Found"))
			return
		}
		mutex.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &Client{limiter: rate.NewLimiter(1, 4)}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mutex.Unlock()
			message := types.RateLimitStruct{
				Status: "Request Failed",
				Body:   "The API is at capacity, try again later.",
			}
			utils.WriteJSON(w, http.StatusTooManyRequests, message)
			return
		}
		mutex.Unlock()
		handlerFunc(w, r)
	}
}

