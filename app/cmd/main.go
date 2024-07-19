package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"user-permissions-api/pkg/api"
	"user-permissions-api/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func initLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout

	file, err := os.OpenFile("logs/server_log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Warn(fmt.Sprintf("Could not open log file: %v\n", err))
	} else {
		log.Out = file
	}
	return log
}

func adminRouters(log *logrus.Logger) *gin.Engine {
	router := gin.Default()

	router.POST("/insertUser", api.InsertHandler(log))
	router.POST("/deleteUser", api.DeleteHandler(log))
	router.POST("/addUserRights", api.AddUserRightsHandler(log))
	router.POST("/deleteUserRights", api.DeleteUserRightsHandler(log))

	return router
}

func clientRouters(log *logrus.Logger) *gin.Engine {
	router := gin.Default()
	router.POST("/checkAccess", api.CheckAccessHandler(log))

	return router
}

func runServer(log *logrus.Logger, router *gin.Engine, port string, wg *sync.WaitGroup) {
	defer wg.Done()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Info("Server exiting")
}

func main() {
	log := initLogger()

	db.ConnectDB(log)
	defer db.DataBase.Close(db.Ctx)
	db.CreateTable()

	wg.Add(1)
	go runServer(log, clientRouters(log), os.Getenv("PORT_CLIENT"), &wg)

	wg.Add(1)
	go runServer(log, adminRouters(log), os.Getenv("PORT_ADMIN"), &wg)

	wg.Wait()
}
