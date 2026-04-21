package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Ticket struct {
	ID		   	int64     `json:"id"`
	Title	  	string    `json:"title"`
	Description string    `json:"description"`
	Status    	string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedBy   int64     `json:"created_by"`
	AssignedTo  *int64    `json:"assigned_to"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TicketReply struct {
	ID        int64     `json:"id"`
	TicketID  int64     `json:"ticket_id"`
	UserID    int64     `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type TicketRepository struct {
	db *pgxpool.Pool
}

func NewTicketRepository(db *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) CreateTicket(ctx context.Context, t *Ticket) error {
	query := `
		INSERT INTO tickets (title, description, status, priority, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, assigned_to, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		t.Title,
		t.Description,
		t.Status,
		t.Priority,
		t.CreatedBy,
	).Scan(
		&t.ID,
		&t.AssignedTo,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
}

func (r *TicketRepository) ListTickets(
	ctx context.Context,
	userID int64,
	role string,
	status string,
	priority string,
	limit int,
	offset int,
) ([]Ticket, error) {

	query := `
			SELECT id, title, description, status, priority,
				created_by, assigned_to, created_at, updated_at
			FROM tickets
			WHERE 
				($1 = 'admin') OR
				($1 = 'agent' AND (assigned_to = $2 OR created_by = $2)) OR
				($1 = 'customer' AND created_by = $2)
	`

	args := []any{role, userID}
	argIdx := 3

	if status != "" {
		query += " AND status = $" + strconv.Itoa(argIdx)
		args = append(args, status)
		argIdx++
	}

	if priority != "" {
		query += " AND priority = $" + strconv.Itoa(argIdx)
		args = append(args, priority)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	query += " LIMIT $" + strconv.Itoa(argIdx)
	args = append(args, limit)
	argIdx++

	query += " OFFSET $" + strconv.Itoa(argIdx)
	args = append(args, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []Ticket

	for rows.Next() {
		var t Ticket
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.CreatedBy,
			&t.AssignedTo,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

func (r *TicketRepository) UpdateTicketStatus(
	ctx context.Context,
	id int64,
	status string,
) (*Ticket, error) {

	query := `
		UPDATE tickets
		SET status = $1, 
			updated_at = NOW()
		WHERE id = $2
		RETURNING id, title, description, status, priority, 
				  created_by, assigned_to, created_at, updated_at
	`

	var t Ticket

	err := r.db.QueryRow(ctx, query, status, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.CreatedBy,
		&t.AssignedTo,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TicketRepository) GetStatusByID(
	ctx context.Context,
	id int64,
) (string, error) {
	var status string

	err := r.db.QueryRow(
		ctx,
		"SELECT status FROM tickets WHERE id = $1",
		id,
	).Scan(&status)

	if err != nil {
		return "", err
	}

	return status, nil
}

func (r *TicketRepository) CreateReply(
	ctx context.Context,
	ticketID int64,
	userID int64,
	message string,
) (*TicketReply, error) {

	query := `
		INSERT INTO ticket_replies (ticket_id, user_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, ticket_id, user_id, message, created_at
	`

	var reply TicketReply

	err := r.db.QueryRow(ctx, query, ticketID, userID, message).Scan(
		&reply.ID,
		&reply.TicketID,
		&reply.UserID,
		&reply.Message,
		&reply.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return &reply, nil
}

func (r *TicketRepository) GetReplies(
	ctx context.Context,
	ticketID int64,
) ([]TicketReply, error) {
	query := `
		SELECT id, ticket_id, user_id, message, created_at
		FROM ticket_replies
		WHERE ticket_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []TicketReply

	for rows.Next() {
		var r TicketReply
		if err := rows.Scan(
			&r.ID,
			&r.TicketID,
			&r.UserID,
			&r.Message,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}
		replies = append(replies, r)
	}

	return replies, rows.Err()
}

func (r *TicketRepository) AssignTicket(
	ctx context.Context,
	id int64,
	assignedTo int64,
) (*Ticket, error) {
	query := `
		UPDATE tickets
		SET assigned_to = $1,
			updated_at = NOW()
		WHERE id = $2
		RETURNING id, title, description, status, priority, 
				  created_by, assigned_to, created_at, updated_at
	`

	var t Ticket

	err := r.db.QueryRow(ctx, query, assignedTo, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.CreatedBy,
		&t.AssignedTo,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TicketRepository) GetTicketByID(
	ctx context.Context,
	id int64,
) (*Ticket, error) {
	query := `
		SELECT id, title, description, status, priority, 
			   created_by, assigned_to, created_at, updated_at
		FROM tickets
		WHERE id = $1
	`

	var t Ticket
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.CreatedBy,
		&t.AssignedTo,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	
	return &t, nil
}