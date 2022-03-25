package model

import (
	"time"
)

type MailingEntry struct {
	Id         int // Primary key
	CustomerId int // Maps many-to-one relationship to Customer.Id
	MailingId  int
	Title      string
	Content    string
	InsertTime time.Time
}
