package main

import (
	"fmt"
	"go/adv-demo/configs"
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/link"
	"go/adv-demo/internal/stat"
	"go/adv-demo/internal/user"
	"go/adv-demo/pkg/db"
	"go/adv-demo/pkg/event"
	"go/adv-demo/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventbus := event.NewEventBus()

	//Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	//Services
	authService := auth.NewUserService(userRepository)
	statService := *stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventbus,
		StatRepository: statRepository,
	})

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		Config:         conf,
		EventBus:       eventbus,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	//Middlewares
	stack := middleware.Chain(middleware.CORS, middleware.Logging)

	go statService.AddClick()

	return stack(router)
}

func main() {

	handler := App()

	server := http.Server{
		Addr:    ":8081",
		Handler: handler,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
