export type ApiErrorResponse = {
  error: {
    code: string;
    message: string;
  };
};

export type LoginRequest = {
  email: string;
  password: string;
};

export type LoginResponse = {
  data: {
    token: string;
  };
};

export type UserRole = "customer" | "agent" | "admin";

export type Ticket = {
  id: number;
  title: string;
  description: string;
  status: "open" | "in_progress" | "resolved";
  priority: "low" | "medium" | "high";
  created_by: number;
  assigned_to: number | null;
  created_at: string;
  updated_at: string;
};

export type ListTicketsResponse = {
  data: Ticket[];
  meta: {
    page: number;
    limit: number;
    count: number;
  };
};

export type GetTicketResponse = {
  data: Ticket;
};

export type TicketReply = {
  id: number;
  ticket_id: number;
  user_id: number;
  message: string;
  created_at: string;
};

export type ListRepliesResponse = {
  data: TicketReply[];
};

export type CreateReplyRequest = {
  message: string;
};