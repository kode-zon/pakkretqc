package frontend

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/zapkub/pakkretqc/internal/conf"
	"github.com/zapkub/pakkretqc/internal/session"
	"github.com/zapkub/pakkretqc/pkg/almsdk"
)

type loginPage struct {
	Username     string            `json:"username"`
	Domains      []*almsdk.Domains `json:"domains"`
	RedirectUrl  string            `json:"redirectUrl"`
	ErrorMessage string            `json:"errorMessage"`
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	method := r.Method
	var (
		loginPage loginPage
		token, _  = r.Cookie(session.CookieKey)
		almclient = almsdk.New(&almsdk.ClientOptions{Endpoint: conf.ALMEndpoint()})
	)
	defer func() {
		q := r.URL.Query()
		loginPage.RedirectUrl = q.Get("then")
		s.servePage(w, "login", loginPage)
	}()

	if token == nil {
		if username := r.FormValue("username"); len(username) > 0 && method == "POST" {
			log.Printf("login with %s", username)
			token := base64.URLEncoding.EncodeToString([]byte(username + ":" + r.FormValue("password")))
			err := almclient.Authenticate(ctx, token)
			if err != nil {
				log.Println(err)
				if errors.Is(err, almsdk.InvalidCredential) {
					loginPage.ErrorMessage = "Cannot login to ALM server. maybe your credential is invalid I guesss 🤔"
				}
				return
			}
			log.Printf("login with %s success", username)

			var cookietoken http.Cookie
			cookietoken.Path = "/"
			cookietoken.HttpOnly = true
			cookietoken.Name = session.CookieKey
			cookietoken.Value = token
			http.SetCookie(w, &cookietoken)

			var cookieUsername http.Cookie
			cookieUsername.Path = "/"
			cookieUsername.HttpOnly = true
			cookieUsername.Name = session.UserName
			cookieUsername.Value = username
			http.SetCookie(w, &cookieUsername)

			if redirect := r.FormValue("redirect"); len(redirect) > 0 {

				log.Printf("redirect to %s", redirect)
				http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
				return
			}

			//load domains
			domains, err := almclient.Domains(ctx)
			if err != nil {
				log.Println(err)
				return
			}
			if domain := r.FormValue("domain"); domain == "" {
				loginPage.Username = username
				loginPage.Domains = domains
			}
			fmt.Printf("%+v", domains)
		}
		return
	}

	if token.Value != "" && method == "POST" {
		switch r.FormValue("action") {
		case "cancel":
			for _, cook := range r.Cookies() {
				cook.MaxAge = -1
				http.SetCookie(w, cook)
			}
			s.servePage(w, "login", loginPage)
			return
		case "proceed":
			currentDomain := r.FormValue("currentDomain")
			http.Redirect(w, r, path.Join("/", "domains", currentDomain), http.StatusTemporaryRedirect)
			return
		}
	}

	token.MaxAge = -1
	http.SetCookie(w, token)

}
