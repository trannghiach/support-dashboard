package dto

type CreateTicketRequest struct {
	Title	    string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
}

type UpdateTicketStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateReplyRequest struct {
	Message string `json:"message" binding:"required"`
}

type AssignTicketRequest struct {
	AssignedTo int64 `json:"assigned_to"`
}