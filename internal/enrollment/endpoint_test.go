package enrollment_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ncostamagna/go_lib_response/response"
	"github.com/stretchr/testify/assert"

	"github.com/ncostamagna/gocourse_enrollment/internal/enrollment"
)

func TestCreateEndpoint(t *testing.T) {

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
}