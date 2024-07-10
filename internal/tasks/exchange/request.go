package exchange

type Request int32

const (
	Validate Request = iota
	Create
	Remove
	Modify
	Execute
	Passive
)

type TaskRequest struct {
	Type    Request
	Changes []string
}
