package enrollment_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ncostamagna/gocourse_enrollment/internal/enrollment"
	"github.com/ncostamagna/gocourse_domain/domain"
	"github.com/ncostamagna/go_lib_response/response"
	"github.com/stretchr/testify/assert"

	userSdk "github.com/ncostamagna/go_course_sdk/user"
	courseSdk "github.com/ncostamagna/go_course_sdk/course"

	"log"
	"io"
	"errors"

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
		tag string
		repositoryMock enrollment.Repository
		userSdkMock userSdk.Transport
		courseSdkMock courseSdk.Transport
		wantErr error
		wantCode int
		wantResponse interface{}
	}{
		{
			tag: "should return an error if user sdk returns an unexpected error",
			userSdkMock: &UserSdkMock{
				GetFunc: func(id string) (*domain.User, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr: errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return an error if user does not exist",
			userSdkMock: &UserSdkMock{
				GetFunc: func(id string) (*domain.User, error) {
					return nil, userSdk.ErrNotFound{ Message: "user not found" }
				},
			},
			wantErr: userSdk.ErrNotFound{ Message: "user not found" },
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
			wantErr: errors.New("unexpected error"),
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
					return nil, courseSdk.ErrNotFound{ Message: "course not found" }
				},
			},
			wantErr: courseSdk.ErrNotFound{ Message: "course not found" },
			wantCode: http.StatusNotFound,
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
				
			}
		})
	}

}