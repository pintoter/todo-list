package dbrepo

import (
	"context"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteById(t *testing.T) {
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

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "Success",
			args: args{
				id: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM notes WHERE id = $1"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failed",
			args: args{
				id: 100,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM notes WHERE id = $1"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.id).
					WillReturnError(errors.New("new error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := r.DeleteNoteById(context.Background(), tt.args.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteNotes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type mockBehavior func()

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func() {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM notes"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WillReturnResult(sqlmock.NewResult(0, 5))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failed",
			mockBehavior: func() {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM notes"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WillReturnError(errors.New("new error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			err := r.DeleteNotes(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
