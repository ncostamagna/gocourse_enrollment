package enrollment_test

import (
	"context"

	"github.com/ncostamagna/gocourse_domain/domain"
)

type RepositoryMock struct {
	CreateFunc func(ctx context.Context, enrollment *domain.Enrollment) error
	GetAllFunc func(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
	UpdateFunc func(ctx context.Context, id string, status *string) error
	CountFunc  func(ctx context.Context, filters Filters) (int, error)
}

func (r *RepositoryMock) Create(ctx context.Context, enrollment *enrollment.Enrollment) error {
	return r.CreateFunc(ctx, enrollment)
}

func (r *RepositoryMock) GetAll(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]enrollment.Enrollment, error) {
	return r.GetAllFunc(ctx, filters, offset, limit)
}

func (r *RepositoryMock) Update(ctx context.Context, id string, status *string) error {
	return r.UpdateFunc(ctx, id, status)
}

func (r *RepositoryMock) Count(ctx context.Context, filters enrollment.Filters) (int, error) {
	return r.CountFunc(ctx, filters)
}