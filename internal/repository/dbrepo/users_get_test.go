package dbrepo

import (
	"context"
	"errors"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pintoter/todo-list/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		id int
	}

	type mockBehavior func(args args)

	users := []entity.User{
		{
			ID:           1,
			Email:        "test@mail.ru",
			Login:        "test",
			Password:     "hashed",
			RegisteredAt: time.Time{},
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantUser     entity.User
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "email", "login", "password", "register_at"}).
					AddRow(users[0].ID, users[0].Email, users[0].Login, users[0].Password, users[0].RegisteredAt)

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.id).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				id: 1,
			},
			wantUser: users[0],
		},
		{
			name: "Failed",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.id).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			args: args{
				id: 100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			got, err := r.GetUserByID(context.Background(), tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUser, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		login string
	}

	type mockBehavior func(args args)

	users := []entity.User{
		{
			ID:           1,
			Email:        "test@mail.ru",
			Login:        "test",
			Password:     "hashed",
			RegisteredAt: time.Time{},
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantUser     entity.User
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "email", "login", "password", "register_at"}).
					AddRow(users[0].ID, users[0].Email, users[0].Login, users[0].Password, users[0].RegisteredAt)

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE login = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.login).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				login: "test",
			},
			wantUser: users[0],
		},
		{
			name: "Failed",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE login = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.login).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			args: args{
				login: "failed",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			got, err := r.GetUserByLogin(context.Background(), tt.args.login)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUser, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		email string
	}

	type mockBehavior func(args args)

	users := []entity.User{
		{
			ID:           1,
			Email:        "test@mail.ru",
			Login:        "test",
			Password:     "hashed",
			RegisteredAt: time.Time{},
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantUser     entity.User
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "email", "login", "password", "register_at"}).
					AddRow(users[0].ID, users[0].Email, users[0].Login, users[0].Password, users[0].RegisteredAt)

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.email).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				email: "test@mail.ru",
			},
			wantUser: users[0],
		},
		{
			name: "Failed",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.email).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			args: args{
				email: "failed",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			got, err := r.GetUserByEmail(context.Background(), tt.args.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUser, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		email    string
		password string
	}

	type mockBehavior func(args args)

	users := []entity.User{
		{
			ID:           1,
			Email:        "test@mail.ru",
			Login:        "test",
			Password:     "hashed",
			RegisteredAt: time.Time{},
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantUser     entity.User
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "email", "login", "password", "register_at"}).
					AddRow(users[0].ID, users[0].Email, users[0].Login, users[0].Password, users[0].RegisteredAt)

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE email = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.email, args.password).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				email:    "test@mail.ru",
				password: "hashed",
			},
			wantUser: users[0],
		},
		{
			name: "Failed",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE email = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.email, args.password).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			args: args{
				email: "failed",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			got, err := r.GetUserByCredentials(context.Background(), tt.args.email, tt.args.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUser, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		refreshToken string
	}

	type mockBehavior func(args args)

	users := []entity.User{
		{
			ID:           1,
			Email:        "test@mail.ru",
			Login:        "test",
			Password:     "hashed",
			RegisteredAt: time.Time{},
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantUser     entity.User
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "email", "login", "password", "register_at"}).
					AddRow(users[0].ID, users[0].Email, users[0].Login, users[0].Password, users[0].RegisteredAt)

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE refresh_token = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.refreshToken).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				refreshToken: "test_refresh_token",
			},
			wantUser: users[0],
		},
		{
			name: "Failed",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectExec := "SELECT id, email, login, password, register_at FROM users WHERE refresh_token = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.refreshToken).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			args: args{
				refreshToken: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			got, err := r.GetUserByRefreshToken(context.Background(), tt.args.refreshToken)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUser, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
