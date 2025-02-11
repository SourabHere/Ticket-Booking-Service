package db

import (
	"context"
	"fmt"
)

type CancelTicketParams struct {
	UserID   int32 `json:"user_id"`
	TicketID int32 `json:"ticket_id"`
	SeatID   int32 `json:"seat_id"`
}

func (store *Store) CancelTicketTx(ctx context.Context, arg CancelTicketParams) error {
	err := store.execTx(ctx, func(q *Queries) error {

		err := q.UpdateTicketStatus(ctx, UpdateTicketStatusParams{
			ID:     arg.TicketID,
			Status: "canceled",
		})
		if err != nil {
			return fmt.Errorf("failed to update ticket status: %w", err)
		}

		err = q.UpdateSeatReservationStatus(ctx, UpdateSeatReservationStatusParams{
			BusSeatID: arg.SeatID,
			Status:    "canceled",
			UserID:    arg.UserID,
		})
		if err != nil {
			return fmt.Errorf("failed to update seat status: %w", err)
		}

		return nil
	})

	return err
}
