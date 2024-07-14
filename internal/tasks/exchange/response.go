package exchange

type Response int32

const (
	Validated Response = iota
	Created
	Removed
	Modified
	Executed
	Skipped
	Failed
	Unsupported
)

type CommandResult struct {
	Cmd string
	Out string
	Err string
	Rc  int
}

func NewCommandResult() *CommandResult {
	return &CommandResult{
		Cmd: "",
		Out: "",
		Err: "",
		Rc:  1,
	}
}

type TaskResponse struct {
	Type          Response
	Changes       []string
	Message       string
	CommandResult *CommandResult
	Error         error
}

func NewTaskResponse() *TaskResponse {
	return &TaskResponse{
		CommandResult: NewCommandResult(),
	}
}
