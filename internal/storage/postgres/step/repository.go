package steps

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
)

type StepRepositorer interface {
	FindMany(ctx context.Context, filter models.StepFilter) ([]models.Step, error)
	Create(ctx context.Context, model *models.Step) error
	Update(ctx context.Context, id int, model models.Step) error
	FindByID(ctx context.Context, id int) (models.Step, error)
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) (models.Step, error)
	Count(ctx context.Context, filter models.StepFilter) (int, error)
}

type stepRepository struct {
	client postgres.Client
}

func (r *stepRepository) Count(ctx context.Context, filter models.StepFilter) (int, error) {
	var (
		query = `
			SELECT COUNT(*) 
			FROM steps 
			WHERE deleted_at IS NULL AND
				name ILIKE $1
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query, "%"+filter.Name+"%").Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count steps: %w", err)
	}
	return count, nil
}

// Create implements StepRepositorer
func (r *stepRepository) Create(ctx context.Context, model *models.Step) error {
	query := `
		INSERT INTO steps (name)
		VALUES ($1)
		RETURNING id
	`
	if err := r.client.QueryRow(ctx, query, model.Name).Scan(&model.ID); err != nil {
		return fmt.Errorf("failed to create step: %w", err)
	}
	return nil
}

// Delete implements StepRepositorer
func (r *stepRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE steps
		SET deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
	`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete step: %w", err)
	}
	return nil
}

// FindByID implements StepRepositorer
func (r *stepRepository) FindByID(ctx context.Context, id int) (models.Step, error) {
	var (
		query = `
			SELECT id, name
			FROM steps
			WHERE id = $1
		`
		model models.Step
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&model.ID, &model.Name); err != nil {
		return models.Step{}, fmt.Errorf("failed to find step: %w", err)
	}
	return model, nil
}

// FindMany implements StepRepositorer
func (r *stepRepository) FindMany(
	ctx context.Context,
	filter models.StepFilter,
) ([]models.Step, error) {
	var (
		query = `
			SELECT id, name
			FROM steps
			WHERE deleted_at IS NULL AND
				name ILIKE $3
			LIMIT $1
			OFFSET $2
		`
		entities []models.Step
	)
	rows, err := r.client.Query(ctx, query, filter.Limit, filter.Offset, "%"+filter.Name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find steps: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var model models.Step
		if err := rows.Scan(&model.ID, &model.Name); err != nil {
			return nil, fmt.Errorf("failed to scan step: %w", err)
		}
		entities = append(entities, model)
	}
	return entities, nil
}

// Restore implements StepRepositorer
func (r *stepRepository) Restore(ctx context.Context, id int) (models.Step, error) {
	var (
		query = `
			UPDATE steps
			SET deleted_at = NULL,
				updated_at = NOW()
			WHERE id = $1
			RETURNING id, name
		`
		model models.Step
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&model.ID, &model.Name); err != nil {
		return models.Step{}, fmt.Errorf("failed to restore step: %w", err)
	}
	return model, nil
}

// Update implements StepRepositorer
func (r *stepRepository) Update(ctx context.Context, id int, model models.Step) error {
	query := `
		UPDATE steps
		SET name = $1,
			updated_at = NOW()
		WHERE id = $2
		RETURNING id, name
	`
	if _, err := r.client.Exec(ctx, query, model.Name, id); err != nil {
		return fmt.Errorf("failed to update step: %w", err)
	}
	return nil
}

func New(client postgres.Client) StepRepositorer {
	return &stepRepository{
		client: client,
	}
}
