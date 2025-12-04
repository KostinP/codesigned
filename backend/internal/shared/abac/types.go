package abac

import (
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
)

type Target struct {
	Resource string `json:"resource"` // course, lesson, analytics, tag, category, *
	Action   string `json:"action"`   // read, create, update, delete, manage, *
}

type Condition struct {
	Attribute string      `json:"attribute"` // user.role, resource.author_id, env.time.hour
	Operator  string      `json:"operator"`  // eq, in, gt, lt, contains
	Value     interface{} `json:"value"`
}

type Policy struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Target     Target      `json:"target"`
	Conditions []Condition `json:"conditions"`
	Effect     string      `json:"effect"` // allow / deny
	Priority   int         `json:"priority"`
}

type Context struct {
	User        *entity.User
	Resource    map[string]interface{} // id, author_id, type, target_author_id и т.д.
	Action      string
	Environment map[string]interface{}
}
