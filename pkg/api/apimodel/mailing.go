package apimodel

import "time"

// MailingEntry defines a mailing entry to create.
type MailingEntry struct {
	MailingId  int       `json:"mailing_id" validate:"required"`            // ID of the mailing list
	Email      string    `json:"email,omitempty" validate:"required,email"` // Email address of the recipient
	Title      string    `json:"title,omitempty" validate:"required"`       // Message title
	Content    string    `json:"content,omitempty" validate:"required"`     // Message con
	InsertTime time.Time `json:"insert_time,omitempty" validate:"required"` // Message creation time
}

// MailingEntryCreated is returned after successfully creating a mailing entry from a MailingEntry.
type MailingEntryCreated struct {
	Id int `json:"id" validation:"required"`
}

// MailingRequest is a request to send mailing entries from a given mailing list.
type MailingRequest struct {
	MailingId int `json:"mailing_id" validate:"required"` // ID of the mailing list
}
