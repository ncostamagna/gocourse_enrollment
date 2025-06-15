package test

import (

	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"io"

	"github.com/joho/godotenv"

	"github.com/ncostamagna/go_http_client/client"

	"testing"
	"github.com/ncostamagna/gocourse_domain/domain"
		"github.com/ncostamagna/gocourse_enrollment/internal/enrollment"
		"github.com/ncostamagna/gocourse_enrollment/pkg/handler"
		"github.com/ncostamagna/gocourse_enrollment/pkg/bootstrap"

)

var cli client.Transport

func TestMain(m *testing.M) {

	_ = godotenv.Load("../.env")
	l := log.New(io.Discard, "", 0)
	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}
	tx := db.Begin()

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	courseTrans := &CourseSdkMock{
		GetFunc: func(id string) (*domain.Course, error) {
			return &domain.Course{ID: id, Name: "Course " + id}, nil
		},
	}
	userTrans := &UserSdkMock{
		GetFunc: func(id string) (*domain.User, error) {
			return &domain.User{ID: id, FirstName: "User " + id, LastName: "Last " + id}, nil
		},
	}

	ctx := context.Background()
	enrollRepo := enrollment.NewRepo(tx, l)
	enrollSrv := enrollment.NewService(l, userTrans, courseTrans, enrollRepo)
	h := handler.NewEnrollmentHTTPServer(ctx, enrollment.MakeEndpoints(enrollSrv, enrollment.Config{LimPageDef: pagLimDef}))
	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	cli = client.New(nil, "http://"+address, 0, false)

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  4 * time.Second,
	}

	errCh := make(chan error)

	go func() {
		l.Println("listen in ", address)
		errCh <- srv.ListenAndServe()
	}()

	r := m.Run()

	if err := srv.Shutdown(context.Background()); err != nil {
		l.Println(err)
	}

	tx.Rollback()

	os.Exit(r)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// UserSdkMock
type UserSdkMock struct {
	GetFunc func(id string) (*domain.User, error)
}

func (r *UserSdkMock) Get(id string) (*domain.User, error) {
	return r.GetFunc(id)
}

// CourseSdkMock
type CourseSdkMock struct {
	GetFunc func(id string) (*domain.Course, error)
}

func (r *CourseSdkMock) Get(id string) (*domain.Course, error) {
	return r.GetFunc(id)
}