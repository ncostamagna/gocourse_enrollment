package enrollment

import (
	"context"
	"errors"

	"github.com/ncostamagna/go_lib_response/response"
	"github.com/ncostamagna/gocourse_meta/meta"
)

//Endpoints struct
type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Update Controller
	}

	CreateReq struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
	}

	GetAllReq struct {
		UserID   string
		CourseID string
		Limit    int
		Page     int
	}

	UpdateReq struct {
		ID     string
		Status *string `json:"status"`
	}

	Config struct {
		LimPageDef string
	}
)

//MakeEndpoints handler endpoints
func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		if req.UserID == "" {
			return nil, response.BadRequest(ErrUserIDRequired.Error())
		}

		if req.CourseID == "" {
			return nil, response.BadRequest(ErrCourseIDRequired.Error())
		}

		enroll, err := s.Create(ctx, req.UserID, req.CourseID)
		if err != nil {

			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", enroll, nil), nil

	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			UserID:   req.UserID,
			CourseID: req.CourseID,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		enrollments, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", enrollments, meta), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.Status != nil && *req.Status == "" {
			return nil, response.BadRequest(ErrStatusRequired.Error())
		}

		if err := s.Update(ctx, req.ID, req.Status); err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}
