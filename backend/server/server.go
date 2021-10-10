package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// initRouter returns created server with controller endpoints.
func initRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())
	router.Use(loggingMiddleware())

	api := router.Group("/api")
	{
		api.POST("/user", CreateUser)
		api.GET("/user/:id", GetUser)

		api.GET("/game_week/:id", GetGameWeek)

		api.POST("/user/:id/next_week", NextWeek)
		api.POST("/user/:id/buy_instrument/:instrument_id", BuyInstrument)
		api.POST("/user/:id/sell_instrument/:instrument_id", SellInstrument)
		api.POST("/user/:id/add_test_answer/:test_answer_id", AddTestAnswer)

		api.GET("/user/:id/game_result", GetGameResult)
	}

	router.GET("/user/:id/dashboard", RenderResult)

	return router
}

// RunForever starts HTTP server on :9584 port.
func RunForever() {
	r := initRouter()

	PORT := ":9584"
	srv := &http.Server{
		Addr:    PORT,
		Handler: r,
	}

	go func() {
		log.Info(fmt.Sprintf("listening on %s", PORT))

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(fmt.Sprintf("listen and serve: %s", err.Error()))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Warn("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Sprintf("server shutdown: %s", err))
	}

	log.Info("server exited")
}