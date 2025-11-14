package stat

import (
	"fmt"
	"go/adv-demo/configs"
	"go/adv-demo/pkg/middleware"
	"go/adv-demo/pkg/response"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandler struct {
	StatRepository *StatRepository
}

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}
	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStatistic(), deps.Config))
}

func (handler *StatHandler) GetStatistic() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Get stat")

		from, err := time.Parse("2006-01-02", req.URL.Query().Get("from"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(from)

		to, err := time.Parse("2006-01-02", req.URL.Query().Get("to"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(to)

		by := req.URL.Query().Get("by")

		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Incorrect by", http.StatusBadRequest)
			return
		}

		fmt.Println(by)

		stats := handler.StatRepository.GetStatistic(by, from, to)

		response.Json(w, stats, http.StatusOK)

	}
}
