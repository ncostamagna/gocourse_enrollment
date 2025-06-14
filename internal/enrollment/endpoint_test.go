package enrollment_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ncostamagna/go_lib_response/response"
	"github.com/ncostamagna/gocourse_domain/domain"
	"github.com/ncostamagna/gocourse_enrollment/internal/enrollment"
	"github.com/stretchr/testify/assert"

	courseSdk "github.com/ncostamagna/go_course_sdk/course"
	userSdk "github.com/ncostamagna/go_course_sdk/user"

	"errors"
	"io"
	"log"
)

func TestCreateEndpoint(t *testing.T) {

	l := log.New(io.Discard, "", 0)

	t.Run("should return an error if user id is required", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		assert.EqualError(t, enrollment.ErrUserIDRequired, resp.Error())
	})

	t.Run("should return an error if course id is required", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{UserID: "1"})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		assert.EqualError(t, enrollment.ErrCourseIDRequired, resp.Error())
	})

	obj := []struct {
		tag            string
		repositoryMock enrollment.Repository
		userSdkMock    userSdk.Transport
		courseSdkMock  courseSdk.Transport
		wantErr        error
		wantCode       int
		wantResponse   *domain.Enrollment
	}{
		{
			tag: "should return an error if user sdk returns an unexpected error",
			userSdkMock: &UserSdkMock{
				GetFunc: func(id string) (*domain.User, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if user does not exist",
			userSdkMock: &UserSdkMock{
				GetFunc: func(id string) (*domain.User, error) {
					return nil, userSdk.ErrNotFound{Message: "user not found"}
				},
			},
			wantErr:  userSdk.ErrNotFound{Message: "user not found"},
			wantCode: http.StatusNotFound,
		},
		{
			tag: "should return an error if course sdk returns an unexpected error",
			userSdkMock: &UserSdkMock{
				GetFunc: func(id string) (*domain.User, error) {
					return &domain.User{}, nil
				},
			},
			courseSdkMock: &CourseSdkMock{
				GetFunc: func(id string) (*domain.Course, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if course does not exist",
			userSdkMock: &UserSdkMock{
				GetFunc: func(id string) (*domain.User, error) {
					return &domain.User{}, nil
				},
			},
			courseSdkMock: &CourseSdkMock{
				GetFunc: func(id string) (*domain.Course, error) {
					return nil, courseSdk.ErrNotFound{Message: "course not found"}
				},
			},
			wantErr:  courseSdk.ErrNotFound{Message: "course not found"},
			wantCode: http.StatusNotFound,
		},
		{
			tag: "should return an error if repository returns an unexpected error",
			userSdkMock: &UserSdkMock{
				GetFunc: func(id string) (*domain.User, error) {
					return &domain.User{}, nil
				},
			},
			courseSdkMock: &CourseSdkMock{
				GetFunc: func(id string) (*domain.Course, error) {
					return &domain.Course{}, nil
				},
			},
			repositoryMock: &RepositoryMock{
				CreateFunc: func(ctx context.Context, enrollment *domain.Enrollment) error {
					return errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should record the enrollment created",
			userSdkMock: &UserSdkMock{
				GetFunc: func(id string) (*domain.User, error) {
					return &domain.User{}, nil
				},
			},
			courseSdkMock: &CourseSdkMock{
				GetFunc: func(id string) (*domain.Course, error) {
					return &domain.Course{}, nil
				},
			},
			repositoryMock: &RepositoryMock{
				CreateFunc: func(ctx context.Context, enrollment *domain.Enrollment) error {
					enrollment.ID = "101021"
					return nil
				},
			},
			wantCode: http.StatusCreated,
			wantResponse: &domain.Enrollment{
				ID:       "101021",
				UserID:   "1",
				CourseID: "4",
				Status:   "P",
			},
		},
	}

	for _, obj := range obj {
		t.Run(obj.tag, func(t *testing.T) {
			service := enrollment.NewService(l, obj.userSdkMock, obj.courseSdkMock, obj.repositoryMock)
			endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
			resp, err := endpoint.Create(context.Background(), enrollment.CreateReq{UserID: "1", CourseID: "4"})
			if obj.wantErr != nil {
				assert.NotNil(t, err)
				assert.Nil(t, resp)

				r := err.(response.Response)
				assert.EqualError(t, obj.wantErr, r.Error())
				assert.Equal(t, obj.wantCode, r.StatusCode())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, resp)

				r := resp.(response.Response)
				assert.Equal(t, obj.wantCode, r.StatusCode())
				assert.Empty(t, r.Error())

				enrollment := r.GetData().(*domain.Enrollment)
				assert.Equal(t, obj.wantResponse.ID, enrollment.ID)
				assert.Equal(t, obj.wantResponse.UserID, enrollment.UserID)
				assert.Equal(t, obj.wantResponse.CourseID, enrollment.CourseID)
				assert.Equal(t, obj.wantResponse.Status, enrollment.Status)
			}
		})
	}

}

func TestGetAllEndpoint(t *testing.T) {

	l := log.New(io.Discard, "", 0)

	t.Run("should return an error if Count repository returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, nil, nil, &RepositoryMock{
			CountFunc: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 0, errors.New("unexpected error")
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
		assert.EqualError(t, wantErr, resp.Error())
	})

	t.Run("should return an error if Meta returns an unexpected error", func(t *testing.T) {
		wantErr := "strconv.Atoi: parsing \"invalid number\": invalid syntax"
		service := enrollment.NewService(l, nil, nil, &RepositoryMock{
			CountFunc: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 3, nil
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimPageDef: "invalid number"})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
		assert.Equal(t, wantErr, resp.Error())
	})

	t.Run("should return an error if GetAll repository returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, nil, nil, &RepositoryMock{
			CountFunc: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 3, nil
			},
			GetAllFunc: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				return nil, errors.New("unexpected error")
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimPageDef: "10"})
		_, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
		assert.EqualError(t, wantErr, resp.Error())
	})

	t.Run("should return the enrollments", func(t *testing.T) {
		wantEnrollments := []domain.Enrollment{
			{ID: "1", UserID: "11", CourseID: "111", Status: "P"},
			{ID: "2", UserID: "22", CourseID: "222", Status: "P"},
			{ID: "3", UserID: "33", CourseID: "333", Status: "P"},
		}
		service := enrollment.NewService(l, nil, nil, &RepositoryMock{
			CountFunc: func(ctx context.Context, filters enrollment.Filters) (int, error) {
				return 3, nil
			},
			GetAllFunc: func(ctx context.Context, filters enrollment.Filters, offset, limit int) ([]domain.Enrollment, error) {
				return []domain.Enrollment{
					{ID: "1", UserID: "11", CourseID: "111", Status: "P"},
					{ID: "2", UserID: "22", CourseID: "222", Status: "P"},
					{ID: "3", UserID: "33", CourseID: "333", Status: "P"},
				}, nil
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{LimPageDef: "10"})
		resp, err := endpoint.GetAll(context.Background(), enrollment.GetAllReq{})
		assert.Nil(t, err)

		r := resp.(response.Response)
		assert.Equal(t, http.StatusOK, r.StatusCode())
		assert.Empty(t, r.Error())

		enrollments := r.GetData().([]domain.Enrollment)
		assert.Equal(t, len(wantEnrollments), len(enrollments))
		for i := range enrollments {
			assert.Equal(t, wantEnrollments[i].ID, enrollments[i].ID)
			assert.Equal(t, wantEnrollments[i].UserID, enrollments[i].UserID)
			assert.Equal(t, wantEnrollments[i].CourseID, enrollments[i].CourseID)
			assert.Equal(t, wantEnrollments[i].Status, enrollments[i].Status)
		}
	})
}

func TestUpdateEndpoint(t *testing.T) {

	l := log.New(io.Discard, "", 0)

	t.Run("should return an error if status is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		status := ""
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		assert.EqualError(t, enrollment.ErrStatusRequired, resp.Error())
	})

	t.Run("should return an error if repository returns a not found error", func(t *testing.T) {
		service := enrollment.NewService(l, nil, nil, &RepositoryMock{
			UpdateFunc: func(ctx context.Context, id string, status *string) error {
				return enrollment.ErrNotFound{EnrollmentsID: id}
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "5"
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode())
		assert.EqualError(t, enrollment.ErrNotFound{EnrollmentsID: "20"}, resp.Error())
	})

	t.Run("should return an error if repository returns an unexpected error", func(t *testing.T) {
		wantErr := errors.New("unexpected error")
		service := enrollment.NewService(l, nil, nil, &RepositoryMock{
			UpdateFunc: func(ctx context.Context, id string, status *string) error {
				return errors.New("unexpected error")
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "5"
		_, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Error(t, err)

		resp := err.(response.Response)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
		assert.EqualError(t, wantErr, resp.Error())
	})

	t.Run("should return success", func(t *testing.T) {
		service := enrollment.NewService(l, nil, nil, &RepositoryMock{
			UpdateFunc: func(ctx context.Context, id string, status *string) error {
				assert.Equal(t, "20", id)
				assert.Equal(t, "5", *status)
				return nil
			},
		})
		endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
		status := "5"
		resp, err := endpoint.Update(context.Background(), enrollment.UpdateReq{ID: "20", Status: &status})
		assert.Nil(t, err)

		r := resp.(response.Response)
		assert.Equal(t, http.StatusOK, r.StatusCode())
		assert.Nil(t, r.GetData())
	})
}