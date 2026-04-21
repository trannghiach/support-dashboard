package ai

import (
	"fmt"
	"strings"

	"github.com/trannghiach/support-dashboard/backend/internal/repository"
)

func BuildTicketAssistPrompt(ticket *repository.Ticket, replies []repository.TicketReply) string {
	var b strings.Builder

	b.WriteString("You are an AI assistant helping a customer support agent.\n")
	b.WriteString("Your task:\n")
	b.WriteString("1. Summarize the ticket briefly and accurately.\n")
	b.WriteString("2. Generate exactly 3 possible next replies for the support agent.\n\n")

	b.WriteString("Rules:\n")
	b.WriteString("- Be concise, professional, and helpful.\n")
	b.WriteString("- Use only the ticket and reply history provided.\n")
	b.WriteString("- Do not invent facts, system actions, investigation results, or internal actions.\n")
	b.WriteString("- Do not claim that any action has already been taken unless it is explicitly stated in the provided context.\n")
	b.WriteString("- Do not say you escalated, triggered, verified, checked, confirmed, refunded, synced, reset, fixed, or changed anything unless that exact action appears in the ticket or reply history.\n")
	b.WriteString("- If information is missing, ask for clarification instead of inventing actions.\n")
	b.WriteString("- Suggested replies must be safe for an agent to send immediately without implying unperformed actions.\n")
	b.WriteString("- The 3 replies must be meaningfully different from each other.\n")
	b.WriteString("- Reply 1 should focus on clarification.\n")
	b.WriteString("- Reply 2 should focus on troubleshooting or next diagnostic steps.\n")
	b.WriteString("- Reply 3 should focus on a careful next-step recommendation, and may mention possible escalation only as a future possibility, not as something already done.\n")
	b.WriteString("- Avoid phrases like: 'I have escalated', 'I have triggered', 'I have verified', 'I checked our system', 'I confirmed', unless explicitly supported by the conversation history.\n")
	b.WriteString("- Return valid JSON only.\n\n")

	b.WriteString("Ticket:\n")
	b.WriteString(fmt.Sprintf("Title: %s\n", ticket.Title))
	b.WriteString(fmt.Sprintf("Description: %s\n", ticket.Description))
	b.WriteString(fmt.Sprintf("Status: %s\n", ticket.Status))
	b.WriteString(fmt.Sprintf("Priority: %s\n", ticket.Priority))

	if ticket.AssignedTo != nil {
		b.WriteString(fmt.Sprintf("AssignedTo: %d\n", *ticket.AssignedTo))
	} else {
		b.WriteString("AssignedTo: null\n")
	}

	b.WriteString("\nConversation history:\n")
	if len(replies) == 0 {
		b.WriteString("No replies yet.\n")
	} else {
		for i, reply := range replies {
			b.WriteString(fmt.Sprintf(
				"%d. User #%d at %s: %s\n",
				i+1,
				reply.UserID,
				reply.CreatedAt.Format("2006-01-02 15:04:05"),
				reply.Message,
			))
		}
	}

	b.WriteString("\nReturn JSON with this exact shape:\n")
	b.WriteString("{\n")
	b.WriteString(`  "summary": "string",` + "\n")
	b.WriteString(`  "suggested_replies": ["string", "string", "string"]` + "\n")
	b.WriteString("}\n\n")

	b.WriteString("Requirements:\n")
	b.WriteString("- summary: 2 to 4 sentences\n")
	b.WriteString("- suggested_replies: exactly 3 different reply drafts\n")
	b.WriteString("- no markdown\n")
	b.WriteString("- no extra explanation outside JSON\n")

	return b.String()
}