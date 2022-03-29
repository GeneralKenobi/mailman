package creator

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/internal/persistence"
	"github.com/GeneralKenobi/mailman/internal/persistence/model"
	apimodel "github.com/GeneralKenobi/mailman/pkg/api/model"
	"testing"
	"time"
)

// Should use the existing customer.
func TestCreateFromDtoCustomerAlreadyExists(t *testing.T) {
	expected := model.MailingEntry{
		Id:         45,
		CustomerId: 33,
		MailingId:  17,
		Title:      "test email",
		Content:    "test content",
		InsertTime: time.Now(),
	}
	input := apimodel.MailingEntryDto{
		MailingId:  expected.MailingId,
		Email:      "test@test.com",
		Title:      expected.Title,
		Content:    expected.Content,
		InsertTime: expected.InsertTime,
	}

	customerCreator := customerCreatorMock{
		createFromEmail: func(ctx context.Context, email string) (model.Customer, error) {
			t.Fatalf("shouldn't be called - customer already exists")
			return model.Customer{}, nil
		},
	}
	repository := repositoryMock{
		findCustomerByEmail: func(ctx context.Context, email string) (model.Customer, error) {
			if email != input.Email {
				t.Fatalf("expected email %q, got %q", input.Email, email)
			}

			customer := model.Customer{
				Id:    expected.CustomerId,
				Email: email,
			}
			return customer, nil
		},
		findMailingEntriesByCustomerIdMailingIdTitleContentInsertTime: func(ctx context.Context, customerId, mailingId int, title, content string, insertTime time.Time) ([]model.MailingEntry, error) {
			queried := model.MailingEntry{
				Id:         expected.Id, // For comparison in assertion below
				CustomerId: customerId,
				MailingId:  mailingId,
				Title:      title,
				Content:    content,
				InsertTime: insertTime,
			}
			if queried != expected {
				t.Fatalf("Expected query for %#v\n, got query for %#v", expected, queried)
			}
			return nil, nil
		},
		insertMailingEntry: func(ctx context.Context, mailingEntry model.MailingEntry) (model.MailingEntry, error) {
			mailingEntry.Id = expected.Id
			return mailingEntry, nil
		},
	}

	testObj := New(repository, customerCreator)
	mailingEntry, err := testObj.CreateFromDto(context.TODO(), input)

	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if mailingEntry != expected {
		t.Errorf("Expected %#v\n, got %#v", expected, mailingEntry)
	}
}

// Should create the customer.
func TestCreateFromDtoCustomerDoesNotExistYet(t *testing.T) {
	expected := model.MailingEntry{
		Id:         45,
		CustomerId: 33,
		MailingId:  17,
		Title:      "test email",
		Content:    "test content",
		InsertTime: time.Now(),
	}
	input := apimodel.MailingEntryDto{
		MailingId:  expected.MailingId,
		Email:      "test@test.com",
		Title:      expected.Title,
		Content:    expected.Content,
		InsertTime: expected.InsertTime,
	}
	customerCreator := customerCreatorMock{
		createFromEmail: func(ctx context.Context, email string) (model.Customer, error) {
			if email != input.Email {
				t.Fatalf("expected email %q, got %q", input.Email, email)
			}

			customer := model.Customer{
				Id:    expected.CustomerId,
				Email: email,
			}
			return customer, nil
		},
	}
	repository := repositoryMock{
		findCustomerByEmail: func(ctx context.Context, email string) (model.Customer, error) {
			if email != input.Email {
				t.Fatalf("expected email %q, got %q", input.Email, email)
			}

			return model.Customer{}, persistence.ErrNoRows
		},
		findMailingEntriesByCustomerIdMailingIdTitleContentInsertTime: func(ctx context.Context, customerId, mailingId int, title, content string, insertTime time.Time) ([]model.MailingEntry, error) {
			queried := model.MailingEntry{
				Id:         expected.Id, // For comparison in assertion below
				CustomerId: customerId,
				MailingId:  mailingId,
				Title:      title,
				Content:    content,
				InsertTime: insertTime,
			}
			if queried != expected {
				t.Fatalf("Expected query for %#v\n, got query for %#v", expected, queried)
			}
			return nil, nil
		},
		insertMailingEntry: func(ctx context.Context, mailingEntry model.MailingEntry) (model.MailingEntry, error) {
			mailingEntry.Id = expected.Id
			return mailingEntry, nil
		},
	}

	testObj := New(repository, customerCreator)
	mailingEntry, err := testObj.CreateFromDto(context.TODO(), input)

	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if mailingEntry != expected {
		t.Errorf("Expected %#v\n, got %#v", expected, mailingEntry)
	}
}

// Should return an error because duplicates are not allowed.
func TestCreateFromDtoMailingEntryAlreadyExists(t *testing.T) {
	expected := model.MailingEntry{
		Id:         45,
		CustomerId: 33,
		MailingId:  17,
		Title:      "test email",
		Content:    "test content",
		InsertTime: time.Now(),
	}
	input := apimodel.MailingEntryDto{
		MailingId:  expected.MailingId,
		Email:      "test@test.com",
		Title:      expected.Title,
		Content:    expected.Content,
		InsertTime: expected.InsertTime,
	}
	customerCreator := customerCreatorMock{
		createFromEmail: func(ctx context.Context, email string) (model.Customer, error) {
			t.Fatalf("shouldn't be called")
			return model.Customer{}, nil
		},
	}
	repository := repositoryMock{
		findCustomerByEmail: func(ctx context.Context, email string) (model.Customer, error) {
			if email != input.Email {
				t.Fatalf("expected email %q, got %q", input.Email, email)
			}

			customer := model.Customer{
				Id:    expected.CustomerId,
				Email: email,
			}
			return customer, nil
		},
		findMailingEntriesByCustomerIdMailingIdTitleContentInsertTime: func(ctx context.Context, customerId, mailingId int, title, content string, insertTime time.Time) ([]model.MailingEntry, error) {
			queried := model.MailingEntry{
				Id:         expected.Id, // For comparison in assertion below
				CustomerId: customerId,
				MailingId:  mailingId,
				Title:      title,
				Content:    content,
				InsertTime: insertTime,
			}
			if queried != expected {
				t.Fatalf("Expected query for %#v\n, got query for %#v", expected, queried)
			}
			return []model.MailingEntry{expected}, nil
		},
		insertMailingEntry: func(ctx context.Context, mailingEntry model.MailingEntry) (model.MailingEntry, error) {
			t.Fatalf("shouldn't be called - should return an error because the mailing entry already exists")
			return model.MailingEntry{}, nil
		},
	}

	testObj := New(repository, customerCreator)
	_, err := testObj.CreateFromDto(context.TODO(), input)

	if err == nil {
		t.Fatalf("Expected error but got none")
	}
	statusErr, ok := err.(api.StatusError)
	if !ok {
		t.Errorf("Expected a StatusError error but got %v (%T)", err, err)
	}
	if statusErr.Status() != api.StatusBadInput {
		t.Errorf("Expected bad input status but got %v in %v", statusErr.Status(), statusErr)
	}
}

type customerCreatorMock struct {
	createFromEmail func(ctx context.Context, email string) (model.Customer, error)
}

func (mock customerCreatorMock) CreateFromEmail(ctx context.Context, email string) (model.Customer, error) {
	return mock.createFromEmail(ctx, email)
}

type repositoryMock struct {
	findCustomerByEmail                                           func(ctx context.Context, email string) (model.Customer, error)
	findMailingEntriesByCustomerIdMailingIdTitleContentInsertTime func(ctx context.Context, customerId, mailingId int, title, content string, insertTime time.Time) ([]model.MailingEntry, error)
	insertMailingEntry                                            func(ctx context.Context, mailingEntry model.MailingEntry) (model.MailingEntry, error)
}

func (mock repositoryMock) FindCustomerByEmail(ctx context.Context, email string) (model.Customer, error) {
	return mock.findCustomerByEmail(ctx, email)
}

func (mock repositoryMock) FindMailingEntriesByCustomerIdMailingIdTitleContentInsertTime(ctx context.Context, customerId, mailingId int, title, content string, insertTime time.Time) ([]model.MailingEntry, error) {
	return mock.findMailingEntriesByCustomerIdMailingIdTitleContentInsertTime(ctx, customerId, mailingId, title, content, insertTime)
}

func (mock repositoryMock) InsertMailingEntry(ctx context.Context, mailingEntry model.MailingEntry) (model.MailingEntry, error) {
	return mock.insertMailingEntry(ctx, mailingEntry)
}
