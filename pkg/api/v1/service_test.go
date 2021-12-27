package v1

import (
	"context"
	"regexp"
	"time"

	"testing"

	"github.com/jmoiron/sqlx"
	v1 "github.com/nedson202/auth-manager/api/proto/v1"
	"github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Test_toDoServiceServer_Signup(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &Config{
		JwtSecret: "jwt_sect",
	}
	tokenService := NewTokenService(*cfg)
	repo := NewPostgresRepository(sqlxDB)
	s := NewAuthServiceServer(repo, tokenService)

	convey.Convey("User signup", t, func() {
		convey.Convey("should return Id and Token on successful registration", func() {
			mock := func() {
				rows := sqlmock.NewRows([]string{"id", "email", "created_at", "updated_at"}).AddRow("931131fe-3252-4530-a4c3-7947c8edfc8b", "test@gmail.com", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO identity.user (id,email,password) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING RETURNING "id", "email", "created_at", "updated_at"`)).
					WillReturnRows(rows)
			}
			mock()
			got, err := s.Signup(ctx, &v1.AuthRequest{
				Email:    "test@gmail.com",
				Password: "random_test_password",
			})

			if err == nil {
				convey.So(got.Id, convey.ShouldNotBeNil)
				convey.So(got.Token, convey.ShouldNotBeNil)
			}
		})

		convey.Convey("should return err if user already exists", func() {
			mock := func() {
				rows := sqlmock.NewRows([]string{"id", "email", "created_at", "updated_at"}).AddRow("", "", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO identity.user (id,email,password) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING RETURNING "id", "email", "created_at", "updated_at"`)).
					WillReturnRows(rows)
			}
			mock()
			_, err := s.Signup(ctx, &v1.AuthRequest{
				Email:    "test@gmail.com",
				Password: "random_test_password",
			})

			convey.So(err.Error(), convey.ShouldEqual, status.Error(codes.AlreadyExists, "please provide another email").Error())
		})
	})
}

func Test_toDoServiceServer_Login(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cfg := &Config{
		JwtSecret: "jwt_sect",
	}
	tokenService := NewTokenService(*cfg)
	repo := NewPostgresRepository(sqlxDB)
	s := NewAuthServiceServer(repo, tokenService)

	convey.Convey("User login", t, func() {
		convey.Convey("should return Id and Token on successful login", func() {
			password := "94jsdk42033-sdrer"
			passwordHash, _ := hashPassword(password)
			mock := func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).AddRow("931131fe-3252-4530-a4c3-7947c8edfc8b", "", "test@gmail.com", passwordHash, time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at, updated_at FROM identity.user WHERE email = $1`)).
					WillReturnRows(rows)
			}
			mock()
			got, _ := s.Login(ctx, &v1.AuthRequest{
				Email:    "test@gmail.com",
				Password: password,
			})

			convey.So(got.Id, convey.ShouldNotBeNil)
			convey.So(got.Token, convey.ShouldNotBeNil)
		})

		convey.Convey("should return err if user does not exist", func() {
			password := "94jsdk42033-sdrer"
			mock := func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).AddRow("", "", "", "", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at, updated_at FROM identity.user WHERE email = $1`)).
					WillReturnRows(rows)
			}
			mock()
			_, err := s.Login(ctx, &v1.AuthRequest{
				Email:    "test@gmail.com",
				Password: password,
			})

			convey.So(err.Error(), convey.ShouldEqual, status.Error(codes.Unauthenticated, "authentication failed due to invalid credentials").Error())
		})

		convey.Convey("should return err if invalid credentials is provided", func() {
			password := "94jsdk42033-sdrer"
			passwordHash, _ := hashPassword(password)
			mock := func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).AddRow("931131fe-3252-4530-a4c3-7947c8edfc8b", "", "test@gmail.com", passwordHash, time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at, updated_at FROM identity.user WHERE email = $1`)).
					WillReturnRows(rows)
			}
			mock()
			_, err := s.Login(ctx, &v1.AuthRequest{
				Email:    "test@gmail.com",
				Password: "94jsdk42033",
			})

			convey.So(err.Error(), convey.ShouldEqual, status.Error(codes.Unauthenticated, "authentication failed due to invalid credentials").Error())
		})
	})
}
