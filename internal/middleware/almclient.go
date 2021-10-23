package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zapkub/pakkretqc/internal/conf"
	"github.com/zapkub/pakkretqc/internal/perrors"
	"github.com/zapkub/pakkretqc/pkg/almsdk"
)

func MustGetALMClient(ctx context.Context) *almsdk.Client {
	var err error
	token, ok := GetSessionToken(ctx)
	if !ok {
		panic(fmt.Errorf("cannot get AML client token notfound: %w", perrors.Unauthenticated))
	}

	var almclient = almsdk.New(&almsdk.ClientOptions{
		Endpoint: conf.ALMEndpoint(),
	})
	err = almclient.Authenticate(ctx, token)
	if err != nil {
		panic(fmt.Errorf("cannot connect to alm %w", err))
	}
	return almclient
}

func ALMClient(n http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		defer func() {
			rco := recover()
			if rco != nil {
				var err error
				switch t := rco.(type) {
				case string:
					log.Printf("found error string from recover :: %v", rco)
					err = errors.New(t)
				case error:
					log.Printf("found error obj from recover :: %v", rco)
					err = t
				default:
					err = errors.New("Unknown error")
				}

				errStr := err.Error()
				log.Printf("errStr :: %v", errStr)
				if strings.Contains(errStr, "token notfound") {
					http.Redirect(rw, r, fmt.Sprintf("/login?then=%s", r.RequestURI), http.StatusTemporaryRedirect)
				} else {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
				}
			}
		}()

		log.Printf("ALMClient :: URL :: %v", r.URL)
		n.ServeHTTP(rw, r)
	})
}
