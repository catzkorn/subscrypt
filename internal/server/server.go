package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Catzkorn/subscrypt/internal/calendar"
	"github.com/Catzkorn/subscrypt/internal/email"
	"github.com/Catzkorn/subscrypt/internal/reminder"
	"github.com/Catzkorn/subscrypt/internal/subscription"
	"github.com/Catzkorn/subscrypt/internal/userprofile"
	"github.com/shopspring/decimal"
)

// Server is the HTTP interface for subscription information
type Server struct {
	dataStore           DataStore
	router              *http.ServeMux
	parsedIndexTemplate *template.Template
	mailer              email.Mailer
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
func NewServer(dataStore DataStore, indexTemplatePath string, mailer email.Mailer) *Server {
	s := &Server{dataStore: dataStore, router: http.NewServeMux()}
	s.router.Handle("/", http.HandlerFunc(s.subscriptionHandler))
	s.router.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	s.router.Handle("/api/reminders", http.HandlerFunc(s.reminderHandler))
	s.router.Handle("/api/subscriptions/", http.HandlerFunc(s.subscriptionsAPIHandler))
	s.router.Handle("/new/user/", http.HandlerFunc(s.userHandler))

	s.parsedIndexTemplate = template.Must(template.New("index.html").ParseFiles(indexTemplatePath))

	s.mailer = mailer

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

// processPostReminder creates an ics file
// TODO: then emails it to the user's email
func (s *Server) processPostReminder(w http.ResponseWriter, r *http.Request) {
	var newReminder reminder.Reminder
	var newSubscription subscription.Subscription
	var userInformation userprofile.Userprofile

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

	newSubscription.Name = subscription.Name
	newSubscription.Amount = subscription.Amount
	newSubscription.DateDue = subscription.DateDue

	user, err := s.dataStore.GetUserDetails()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userInformation.Name = user.Name
	userInformation.Email = user.Email

	newReminder = reminder.Reminder{
		Email:          userInformation.Email,
		SubscriptionID: newSubscription.ID,
		ReminderDate:   newSubscription.DateDue.AddDate(0, 0, -5),
	}

	cal := calendar.CreateReminderInvite(newSubscription, newReminder)

	email.SendEmail(newReminder, userInformation, cal, s.mailer, s.dataStore)

	w.WriteHeader(http.StatusOK)

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
