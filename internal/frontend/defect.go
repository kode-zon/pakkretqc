package frontend

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zapkub/pakkretqc/internal/middleware"
	"github.com/zapkub/pakkretqc/internal/session"
	"github.com/zapkub/pakkretqc/pkg/almsdk"
)

type defectPage struct {
	Defect       *almsdk.Defect       `json:"defect"`
	Attachment   []*almsdk.Attachment `json:"attachment"`
	Project      string               `json:"project"`
	Domain       string               `json:"domain"`
	Username     string               `json:"username"`
	UserFullName string               `json:"userfullname"`
}

func (s *Server) defectPageHandler(w http.ResponseWriter, r *http.Request) {
	var (
		page      defectPage
		vars      = mux.Vars(r)
		domain    = vars["domain"]
		project   = vars["project"]
		id        = vars["id"]
		ctx       = r.Context()
		almclient = middleware.MustGetALMClient(ctx)
	)

	if r.Method == "PUT" {
		err := almclient.PutDefect(ctx, domain, project, id, r, w)
		if err != nil {
			panic(err)
		}
	} else {
		deflect, err := almclient.Defect(ctx, domain, project, id, r)
		if err != nil {
			panic(err)
		}
		page.Defect = deflect

		attachment, err := almclient.Attachments(ctx, domain, project, fmt.Sprintf("parent-id = %s ; parent-type = '%s'", id, "defect"), 10, 0)
		if err != nil {
			panic(err)
		}

		//TODO: move this to react instead
		currentUserCookie, err := r.Cookie(session.UserName)
		if err != nil {
			panic(err)
		}
		currentUser, err := almclient.UserDetail(ctx, domain, project, currentUserCookie.Value, r, nil)
		if err != nil {
			panic(err)
		}

		page.Attachment = attachment
		page.Domain = domain
		page.Project = project
		page.Username = currentUser.Name
		page.UserFullName = currentUser.FullName

		s.servePage(w, "defect", page)
	}
}

func (s *Server) attachmentDownloadHandler(w http.ResponseWriter, r *http.Request) {

	var (
		ctx       = r.Context()
		vars      = mux.Vars(r)
		domain    = vars["domain"]
		project   = vars["project"]
		id        = vars["id"] //attach id
		almclient = middleware.MustGetALMClient(ctx)
	)

	err := almclient.Attachment(ctx, domain, project, id, r, w)
	if err != nil {
		log.Printf("ERROR: %+v", err)
	}

}

func (s *Server) defectAttachmentUploadHandler(w http.ResponseWriter, r *http.Request) {

	var (
		ctx       = r.Context()
		vars      = mux.Vars(r)
		domain    = vars["domain"]
		project   = vars["project"]
		id        = vars["id"]
		almclient = middleware.MustGetALMClient(ctx)
	)

	err := almclient.PostAttachment(ctx, domain, project, "defects", id, r, w)
	if err != nil {
		log.Printf("ERROR: %+v", err)
		panic(err)
	}

}
