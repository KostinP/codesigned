// internal/category/repository/category_navigation_repository.go
package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/category/entity"
	searchEntity "github.com/kostinp/edu-platform-backend/internal/search/entity"
	tagEntity "github.com/kostinp/edu-platform-backend/internal/tag/entity"
)

type CategoryNavigationRepository interface {
	// Получение дерева категорий
	GetTree(ctx context.Context) ([]*entity.TreeCategory, error)
	GetSubtree(ctx context.Context, categoryID uuid.UUID) ([]*entity.TreeCategory, error)
	GetBreadcrumbs(ctx context.Context, categoryID uuid.UUID) ([]*entity.Category, error)
	GetCategoryWithStats(ctx context.Context, categoryID uuid.UUID) (*entity.CategoryWithStats, error)
	// Поиск по категории
	SearchInCategory(ctx context.Context, categoryID uuid.UUID, query string) ([]searchEntity.SearchEntity, error)
	GetContentByCategory(ctx context.Context, categoryID uuid.UUID, contentType string) ([]searchEntity.SearchEntity, error)
	// Получение связанных категорий
	GetSiblingCategories(ctx context.Context, categoryID uuid.UUID) ([]*entity.Category, error)
	GetPathCategories(ctx context.Context, categoryID uuid.UUID) ([]*entity.Category, error)
}

type PostgresCategoryNavigationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCategoryNavigationRepository(db *pgxpool.Pool) *PostgresCategoryNavigationRepository {
	return &PostgresCategoryNavigationRepository{db: db}
}

// GetTree возвращает полное дерево категорий
func (r *PostgresCategoryNavigationRepository) GetTree(ctx context.Context) ([]*entity.TreeCategory, error) {
	// Сначала получаем все категории
	query := `
		SELECT
			c.id, c.name, c.slug, c.description, c.parent_id,
			c.author_id, c.sort_order, c.is_visible, c.created_at, c.updated_at,
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'id', t.id,
				'name', t.name,
				'description', t.description
			)) FILTER (WHERE t.id IS NOT NULL), '[]') as tags
		FROM categories c
		LEFT JOIN tags t ON t.category_id = c.id
		WHERE c.is_visible = true
		GROUP BY c.id
		ORDER BY c.sort_order, c.name
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*entity.TreeCategory
	for rows.Next() {
		var cat entity.TreeCategory
		var tagsJSON []byte
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &cat.Description, &cat.ParentID,
			&cat.AuthorID, &cat.SortOrder, &cat.IsVisible, &cat.CreatedAt, &cat.UpdatedAt,
			&tagsJSON,
		)
		if err != nil {
			return nil, err
		}
		// Парсим теги
		cat.Tags = parseTagsJSON(tagsJSON)
		categories = append(categories, &cat)
	}
	// Строим дерево
	return r.buildTree(categories, nil), nil
}

func (r *PostgresCategoryNavigationRepository) buildTree(
	categories []*entity.TreeCategory,
	parentID *uuid.UUID,
) []*entity.TreeCategory {
	var tree []*entity.TreeCategory
	for _, category := range categories {
		if (parentID == nil && category.ParentID == nil) ||
			(parentID != nil && category.ParentID != nil && *category.ParentID == *parentID) {
			category.Children = r.buildTree(categories, &category.ID)
			tree = append(tree, category)
		}
	}
	return tree
}

// GetSubtree возвращает поддерево категорий
func (r *PostgresCategoryNavigationRepository) GetSubtree(ctx context.Context, categoryID uuid.UUID) ([]*entity.TreeCategory, error) {
	// Аналогично GetTree, но начиная с указанной категории
	// Сначала получаем все подкатегории рекурсивно
	query := `
		WITH RECURSIVE sub_categories AS (
		    SELECT * FROM categories WHERE id = $1
		    UNION
		    SELECT c.* FROM categories c
		    INNER JOIN sub_categories sc ON c.parent_id = sc.id
		)
		SELECT
			sc.id, sc.name, sc.slug, sc.description, sc.parent_id,
			sc.author_id, sc.sort_order, sc.is_visible, sc.created_at, sc.updated_at,
			COALESCE(json_agg(DISTINCT jsonb_build_object(
				'id', t.id,
				'name', t.name,
				'description', t.description
			)) FILTER (WHERE t.id IS NOT NULL), '[]') as tags
		FROM sub_categories sc
		LEFT JOIN tags t ON t.category_id = sc.id
		WHERE sc.is_visible = true
		GROUP BY sc.id
		ORDER BY sc.sort_order, sc.name
	`
	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*entity.TreeCategory
	for rows.Next() {
		var cat entity.TreeCategory
		var tagsJSON []byte
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &cat.Description, &cat.ParentID,
			&cat.AuthorID, &cat.SortOrder, &cat.IsVisible, &cat.CreatedAt, &cat.UpdatedAt,
			&tagsJSON,
		)
		if err != nil {
			return nil, err
		}
		// Парсим теги
		cat.Tags = parseTagsJSON(tagsJSON)
		categories = append(categories, &cat)
	}
	// Строим дерево, начиная с указанной категории
	return r.buildTree(categories, &categoryID), nil
}

// GetBreadcrumbs возвращает хлебные крошки для категории
func (r *PostgresCategoryNavigationRepository) GetBreadcrumbs(ctx context.Context, categoryID uuid.UUID) ([]*entity.Category, error) {
	// Используем рекурсивный CTE для получения пути
	query := `
		WITH RECURSIVE category_path AS (
			SELECT id, name, slug, parent_id, 1 as depth
			FROM categories
			WHERE id = $1
			UNION ALL
			SELECT c.id, c.name, c.slug, c.parent_id, cp.depth + 1
			FROM categories c
			INNER JOIN category_path cp ON c.id = cp.parent_id
		)
		SELECT id, name, slug, parent_id
		FROM category_path
		ORDER BY depth DESC
	`
	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var breadcrumbs []*entity.Category
	for rows.Next() {
		var cat entity.Category
		var parentID *uuid.UUID
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug, &parentID)
		if err != nil {
			return nil, err
		}
		cat.ParentID = parentID
		breadcrumbs = append(breadcrumbs, &cat)
	}
	return breadcrumbs, nil
}

// GetCategoryWithStats возвращает категорию со статистикой
func (r *PostgresCategoryNavigationRepository) GetCategoryWithStats(ctx context.Context, categoryID uuid.UUID) (*entity.CategoryWithStats, error) {
	query := `
		SELECT
			c.*,
			(SELECT COUNT(*) FROM category_assignments ca WHERE ca.category_id = c.id AND ca.target_type = 'course') as courses_count,
			(SELECT COUNT(*) FROM category_assignments ca WHERE ca.category_id = c.id AND ca.target_type = 'lesson') as lessons_count,
			(SELECT COUNT(*) FROM category_assignments ca WHERE ca.category_id = c.id AND ca.target_type = 'module') as modules_count
		FROM categories c
		WHERE c.id = $1
	`
	row := r.db.QueryRow(ctx, query, categoryID)
	var stats entity.CategoryWithStats
	err := row.Scan(
		&stats.ID, &stats.Name, &stats.Slug, &stats.Description, &stats.ParentID,
		&stats.AuthorID, &stats.SortOrder, &stats.IsVisible, &stats.CreatedAt, &stats.UpdatedAt,
		&stats.CoursesCount, &stats.LessonsCount, &stats.ModulesCount,
	)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// SearchInCategory выполняет поиск внутри категории
func (r *PostgresCategoryNavigationRepository) SearchInCategory(ctx context.Context, categoryID uuid.UUID, query string) ([]searchEntity.SearchEntity, error) {
	// Упрощённая реализация: поиск по курсам и урокам в категории
	searchQuery := `
		SELECT
			c.id, c.title, c.description, 'course' as type, c.created_at, c.updated_at, 1.0 as relevance
		FROM courses c
		INNER JOIN category_assignments ca ON ca.target_id = c.id AND ca.target_type = 'course'
		WHERE ca.category_id = $1 AND (c.title ILIKE $2 OR c.description ILIKE $2) AND c.deleted_at IS NULL
		UNION ALL
		SELECT
			l.id, l.title, l.content as description, 'lesson' as type, l.created_at, l.updated_at, 1.0 as relevance
		FROM lessons l
		INNER JOIN category_assignments ca ON ca.target_id = l.id AND ca.target_type = 'lesson'
		WHERE ca.category_id = $1 AND (l.title ILIKE $2 OR l.content ILIKE $2) AND l.deleted_at IS NULL
		ORDER BY relevance DESC, created_at DESC
	`
	rows, err := r.db.Query(ctx, searchQuery, categoryID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []searchEntity.SearchEntity
	for rows.Next() {
		var res searchEntity.SearchEntity
		err := rows.Scan(&res.ID, &res.Title, &res.Description, &res.Type, &res.CreatedAt, &res.UpdatedAt, &res.Relevance)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}

// GetContentByCategory возвращает контент по категории
func (r *PostgresCategoryNavigationRepository) GetContentByCategory(
	ctx context.Context,
	categoryID uuid.UUID,
	contentType string,
) ([]searchEntity.SearchEntity, error) {
	var query string
	switch contentType {
	case "course":
		query = `
			SELECT
				c.id,
				c.title,
				c.description,
				'course' as type,
				c.created_at,
				json_agg(DISTINCT jsonb_build_object(
					'id', t.id,
					'name', t.name
				)) as tags,
				COALESCE(json_agg(DISTINCT jsonb_build_object(
					'id', cat.id,
					'name', cat.name,
					'slug', cat.slug
				)) FILTER (WHERE cat.id IS NOT NULL), '[]') as categories
			FROM courses c
			INNER JOIN category_assignments ca ON ca.target_id = c.id AND ca.target_type = 'course'
			LEFT JOIN tag_assignments ta ON ta.entity_id = c.id AND ta.entity_type = 'course'
			LEFT JOIN tags t ON t.id = ta.tag_id
			LEFT JOIN categories cat ON cat.id = ca.category_id
			WHERE ca.category_id = $1 AND c.deleted_at IS NULL
			GROUP BY c.id, c.title, c.description, c.created_at
			ORDER BY c.created_at DESC
		`
	case "lesson":
		query = `
			SELECT
				l.id,
				l.title,
				l.content as description,
				'lesson' as type,
				l.created_at,
				json_agg(DISTINCT jsonb_build_object(
					'id', t.id,
					'name', t.name
				)) as tags,
				COALESCE(json_agg(DISTINCT jsonb_build_object(
					'id', cat.id,
					'name', cat.name,
					'slug', cat.slug
				)) FILTER (WHERE cat.id IS NOT NULL), '[]') as categories
			FROM lessons l
			INNER JOIN category_assignments ca ON ca.target_id = l.id AND ca.target_type = 'lesson'
			LEFT JOIN tag_assignments ta ON ta.entity_id = l.id AND ta.entity_type = 'lesson'
			LEFT JOIN tags t ON t.id = ta.tag_id
			LEFT JOIN categories cat ON cat.id = ca.category_id
			WHERE ca.category_id = $1 AND l.deleted_at IS NULL
			GROUP BY l.id, l.title, l.content, l.created_at
			ORDER BY l.created_at DESC
		`
	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}
	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []searchEntity.SearchEntity
	for rows.Next() {
		var e searchEntity.SearchEntity
		var tagsJSON, categoriesJSON []byte
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.Type,
			&e.CreatedAt, &tagsJSON, &categoriesJSON,
		)
		if err != nil {
			return nil, err
		}
		e.Tags = parseTagsJSON(tagsJSON)
		e.Categories = parseCategoriesJSON(categoriesJSON)
		entities = append(entities, e)
	}
	return entities, nil
}

// GetSiblingCategories возвращает соседние категории
func (r *PostgresCategoryNavigationRepository) GetSiblingCategories(ctx context.Context, categoryID uuid.UUID) ([]*entity.Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, author_id, sort_order, is_visible, created_at, updated_at
		FROM categories
		WHERE parent_id = (SELECT parent_id FROM categories WHERE id = $1)
		AND id != $1
		AND is_visible = true
		ORDER BY sort_order, name
	`
	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var siblings []*entity.Category
	for rows.Next() {
		var cat entity.Category
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug, &cat.Description, &cat.ParentID, &cat.AuthorID, &cat.SortOrder, &cat.IsVisible, &cat.CreatedAt, &cat.UpdatedAt)
		if err != nil {
			return nil, err
		}
		siblings = append(siblings, &cat)
	}
	return siblings, nil
}

// GetPathCategories возвращает путь категорий от корня
func (r *PostgresCategoryNavigationRepository) GetPathCategories(ctx context.Context, categoryID uuid.UUID) ([]*entity.Category, error) {
	// Аналогично breadcrumbs, но с полными данными
	query := `
		WITH RECURSIVE category_path AS (
			SELECT id, name, slug, description, parent_id, author_id, sort_order, is_visible, created_at, updated_at, 1 as depth
			FROM categories
			WHERE id = $1
			UNION ALL
			SELECT c.id, c.name, c.slug, c.description, c.parent_id, c.author_id, c.sort_order, c.is_visible, c.created_at, c.updated_at, cp.depth + 1
			FROM categories c
			INNER JOIN category_path cp ON c.id = cp.parent_id
		)
		SELECT id, name, slug, description, parent_id, author_id, sort_order, is_visible, created_at, updated_at
		FROM category_path
		ORDER BY depth DESC
	`
	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var path []*entity.Category
	for rows.Next() {
		var cat entity.Category
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug, &cat.Description, &cat.ParentID, &cat.AuthorID, &cat.SortOrder, &cat.IsVisible, &cat.CreatedAt, &cat.UpdatedAt)
		if err != nil {
			return nil, err
		}
		path = append(path, &cat)
	}
	return path, nil
}

func parseTagsJSON(data []byte) []*tagEntity.Tag {
	if len(data) == 0 || string(data) == "[]" {
		return nil
	}
	var tags []struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
	}
	err := json.Unmarshal(data, &tags)
	if err != nil {
		// Логировать ошибку, но вернуть nil для простоты
		return nil
	}
	var result []*tagEntity.Tag
	for _, t := range tags {
		result = append(result, &tagEntity.Tag{
			Name:        t.Name,
			Description: t.Description,
		})
	}
	return result
}

func parseCategoriesJSON(data []byte) []*entity.Category {
	if len(data) == 0 || string(data) == "[]" {
		return nil
	}
	var categories []struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
		Slug string    `json:"slug"`
	}
	err := json.Unmarshal(data, &categories)
	if err != nil {
		// Логировать ошибку, но вернуть nil для простоты
		return nil
	}
	var result []*entity.Category
	for _, c := range categories {
		result = append(result, &entity.Category{
			ID:   c.ID,
			Name: c.Name,
			Slug: c.Slug,
		})
	}
	return result
}
