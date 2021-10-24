package frontend

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zapkub/pakkretqc/internal/middleware"
)

func (s *Server) testInstancesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx       = r.Context()
		domain    = vars["domain"]
		project   = vars["project"]
		id        = vars["id"]
		almclient = middleware.MustGetALMClient(ctx)
	)

	err := almclient.TestInstances(ctx, domain, project, id, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
	}

}
