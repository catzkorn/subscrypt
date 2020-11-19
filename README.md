# Subscrypt

## About Subscrypt

 Subscrypt was born from the frustration of forgetting when subscriptions are due to be renewed, or forgotten completely, and wasting money on unwanted services. To counter this, we made Subscrypt, a subscription manager application. We allow users to import transactions from their bank via the open banking API and have reoccurring monthly subscriptions filtered into Subscrypt. Users are able to view all their subscriptions and request a calendar reminder for 5 days prior to subscription renewal which is emailed to them. The calendar reminder is a .ics file which can be imported into any popular calendar application.

Users are also able to manually add subscriptions to Subscrypt if they do not wish to integrate their bank to the app or if a subscription is not detected by our filters.

Note: Subscrypt is a proof of concept product, using the Plaid open banking API sandbox and users bank accounts are not able to connect to the app currently.


# Go, Team!

<img src="https://imgur.com/oePX1Wo.png" width="200" height="200"> <img src="https://imgur.com/YKZhxGt.png" width="200" height="200"> <img src="https://imgur.com/zBtpZ4o.png" width="200" height="200"> 

Charlotte Brandhorst-Satzkorn ----- Farhaan Ali ------ Veronica Lee
- Gophers by [Gopherize.me](https://gopherize.me/)

Go, Team! and the Subscrypt project are the final engineering project for the [Makers Academy](https://makers.tech) Bootcamp for the September 2020 cohort. 

## Group goals

The collective aim of the group was to build a project from scratch in a new language. Go was settled on as a language to gain experience with strongly typed languages, the extensive standard library and well documented resources, and employment possibilities for Go engineers.

Our team charter can be viewed [here](https://github.com/Catzkorn/subscrypt/wiki/Team-Charter).

# Tech Stack

|      Area      | Technology  |
| :------------- | :----------: | 
|  Languages | Go, JavaScript  |
|  APIs | [Plaid](https://plaid.com/uk/), [SendGrid](https://sendgrid.com/)  |
|  Database | [Postgresql](https://www.postgresql.org/) |
| Testing & Coverage | [Go Tests](https://golang.org/pkg/testing/)  |   
|  Linting | [Golangci-lint](https://golangci-lint.run/) |  
| CI/CD   | [Github Actions, Heroku](https://github.com/Catzkorn/subscrypt/wiki/CI-and-CD) | 

# Using Subscrypt

## Initial Setup
Clone the repository: `https://github.com/Catzkorn/subscrypt.git`

For full functionality you will need to either replace code directly, or store specific information as ENV Variables. 

|      Service      | ENV Key Name  | Example |
| :------------- | :----------: | :----------: | 
|  SendGrid | SENDGRID_API_KEY  | [Documentation](https://sendgrid.com/docs/API_Reference/api_getting_started.html)
|  Database | DATABASE_CONN_STRING | "user={your_name}  host=localhost port=5432 database=subscryptdb sslmode=disable" 
| Plaid API | SECRET  |   [Documentation](https://plaid.com/docs/api/)
|  Plaid API | CLIENT_ID  |  [Documentation](https://plaid.com/docs/api/)
|  Email Address | EMAIL  |  "test@test.com"


### Database setup

[Postgresql](https://www.postgresql.org/) is required for this setup.

```Go
// Run postgres
psql -c 'CREATE DATABASE subscryptdb;'
psql -d subscryptdb -f db/migrations/database_setup.sql
```

## How to Run
```Go
$ go run ./cmd/subscrypt/main.go
```

Navigate to `http://localhost:5000/`

## How to Use

### Add Name and Email

To access the subscription manager, add an email and password (this is not [user authentication](https://github.com/Catzkorn/subscrypt/blob/main/README.md#user-authentication)).

#### Edit Name and Email

As a user, you are able to edit the name and the email the subscription reminders go to by clicking the pencil icon next to the name and email fields and replacing the contents with the updated versions. 

### Import Subscriptions 

Instructions here

Demonstration of the open banking api integration

### Add Subscription Manually

Instructions here

### Receive a Calendar Reminder

To receive a calendar reminder, click the envelope next to the desired subscription and an email will be sent to the specified user email address with a .ics file attachment that can be added to the desired calendar application. 

### Delete a Subscription

To delete a subscription, press the bin icon next to the desired subscription.

## Testing

Testing for the project is handled by the [Go standard library testing package](https://golang.org/pkg/testing/). 

### How to Test

```Go
// Run tests
$ go test ./...
```

### Test Coverage

```Go
// View test coverage per file
$ go test ./... -cover

// Generate a test coverage profile
$ go test ./... -coverprofile=coverage.out

// View coverage report (opens broswer window)
$ go tool cover -html=coverage.out
```


# Additional Information

## Future Goals

### User Authentication

The current state of the product does not allow for individual user accounts and is limited to a singular user who has the ability to edit their name and email. This feature was left out due to the time constraints of the project, and priority was given to integration of the Open Banking API and Email/Calendar invite API.

Future versions of this product would include users being able to sign up, log in , manage their details and have the ability to delete their account if they wished to. 

### Subscription Frequencies

We are currently only able to handle monthly subscriptions, which excludes being able to handle users weekly, bi-weekly, quarterly, half or yearly subscriptions. Adding this functionality would provide the user with a more complete experience. 

### Customisable Reminder Time Frame

The current configuration defaults all calendar reminders to be stored as calendar events 5 days before the subscription is due to renew. Future iterations would allow for the user to determine the time frame in which a reminder would appear before the subscription renewal is due. 

### Categories and Cost

Subscription categories such as Fitness, Food, Entertainment, etc. could be introduced to allow the user to see a breakdown of subscriptions per category. Additionally, a user would be able to see a breakdown of how much the total of their subscriptions cost for each category, as well as for all of their subscriptions. 

### Frontend Testing

At present our frontend is only manually tested due to time constraints and a late decision to move to JavaScript/JSON API. Future iterations of the project would include testing these aspects to ensure full functionality and consistent user experience.


## Attributions

### [Hero Icons](https://heroicons.dev/)

Thank you to Heroicons for providing MIT-licensed SVG icons which have been used as a part of our frontend design. 

### [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)

A massive thank you to [Chris James](https://github.com/quii) for Learn Go with Tests which was extensively used by Go, Team! at the start of the project and an excellent learning resource for diving into the Go world. 

### [GoTest](https://github.com/rakyll/gotest)

Thank you [Jaana Dogan](https://github.com/rakyll) who published GoTest, which improves the accessibility of go tests.
