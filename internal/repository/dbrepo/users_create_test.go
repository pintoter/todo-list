package dbrepo

import (
	"context"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pintoter/todo-list/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		user entity.User
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectExec := "INSERT INTO users (email,login,password,register_at) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectExec)).
					WithArgs(args.user.Email, args.user.Login, args.user.Password, args.user.RegisteredAt).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			args: args{
				user: entity.User{
					Email:        "testy",
					Login:        "test",
					Password:     "hashedpw",
					RegisteredAt: time.Now().Round(time.Second),
				},
			},
			id: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			got, err := r.CreateUser(context.Background(), tt.args.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
