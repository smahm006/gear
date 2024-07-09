package requonse

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
}

func NewTaskResponse() *TaskResponse {
	return &TaskResponse{
		CommandResult: NewCommandResult(),
	}
}
