package abac

import (
	"fmt"
	"strings"
	"sync"
)

type Engine struct {
	Policies []Policy
	mu       sync.RWMutex
}

func NewABACEngine() *Engine {
	return &Engine{
		Policies: make([]Policy, 0),
	}
}

func (e *Engine) AddPolicy(p Policy) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Policies = append(e.Policies, p)
}

func (e *Engine) Evaluate(ctx Context) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var allowPolicies []Policy
	var denyPolicies []Policy

	for _, policy := range e.Policies {
		if !e.matchTarget(policy.Target, ctx) {
			continue
		}
		if e.matchConditions(policy.Conditions, ctx) {
			if policy.Effect == "allow" {
				allowPolicies = append(allowPolicies, policy)
			} else if policy.Effect == "deny" {
				denyPolicies = append(denyPolicies, policy)
			}
		}
	}

	// Находим политику с максимальным приоритетом среди deny и allow
	var bestAllow, bestDeny *Policy
	for _, p := range allowPolicies {
		if bestAllow == nil || p.Priority > bestAllow.Priority {
			bestAllow = &p
		}
	}
	for _, p := range denyPolicies {
		if bestDeny == nil || p.Priority > bestDeny.Priority {
			bestDeny = &p
		}
	}

	if bestDeny != nil && (bestAllow == nil || bestDeny.Priority >= bestAllow.Priority) {
		return false, nil
	}
	if bestAllow != nil {
		return true, nil
	}
	return false, nil // по умолчанию — deny
}

func (e *Engine) matchTarget(t Target, ctx Context) bool {
	if t.Resource != "*" && t.Resource != getString(ctx.Resource, "type") {
		return false
	}
	if t.Action != "*" && t.Action != ctx.Action {
		return false
	}
	return true
}

func (e *Engine) matchConditions(conds []Condition, ctx Context) bool {
	for _, c := range conds {
		if !e.evalCondition(c, ctx) {
			return false
		}
	}
	return true
}

func (e *Engine) evalCondition(c Condition, ctx Context) bool {
	var actual interface{}

	switch c.Attribute {
	case "user.id":
		if ctx.User != nil {
			actual = ctx.User.ID.String()
		}
	case "user.role":
		if ctx.User != nil {
			actual = string(ctx.User.Role)
		}
	default:
		// resource.xxx или env.xxx
		if attr, ok := extractAttribute(c.Attribute, ctx); ok {
			actual = attr
		}
	}

	if actual == nil {
		return false
	}

	switch c.Operator {
	case "eq":
		return compareEqual(actual, c.Value)
	case "in":
		return compareIn(actual, c.Value)
	case "gt", "lt", "gte", "lte":
		return compareNumeric(actual, c.Value, c.Operator)
	}
	return false
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func extractAttribute(attr string, ctx Context) (interface{}, bool) {
	parts := strings.Split(attr, ".")
	if len(parts) != 2 {
		return nil, false
	}

	switch parts[0] {
	case "user":
		if ctx.User == nil {
			return nil, false
		}
		switch parts[1] {
		case "id":
			return ctx.User.ID.String(), true
		case "role":
			return string(ctx.User.Role), true
		}
	case "resource":
		if val, ok := ctx.Resource[parts[1]]; ok {
			return val, true
		}
	case "env", "environment":
		if val, ok := ctx.Environment[parts[1]]; ok {
			return val, true
		}
	}

	return nil, false
}

func compareEqual(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func compareIn(a interface{}, values interface{}) bool {
	slice, ok := values.([]string)
	if !ok {
		return false
	}
	aval := fmt.Sprintf("%v", a)
	for _, v := range slice {
		if aval == v {
			return true
		}
	}
	return false
}

func compareNumeric(a, b interface{}, op string) bool {
	// упрощённо, можно улучшить через reflect
	return false
}
