package entity

import (
	"context"
	"errors"
	"fmt"
	"github.com/Stuko0/scarlet-backend/internal/database"
	"github.com/jackc/pgx/v5"
)

type EntityRepositoryInterface interface {
	CreateEntity(ctx context.Context, entitie *Entity) (*Entity, error)
	GetEntity(ctx context.Context, entityId int64) (*Entity, error)
	UpdateEntity(ctx context.Context, entity *Entity) (*Entity, error)
	DeleteEntity(ctx context.Context, entityId int64) error
}

type EntitieRepository struct {
	db *database.Postgres
}

func NewEntitieRepository(db *database.Postgres) *EntitieRepository {
	return &EntitieRepository{db: db}
}

func (r *EntitieRepository) CreateEntity(ctx context.Context, entity *Entity) (*Entity, error) {
	query := `INSERT INTO scarlet.Entitys (name, entity_type, address, phone, email, image) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING entity_id`
	var createdEntity Entity
	err := r.db.Pool.QueryRow(ctx, query,
		entity.Name,
		entity.EntityType,
		entity.Address,
		entity.Phone,
		entity.Email,
		entity.ImageUrl,).Scan(&createdEntity.EntityId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no rows found: %w", err)
		}
		return nil, err
	}
	return &createdEntity, nil
}

func (r *EntitieRepository) GetEntity(ctx context.Context, entityId int64) (*Entity, error) {
	query := `SELECT entity_id, name, entity_type, address, phone, email, image, created_at, updated_at, active 
	FROM scarlet.Entitys WHERE entity_id = $1`
	var entity Entity
	err := r.db.Pool.QueryRow(ctx, query, entityId).Scan(
		&entity.EntityId,
		&entity.Name,
		&entity.EntityType,
		&entity.Address,
		&entity.Phone,
		&entity.Email,
		&entity.ImageUrl,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.Active,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no entity found with id %d: %w", entityId, err)
		}
		return nil, err
	}
	return &entity, nil
}

func (r *EntitieRepository) UpdateEntity(ctx context.Context, entity *Entity) (*Entity, error) {
	query := `UPDATE scarlet.Entitys SET name = $1, entity_type = $2, address = $3, phone = $4, email = $5, image = $6, 
	updated_at = NOW() WHERE entity_id = $7 RETURNING entity_id`
	var updatedEntity Entity
	err := r.db.Pool.QueryRow(ctx, query,
		entity.Name,
		entity.EntityType,
		entity.Address,
		entity.Phone,
		entity.Email,
		entity.ImageUrl,
		entity.EntityId).Scan(&updatedEntity.EntityId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no entity found with id %d: %w", entity.EntityId, err)
		}
		return nil, err
	}
	updatedEntity.Name = entity.Name
	updatedEntity.EntityType = entity.EntityType
	updatedEntity.Address = entity.Address
	updatedEntity.Phone = entity.Phone
	updatedEntity.Email = entity.Email
	updatedEntity.ImageUrl = entity.ImageUrl
	updatedEntity.CreatedAt = entity.CreatedAt
	updatedEntity.UpdatedAt = entity.UpdatedAt
	updatedEntity.Active = entity.Active
	return &updatedEntity, nil
}

func (r *EntitieRepository) DeleteEntity(ctx context.Context, entityId int64) error {
	query := `DELETE FROM scarlet.Entitys WHERE entity_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, entityId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("no entity found with id %d: %w", entityId, err)
		}
		return err
	}
	return nil
}