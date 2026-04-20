package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/trannghiach/support-dashboard/backend/internal/dto"
	"github.com/trannghiach/support-dashboard/backend/internal/repository"
)

type TicketService struct {
	ticketRepo *repository.TicketRepository
	userRepo   *repository.UserRepository
}

func NewTicketService(ticketRepo *repository.TicketRepository, userRepo *repository.UserRepository) *TicketService {
	return &TicketService{
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
	}
}

func (s *TicketService) CreateTicket(
	ctx context.Context,
	userID int64,
	role string,
	req dto.CreateTicketRequest,
) (*repository.Ticket, error) {
	if role != "customer" {
		return nil, errors.New("only customers can create tickets")
	}

	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)
	priority := strings.TrimSpace(req.Priority)

	if title == "" {
		return nil, errors.New("title is required")
	}

	if description == "" {
		return nil, errors.New("description is required")
	}
	
	switch priority {
		case "low", "medium", "high":
		default:
			return nil, errors.New("priority must be one of: low, medium, high")
	}

	ticket := &repository.Ticket{
		Title:       title,
		Description: description,
		Status:      "open",
		Priority:    priority,
		CreatedBy:   userID,
	}

	if err := s.ticketRepo.CreateTicket(ctx, ticket); err != nil {
		return nil, err
	}
	
	return ticket, nil
}

func (s *TicketService) ListTickets(
	ctx context.Context,
	status string,
	priority string,
	page int,
	limit int,
) ([]repository.Ticket, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	tickets, err := s.ticketRepo.ListTickets(ctx, status, priority, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return tickets, len(tickets), nil
}

func (s *TicketService) UpdateTicketStatus(
	ctx context.Context,
	id int64,
	userID int64,
	role string,
	req dto.UpdateTicketStatusRequest,
) (*repository.Ticket, error) {
	
	// Validate enum
	switch req.Status {
		case "open", "in_progress", "resolved":
		default:
			return nil, errors.New("invalid status")
	}

	// Get ticket
	ticket, err := s.ticketRepo.GetTicketByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}

	// RBAC
	switch role {
	case "customer":
		return nil, errors.New("customers cannot update ticket status")
		
	case "agent":
		if ticket.AssignedTo == nil || *ticket.AssignedTo != userID {
			return nil, errors.New("agents can only update status of tickets assigned to them")
		}

	case "admin":
		// full access, do nothing

	default:
		return nil, errors.New("invalid role")
	}

	// Enforce transition
	if !isValidTransition(ticket.Status, req.Status) {
		return nil, errors.New("invalid status transition")
	}

	// Update status
	updated, err := s.ticketRepo.UpdateTicketStatus(ctx, id, req.Status)
	if err != nil {
		return nil, err
	}
	
	return updated, nil
}

func isValidTransition(from, to string) bool {
	switch from {
	case "open":
		return to == "in_progress"

	case "in_progress":
		return to == "resolved"

	case "resolved":
		return false

	default:
		return false
	}
}

func (s *TicketService) CreateReply(
	ctx context.Context,
	ticketID int64,
	userID int64,
	role string,
	req dto.CreateReplyRequest,
) (*repository.TicketReply, error) {
	
	message := strings.TrimSpace(req.Message)
	if message == "" {
		return nil, errors.New("message is required")
	}

	ticket, err := s.ticketRepo.GetTicketByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}

	// RBAC
	if role == "customer" && ticket.CreatedBy != userID {
		return nil, errors.New("customers can only reply to their own tickets")
	}

	// admins and agents can reply to any ticket temporarily, we can add more rules later if needed
	reply, err := s.ticketRepo.CreateReply(ctx, ticketID, userID, message)
	if err != nil {
		return nil, err
	}
	
	return reply, nil
}

func (s *TicketService) GetReplies(
	ctx context.Context,
	ticketID int64,
) ([]repository.TicketReply, error) {
	replies, err := s.ticketRepo.GetReplies(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	return replies, nil
}

func (s *TicketService) AssignTicket(
	ctx context.Context,
	id int64,
	userID int64,
	role string,
	req dto.AssignTicketRequest,
) (*repository.Ticket, error) {
	if role != "admin" {
		return nil, errors.New("only admins can assign tickets")
	}

	if req.AssignedTo <= 0 {
		return nil, errors.New("assigned_to must be a positive integer")
	}

	assignee, err := s.userRepo.GetByID(ctx, req.AssignedTo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("assignee not found")
		}
		return nil, err
	}

	if assignee.Role != "agent" {
		return nil, errors.New("can only assign tickets to agents")
	}

	ticket, err := s.ticketRepo.AssignTicket(ctx, id, req.AssignedTo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}

	return ticket, nil
}