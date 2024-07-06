package common

type Response int32

const (
	Validated Response = iota
	Created
	Removed
	Modified
	Executed
	Skipped
	Failed
)

type TaskResponse struct {
	Type    Response
	Changes []string
	Message string
}
