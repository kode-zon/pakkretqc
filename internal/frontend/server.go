package frontend

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/zapkub/pakkretqc/internal/fsutil"
	"github.com/zapkub/pakkretqc/internal/session"
)

type Server struct {
	apptemplate map[string]*template.Template
}

func (s *Server) Install(handle func(string, http.Handler)) {
	//handle by :: qcbin/rest
	handle("/domains/{domain}/projects/{project}/customization/entities/{entity_name}/fields", http.HandlerFunc(s.customizationFieldsCollectionHandler))
	handle("/domains/{domain}/projects/{project}/customization/entities/{entity_name}/lists", http.HandlerFunc(s.customizationListsRelatedToEntityHandler))
	handle("/domains/{domain}/projects/{project}/customization/entities/{entity_name}/permissions", http.HandlerFunc(s.customizationPermissionCollectionHandler))
	handle("/domains/{domain}/projects/{project}/customization/entities/{entity_name}/relations", http.HandlerFunc(s.customizationRelationsCollectionHandler))
	handle("/domains/{domain}/projects/{project}/customization/entities/{entity_name}/types", http.HandlerFunc(s.customizationTypesCollectionHandler))
	handle("/domains/{domain}/projects/{project}/customization/entities/{entity_name}/types/{subtype_id}/fields", http.HandlerFunc(s.customizationTypeSubtypesFieldCollectionHandler))
	handle("/domains/{domain}/projects/{project}/customization/entities/{entity_name}", http.HandlerFunc(s.customizationEntityMetadataHandler))
	handle("/domains/{domain}/projects/{project}/customization/entities", http.HandlerFunc(s.customizationEntitiesCollectionHandler))
	handle("/domains/{domain}/projects/{project}/customization/used-lists", http.HandlerFunc(s.customizationUsedListsCollectionHandler))
	handle("/domains/{domain}/projects/{project}/customization/users", http.HandlerFunc(s.customizationUsersCollectionHandler))
	//handle by :: qcbin/api
	handle("/login", http.HandlerFunc(s.loginHandler))
	handle("/domains/{domain}/projects/{project}/{collection_name}/$metadata/fields", http.HandlerFunc(s.metadataHandler))
	handle("/domains/{domain}/projects/{project}/attachments/{id}", http.HandlerFunc(s.attachmentDownloadHandler))
	handle("/domains/{domain}/projects/{project}/defects/{id}", http.HandlerFunc(s.defectPageHandler))
	handle("/domains/{domain}/projects/{project}/test-instances/{id}", http.HandlerFunc(s.testInstancesHandler))
	handle("/domains/{domain}/projects/{project}/list-items", http.HandlerFunc(s.listItemsHandler))
	handle("/domains/{domain}/projects/{project}", http.HandlerFunc(s.projectHandler))
	handle("/domains/{domain}", http.HandlerFunc(s.domainHandler))
	//self handle
	handle("/", http.HandlerFunc(s.indexHandler))
}

func parseTemplates(filename string) *template.Template {
	var err error
	tmpl, err := template.
		New("base.html").
		Funcs(template.FuncMap{
			"toJSON": EncodeJSON,
		}).
		ParseFiles(
			fsutil.PathFromWebDir("common/base.html"),
			filename,
		)
	if err != nil {
		panic(fmt.Sprintf("BUG: cannot parse template %+v", err))
	}
	return tmpl
}

func (s *Server) servePage(w http.ResponseWriter, pagename string, page interface{}) {
	log.Println("serve page", pagename)

	err := s.apptemplate[pagename].Execute(w, page)
	if err != nil {
		panic(fmt.Sprintf("BUG: cannot serve page %+v", err))
	}
}

func New() *Server {
	return &Server{
		apptemplate: map[string]*template.Template{
			"index":   parseTemplates(fsutil.PathFromWebDir("index.html")),
			"login":   parseTemplates(fsutil.PathFromWebDir("login.html")),
			"domain":  parseTemplates(fsutil.PathFromWebDir("domain.html")),
			"project": parseTemplates(fsutil.PathFromWebDir("project.html")),
			"defect":  parseTemplates(fsutil.PathFromWebDir("defect.html")),
		},
	}
}

func UserName(r *http.Request) string {
	if usernamecookie, err := r.Cookie(session.UserName); err == nil {
		return usernamecookie.Value
	}

	return ""
}
