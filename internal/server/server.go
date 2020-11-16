package server

import (
	"fmt"
	"github.com/Catzkorn/subscrypt/internal/plaid"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Server is the HTTP interface for subscription information
type Server struct {
	dataStore DataStore
	router    *http.ServeMux
	transactionAPI TransactionAPI
}

type TransactionAPI interface {
	GetTransactions() (plaid.TransactionList, error)
}

type IndexPageData struct {
	PageTitle string
	Subscriptions []subscription.Subscription
}

// DataStore provides an interface to store information about individual subscriptions
type DataStore interface {
	GetSubscriptions() ([]subscription.Subscription, error)
	RecordSubscription(subscription subscription.Subscription) (*subscription.Subscription, error)
	DeleteSubscription(ID int) error
	GetSubscription(ID int) (*subscription.Subscription, error)
}

// NewServer returns a instance of a Server
func NewServer(dataStore DataStore, transactionAPI TransactionAPI) *Server {
	s := &Server{dataStore: dataStore, router: http.NewServeMux(), transactionAPI: transactionAPI}
	s.router.Handle("/", http.HandlerFunc(s.subscriptionHandler))
	s.router.Handle("/api/subscriptions/", http.HandlerFunc(s.subscriptionsAPIHandler))
	s.router.Handle("/api/transactions/", http.HandlerFunc(s.transactionAPIHandler))

	return s
}

func (s *Server) transactionAPIHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		_ , err := s.transactionAPI.GetTransactions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// ServeHTTP implements the http handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// subscriptionHandler handles the routing logic for the index
func (s *Server) subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := s.processGetSubscription(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		s.processPostSubscription(w, r)
	}
}

// subscriptionsAPIHandler handles the routing logic for the '/api/subscriptions' paths
func (s *Server) subscriptionsAPIHandler(w http.ResponseWriter, r *http.Request) {
	urlID := strings.TrimPrefix(r.URL.Path, "/api/subscriptions/")
	ID, err := strconv.Atoi(urlID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete{
		s.processDeleteSubscription(w, ID)
	}
}

// JsonContentType defines application/json
const JsonContentType = "application/json"

// processGetSubscription processes the GET subscription request, returning the store subscriptions as json
func (s *Server) processGetSubscription(w http.ResponseWriter) error {

	subscriptions, err := s.dataStore.GetSubscriptions()

	if err != nil {
		return err
	}

	data := IndexPageData{
		PageTitle: "My Subscriptions List",
		Subscriptions: subscriptions,
	}

	err = ParsedIndexTemplate.Execute(w, data)

	if err != nil {
		return err
	}

	return nil
}

// processPostSubscription tells the SubscriptionStore to record the subscription from the post body
func (s *Server) processPostSubscription(w http.ResponseWriter, r *http.Request) {

	amount, _ := decimal.NewFromString(r.FormValue("amount"))

	layout := "2006-01-02"
	str := r.FormValue("date")

	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}

	entry := subscription.Subscription{
		Name:   r.FormValue("name"),
		Amount: amount,
		DateDue: t,
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = s.dataStore.RecordSubscription(entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// processDeleteSubscription tells the SubscriptionStore to delete the subscription with the given ID
func (s *Server) processDeleteSubscription(w http.ResponseWriter, ID int) {

	retrievedSubscription, err := s.dataStore.GetSubscription(ID)

	switch {
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	case retrievedSubscription == nil:
		errorMessage := "Failed to delete subscription - subscription not found"
		http.Error(w, errorMessage, http.StatusNotFound)
		return
	default:
		err = s.dataStore.DeleteSubscription(ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
