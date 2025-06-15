package test

import (
	"net/http"
	"testing"

	"github.com/ncostamagna/gocourse_domain/domain"
	"github.com/ncostamagna/gocourse_enrollment/internal/enrollment"
	"github.com/stretchr/testify/assert"
)

type dataResponse struct {
	Message string      `json:"message"`
	Status    int         `json:"status"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}

func TestEnrollments(t *testing.T) {

	t.Run("create an enrollment and get it", func(t *testing.T) {
		bodyRequest := enrollment.CreateReq{
			UserID:   "11",
			CourseID: "22",
		}
		resp := cli.Post("/enrollments", bodyRequest)
		
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		dataCreated := domain.Enrollment{}
		dRespCreated := dataResponse{Data: &dataCreated}

		err := resp.FillUp(&dRespCreated)
		assert.Nil(t, err)

		assert.Equal(t, "success", dRespCreated.Message)
		assert.Equal(t, http.StatusCreated, dRespCreated.Status)

		assert.NotEmpty(t, dataCreated.ID)
		assert.Equal(t, "11", dataCreated.UserID)
		assert.Equal(t, "22", dataCreated.CourseID)


		resp = cli.Get("/enrollments?user_id=" + dataCreated.UserID + "&course_id=" + dataCreated.CourseID)

		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var dataGetAll []domain.Enrollment
		dRespGetAll := dataResponse{Data: &dataGetAll}

		err = resp.FillUp(&dRespGetAll)
		assert.Nil(t, err)
		assert.Equal(t, "success", dRespGetAll.Message)
		assert.Equal(t, http.StatusOK, dRespGetAll.Status)
		assert.Equal(t, 1, len(dataGetAll))

		assert.Equal(t, dataCreated.ID, dataGetAll[0].ID)
		assert.Equal(t, dataCreated.UserID, dataGetAll[0].UserID)
		assert.Equal(t, dataCreated.CourseID, dataGetAll[0].CourseID)
		assert.Equal(t, "P", dataGetAll[0].Status)
	})

	t.Run("update an enrollment", func(t *testing.T) {
		bodyRequest := enrollment.CreateReq{
			UserID:   "22",
			CourseID: "33",
		}
		resp := cli.Post("/enrollments", bodyRequest)
		
		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		dataCreated := domain.Enrollment{}
		dRespCreated := dataResponse{Data: &dataCreated}

		err := resp.FillUp(&dRespCreated)
		assert.Nil(t, err)

		status := "A"
		resp = cli.Patch("/enrollments/" + dataCreated.ID, enrollment.UpdateReq{
			Status: &status,
		})

		resp = cli.Get("/enrollments?user_id=" + dataCreated.UserID + "&course_id=" + dataCreated.CourseID)

		assert.Nil(t, resp.Err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var dataGetAll []domain.Enrollment
		dRespGetAll := dataResponse{Data: &dataGetAll}

		err = resp.FillUp(&dRespGetAll)
		assert.Nil(t, err)
		assert.Equal(t, "success", dRespGetAll.Message)
		assert.Equal(t, http.StatusOK, dRespGetAll.Status)
		assert.Equal(t, 1, len(dataGetAll))

		assert.Equal(t, dataCreated.ID, dataGetAll[0].ID)
		assert.Equal(t, dataCreated.UserID, dataGetAll[0].UserID)
		assert.Equal(t, dataCreated.CourseID, dataGetAll[0].CourseID)
		assert.Equal(t, "A", dataGetAll[0].Status)
	})
}