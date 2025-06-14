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

		var want error = errors.New("my expected error")
		var wantCounter int = 1
		var counter int = 0
		repo := &RepositoryMock{
			GetAllFunc: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				counter++
				return nil, want
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)

		enrollments, err := service.GetAll(context.Background(), enrollment.Filters{}, 0, 10)
		assert.ErrorIs(t, err, want)
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
				return want, nil
			},
		}

		service := enrollment.NewService(l, nil, nil, repo)

		enrollments, err := service.GetAll(context.Background(), enrollment.Filters{}, 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, want, enrollments)
		assert.Equal(t, wantCounter, counter)
	})
	
	
}