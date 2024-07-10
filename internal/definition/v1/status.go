package v1

const (
	Create = "create"
	Update = "update"
	Delete = "delete"
	None   = ""

	Pending = "pending"

	Applied = "applied"
	Failed  = "failed"
)

type Status struct {
	Current string `json:"current" yaml:"current" bson:"current"`
	Desired string `json:"desired,omitempty" yaml:"desired,omitempty" bson:"desired,omitempty"`

	CreatedAt          string `json:"createdAt" yaml:"createdAt" bson:"createdAt"`
	UpdatedAt          string `json:"updatedAt" yaml:"updatedAt" bson:"updatedAt"`
	MaxPendingDuration int    `json:"maxPendingDuration,omitempty" yaml:"maxPendingDuration,omitempty" bson:"maxPendingDuration,omitempty"`
}

func (s *Status) ClearDesired() {
	s.Desired = None
}

func (s *Status) SetCurrentToApplied() {
	s.Current = Applied
}

func (s *Status) SetCurrentToPending() {
	s.Current = Pending
}

func (s *Status) SetDesiredToUpdate() {
	s.Desired = Update
}

func (s *Status) SetDesiredToDelete() {
	s.Desired = Delete
}
