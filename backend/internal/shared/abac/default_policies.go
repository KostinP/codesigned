// internal/shared/abac/default_policies.go
package abac

func GetDefaultPolicies() []Policy {
	return []Policy{
		// 0. Админы могут всё
		{
			ID:         "admin_full_access",
			Name:       "Admin Full Access",
			Target:     Target{Resource: "*", Action: "*"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "eq", Value: "admin"}},
			Effect:     "allow",
			Priority:   1000,
		},
		// 1. Сессии — только свои
		{
			ID:     "session_manage_own",
			Name:   "Manage Own Sessions",
			Target: Target{Resource: "user_sessions", Action: "*"},
			Conditions: []Condition{
				{Attribute: "user.role", Operator: "in", Value: []string{"student", "teacher", "admin"}},
				{Attribute: "resource.user_id", Operator: "eq", Value: "user.id"},
			},
			Effect:   "allow",
			Priority: 100,
		},
		// 2. Аналитика — только админы
		{
			ID:         "analytics_read",
			Name:       "Read Analytics",
			Target:     Target{Resource: "analytics", Action: "read"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "eq", Value: "admin"}},
			Effect:     "allow",
			Priority:   200,
		},
		// ========== КУРСЫ ==========
		// 2.1 Создание курсов — teacher/admin
		{
			ID:         "course_create",
			Name:       "Create Courses",
			Target:     Target{Resource: "course", Action: "create"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "in", Value: []string{"teacher", "admin"}}},
			Effect:     "allow",
			Priority:   100,
		},
		// 2.2 Чтение всех курсов — все авторизованные
		{
			ID:         "course_read_all",
			Name:       "Read Any Course",
			Target:     Target{Resource: "course", Action: "read"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "in", Value: []string{"student", "teacher", "admin"}}},
			Effect:     "allow",
			Priority:   50,
		},
		// 2.3 Редактирование/удаление — только автор или админ
		{
			ID:         "course_manage_own",
			Name:       "Manage Own Course",
			Target:     Target{Resource: "course", Action: "*"},
			Conditions: []Condition{{Attribute: "resource.author_id", Operator: "eq", Value: "user.id"}},
			Effect:     "allow",
			Priority:   150,
		},
		// ========== МОДУЛИ ==========
		{
			ID:         "module_create_update_delete",
			Name:       "Manage Own Module",
			Target:     Target{Resource: "module", Action: "*"},
			Conditions: []Condition{{Attribute: "resource.author_id", Operator: "eq", Value: "user.id"}},
			Effect:     "allow",
			Priority:   150,
		},
		{
			ID:         "module_read",
			Name:       "Read Modules",
			Target:     Target{Resource: "module", Action: "read"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "in", Value: []string{"student", "teacher", "admin"}}},
			Effect:     "allow",
			Priority:   50,
		},
		// ========== УРОКИ ==========
		{
			ID:         "lesson_manage_own",
			Name:       "Manage Own Lesson",
			Target:     Target{Resource: "lesson", Action: "*"},
			Conditions: []Condition{{Attribute: "resource.author_id", Operator: "eq", Value: "user.id"}},
			Effect:     "allow",
			Priority:   150,
		},
		{
			ID:         "lesson_read",
			Name:       "Read Lessons",
			Target:     Target{Resource: "lesson", Action: "read"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "in", Value: []string{"student", "teacher", "admin"}}},
			Effect:     "allow",
			Priority:   50,
		},
		// ========== КАТЕГОРИИ ==========
		// 7.1 Создание категорий — teacher/admin
		{
			ID:         "category_create",
			Name:       "Create Categories",
			Target:     Target{Resource: "category", Action: "create"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "in", Value: []string{"teacher", "admin"}}},
			Effect:     "allow",
			Priority:   100,
		},
		// 7.2 Управление своими категориями
		{
			ID:         "category_manage_own",
			Name:       "Manage Own Category",
			Target:     Target{Resource: "category", Action: "*"},
			Conditions: []Condition{{Attribute: "resource.author_id", Operator: "eq", Value: "user.id"}},
			Effect:     "allow",
			Priority:   150,
		},
		// 7.3 Чтение всех категорий
		{
			ID:         "category_read",
			Name:       "Read Categories",
			Target:     Target{Resource: "category", Action: "read"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "in", Value: []string{"student", "teacher", "admin"}}},
			Effect:     "allow",
			Priority:   50,
		},
		// ========== ТЕГИ (уже были, но обновлённые) ==========
		{
			ID:         "tag_create",
			Name:       "Create Tags",
			Target:     Target{Resource: "tag", Action: "create"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "in", Value: []string{"teacher", "admin"}}},
			Effect:     "allow",
			Priority:   100,
		},
		{
			ID:         "tag_manage_own",
			Name:       "Manage Own Tags",
			Target:     Target{Resource: "tag", Action: "*"},
			Conditions: []Condition{{Attribute: "resource.author_id", Operator: "eq", Value: "user.id"}},
			Effect:     "allow",
			Priority:   150,
		},
		{
			ID:         "tag_read",
			Name:       "Read Tags",
			Target:     Target{Resource: "tag", Action: "read"},
			Conditions: []Condition{{Attribute: "user.role", Operator: "in", Value: []string{"student", "teacher", "admin"}}},
			Effect:     "allow",
			Priority:   50,
		},
	}
}
