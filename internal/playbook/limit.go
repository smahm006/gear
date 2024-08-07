package playbook

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/smahm006/gear/internal/inventory"
	"github.com/smahm006/gear/internal/playbook/state"
)

type LimitValidationError struct {
	Limit string
	Err   error
}

func (l *LimitValidationError) Error() string {
	return fmt.Sprintf("invalid limit %q\n%v", l.Limit, l.Err)
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

func getAllHosts(state *state.RunState) []string {
	var hosts []string
	for host_name, _ := range state.Inventory.GroupHostsMembership.HostsToGroup {
		hosts = append(hosts, host_name)
	}
	return hosts
}

func getHostsByName(name string, state *state.RunState, play *Play) ([]string, error) {
	if slices.Contains(play.Groups, name) {
		if hosts, exists := state.Inventory.GroupHostsMembership.GroupToHosts[name]; exists {
			return hosts, nil
		}
	}
	if groups, exists := state.Inventory.GroupHostsMembership.HostsToGroup[name]; exists {
		for _, gname := range groups {
			if slices.Contains(play.Groups, gname) {
				return []string{name}, nil
			}
		}
	}
	return nil, fmt.Errorf("no group or hosts found for (%s) referenced in limit found", name)
}

func divideByToken(limit string, token Token) []string {
	var r *regexp.Regexp
	switch token {
	case Not:
		r = regexp.MustCompile(`NOT`)
	case Or:
		r = regexp.MustCompile(`OR`)
	case And:
		r = regexp.MustCompile(`AND`)
	}
	parts := r.Split(limit, 2)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func expressionParse(limit string, state *state.RunState, play *Play) ([]string, error) {
	limit_err := &LimitValidationError{Limit: limit}
	var hosts_limited []string
	tokens := []Token{Not, Or, And}

	var applyOperations func(string, bool) ([]string, error)
	applyOperations = func(expression string, started_with_not bool) ([]string, error) {
		var result []string
		re := regexp.MustCompile(`NOT|OR|AND`)
		if len(expression) == 0 && started_with_not {
			all_hosts := getAllHosts(state)
			return append(result, all_hosts...), nil
		} else if !re.MatchString(expression) {
			sole_result, err := getHostsByName(expression, state, play)
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
	hosts_limited, err := applyOperations(limit, re_not.MatchString(limit))
	if err != nil {
		limit_err.Err = err
		return nil, limit_err
	}

	return hosts_limited, nil
}

func getHostsGivenLimit(limit string, state *state.RunState, play *Play) (map[string]*inventory.Host, error) {
	hosts_limited := make(map[string]*inventory.Host)
	hosts, err := expressionParse(limit, state, play)
	if err != nil {
		return nil, err
	}
	for _, host_name := range hosts {
		hosts_limited[host_name], _ = state.Inventory.GetHost(host_name)
	}
	return hosts_limited, nil
}

func collectHosts(state *state.RunState, play *Play) (map[string]*inventory.Host, error) {
	var err error
	collective_hosts := make(map[string]*inventory.Host)
	limited := len(state.ParsedFlags.Limit) != 0
	if !limited {
		for _, group_name := range play.Groups {
			hosts := state.Inventory.GroupHostsMembership.GroupToHosts[group_name]
			for _, host_name := range hosts {
				host, _ := state.Inventory.GetHost(host_name)
				collective_hosts[host_name] = host
			}
		}
		return collective_hosts, nil
	} else {
		collective_hosts, err = getHostsGivenLimit(state.ParsedFlags.Limit, state, play)
		if err != nil {
			return nil, err
		}
	}
	return collective_hosts, nil
}
