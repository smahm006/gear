package tasks

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/smahm006/gear/internal/playbook/state"
)

type TagsValidationError struct {
	Tags string
	Err  error
}

func (l *TagsValidationError) Error() string {
	return fmt.Sprintf("invalid tags %q\n%v", l.Tags, l.Err)
}

type Token int

const (
	Not Token = iota
	Or
	And
)

func (t Token) Apply(slice1 []string, slice2 []string) []string {
	hash := make(map[string]bool)
	var inter []string
	switch t {
	case Not:
		// Removal
		for _, e := range slice2 {
			hash[e] = true
		}
		for _, val := range slice1 {
			if !hash[val] {
				inter = append(inter, val)
			}
		}
		return inter
	case Or:
		// Union
		for _, e := range slice1 {
			hash[e] = true
		}
		for _, e := range slice2 {
			hash[e] = true
		}
		for k := range hash {
			inter = append(inter, k)
		}
		return inter
	case And:
		// Intersection
		for _, e := range slice1 {
			hash[e] = true
		}
		for _, e := range slice2 {
			if hash[e] {
				inter = append(inter, e)
			}
		}
		return inter
	}
	return nil
}

func getAllTasks(tasks []*Task) []string {
	var tasks_list []string
	for _, task := range tasks {
		tasks_list = append(tasks_list, task.Name)
	}
	return tasks_list
}

func getTasksByName(name string, tasks []*Task) ([]string, error) {
	var tasks_list []string
	for _, task := range tasks {
		if task.With != nil {
			if slices.Contains(task.With.Tags, name) {
				tasks_list = append(tasks_list, task.Name)
			}
		}
	}
	if len(tasks_list) != 0 {
		return tasks_list, nil
	}
	return nil, fmt.Errorf("no tasks found for (%s) referenced in tags found", name)
}

func divideByToken(tags string, token Token) []string {
	var r *regexp.Regexp
	switch token {
	case Not:
		r = regexp.MustCompile(`NOT`)
	case Or:
		r = regexp.MustCompile(`OR`)
	case And:
		r = regexp.MustCompile(`AND`)
	}
	parts := r.Split(tags, 2)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func expressionParse(tags string, tasks []*Task) ([]string, error) {
	limit_err := &TagsValidationError{Tags: tags}
	var tasks_limited []string
	tokens := []Token{Not, Or, And}

	var applyOperations func(string, bool) ([]string, error)
	applyOperations = func(expression string, started_with_not bool) ([]string, error) {
		var result []string
		re := regexp.MustCompile(`NOT|OR|AND`)
		if len(expression) == 0 && started_with_not {
			all_tasks := getAllTasks(tasks)
			return append(result, all_tasks...), nil
		} else if len(expression) != 0 && !re.MatchString(expression) {
			sole_result, err := getTasksByName(expression, tasks)
			if err != nil {
				return nil, err
			}
			return append(result, sole_result...), nil
		}
		for _, token := range tokens {
			sub_exprs := divideByToken(expression, token)
			switch token {
			case Not:
				if len(sub_exprs) > 1 {
					sub_result_first, err := applyOperations(sub_exprs[0], started_with_not)
					if err != nil {
						return nil, err
					}
					sub_result_remainder, err := applyOperations(sub_exprs[1], started_with_not)
					if err != nil {
						return nil, err
					}
					result = Not.Apply(sub_result_first, sub_result_remainder)
					return result, nil
				}
			case Or:
				if len(sub_exprs) > 1 {
					sub_result_first, err := applyOperations(sub_exprs[0], started_with_not)
					if err != nil {
						return nil, err
					}
					sub_result_remainder, err := applyOperations(sub_exprs[1], started_with_not)
					if err != nil {
						return nil, err
					}
					result = Or.Apply(sub_result_first, sub_result_remainder)
					return result, nil
				}
			case And:
				if len(sub_exprs) > 1 {
					sub_result_first, err := applyOperations(sub_exprs[0], started_with_not)
					if err != nil {
						return nil, err
					}
					sub_result_remainder, err := applyOperations(sub_exprs[1], started_with_not)
					if err != nil {
						return nil, err
					}
					result = And.Apply(sub_result_first, sub_result_remainder)
					return result, nil
				}
			}
		}
		return result, nil
	}

	// Check if limit started with NOT
	re_not := regexp.MustCompile(`^NOT.*`)

	// Call the recursive function to apply operations
	tasks_limited, err := applyOperations(tags, re_not.MatchString(tags))
	if err != nil {
		limit_err.Err = err
		return nil, limit_err
	}

	return tasks_limited, nil
}

func getTasksGivenTags(tags string, tasks []*Task) ([]*Task, error) {
	var tasks_limited []*Task
	tasks_list, err := expressionParse(tags, tasks)
	for _, task_name := range tasks_list {
		for _, task := range tasks {
			if task.Name == task_name {
				tasks_limited = append(tasks_limited, task)
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return tasks_limited, nil
}

func collectTasks(state *state.RunState, tasks []*Task) ([]*Task, error) {
	var err error
	var collective_tasks []*Task = tasks
	limited := len(state.ParsedFlags.Tags) != 0
	if !limited {
		return collective_tasks, nil
	} else {
		collective_tasks, err = getTasksGivenTags(state.ParsedFlags.Tags, tasks)
		if err != nil {
			return nil, err
		}
	}
	return collective_tasks, nil
}
