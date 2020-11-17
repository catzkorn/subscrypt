package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/Catzkorn/subscrypt/internal/calendar"
	"github.com/Catzkorn/subscrypt/internal/email"
	"github.com/Catzkorn/subscrypt/internal/plaid"
	"github.com/Catzkorn/subscrypt/internal/reminder"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/Catzkorn/subscrypt/internal/userprofile"
)

// Server is the HTTP interface for subscription information
type Server struct {
	dataStore           DataStore
	router              *http.ServeMux
	parsedIndexTemplate *template.Template
	mailer              email.Mailer
	transactionAPI      TransactionAPI
}

// TransactionAPI defines the transaction api interface
type TransactionAPI interface {
	GetTransactions() (plaid.TransactionList, error)
}

// IndexPageData defines data shown on the page
type IndexPageData struct {
	PageTitle     string
	Subscriptions []subscription.Subscription
	Userprofile   *userprofile.Userprofile
}

// DataStore provides an interface to store information about individual subscriptions
type DataStore interface {
	GetSubscriptions() ([]subscription.Subscription, error)
	RecordSubscription(subscription subscription.Subscription) (*subscription.Subscription, error)
	DeleteSubscription(ID int) error
	GetSubscription(ID int) (*subscription.Subscription, error)
	RecordUserDetails(name string, email string) (*userprofile.Userprofile, error)
	GetUserDetails() (*userprofile.Userprofile, error)
}

// NewServer returns a instance of a Server
func NewServer(dataStore DataStore, indexTemplatePath string, mailer email.Mailer, transactionAPI TransactionAPI) *Server {
	s := &Server{dataStore: dataStore, router: http.NewServeMux(), transactionAPI: transactionAPI}
	s.router.Handle("/", http.HandlerFunc(s.indexHandler))
	s.router.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	s.router.Handle("/api/reminders", http.HandlerFunc(s.reminderHandler))
	s.router.Handle("/api/subscriptions", http.HandlerFunc(s.subscriptionsAPIHandler))
	s.router.Handle("/api/subscriptions/", http.HandlerFunc(s.subscriptionIDAPIHandler))
	s.router.Handle("/api/transactions/", http.HandlerFunc(s.transactionAPIHandler))
	s.router.Handle("/api/users", http.HandlerFunc(s.userHandler))

	s.parsedIndexTemplate = template.Must(template.New("index.html").ParseFiles(indexTemplatePath))

	s.mailer = mailer

	return s
}

func (s *Server) transactionAPIHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		transactions, err := s.transactionAPI.GetTransactions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		subscriptions := subscription.ProcessTransactions(transactions)
		for _, entry := range subscriptions {
			_, _ = s.dataStore.RecordSubscription(entry)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// ServeHTTP implements the http handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// indexHandler handles the routing logic for the index
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := s.processGetIndex(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) reminderHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.processPostReminder(w, r)
	}
}

// processPostReminder creates an ics file
// TODO: then emails it to the user's email
func (s *Server) processPostReminder(w http.ResponseWriter, r *http.Request) {
	var newReminder reminder.Reminder
	var newSubscription subscription.Subscription

	err := json.NewDecoder(r.Body).Decode(&newSubscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscription, err := s.dataStore.GetSubscription(newSubscription.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.dataStore.GetUserDetails()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newReminder = reminder.Reminder{
		Email:          user.Email,
		SubscriptionID: subscription.ID,
		ReminderDate:   subscription.DateDue.AddDate(0, 0, -5),
	}

	cal := calendar.CreateReminderInvite(*subscription, newReminder)

	err = email.SendEmail(newReminder, *user, cal, s.mailer, s.dataStore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// subscriptionsAPIHandler handles the routing logic for the '/api/subscriptions' paths
func (s *Server) subscriptionsAPIHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		s.processGetSubscriptions(w)
	case http.MethodPost:
		s.processPostSubscription(w, r)

	}

}

// subscriptionIDAPIHandler handles the routing logic for the '/api/subscriptions/:id' paths
func (s *Server) subscriptionIDAPIHandler(w http.ResponseWriter, r *http.Request) {
	urlID := strings.TrimPrefix(r.URL.Path, "/api/subscriptions/")
	ID, err := strconv.Atoi(urlID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		s.processDeleteSubscription(w, ID)
	}
}

func (s *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.processPostUser(w, r)

	}
}

func (s *Server) processPostUser(w http.ResponseWriter, r *http.Request) {
	var userProfile userprofile.Userprofile

	err := json.NewDecoder(r.Body).Decode(&userProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = s.dataStore.RecordUserDetails(userProfile.Name, userProfile.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// JSONContentType defines application/json
const JSONContentType = "application/json"

// processGetIndex processes the GET / request, returning the index page html
func (s *Server) processGetIndex(w http.ResponseWriter) error {

	userInfo, err := s.dataStore.GetUserDetails()
	if err != nil {
		return err
	}

	data := IndexPageData{
		Userprofile:   userInfo,
	}

	err = s.parsedIndexTemplate.Execute(w, data)

	if err != nil {
		return err
	}

	return nil
}

// processGetSubscriptions processes the GET /api/subscriptions request
// It returns the stored subscriptions as json
func (s *Server) processGetSubscriptions(w http.ResponseWriter) {
	w.Header().Set("content-type", JSONContentType)
	subscriptions, err := s.dataStore.GetSubscriptions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(subscriptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// processPostSubscription tells the SubscriptionStore to record the subscription from the post body
func (s *Server) processPostSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription subscription.Subscription
	err := json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = s.dataStore.RecordSubscription(subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
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
