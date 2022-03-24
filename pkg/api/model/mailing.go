package model

import "time"

type MailingEntryDto struct {
	MailingId  int       `json:"mailing_id"`
	Email      string    `json:"email,omitempty"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content,omitempty"`
	InsertTime time.Time `json:"insert_time,omitempty"`
}

type MailingRequestDto struct {
	MailingId int `json:"mailing_id"`
}
