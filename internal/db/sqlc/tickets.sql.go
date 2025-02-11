// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: tickets.sql

package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPenalty = `-- name: CreatePenalty :one
INSERT INTO penalties (bus_id, actual_hours_before, hours_before, percent, custom_text)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, bus_id, actual_hours_before, hours_before, percent, custom_text
`

type CreatePenaltyParams struct {
	BusID             int32         `json:"bus_id"`
	ActualHoursBefore pgtype.Float8 `json:"actual_hours_before"`
	HoursBefore       pgtype.Float8 `json:"hours_before"`
	Percent           int32         `json:"percent"`
	CustomText        pgtype.Text   `json:"custom_text"`
}

func (q *Queries) CreatePenalty(ctx context.Context, arg CreatePenaltyParams) (Penalty, error) {
	row := q.db.QueryRow(ctx, createPenalty,
		arg.BusID,
		arg.ActualHoursBefore,
		arg.HoursBefore,
		arg.Percent,
		arg.CustomText,
	)
	var i Penalty
	err := row.Scan(
		&i.ID,
		&i.BusID,
		&i.ActualHoursBefore,
		&i.HoursBefore,
		&i.Percent,
		&i.CustomText,
	)
	return i, err
}

const deleteTicket = `-- name: DeleteTicket :exec
DELETE FROM tickets
WHERE id = $1
`

func (q *Queries) DeleteTicket(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteTicket, id)
	return err
}

const getBusPenalties = `-- name: GetBusPenalties :many
SELECT id, bus_id, actual_hours_before, hours_before, percent, custom_text
FROM penalties
WHERE bus_id = $1
`

func (q *Queries) GetBusPenalties(ctx context.Context, busID int32) ([]Penalty, error) {
	rows, err := q.db.Query(ctx, getBusPenalties, busID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Penalty{}
	for rows.Next() {
		var i Penalty
		if err := rows.Scan(
			&i.ID,
			&i.BusID,
			&i.ActualHoursBefore,
			&i.HoursBefore,
			&i.Percent,
			&i.CustomText,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReservedTicketsCount = `-- name: GetReservedTicketsCount :one
SELECT COUNT(*)
FROM tickets
WHERE bus_id = $1
`

func (q *Queries) GetReservedTicketsCount(ctx context.Context, busID int32) (int64, error) {
	row := q.db.QueryRow(ctx, getReservedTicketsCount, busID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTicketByID = `-- name: GetTicketByID :one
SELECT id, user_id, bus_id, reserved_at, status, seat_reservation_id
FROM tickets
WHERE id = $1
`

type GetTicketByIDRow struct {
	ID                int32              `json:"id"`
	UserID            int32              `json:"user_id"`
	BusID             int32              `json:"bus_id"`
	ReservedAt        pgtype.Timestamptz `json:"reserved_at"`
	Status            string             `json:"status"`
	SeatReservationID int32              `json:"seat_reservation_id"`
}

func (q *Queries) GetTicketByID(ctx context.Context, id int32) (GetTicketByIDRow, error) {
	row := q.db.QueryRow(ctx, getTicketByID, id)
	var i GetTicketByIDRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.BusID,
		&i.ReservedAt,
		&i.Status,
		&i.SeatReservationID,
	)
	return i, err
}

const getUserTickets = `-- name: GetUserTickets :many
SELECT t.id, b.route_id, b.departure_time, b.arrival_time, b.capacity, b.price, b.bus_type, b.corporation, b.super_corporation, b.service_number, b.is_vip
FROM tickets t
JOIN buses b ON t.bus_id = b.id
WHERE t.user_id = $1
`

type GetUserTicketsRow struct {
	ID               int32       `json:"id"`
	RouteID          int32       `json:"route_id"`
	DepartureTime    time.Time   `json:"departure_time"`
	ArrivalTime      time.Time   `json:"arrival_time"`
	Capacity         int32       `json:"capacity"`
	Price            int32       `json:"price"`
	BusType          string      `json:"bus_type"`
	Corporation      pgtype.Text `json:"corporation"`
	SuperCorporation pgtype.Text `json:"super_corporation"`
	ServiceNumber    pgtype.Text `json:"service_number"`
	IsVip            pgtype.Bool `json:"is_vip"`
}

func (q *Queries) GetUserTickets(ctx context.Context, userID int32) ([]GetUserTicketsRow, error) {
	rows, err := q.db.Query(ctx, getUserTickets, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserTicketsRow{}
	for rows.Next() {
		var i GetUserTicketsRow
		if err := rows.Scan(
			&i.ID,
			&i.RouteID,
			&i.DepartureTime,
			&i.ArrivalTime,
			&i.Capacity,
			&i.Price,
			&i.BusType,
			&i.Corporation,
			&i.SuperCorporation,
			&i.ServiceNumber,
			&i.IsVip,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUserTickets = `-- name: ListUserTickets :many
SELECT 
    t.id AS ticket_id,
    t.bus_id,
    sr.bus_seat_id AS seat_id,
    t.reserved_at,
    b.departure_time,
    b.arrival_time,
    b.price,
    s.seat_number,
    sr.status AS reservation_status
FROM 
    tickets t
JOIN 
    buses b ON t.bus_id = b.id
JOIN 
    seat_reservations sr ON t.seat_reservation_id = sr.id
JOIN 
    bus_seats s ON sr.bus_seat_id = s.id
WHERE 
    t.user_id = $1
ORDER BY 
    t.reserved_at DESC
`

type ListUserTicketsRow struct {
	TicketID          int32              `json:"ticket_id"`
	BusID             int32              `json:"bus_id"`
	SeatID            int32              `json:"seat_id"`
	ReservedAt        pgtype.Timestamptz `json:"reserved_at"`
	DepartureTime     time.Time          `json:"departure_time"`
	ArrivalTime       time.Time          `json:"arrival_time"`
	Price             int32              `json:"price"`
	SeatNumber        int32              `json:"seat_number"`
	ReservationStatus string             `json:"reservation_status"`
}

func (q *Queries) ListUserTickets(ctx context.Context, userID int32) ([]ListUserTicketsRow, error) {
	rows, err := q.db.Query(ctx, listUserTickets, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListUserTicketsRow{}
	for rows.Next() {
		var i ListUserTicketsRow
		if err := rows.Scan(
			&i.TicketID,
			&i.BusID,
			&i.SeatID,
			&i.ReservedAt,
			&i.DepartureTime,
			&i.ArrivalTime,
			&i.Price,
			&i.SeatNumber,
			&i.ReservationStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const purchaseTicket = `-- name: PurchaseTicket :one
INSERT INTO tickets (user_id, bus_id, seat_reservation_id, status, purchased_at)
VALUES ($1, $2, $3, 'purchased', NOW())
RETURNING id, user_id, bus_id, seat_reservation_id, status, purchased_at
`

type PurchaseTicketParams struct {
	UserID            int32 `json:"user_id"`
	BusID             int32 `json:"bus_id"`
	SeatReservationID int32 `json:"seat_reservation_id"`
}

type PurchaseTicketRow struct {
	ID                int32              `json:"id"`
	UserID            int32              `json:"user_id"`
	BusID             int32              `json:"bus_id"`
	SeatReservationID int32              `json:"seat_reservation_id"`
	Status            string             `json:"status"`
	PurchasedAt       pgtype.Timestamptz `json:"purchased_at"`
}

func (q *Queries) PurchaseTicket(ctx context.Context, arg PurchaseTicketParams) (PurchaseTicketRow, error) {
	row := q.db.QueryRow(ctx, purchaseTicket, arg.UserID, arg.BusID, arg.SeatReservationID)
	var i PurchaseTicketRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.BusID,
		&i.SeatReservationID,
		&i.Status,
		&i.PurchasedAt,
	)
	return i, err
}

const reserveTicket = `-- name: ReserveTicket :one
INSERT INTO tickets (user_id, bus_id, seat_reservation_id, status, reserved_at)
VALUES ($1, $2, $3, 'reserved', NOW())
RETURNING id, user_id, bus_id, seat_reservation_id, status, reserved_at
`

type ReserveTicketParams struct {
	UserID            int32 `json:"user_id"`
	BusID             int32 `json:"bus_id"`
	SeatReservationID int32 `json:"seat_reservation_id"`
}

type ReserveTicketRow struct {
	ID                int32              `json:"id"`
	UserID            int32              `json:"user_id"`
	BusID             int32              `json:"bus_id"`
	SeatReservationID int32              `json:"seat_reservation_id"`
	Status            string             `json:"status"`
	ReservedAt        pgtype.Timestamptz `json:"reserved_at"`
}

func (q *Queries) ReserveTicket(ctx context.Context, arg ReserveTicketParams) (ReserveTicketRow, error) {
	row := q.db.QueryRow(ctx, reserveTicket, arg.UserID, arg.BusID, arg.SeatReservationID)
	var i ReserveTicketRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.BusID,
		&i.SeatReservationID,
		&i.Status,
		&i.ReservedAt,
	)
	return i, err
}

const updateTicketStatus = `-- name: UpdateTicketStatus :exec
UPDATE tickets
SET status = $2
WHERE id = $1
`

type UpdateTicketStatusParams struct {
	ID     int32  `json:"id"`
	Status string `json:"status"`
}

func (q *Queries) UpdateTicketStatus(ctx context.Context, arg UpdateTicketStatusParams) error {
	_, err := q.db.Exec(ctx, updateTicketStatus, arg.ID, arg.Status)
	return err
}
