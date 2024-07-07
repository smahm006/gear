/*
   1. NOT - (webservers OR dbservers AND staging) NOT (phoenix)
   2. OR  ((webservers) OR (dbservers AND staging)) NOT (phoenix)
   3. AND ((webservers) OR ((dbservers) AND (staging))) NOT (phoenix)

   “all machines in the groups ‘webservers’ and ‘dbservers’ are to be managed if they are in the group ‘staging’ also, but the machines are not to be managed if they are in the group ‘phoenix’"

   webservers:dbservers:&staging:!phoenix

   gear --limit "webservers OR dbservers AND staging AND fire NOT phoenix OR test"

   gear --limit "webservers OR dbservers AND staging NOT phoenix"

   OR = (?:^(?:([^(?:NOT|AND)]*)OR([^(?:NOT|AND)]*))|OR([^(?:NOT|AND)]*))
   NOT = (?:^(?:([^(?:AND|OR)]*)NOT([^(?:AND|OR)]*))|NOT([^(?:AND|OR)]*))
   AND = (?:^(?:([^(?:ORNOTAND)]*)AND([^(?:ORNOTAND)]*))|AND([^(?:ORNOTAND)]*))

*/

package playbook

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/smahm006/gear/src/inventory"
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
		fmt.Printf("REMOVAL OF\n%v, %v\n%v\n", slice1, slice2, inter)
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
		fmt.Printf("UNION OF\n%v, %v\n%v\n", slice1, slice2, inter)
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
		fmt.Printf("INTERSECTION OF\n%v, %v\n%v\n", slice1, slice2, inter)
		return inter
	}
	return nil
}

func getHostsByName(name string, groups map[string]*inventory.Group) ([]string, error) {
	var hosts []string
	if group, exists := groups[name]; exists {
		for host_name, _ := range group.Hosts {
			hosts = append(hosts, host_name)
		}
		return hosts, nil
	}
	for _, group := range groups {
		if group.SubGroups != nil {
			sub_group_hosts, err := getHostsByName(name, group.SubGroups)
			if err == nil {
				hosts = append(hosts, sub_group_hosts...)
				return hosts, nil
			}
		}
		if host, exists := group.GetHost(name); exists {
			hosts = append(hosts, host.Name)
			return hosts, nil
		}
	}
	return nil, fmt.Errorf("no group or host (%s) referenced in limit found", name)
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

// func expressionParse(limit string, groups map[string]*inventory.Group) ([]string, error) {
// 	limit_err := &LimitValidationError{Limit: limit}
// 	var hosts_limited []string
// 	not_groups := divideByToken(limit, Not)
// 	for i := len(not_groups) - 1; i >= 1; i-- {
// 		or_groups := divideByToken(not_groups[i], Or)
// 		for i := len(or_groups) - 1; i >= 1; i-- {
// 			and_groups := divideByToken(or_groups[i], And)
// 			fmt.Println(and_groups)
// 			and_intersection, err := getHostsByName(and_groups[len(and_groups)-1], groups)
// 			if err != nil {
// 				limit_err.Err = err
// 				return nil, limit_err
// 			}
// 			for i := len(and_groups) - 2; i >= 0; i-- {
// 				and_host, err := getHostsByName(and_groups[i], groups)
// 				if err != nil {
// 					limit_err.Err = err
// 					return nil, limit_err
// 				}
// 				and_intersection = And.Apply(and_intersection, and_host, 0)
// 			}
// 		}
// 	}
// 	return hosts_limited, nil
// }

func expressionParse(limit string, groups map[string]*inventory.Group) ([]string, error) {
	limitErr := &LimitValidationError{Limit: limit}
	var hostsLimited []string
	tokens := []Token{Not, Or, And}

	var applyOperations func(string) ([]string, error)
	applyOperations = func(expression string) ([]string, error) {
		fmt.Println("EXPRESSION: ", expression)
		var result []string
		var re = regexp.MustCompile(`NOT|OR|AND`)
		if !re.MatchString(expression) {
			sole_name := expression
			sole_hosts, err := getHostsByName(sole_name, groups)
			if err != nil {
				return nil, err
			}
			result = append(result, sole_hosts...)
			return result, nil
		}
		for _, token := range tokens {
			sub_exprs := divideByToken(expression, token)
			switch token {
			case Not:
				if len(sub_exprs) > 1 {
					not_name := sub_exprs[0]
					not_hosts, err := getHostsByName(not_name, groups)
					if err != nil {
						return nil, err
					}
					sub_result, err := applyOperations(sub_exprs[1])
					if err != nil {
						return nil, err
					}
					result = Not.Apply(not_hosts, sub_result)
				}
			case Or:
				if len(sub_exprs) > 1 {
					or_name := sub_exprs[0]
					or_hosts, err := getHostsByName(or_name, groups)
					if err != nil {
						return nil, err
					}
					sub_result, err := applyOperations(sub_exprs[1])
					if err != nil {
						return nil, err
					}
					result = Or.Apply(or_hosts, sub_result)
				}
			case And:
				if len(sub_exprs) > 1 {
					and_name := sub_exprs[0]
					and_hosts, err := getHostsByName(and_name, groups)
					if err != nil {
						return nil, err
					}
					sub_result, err := applyOperations(sub_exprs[1])
					if err != nil {
						return nil, err
					}
					result = And.Apply(and_hosts, sub_result)
				}
			}
		}
		return result, nil
	}

	// Call the recursive function to apply operations
	hostsLimited, err := applyOperations(limit)
	if err != nil {
		limitErr.Err = err
		return nil, limitErr
	}

	return hostsLimited, nil
}

func getHostsGivenLimit(limit string, groups map[string]*inventory.Group) (map[string]*inventory.Host, error) {
	var matching_hosts map[string]*inventory.Host
	_, err := expressionParse(limit, groups)
	if err != nil {
		return nil, err
	}
	return matching_hosts, nil
}
