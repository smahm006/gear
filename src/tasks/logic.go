package tasks

type PreTaskLogic struct {
	conditions string
	sudo       string
	items      []string
	tags       []string
	delegateTo string
}

type PostTaskLogic struct {
	ignore_errors string
	retry         int
	delay         int
	notify        string
}
