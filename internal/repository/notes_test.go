package repository

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

func TestNote_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		note entity.Note
	}

	type mockBehavior func(args args)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "Success test",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			args: args{
				note: entity.Note{
					Title:       "Test title",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      entity.StatusDone,
				},
			},
			id: 1,
		},
		{
			name: "Success with empty date",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			args: args{
				note: entity.Note{
					Title:       "Test title",
					Description: "Test describstion",
					Status:      entity.StatusDone,
				},
			},
			id: 1,
		},
		{
			name: "Success with empty description",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			args: args{
				note: entity.Note{
					Title:  "Test title",
					Date:   time.Now().Round(time.Second),
					Status: entity.StatusDone,
				},
			},
			id: 1,
		},
		{
			name: "Empty title",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnError(errors.New("empty title"))

				mock.ExpectRollback()
			},
			args: args{
				note: entity.Note{
					Title:       "",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      entity.StatusDone,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid status",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnError(errors.New("empty title"))

				mock.ExpectRollback()
			},
			args: args{
				note: entity.Note{
					Title:       "Title",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      "not dont",
				},
			},
			wantErr: true,
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.Create(context.Background(), testCase.args.note)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
/*
func TestNote_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		id int
	}

	testNoteId := 1
	testNote := entity.Note{
		Title: "Test title",
		Description: "Test description",
		Date: time.Now().Round(time.Second),
		Status: entity.StatusDone,
	}

	type mockBehavior func(args args)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
	}{
		{
			name: "Success test",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			args: args{
				note: entity.Note{
					Title:       "Test title",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      entity.StatusDone,
				},
			},
			id: 1,
		},
		{
			name: "Success with empty date",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			args: args{
				note: entity.Note{
					Title:       "Test title",
					Description: "Test describstion",
					Status:      entity.StatusDone,
				},
			},
			id: 1,
		},
		{
			name: "Success with empty description",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			args: args{
				note: entity.Note{
					Title:  "Test title",
					Date:   time.Now().Round(time.Second),
					Status: entity.StatusDone,
				},
			},
			id: 1,
		},
		{
			name: "Empty title",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnError(errors.New("empty title"))

				mock.ExpectRollback()
			},
			args: args{
				note: entity.Note{
					Title:       "",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      entity.StatusDone,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid status",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (title,description,date,status) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnError(errors.New("empty title"))

				mock.ExpectRollback()
			},
			args: args{
				note: entity.Note{
					Title:       "Title",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      "not dont",
				},
			},
			wantErr: true,
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.Create(context.Background(), testCase.args.note)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
*/