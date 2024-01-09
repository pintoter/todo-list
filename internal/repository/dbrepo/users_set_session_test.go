package dbrepo

import (
	"context"
	"database/sql/driver"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pintoter/todo-list/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestSetUserSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		id           int
		refreshToken string
		expiresAt    time.Time
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectExec := "UPDATE users SET refresh_token = $1, expires_at = $2 WHERE user_id = $3"
				mock.ExpectExec(regexp.QuoteMeta(expectExec)).
					WithArgs(args.refreshToken, args.expiresAt, args.id).
					WillReturnResult(driver.RowsAffected(0))

				mock.ExpectCommit()
			},
			args: args{
				id:           0,
				refreshToken: "refresh_token_test",
				expiresAt:    time.Time{},
			},
		},
		{
			name: "Failed",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectExec := "UPDATE users SET refresh_token = $1, expires_at = $2 WHERE user_id = $3"
				mock.ExpectExec(regexp.QuoteMeta(expectExec)).WithArgs(args.refreshToken, args.expiresAt, args.id)

				mock.ExpectRollback()
			},
			args:    args{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			err := r.SetSession(context.Background(), tt.args.id, entity.Session{
				RefreshToken: tt.args.refreshToken,
				ExpiresAt:    tt.args.expiresAt,
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}
