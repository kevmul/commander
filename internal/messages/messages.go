package messages

// =====================================
// Modal messages
// =====================================

type ModalClosedMsg struct{}

type ShowModalMsg struct {
	ModalType string // "entry", "help", etc
}

type ItemDeletedMsg struct {
	ID   string
	Type string
}

