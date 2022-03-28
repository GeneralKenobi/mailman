package remover

import (
	"context"
	"errors"
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/internal/persistence"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
)

type Repository interface {
	DeleteMailingEntryById(ctx context.Context, id int) error
}

func New(repository Repository) *Remover {
	return &Remover{repository: repository}
}

type Remover struct {
	repository Repository
}

func (remover *Remover) Remove(ctx context.Context, id int) error {
	mdctx.Infof(ctx, "Deleting mailing entry %d", id)
	err := remover.repository.DeleteMailingEntryById(ctx, id)
	if err != nil && errors.Is(err, persistence.ErrNoRows) {
		return api.StatusNotFound.WithMessageAndCause(err, "mailing entry with ID %d doesn't exist", id)
	}
	return err
}
