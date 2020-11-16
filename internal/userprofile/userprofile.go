package userprofile

import "github.com/jackc/pgtype"

// Userprofile defines a users details
type Userprofile struct {
	ID    pgtype.UUID
	Name  string
	Email string
}
