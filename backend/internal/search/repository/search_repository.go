// internal/search/repository/search_repository.go
package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/search/entity"
)

type SearchRepository interface {
	Search(ctx context.Context, filters entity.SearchFilters) (*entity.SearchResult, error)
	SearchAdvanced(ctx context.Context, filters entity.SearchFilters) (*entity.SearchResult, error)
	GetAutocomplete(ctx context.Context, query string, limit int) ([]string, error)
	GetPopularSearches(ctx context.Context, limit int) ([]string, error)
}

type PostgresSearchRepository struct {
	db *pgxpool.Pool
}

func NewPostgresSearchRepository(db *pgxpool.Pool) *PostgresSearchRepository {
	return &PostgresSearchRepository{db: db}
}

func (r *PostgresSearchRepository) Search(ctx context.Context, filters entity.SearchFilters) (*entity.SearchResult, error) {
	// Собираем WHERE условия
	var conditions []string
	var args []interface{}
	argIndex := 1
	// Базовые условия
	conditions = append(conditions, "deleted_at IS NULL")
	// Поиск по тексту
	if filters.Query != "" {
		conditions = append(conditions,
			fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+filters.Query+"%")
		argIndex++
	}
	// Поиск по категории
	if filters.CategoryID != nil {
		conditions = append(conditions,
			fmt.Sprintf("id IN (SELECT target_id FROM category_assignments WHERE category_id = $%d)", argIndex))
		args = append(args, *filters.CategoryID)
		argIndex++
	}
	// Поиск по тегам
	if len(filters.Tags) > 0 {
		placeholders := make([]string, len(filters.Tags))
		for i := range filters.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argIndex+i)
		}
		conditions = append(conditions,
			fmt.Sprintf("id IN (SELECT entity_id FROM tag_assignments WHERE entity_type = 'course' AND tag_id IN (SELECT id FROM tags WHERE name IN (%s)))",
				strings.Join(placeholders, ",")))
		for _, tag := range filters.Tags { // Исправлено: ручная конвертация []string в append
			args = append(args, tag)
		}
		argIndex += len(filters.Tags)
	}
	// Поиск по цене
	if filters.PriceMin != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, *filters.PriceMin)
		argIndex++
	}
	if filters.PriceMax != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, *filters.PriceMax)
		argIndex++
	}
	// Определяем таблицу для поиска
	// Упрощенный пример - ищем только курсы
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	// Запрос для подсчета общего количества
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM courses %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}
	// Запрос данных с пагинацией
	query := fmt.Sprintf(`
		SELECT
			id, title, description, price, author_id, created_at, updated_at,
			'course' as type,
			1.0 as relevance
		FROM courses
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filters.Limit, filters.Offset)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*entity.SearchEntity
	for rows.Next() {
		var e entity.SearchEntity
		var price *int
		var authorID *uuid.UUID
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &price, &authorID,
			&e.CreatedAt, &e.UpdatedAt, &e.Type, &e.Relevance,
		)
		if err != nil {
			return nil, err
		}
		e.Price = price
		e.AuthorID = authorID
		entities = append(entities, &e)
	}
	return &entity.SearchResult{
		Entities: entities,
		Total:    total,
		Limit:    filters.Limit,
		Offset:   filters.Offset,
	}, nil
}

func (r *PostgresSearchRepository) SearchAdvanced(ctx context.Context, filters entity.SearchFilters) (*entity.SearchResult, error) {
	// Расширенный поиск по всем типам контента
	// Используем UNION для поиска по курсам, урокам и модулям
	var queries []string
	var args []interface{}
	argIndex := 1
	// Базовые условия
	baseConditions := []string{"deleted_at IS NULL"}
	// Добавляем условия по запросу
	if filters.Query != "" {
		baseConditions = append(baseConditions,
			fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+filters.Query+"%")
		argIndex++
	}
	// Курсы
	coursesQuery := fmt.Sprintf(`
		SELECT
			id, title, description, price as price, author_id,
			created_at, updated_at, 'course' as type, 1.0 as relevance
		FROM courses
		WHERE %s
	`, strings.Join(baseConditions, " AND "))
	queries = append(queries, coursesQuery)
	// Уроки
	lessonsQuery := fmt.Sprintf(`
		SELECT
			id, title, content as description, NULL as price, author_id,
			created_at, updated_at, 'lesson' as type, 1.0 as relevance
		FROM lessons
		WHERE %s
	`, strings.Join(baseConditions, " AND "))
	queries = append(queries, lessonsQuery)
	// Объединяем запросы
	fullQuery := "(" + strings.Join(queries, ") UNION ALL (") + ")"
	fullQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, filters.Offset)
	rows, err := r.db.Query(ctx, fullQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*entity.SearchEntity
	for rows.Next() {
		var e entity.SearchEntity
		var price *int
		var authorID *uuid.UUID
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &price, &authorID,
			&e.CreatedAt, &e.UpdatedAt, &e.Type, &e.Relevance,
		)
		if err != nil {
			return nil, err
		}
		e.Price = price
		e.AuthorID = authorID
		entities = append(entities, &e)
	}
	// TODO: Получить общее количество
	total := len(entities) // Упрощенно
	return &entity.SearchResult{
		Entities: entities,
		Total:    total,
		Limit:    filters.Limit,
		Offset:   filters.Offset,
	}, nil
}

func (r *PostgresSearchRepository) GetAutocomplete(ctx context.Context, query string, limit int) ([]string, error) {
	// TODO: реализовать автодополнение, например, с использованием LIKE или полнотекстового поиска
	return []string{}, nil
}

func (r *PostgresSearchRepository) GetPopularSearches(ctx context.Context, limit int) ([]string, error) {
	// TODO: реализовать, если есть таблица популярных запросов
	return []string{}, nil
}
