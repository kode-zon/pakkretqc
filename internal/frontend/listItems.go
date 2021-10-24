package frontend

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zapkub/pakkretqc/internal/middleware"
)

func (s *Server) listItemsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx       = r.Context()
		domain    = vars["domain"]
		project   = vars["project"]
		almclient = middleware.MustGetALMClient(ctx)
	)

	err := almclient.ListItems(ctx, domain, project, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
	}

}
