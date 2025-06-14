package enrollment_test

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/ncostamagna/gocourse_domain/domain"
	"github.com/ncostamagna/gocourse_enrollment/internal/enrollment"
	"github.com/stretchr/testify/assert"
)

func TestService_GetAll(t *testing.T) {

	l := log.New(io.Discard, "", 0)

	t.Run("should return an error", func(t *testing.T) {

		var want string = "my expected error"
		var wantCounter int = 1
		var counter int = 0
		repo := &RepositoryMock{
			GetAllFunc: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				counter++
				return nil, errors.New("my expected error")
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)

		enrollments, err := service.GetAll(context.Background(), enrollment.Filters{}, 0, 10)
		assert.Error(t, err)
		assert.EqualError(t, err, want)
		assert.Nil(t, enrollments)
		assert.Equal(t, wantCounter, counter)
	})

	t.Run("should return all enrollments", func(t *testing.T) {
		want := []domain.Enrollment{
			{
				ID: "1",
				UserID: "1",
				CourseID: "1",
				Status: "active",
			},
		}
		var wantCounter int = 1
		var counter int = 0
		repo := &RepositoryMock{
			GetAllFunc: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				counter++
				return []domain.Enrollment{
					{
						ID: "1",
						UserID: "1",
						CourseID: "1",
						Status: "active",
					},
				}, nil
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)

		enrollments, err := service.GetAll(context.Background(), enrollment.Filters{}, 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, want, enrollments)
		assert.Equal(t, wantCounter, counter)
	})
	
	
}

func TestService_Update(t *testing.T) {

	l := log.New(io.Discard, "", 0)

	t.Run("should return an error", func(t *testing.T) {
		var want string = "my expected error"
		var wantCounter int = 1
		var counter int = 0
		var wantStatus string = "active"
		var wantID string = "1"
		repo := &RepositoryMock{
			UpdateFunc: func(ctx context.Context, id string, status *string) error {
				counter++
				assert.Equal(t, wantID, id)
				assert.NotNil(t, status)
				assert.Equal(t, wantStatus, *status)
				return errors.New("my expected error")
			},
		}
		service := enrollment.NewService(l, nil, nil, repo)

		status := "active"
		err := service.Update(context.Background(), "1", &status)
		
		assert.Error(t, err)
		assert.EqualError(t, err, want)
		assert.Equal(t, wantCounter, counter)
	})

	t.Run("should update an enrollment", func(t *testing.T) {
		var wantCounter int = 1
		var counter int = 0
		var wantStatus string = "active"
		var wantID string = "1"
		repo := &RepositoryMock{
			UpdateFunc: func(ctx context.Context, id string, status *string) error {
				counter++
				assert.Equal(t, wantID, id)
				assert.NotNil(t, status)
				assert.Equal(t, wantStatus, *status)
				return nil
			},
		}
		service := enrollment.NewService(l, nil, nil, repo)

		status := "active"
		err := service.Update(context.Background(), "1", &status)
		assert.NoError(t, err)
		assert.Equal(t, wantCounter, counter)
	})
}

func TestService_Count(t *testing.T) {

	l := log.New(io.Discard, "", 0)

	t.Run("should return an error", func(t *testing.T) {
		var want string = "my expected error"
		var wantCounter int = 1
		var counter int = 0
		repo := &RepositoryMock{
			CountFunc: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				counter++
				return 0, errors.New("my expected error")
			},
		}
		service := enrollment.NewService(l, nil, nil, repo)

		count, err := service.Count(context.Background(), enrollment.Filters{})
		assert.Error(t, err)
		assert.EqualError(t, err, want)
		assert.Zero(t, count)
		assert.Equal(t, wantCounter, counter)
	})

	t.Run("should return the count of enrollments", func(t *testing.T) {
		var want int = 5
		var wantCounter int = 1
		var counter int = 0

		repo := &RepositoryMock{
			CountFunc: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				counter++
				return 5, nil
			},
		}
		service := enrollment.NewService(l, nil, nil, repo)

		count, err := service.Count(context.Background(), enrollment.Filters{})
		assert.NoError(t, err)
		assert.Equal(t, want, count)
		assert.Equal(t, wantCounter, counter)
	})
}
