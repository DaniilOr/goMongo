package app

import (
	"context"
	"encoding/json"
	"github.com/DaniilOr/goMongo/cmd/service/app/dtos"
	"github.com/DaniilOr/goMongo/cmd/service/app/middleware/authenticator"
	"github.com/DaniilOr/goMongo/cmd/service/app/middleware/authorizator"
	"github.com/DaniilOr/goMongo/cmd/service/app/middleware/cacher"
	"github.com/DaniilOr/goMongo/cmd/service/app/middleware/identificator"
	"github.com/DaniilOr/goMongo/cmd/service/app/middleware/logger"
	"github.com/DaniilOr/goMongo/pkg/cache"
	"github.com/DaniilOr/goMongo/pkg/payments"
	"github.com/DaniilOr/goMongo/pkg/security"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	securitySvc *security.Service
	paymentsSvc *payments.Service
	router      chi.Router
	cacheSvc 	*cache.Service
}

func NewServer(securitySvc *security.Service, paymentsSvc *payments.Service, router chi.Router, cacheSvc *cache.Service) *Server {
	return &Server{securitySvc: securitySvc, paymentsSvc: paymentsSvc, router: router, cacheSvc: cacheSvc}
}

func (s *Server) Init() error {
	s.router.Put("/user", s.handleLogin)
	identificatorMd := identificator.Identificator
	authenticatorMd := authenticator.Authenticator(
		identificator.Identifier, s.securitySvc.UserDetails,
	)
	roleChecker := func(ctx context.Context, roles ...string) bool {
		userDetails, err := authenticator.Authentication(ctx)
		if err != nil {
			return false
		}
		return s.securitySvc.HasAnyRole(ctx, userDetails, roles...)
	}
	serviceRoleMd := authorizator.Authorizator(roleChecker, security.RoleService)
	userRoleMd := authorizator.Authorizator(roleChecker, security.RoleUser)
	logger := logger.Logger
	cacher := cacher.Cache(s.cacheSvc.FromCache, s.cacheSvc.ToCache)
	s.router.With(identificatorMd, authenticatorMd, serviceRoleMd, logger).Post("/service/add/suggestion/{id}", s.handleAdd)
	s.router.With(identificatorMd, authenticatorMd, userRoleMd, logger, cacher).Get("/user/get/suggestions/{id}", s.handleGet)
	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleLogin(writer http.ResponseWriter, request *http.Request) {
	login := request.PostFormValue("login")
	if login == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	password := request.PostFormValue("password")
	if password == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.securitySvc.Login(request.Context(), login, password)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &dtos.TokenDTO{Token: token}
	respBody, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(respBody)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	suggestedPayments, err := s.paymentsSvc.GetPayments(r, id)
	response := dtos.Response{SuggestedPayments: suggestedPayments}
	body, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var payment dtos.Payment
	err = json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.paymentsSvc.AddPredictedPayment(r, id, payment)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := dtos.ResultDTO{Result: "added"}
	body, err := json.Marshal(result)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
}
