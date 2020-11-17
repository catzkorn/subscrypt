package server

import (
	"fmt"
	"github.com/Catzkorn/subscrypt/internal/plaid"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/shopspring/decimal"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Catzkorn/subscrypt/internal/reminder"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/Catzkorn/subscrypt/internal/userprofile"
	"github.com/shopspring/decimal"
)

// Server is the HTTP interface for subscription information
type Server struct {
	dataStore DataStore
	router    *http.ServeMux
	parsedIndexTemplate *template.Template
	transactionAPI TransactionAPI
}

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
func NewServer(dataStore DataStore, indexTemplatePath string, transactionAPI TransactionAPI) *Server {
	s := &Server{dataStore: dataStore, router: http.NewServeMux(), transactionAPI: transactionAPI}
	s.router.Handle("/", http.HandlerFunc(s.subscriptionHandler))
	s.router.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	s.router.Handle("/reminder", http.HandlerFunc(s.reminderHandler))
	s.router.Handle("/api/subscriptions/", http.HandlerFunc(s.subscriptionsAPIHandler))
	s.router.Handle("/api/transactions/", http.HandlerFunc(s.transactionAPIHandler))
	s.router.Handle("/new/user/", http.HandlerFunc(s.userHandler))

	s.parsedIndexTemplate = template.Must(template.New("index.html").ParseFiles(indexTemplatePath))

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
			return
		}
	case http.MethodPost:
		s.processPostSubscription(w, r)
	}
}

func (s *Server) reminderHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.processPostReminder(w, r)
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

	_, err := s.dataStore.RecordUserDetails(r.FormValue("username"), r.FormValue("email"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// JsonContentType defines application/json
const JsonContentType = "application/json"

// processGetSubscription processes the GET subscription request, returning the store subscriptions as json
func (s *Server) processGetSubscription(w http.ResponseWriter) error {

	subscriptions, err := s.dataStore.GetSubscriptions()
	if err != nil {
		return err
	}

	userInfo, err := s.dataStore.GetUserDetails()
	if err != nil {
		return err
	}

	data := IndexPageData{
		PageTitle:     "My Subscriptions List",
		Subscriptions: subscriptions,
		Userprofile:   userInfo,
	}

	err = s.parsedIndexTemplate.Execute(w, data)

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
		Name:    r.FormValue("name"),
		Amount:  amount,
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

// processPostReminder creates an ics file
// TODO: then emails it to the user's email
func (s *Server) processPostReminder(w http.ResponseWriter, r *http.Request) {
	var newReminder reminder.Reminder

	subscriptionID, err := strconv.Atoi(r.FormValue("subscriptionID"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reminderDate, err := time.Parse("2006-01-02", r.FormValue("reminderDate"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newReminder = reminder.Reminder{
		Email:          r.FormValue("email"),
		SubscriptionID: subscriptionID,
		ReminderDate:   reminderDate,
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(newReminder)

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
