package subscription

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

const JsonContentType = "application/json"

// NewSubscriptionServer returns a instance of a SubscriptionServer
func NewSubscriptionServer(store SubscriptionStore) *SubscriptionServer {
	s := new(SubscriptionServer)
	s.store = store
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(s.subscriptionHandler))
	s.Handler = router
	return s
}

// subscriptionHandler handles the routing logic for the index
func (s *SubscriptionServer) subscriptionHandler(w http.ResponseWriter, r *http.Request) {
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

// processGetSubscription processes the GET subscription request, returning the store subscriptions as json
func (s *SubscriptionServer) processGetSubscription(w http.ResponseWriter) error {
	w.Header().Set("content-type", JsonContentType)
	subscriptions, err := s.store.GetSubscriptions()
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
func (s *SubscriptionServer) processPostSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription Subscription

	err := json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		log.Fatalln(err)
	}
	s.store.RecordSubscription(subscription)
	w.WriteHeader(http.StatusAccepted)
}

// SubscriptionServer is the HTTP interface for subscription information
type SubscriptionServer struct {
	store SubscriptionStore
	http.Handler
}

// SubscriptionStore provides an interface to store information about individual subscriptions
type SubscriptionStore interface {
	GetSubscriptions() ([]Subscription, error)
	RecordSubscription(subscription Subscription) error
}

// Subscription defines a subscription. ID is unique per subscription.
// Name is the name of the subscription stored as a string.
// Amount is the cost of the subscription, stored as a decimal.
// DateDue is the date that the subscription is due on, stored as a date.
type Subscription struct {
	ID      int
	Name    string
	Amount  decimal.Decimal
	DateDue time.Time
}
