package middleware

import (
	"avito_hr/pkg/session"
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type pair struct {
	method string
	path   string
}

type Middleware struct {
	SessionsManager    session.SessionsPacker
	mu                 *sync.RWMutex
	unmonitoredQueries map[pair]struct{}
}

func NewMiddleware(packer session.SessionsPacker) *Middleware {
	return &Middleware{
		SessionsManager:    packer,
		mu:                 &sync.RWMutex{},
		unmonitoredQueries: make(map[pair]struct{}),
	}
}

func (middle *Middleware) AddUnmonitoredQuery(method, path string) {
	middle.mu.Lock()
	defer middle.mu.Unlock()

	middle.unmonitoredQueries[pair{
		method: method,
		path:   path,
	}] = struct{}{}
}

func (middle *Middleware) checkUnmonitoredQuery(method, path string) bool {
	middle.mu.RLock()
	defer middle.mu.RUnlock()

	_, ok := middle.unmonitoredQueries[pair{
		method: method,
		path:   path,
	}]

	return ok
}

func (middle *Middleware) CheckRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middle.checkUnmonitoredQuery(r.Method, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("token")

		if token == "" {
			log.Error("no token in header for request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		s, err := middle.SessionsManager.Unpack(token)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "role", s.Sub)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (middle *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"Path":   r.URL.Path,
			"Method": r.Method,
			"Host":   r.Host,
			"Header": r.Header,
		}).Info("trying handle query")

		start := time.Now()
		next.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"Time spent": time.Since(start),
		}).Info("Handled query")
	})
}

func (middle *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"Panic message": err,
				}).Error("Catching panic...")
			}
		}()*/

		next.ServeHTTP(w, r)
	})
}

func (middle *Middleware) PackInMiddleware(next http.Handler) http.Handler {
	return middle.RecoverPanic(
		middle.Logging(
			middle.CheckRole(
				next,
			),
		),
	)
}
