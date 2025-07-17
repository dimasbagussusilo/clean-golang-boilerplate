package grpcs

import (
	"appsku-golang/app/models"
	"appsku-golang/app/usecases"
	pb "appsku-golang/files/grpc-protos"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type StoreGrpc struct {
	pb.StoreServiceServer
	StoreUseCase usecases.IStoreUseCase
}

func NewStoreGrpc(StoreUseCase usecases.IStoreUseCase) *StoreGrpc {
	return &StoreGrpc{
		StoreUseCase: StoreUseCase,
	}
}

func (g *StoreGrpc) GetStore(ctx context.Context, req *pb.GetStoreByIDRequest) (*pb.GetStoreResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &pb.GetStoreResponse{
			StatusCode: 400,
			Error: &pb.StoreErrorResponse{
				Message:    "Invalid ID format",
				StatusCode: 400,
			},
		}, nil
	}

	store, err := g.StoreUseCase.GetById(ctx, objectID)
	if err != nil {
		return &pb.GetStoreResponse{
			StatusCode: 500,
			Error: &pb.StoreErrorResponse{
				Message:    "Failed to get store",
				StatusCode: 500,
			},
		}, nil
	}

	if store == nil {
		return &pb.GetStoreResponse{
			StatusCode: 404,
			Error: &pb.StoreErrorResponse{
				Message:    "Store not found",
				StatusCode: 404,
			},
		}, nil
	}

	var createdAt, updatedAt, deletedAt string
	if store.CreatedAt != nil {
		createdAt = store.CreatedAt.Format(time.RFC3339)
	}
	if store.UpdatedAt != nil {
		updatedAt = store.UpdatedAt.Format(time.RFC3339)
	}
	if store.DeletedAt != nil {
		deletedAt = store.DeletedAt.Format(time.RFC3339)
	}

	return &pb.GetStoreResponse{
		StatusCode: 200,
		Data: &pb.StoreResponse{
			Id:          store.ID.Hex(),
			Name:        store.Name,
			Description: store.Description,
			Type:        store.Type,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			DeletedAt:   deletedAt,
		},
	}, nil
}

func (g *StoreGrpc) ListStores(ctx context.Context, req *pb.ListStoresRequest) (*pb.ListStoresResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	stores, total, err := g.StoreUseCase.GetAll(ctx, map[string]interface{}{}, page, limit)
	if err != nil {
		return &pb.ListStoresResponse{
			StatusCode: 500,
			Error: &pb.StoreErrorResponse{
				Message:    "Failed to list stores",
				StatusCode: 500,
			},
		}, nil
	}

	var storeResponses []*pb.StoreResponse
	for _, store := range stores {
		var createdAt, updatedAt, deletedAt string
		if store.CreatedAt != nil {
			createdAt = store.CreatedAt.Format(time.RFC3339)
		}
		if store.UpdatedAt != nil {
			updatedAt = store.UpdatedAt.Format(time.RFC3339)
		}
		if store.DeletedAt != nil {
			deletedAt = store.DeletedAt.Format(time.RFC3339)
		}

		storeResponses = append(storeResponses, &pb.StoreResponse{
			Id:          store.ID.Hex(),
			Name:        store.Name,
			Description: store.Description,
			Type:        store.Type,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			DeletedAt:   deletedAt,
		})
	}

	return &pb.ListStoresResponse{
		StatusCode: 200,
		Data:       storeResponses,
		Total:      int32(total),
	}, nil
}

func (g *StoreGrpc) CreateStore(ctx context.Context, req *pb.CreateStoreRequest) (*pb.CreateStoreResponse, error) {
	store := &models.Store{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}

	createdStore, err := g.StoreUseCase.Insert(ctx, store)
	if err != nil {
		return &pb.CreateStoreResponse{
			StatusCode: 500,
			Error: &pb.StoreErrorResponse{
				Message:    "Failed to create store",
				StatusCode: 500,
			},
		}, nil
	}

	var createdAt, updatedAt, deletedAt string
	if createdStore.CreatedAt != nil {
		createdAt = createdStore.CreatedAt.Format(time.RFC3339)
	}
	if createdStore.UpdatedAt != nil {
		updatedAt = createdStore.UpdatedAt.Format(time.RFC3339)
	}
	if createdStore.DeletedAt != nil {
		deletedAt = createdStore.DeletedAt.Format(time.RFC3339)
	}

	return &pb.CreateStoreResponse{
		StatusCode: 201,
		Data: &pb.StoreResponse{
			Id:          createdStore.ID.Hex(),
			Name:        createdStore.Name,
			Description: createdStore.Description,
			Type:        createdStore.Type,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			DeletedAt:   deletedAt,
		},
	}, nil
}

func (g *StoreGrpc) UpdateStore(ctx context.Context, req *pb.UpdateStoreRequest) (*pb.UpdateStoreResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &pb.UpdateStoreResponse{
			StatusCode: 400,
			Error: &pb.StoreErrorResponse{
				Message:    "Invalid ID format",
				StatusCode: 400,
			},
		}, nil
	}

	existingStore, err := g.StoreUseCase.GetById(ctx, objectID)
	if err != nil {
		return &pb.UpdateStoreResponse{
			StatusCode: 500,
			Error: &pb.StoreErrorResponse{
				Message:    "Failed to get store",
				StatusCode: 500,
			},
		}, nil
	}

	if existingStore == nil {
		return &pb.UpdateStoreResponse{
			StatusCode: 404,
			Error: &pb.StoreErrorResponse{
				Message:    "Store not found",
				StatusCode: 404,
			},
		}, nil
	}

	now := time.Now()
	updateData := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"type":        req.Type,
		"updated_at":  now,
	}

	err = g.StoreUseCase.Update(ctx, objectID, updateData)
	if err != nil {
		return &pb.UpdateStoreResponse{
			StatusCode: 500,
			Error: &pb.StoreErrorResponse{
				Message:    "Failed to update store",
				StatusCode: 500,
			},
		}, nil
	}

	updatedStore, err := g.StoreUseCase.GetById(ctx, objectID)
	if err != nil {
		return &pb.UpdateStoreResponse{
			StatusCode: 500,
			Error: &pb.StoreErrorResponse{
				Message:    "Failed to get updated store",
				StatusCode: 500,
			},
		}, nil
	}

	var createdAt, updatedAt, deletedAt string
	if updatedStore.CreatedAt != nil {
		createdAt = updatedStore.CreatedAt.Format(time.RFC3339)
	}
	if updatedStore.UpdatedAt != nil {
		updatedAt = updatedStore.UpdatedAt.Format(time.RFC3339)
	}
	if updatedStore.DeletedAt != nil {
		deletedAt = updatedStore.DeletedAt.Format(time.RFC3339)
	}

	return &pb.UpdateStoreResponse{
		StatusCode: 200,
		Data: &pb.StoreResponse{
			Id:          updatedStore.ID.Hex(),
			Name:        updatedStore.Name,
			Description: updatedStore.Description,
			Type:        updatedStore.Type,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			DeletedAt:   deletedAt,
		},
	}, nil
}

func (g *StoreGrpc) DeleteStore(ctx context.Context, req *pb.DeleteStoreRequest) (*pb.DeleteStoreResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return &pb.DeleteStoreResponse{
			StatusCode: 400,
			Error: &pb.StoreErrorResponse{
				Message:    "Invalid ID format",
				StatusCode: 400,
			},
		}, nil
	}

	existingStore, err := g.StoreUseCase.GetById(ctx, objectID)
	if err != nil {
		return &pb.DeleteStoreResponse{
			StatusCode: 500,
			Error: &pb.StoreErrorResponse{
				Message:    "Failed to get store",
				StatusCode: 500,
			},
		}, nil
	}

	if existingStore == nil {
		return &pb.DeleteStoreResponse{
			StatusCode: 404,
			Error: &pb.StoreErrorResponse{
				Message:    "Store not found",
				StatusCode: 404,
			},
		}, nil
	}

	// Soft delete the store
	err = g.StoreUseCase.Delete(ctx, objectID, false)
	if err != nil {
		return &pb.DeleteStoreResponse{
			StatusCode: 500,
			Error: &pb.StoreErrorResponse{
				Message:    "Failed to delete store",
				StatusCode: 500,
			},
		}, nil
	}

	return &pb.DeleteStoreResponse{
		StatusCode: 200,
	}, nil
}
