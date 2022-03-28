package repository

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/persistence/model"
	"time"
)

func (repository *Repository) FindMailingEntriesByMailingId(ctx context.Context, mailingId int) ([]model.MailingEntry, error) {
	return selectingAll(ctx, "find mailing entries by mailing ID", repository.sql, mailingEntryRowScanSupplier,
		"SELECT id, customer_id, mailing_id, title, content, insert_time FROM mailmandb.mailing_entry WHERE mailing_id = $1", mailingId)
}

func (repository *Repository) FindMailingEntriesOlderThan(ctx context.Context, olderThan time.Duration) ([]model.MailingEntry, error) {
	return selectingAll(ctx, "find mailing entries older than", repository.sql, mailingEntryRowScanSupplier,
		"SELECT id, customer_id, mailing_id, title, content, insert_time FROM mailmandb.mailing_entry WHERE insert_time < $1", time.Now().Add(-olderThan))
}

func (repository *Repository) FindMailingEntriesByMailingIdOlderThan(ctx context.Context, mailingId int, olderThan time.Duration) ([]model.MailingEntry, error) {
	return selectingAll(ctx, "find mailing entries by mailing ID older than", repository.sql, mailingEntryRowScanSupplier,
		"SELECT id, customer_id, mailing_id, title, content, insert_time FROM mailmandb.mailing_entry WHERE mailing_entry.mailing_id = $1 AND insert_time < $2", mailingId, time.Now().Add(-olderThan))
}

func (repository *Repository) FindMailingEntriesByCustomerId(ctx context.Context, customerId int) ([]model.MailingEntry, error) {
	return selectingAll(ctx, "find mailing by customer ID", repository.sql, mailingEntryRowScanSupplier,
		"SELECT id, customer_id, mailing_id, title, content, insert_time FROM mailmandb.mailing_entry WHERE mailing_entry.customer_id = $1", customerId)
}

func (repository *Repository) FindMailingEntriesByCustomerIdMailingIdTitleContentInsertTime(
	ctx context.Context, customerId, mailingId int, title, content string, insertTime time.Time) ([]model.MailingEntry, error) {

	return selectingAll(ctx, "find mailing entries by customer ID, mailing ID, title, content and insert time",
		repository.sql, mailingEntryRowScanSupplier,
		"SELECT id, customer_id, mailing_id, title, content, insert_time FROM mailmandb.mailing_entry WHERE customer_id = $1 AND mailing_id = $2 AND title = $3 AND content = $4 AND insert_time = $5",
		customerId, mailingId, title, content, insertTime)
}

func (repository *Repository) InsertMailingEntry(ctx context.Context, mailingEntry model.MailingEntry) (model.MailingEntry, error) {
	return selectingOne(ctx, "insert mailing entry", repository.sql, mailingEntryRowScanSupplier,
		"INSERT INTO mailmandb.mailing_entry(customer_id, mailing_id, title, content, insert_time) VALUES ($1, $2, $3, $4, $5) RETURNING id, customer_id, mailing_id, title, content, insert_time",
		mailingEntry.CustomerId, mailingEntry.MailingId, mailingEntry.Title, mailingEntry.Content, mailingEntry.InsertTime)
}

func (repository *Repository) DeleteMailingEntryById(ctx context.Context, id int) error {
	return affectingOne(ctx, "delete mailing entry by ID", repository.sql,
		"DELETE FROM mailmandb.mailing_entry WHERE id = $1", id)
}

func mailingEntryRowScanSupplier() (*model.MailingEntry, []any) {
	var mailingEntry model.MailingEntry
	return &mailingEntry, []any{
		&mailingEntry.Id,
		&mailingEntry.CustomerId,
		&mailingEntry.MailingId,
		&mailingEntry.Title,
		&mailingEntry.Content,
		&mailingEntry.InsertTime,
	}
}
