package frontend

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zapkub/pakkretqc/internal/middleware"
)

func (s *Server) customizationFieldsCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx         = r.Context()
		domain      = vars["domain"]
		project     = vars["project"]
		entity_name = vars["entity_name"] // example "Defect"
		almclient   = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.FieldsCollection(ctx, domain, project, entity_name, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}

}

func (s *Server) customizationListsRelatedToEntityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx         = r.Context()
		domain      = vars["domain"]
		project     = vars["project"]
		entity_name = vars["entity_name"] // example "Defect"
		almclient   = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.ListsRelatedToEntity(ctx, domain, project, entity_name, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}
}

func (s *Server) customizationPermissionCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx         = r.Context()
		domain      = vars["domain"]
		project     = vars["project"]
		entity_name = vars["entity_name"] // example "Defect"
		almclient   = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.PermissionCollection(ctx, domain, project, entity_name, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}
}

func (s *Server) customizationRelationsCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx         = r.Context()
		domain      = vars["domain"]
		project     = vars["project"]
		entity_name = vars["entity_name"] // example "Defect"
		almclient   = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.RelationsCollection(ctx, domain, project, entity_name, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}
}

func (s *Server) customizationTypesCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx         = r.Context()
		domain      = vars["domain"]
		project     = vars["project"]
		entity_name = vars["entity_name"] // example "Defect"
		almclient   = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.TypesCollection(ctx, domain, project, entity_name, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}
}

func (s *Server) customizationTypeSubtypesFieldCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx         = r.Context()
		domain      = vars["domain"]
		project     = vars["project"]
		entity_name = vars["entity_name"] // example "Defect"
		subtype_id  = vars["subtype_id"]  // id from types_Collection
		almclient   = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.TypeSubtypesFieldCollection(ctx, domain, project, entity_name, subtype_id, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}
}

func (s *Server) customizationEntityMetadataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx         = r.Context()
		domain      = vars["domain"]
		project     = vars["project"]
		entity_name = vars["entity_name"] // example "Defect"
		almclient   = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.EntityMetadata(ctx, domain, project, entity_name, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}
}

func (s *Server) customizationEntitiesCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx       = r.Context()
		domain    = vars["domain"]
		project   = vars["project"]
		almclient = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.EntitiesCollection(ctx, domain, project, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}
}

func (s *Server) customizationUsedListsCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx       = r.Context()
		domain    = vars["domain"]
		project   = vars["project"]
		almclient = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.UsedListsCollection(ctx, domain, project, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}
}

func (s *Server) customizationUsersCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		ctx       = r.Context()
		domain    = vars["domain"]
		project   = vars["project"]
		almclient = middleware.MustGetALMClient(ctx)
	)

	_, err := almclient.UserCollection(ctx, domain, project, r, w)
	if err != nil {
		log.Panic(fmt.Errorf("ALM Error :: %w", err))
		return
	}
}
