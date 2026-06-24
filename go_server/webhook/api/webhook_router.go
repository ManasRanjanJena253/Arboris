package api

import (
	"Arboris/go_server/config"
	"Arboris/go_server/webhook/middleware"
	"Arboris/go_server/webhook/queue"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

func NewHookRouter(envVar config.Config, redisClient *redis.Client) (*chi.Mux, error) {
	hookRouter := chi.NewRouter()

	rateLimit := rate.Limit(envVar.WebHook.RateLimit)
	hookHandler := WebhookHandlerStruct{queue: queue.NewJobQueue(100)}
	cache := &sync.Map{}

	hookRouter.Use(middleware.MaxBodySize(int64(envVar.WebHook.PayloadMaxSize)))
	hookRouter.Use(middleware.Recovery)
	hookRouter.Use(middleware.VerifyHMAC(envVar.WebHook.Secret))
	hookRouter.Use(middleware.ExtractInstallID)
	hookRouter.Use(middleware.PreventReplay(redisClient))
	hookRouter.Use(middleware.RateLimiter(cache, rateLimit, envVar.WebHook.Burst))
	hookRouter.Use(middleware.Logger)

	hookRouter.Post("/webhook/github", hookHandler.WebhookHandler)

	return hookRouter, nil
}
