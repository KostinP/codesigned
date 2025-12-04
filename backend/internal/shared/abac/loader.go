// internal/shared/abac/loader.go
package abac

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbPolicy struct {
	ID             uuid.UUID       `db:"id"`
	Name           string          `db:"name"`
	TargetResource string          `db:"target_resource"`
	TargetAction   string          `db:"target_action"`
	Conditions     json.RawMessage `db:"conditions"`
	Effect         string          `db:"effect"`
	Priority       int             `db:"priority"`
}

func LoadPoliciesFromDB(engine *Engine, pool *pgxpool.Pool) error {
	rows, err := pool.Query(context.Background(), `
        SELECT id, name, target_resource, target_action, conditions, effect, priority
        FROM abac_policies
        ORDER BY priority DESC
    `)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var dbp dbPolicy
		var condsJSON json.RawMessage
		if err := rows.Scan(&dbp.ID, &dbp.Name, &dbp.TargetResource, &dbp.TargetAction, &condsJSON, &dbp.Effect, &dbp.Priority); err != nil {
			continue // или return err
		}

		var conditions []Condition
		if err := json.Unmarshal(condsJSON, &conditions); err != nil {
			continue
		}

		policy := Policy{
			ID:   dbp.ID.String(),
			Name: dbp.Name,
			Target: Target{
				Resource: dbp.TargetResource,
				Action:   dbp.TargetAction,
			},
			Conditions: conditions,
			Effect:     dbp.Effect,
			Priority:   dbp.Priority,
		}

		engine.AddPolicy(policy)
	}

	return rows.Err()
}
