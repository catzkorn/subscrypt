package integration
//
//import (
//	"github.com/Catzkorn/subscrypt/internal/database"
//	"github.com/Catzkorn/subscrypt/internal/plaid"
//	"github.com/Catzkorn/subscrypt/internal/server"
//	"github.com/sendgrid/rest"
//	"github.com/sendgrid/sendgrid-go/helpers/mail"
//	"io/ioutil"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)
//
//func TestCallingAPIAndSeeingSubscriptions(t *testing.T) {
//	store := database.NewInMemorySubscriptionStore()
//
//	api := &plaid.PlaidAPI{}
//
//	testServer := server.NewServer(store, indexTemplatePath, &StubMailer{}, api)
//
//
//	getRequest := newTransactionAPIRequest()
//	response := httptest.NewRecorder()
//	testServer.ServeHTTP(response, getRequest)
//
//	getRequest = newGetSubscriptionRequest()
//	response = httptest.NewRecorder()
//	testServer.ServeHTTP(response, getRequest)
//
//	body, err := ioutil.ReadAll(response.Body)
//
//	if err != nil {
//		t.Errorf("unexpected error: %w", err)
//	}
//
//	bodyString := string(body)
//	got := bodyString
//
//	res := strings.Contains(got, "Touchstone Climbing")
//
//	if res != true {
//		t.Errorf("webpage did not contain subscription of name %v", "Touchstone Climbing")
//	}
//}
//
//type StubMailer struct {
//	sentEmail *mail.SGMailV3
//}
//
//func (s *StubMailer) Send(email *mail.SGMailV3) (*rest.Response, error) {
//	s.sentEmail = email
//	return &rest.Response{StatusCode: http.StatusAccepted}, nil
//}
//
//const indexTemplatePath = "../../web/index.html"
//
//func newTransactionAPIRequest() *http.Request {
//	req, _ := http.NewRequest(http.MethodGet, "/api/transactions/", nil)
//	return req
//}
//
//func newGetSubscriptionRequest() *http.Request {
//	req, _ := http.NewRequest(http.MethodGet, "/", nil)
//	return req
//}
