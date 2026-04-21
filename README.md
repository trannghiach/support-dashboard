# Support Dashboard (Fullstack + AI Assistant)

Ticket management system with role-based access control and an AI assistant that summarizes conversations and suggests replies for support agents.

---

<img width="1856" height="1024" alt="Recording 2026-04-21 141444" src="https://github.com/user-attachments/assets/c73be1a4-8328-4780-ad88-7c7b8fbd018b" />

---

## Demo

- Customer: alice@test.com / alice123  
- Agent: bob@test.com / bob123  
- Admin: admin@test.com / admin123  

Frontend: http://localhost:3000  
Backend: http://localhost:8080  

---

## What this project shows

- Fullstack system design (Go + Next.js)
- RBAC enforced in backend (not just UI)
- Real workflow: ticket → assignment → replies → status transitions
- Practical AI integration (human-in-the-loop, validated output)

---

## Key Features

- Customers create and view their own tickets  
- Agents can only access tickets assigned to them  
- Admin assigns tickets and controls workflow  
- Status flow: `open → in_progress → resolved` (validated transitions)  
- Threaded replies with permission checks  

---

## Architecture

Backend follows a layered design:

- **Handler** → HTTP layer  
- **Service** → business logic + RBAC enforcement  
- **Repository** → database access (PostgreSQL)

Design decisions:
- RBAC enforced in service layer to ensure consistency across all endpoints  
- Nested resources (e.g. replies) reuse ticket-level access checks to prevent data leaks  
- Status transitions validated centrally instead of ad-hoc checks  

---

## Authentication & RBAC

- JWT-based authentication (`/auth/login`)
- Token includes `user_id` and `role`
- Middleware validates token and injects context into request

RBAC rules:
- Customer → own tickets only  
- Agent → assigned tickets only  
- Admin → full access  

Enforced in service layer before any data is returned.

---

## AI Assistant

- Summarizes ticket + conversation history  
- Generates **exactly 3 reply options**:
  - clarification  
  - troubleshooting  
  - next-step guidance  

Design:
- Prompt built on backend (not frontend)
- Structured JSON output enforced via schema

Purpose:
> Assist agents, not auto-reply (human-in-the-loop)

---

## Tech Stack

Frontend: Next.js, React, Ant Design  
Backend: Go (Gin), PostgreSQL  
AI: Gemini API  

---

## Run Locally

```bash
docker compose up -d
cd backend && go run cmd/api/main.go
cd frontend && npm install && npm run dev