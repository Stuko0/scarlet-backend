package entity

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	entityv1 "github.com/Stuko0/scarlet-backend/gen/proto/entity/v1"
)

type EntityService struct {
	repo EntityRepositoryInterface
}

func entitieToProto(entitie *Entity) *entityv1.Entity {
	return &entityv1.Entity{
		EntityId:   int64(entitie.EntityId),
		Name:        entitie.Name,
		EntityType: entitie.EntityType,
		Address:     entitie.Address,
		Phone:       entitie.Phone,
		Email:       entitie.Email,
		Image:       entitie.ImageUrl,
		Active:      entitie.Active,
		CreatedAt:   entitie.CreatedAt,
		UpdatedAt:   entitie.UpdatedAt,
	}
}

func NewEntityService(repo EntityRepositoryInterface) *EntityService {
	return &EntityService{
		repo: repo,
	}
}

func (s *EntityService) CreateEntity(ctx context.Context, req *connect.Request[entityv1.CreateEntityRequest]) (*connect.Response[entityv1.EntityResponse], error) {
	if req.Msg.Entity.Name == "" || req.Msg.Entity.EntityType == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, connect.NewError(connect.CodeInvalidArgument, errors.New("name and type are required")))
	}

	entitieModel := &Entity{
		Name:       req.Msg.Entity.Name,
		EntityType: req.Msg.Entity.EntityType,
		Address:    req.Msg.Entity.Address,
		Phone:      req.Msg.Entity.Phone,
		Email:      req.Msg.Entity.Email,
		ImageUrl:      req.Msg.Entity.Image,
		Active:     req.Msg.Entity.Active,
		CreatedAt:  req.Msg.Entity.CreatedAt,
		UpdatedAt:  req.Msg.Entity.UpdatedAt,
	}

	createdEntity, err := s.repo.CreateEntity(ctx, entitieModel)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := connect.NewResponse(&entityv1.EntityResponse{
		Entity: entitieToProto(createdEntity),
	})

	return res, nil
}

func (s *EntityService) GetEntity(ctx context.Context, req *connect.Request[entityv1.GetEntityRequest]) (*connect.Response[entityv1.EntityResponse], error) {
	if req.Msg.EntityId <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("entity_id must be greater than 0"))
	}

	entity, err := s.repo.GetEntity(ctx, req.Msg.EntityId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&entityv1.EntityResponse{
		Entity: entitieToProto(entity),
	})

	return res, nil
}

func (s *EntityService) UpdateEntity(ctx context.Context, req *connect.Request[entityv1.UpdateEntityRequest]) (*connect.Response[entityv1.EntityResponse], error) {
	if req.Msg.Entity.EntityId <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("entity_id must be greater than 0"))
	}

	entitieModel := &Entity{
		EntityId:   req.Msg.Entity.EntityId,
		Name:       req.Msg.Entity.Name,
		EntityType: req.Msg.Entity.EntityType,
		Address:    req.Msg.Entity.Address,
		Phone:      req.Msg.Entity.Phone,
		Email:      req.Msg.Entity.Email,
		ImageUrl:      req.Msg.Entity.Image,
		Active:     req.Msg.Entity.Active,
		CreatedAt:  req.Msg.Entity.CreatedAt,
		UpdatedAt:  req.Msg.Entity.UpdatedAt,
	}

	updatedEntity, err := s.repo.UpdateEntity(ctx, entitieModel)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&entityv1.EntityResponse{
		Entity: entitieToProto(updatedEntity),
	})

	return res, nil
}

func (s *EntityService) DeleteEntity (ctx context.Context, req *connect.Request[entityv1.DeleteEntityRequest]) (*connect.Response[entityv1.DeleteEntityResponse], error) {
	if req.Msg.EntityId <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("entity_id must be greater than 0"))
	}

	err := s.repo.DeleteEntity(ctx, req.Msg.EntityId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&entityv1.DeleteEntityResponse{
		Message: "Entity deleted successfully",
	})

	return res, nil
}