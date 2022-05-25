package promise

type Status string

var (
	Pending   Status = "pending"
	Fulfilled Status = "fulfilled"
	Rejected  Status = "rejected"
)

func (s Status) Valid() bool {
	switch s {
	case Pending,
		Fulfilled,
		Rejected:
		return true
	default:
		return false
	}
}
