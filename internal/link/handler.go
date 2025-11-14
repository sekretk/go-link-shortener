package link

import (
	"fmt"
	"go/adv-demo/configs"
	"go/adv-demo/pkg/event"
	"go/adv-demo/pkg/middleware"
	"go/adv-demo/pkg/request"
	"go/adv-demo/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	Config         *configs.Config
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}
	router.Handle("POST /link", middleware.IsAuthed(handler.Create(), deps.Config))
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.Delete(), deps.Config))
	router.Handle("GET /link", middleware.IsAuthed(handler.GetList(), deps.Config))
	router.HandleFunc("GET /{hash}", handler.GoTo())
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if email, ok := req.Context().Value(middleware.ContextEmailKey).(string); ok {
			fmt.Printf("Create %s\n", email)
		}

		body, err := request.HandleBody[LinkCreateRequest](&w, req)

		if err != nil {
			return
		}

		link := &Link{
			Url:  body.Url,
			Hash: "FirstL",
		}

		link.GenerateHash()

		for {
			isUnique, err := handler.LinkRepository.IsUniqueHash(link.Hash)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if isUnique {
				break
			}

			link.GenerateHash()
		}

		createdLink, err := handler.LinkRepository.Create(link)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, createdLink, http.StatusCreated)
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("GoTo")

		hash := req.PathValue("hash")

		link, err := handler.LinkRepository.GetByHash(hash)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})

		http.Redirect(w, req, link.Url, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) GetList() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Read")

		limit, err := strconv.Atoi(req.URL.Query().Get("limit"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(req.URL.Query().Get("offset"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		links := handler.LinkRepository.GetList(uint(limit), uint(offset))

		count := handler.LinkRepository.Count()

		response.Json(w, GetLinksListResponse{
			Links: links,
			Count: count,
		}, http.StatusOK)

	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if email, ok := req.Context().Value(middleware.ContextEmailKey).(string); ok {
			fmt.Printf("Update %s\n", email)
		}

		body, err := request.HandleBody[LinkUpdateRequest](&w, req)
		if err != nil {
			return
		}
		idStr := req.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, link, 201)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Delete")
		idStr := req.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.LinkRepository.GetById(uint(id))

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Json(w, nil, 200)
	}
}
