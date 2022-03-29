package model

import "time"

type MailingEntryDto struct {
	MailingId  int       `json:"mailing_id" validate:"required"`
	Email      string    `json:"email,omitempty" validate:"required,email"`
	Title      string    `json:"title,omitempty" validate:"required"`
	Content    string    `json:"content,omitempty" validate:"required"`
	InsertTime time.Time `json:"insert_time,omitempty" validate:"required"`
}

type MailingRequestDto struct {
	MailingId int `json:"mailing_id" validate:"required"`
}
