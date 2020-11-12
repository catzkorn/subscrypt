package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Catzkorn/subscrypt/internal/subscription"
)

// SubscriptionServer is the HTTP interface for subscription information
type Server struct {
	dataStore DataStore
	http.Handler
}

// DataStore provides an interface to store information about individual subscriptions
type DataStore interface {
	GetSubscriptions() ([]subscription.Subscription, error)
	RecordSubscription(subscription subscription.Subscription) (*subscription.Subscription, error)
}

// NewSubscriptionServer returns a instance of a SubscriptionServer
func NewServer(dataStore DataStore) *Server {
	s := new(Server)
	s.dataStore = dataStore
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(s.subscriptionHandler))
	s.Handler = router
	return s
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

const JsonContentType = "application/json"

// processGetSubscription processes the GET subscription request, returning the store subscriptions as json
func (s *Server) processGetSubscription(w http.ResponseWriter) error {
	w.Header().Set("content-type", JsonContentType)
	subscriptions, err := s.dataStore.GetSubscriptions()
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(subscriptions)
	if err != nil {
		return err
	}
	return nil
}

// processPostSubscription tells the SubscriptionStore to record the subscription from the post body
func (s *Server) processPostSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription subscription.Subscription

	err := json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		log.Fatalln(err)
	}
	s.dataStore.RecordSubscription(subscription)
	w.WriteHeader(http.StatusAccepted)
}
