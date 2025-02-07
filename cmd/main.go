package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"template/internal/api"
	"template/internal/common"
	"template/internal/config"
	"template/internal/db"
	"template/internal/handler"
	"template/internal/phonebook"
	"template/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

func Execute() {
	common.SetLogger(common.NewLogrusLogger())

	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("error connect DB: %s", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal()
		}
	}()

	repo := repository.NewPostgreSQLRepository(db)

	userSvc := phonebook.NewUserService(repo)
	addressSvc := phonebook.NewAddressService(repo)

	handler := handler.NewRESTHandler(userSvc, addressSvc)

	r := api.Setup(
		api.Route{Method: "POST", Path: "/register", Handler: []gin.HandlerFunc{handler.Register}},
		api.Route{Method: "POST", Path: "/login", Handler: []gin.HandlerFunc{handler.Login}},

		api.Route{Method: "POST", Path: "/addresses", Handler: []gin.HandlerFunc{api.Authentication(), handler.NewAddress}},
		api.Route{Method: "GET", Path: "/addresses", Handler: []gin.HandlerFunc{api.Authentication(), handler.Addresses}},
		api.Route{Method: "GET", Path: "/addresses/user", Handler: []gin.HandlerFunc{api.Authentication(), handler.GetAddressesByUserID}},
		api.Route{Method: "GET", Path: "/addresses/:id", Handler: []gin.HandlerFunc{api.Authentication(), handler.GetAddressByID}},
		api.Route{Method: "PUT", Path: "/addresses/:id", Handler: []gin.HandlerFunc{api.Authentication(), handler.UpdateAddress}},
		api.Route{Method: "DELETE", Path: "/addresses/:id", Handler: []gin.HandlerFunc{api.Authentication(), handler.DeleteAddress}},
	)

	srv := http.Server{
		Addr:    config.PORT,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	<-ctx.Done()

	log.Println("Server exited gracefully")
}
