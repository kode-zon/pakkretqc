package frontend

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zapkub/pakkretqc/internal/middleware"
)

func (s *Server) metadataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx             = r.Context()
		domain          = vars["domain"]
		project         = vars["project"]
		collection_name = vars["collection_name"] // example "defects"
		almclient       = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.Metadata(ctx, domain, project, collection_name, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}

}
