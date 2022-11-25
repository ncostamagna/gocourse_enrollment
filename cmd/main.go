package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/ncostamagna/gocourse_enrollment/internal/enrollment"
	"github.com/ncostamagna/gocourse_enrollment/pkg/bootstrap"
	"github.com/ncostamagna/gocourse_enrollment/pkg/handler"

	courseSdk "github.com/ncostamagna/go_course_sdk/course"
	userSdk "github.com/ncostamagna/go_course_sdk/user"
)

func main() {

	_ = godotenv.Load()
	l := bootstrap.InitLogger()
	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	courseTrans := courseSdk.NewHttpClient(os.Getenv("API_COURSE_URL"), "")
	userTrans := userSdk.NewHttpClient(os.Getenv("API_USER_URL"), "")

	ctx := context.Background()
	enrollRepo := enrollment.NewRepo(db, l)
	enrollSrv := enrollment.NewService(l, userTrans, courseTrans, enrollRepo)
	h := handler.NewEnrollmentHTTPServer(ctx, enrollment.MakeEndpoints(enrollSrv, enrollment.Config{LimPageDef: pagLimDef}))
	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)
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

	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}

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
