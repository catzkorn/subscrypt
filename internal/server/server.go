package server

import (
	"fmt"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/shopspring/decimal"
	"html/template"
	"net/http"
	"time"
)

// Server is the HTTP interface for subscription information
type Server struct {
	dataStore DataStore
	router    *http.ServeMux
}

type IndexPageData struct {
	PageTitle string
	Subscriptions []subscription.Subscription
}

// DataStore provides an interface to store information about individual subscriptions
type DataStore interface {
	GetSubscriptions() ([]subscription.Subscription, error)
	RecordSubscription(subscription subscription.Subscription) (*subscription.Subscription, error)
}

// NewServer returns a instance of a Server
func NewServer(dataStore DataStore) *Server {
	s := &Server{dataStore: dataStore, router: http.NewServeMux()}
	s.router.Handle("/", http.HandlerFunc(s.subscriptionHandler))

	return s
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

// JsonContentType defines application/json
const JsonContentType = "application/json"

// processGetSubscription processes the GET subscription request, returning the store subscriptions as json
func (s *Server) processGetSubscription(w http.ResponseWriter) error {
	tmpl := template.Must(template.ParseFiles("./web/templates/index.html"))

	subscriptions, err := s.dataStore.GetSubscriptions()

	if err != nil {
		return err
	}

	data := IndexPageData{
		PageTitle: "My Subscriptions List",
		Subscriptions: subscriptions,
	}

	err = tmpl.Execute(w, data)

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
	w.WriteHeader(http.StatusAccepted)

}
