package messages

// =====================================
// Modal messages
// =====================================

type ModalClosedMsg struct{}

type ShowModalMsg struct {
	ModalType string // "workflow", "command", etc.
}

type ItemDeletedMsg struct {
	ID   string
	Type string
}

// =====================================
// Workflow messages
// =====================================

type WorkflowCreatedMsg struct {
	WorkflowID string
}

type WorkflowUpdatedMsg struct {
	WorkflowID string
}

type WorkflowDeletedMsg struct {
	WorkflowID string
}

// =====================================
// Command messages
// =====================================

type CommandCreatedMsg struct {
	CommandID string
}
