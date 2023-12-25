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

	type mockBehavior func(args args)

	id := 1
	notes := []entity.Note{
		{
			Title:       "Test title",
			Description: "Test description",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantNote     entity.Note
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"title", "description", "date", "status"}).
					AddRow(notes[0].Title, notes[0].Description, notes[0].Date, notes[0].Status)

				expectedQuery := "SELECT title, description, date, status FROM notes WHERE id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args:     args{id: id},
			wantNote: notes[0],
		},
		{
			name: "Failed_NotFound",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"title", "description", "date", "status"})

				expectedQuery := "SELECT title, description, date, status FROM notes WHERE id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id).WillReturnError(errors.New("failed test")).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args:     args{id: id},
			wantNote: entity.Note{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotNote, err := r.GetNoteById(context.Background(), tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantNote, gotNote)
		})
	}
}

func TestNote_GetByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		title string
	}

	type mockBehavior func(args args)

	title := "Test title"
	notes := []entity.Note{
		{
			Title:       title,
			Description: "Test description",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantNote     entity.Note
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"title", "description", "date", "status"}).
					AddRow(notes[0].Title, notes[0].Description, notes[0].Date, notes[0].Status)

				expectedQuery := "SELECT title, description, date, status FROM notes WHERE title = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.title).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args:     args{title: title},
			wantNote: notes[0],
		},
		{
			name: "Failed_NotFound",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"title", "description", "date", "status"})

				expectedQuery := "SELECT title, description, date, status FROM notes WHERE title = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.title).WillReturnError(errors.New("failed test")).WillReturnRows(rows)
			},
			args:     args{title: title},
			wantNote: entity.Note{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotNote, err := r.GetNoteByTitle(context.Background(), tt.args.title)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantNote, gotNote)
		})
	}
}

func TestGetNotes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type mockBehavior func()

	notes := []entity.Note{
		{
			Title:       "Test title 1",
			Description: "Test description 1",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
		{
			Title:       "Test title 2",
			Description: "Test description 2",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
		{
			Title:       "Test title 3",
			Description: "Test description 3",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
		{
			Title:       "Test title 4",
			Description: "Test description 4",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		wantNotes    []entity.Note
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func() {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "status"}).
					AddRow(0, notes[0].Title, notes[0].Description, notes[0].Date, notes[0].Status).
					AddRow(0, notes[1].Title, notes[1].Description, notes[1].Date, notes[1].Status).
					AddRow(0, notes[2].Title, notes[2].Description, notes[2].Date, notes[2].Status).
					AddRow(0, notes[3].Title, notes[3].Description, notes[3].Date, notes[3].Status)

				expectedQuery := "SELECT id, title, description, date, status FROM notes ORDER BY id ASC"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantNotes: []entity.Note{notes[0], notes[1], notes[2], notes[3]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			got, err := r.GetNotes(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantNotes, got)
			}
		})
	}
}

func TestGetNotesExtended(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		limit  int
		offset int
		status string
		date   time.Time
	}

	type mockBehavior func(args args)

	dateFormatted, _ := time.Parse("2006-01-02", "2020-04-18")

	notes := []entity.Note{
		{
			ID:          1,
			Title:       "Test title 1",
			Description: "Test description 1",
			Date:        time.Time{},
			Status:      entity.StatusNotDone,
		},
		{
			ID:          2,
			Title:       "Test title 2",
			Description: "Test description 2",
			Date:        time.Time{},
			Status:      entity.StatusDone,
		},
		{
			ID:          3,
			Title:       "Test title 3",
			Description: "Test description 3",
			Date:        time.Time{},
			Status:      entity.StatusNotDone,
		},
		{
			ID:          4,
			Title:       "Test title 4",
			Description: "Test description 4",
			Date:        time.Time{},
			Status:      entity.StatusDone,
		},
	}

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantNotes    []entity.Note
		wantErr      bool
	}{
		{
			name: "SuccessWithoutStatusAndDate",
			args: args{
				limit:  5,
				offset: 5,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "status"}).
					AddRow(notes[0].ID, notes[0].Title, notes[0].Description, notes[0].Date, notes[0].Status).
					AddRow(notes[1].ID, notes[1].Title, notes[1].Description, notes[1].Date, notes[1].Status).
					AddRow(notes[2].ID, notes[2].Title, notes[2].Description, notes[2].Date, notes[2].Status).
					AddRow(notes[3].ID, notes[3].Title, notes[3].Description, notes[3].Date, notes[3].Status)

				expectedQuery := "SELECT id, title, description, date, status FROM notes ORDER BY id ASC LIMIT 5 OFFSET 5"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantNotes: []entity.Note{notes[0], notes[1], notes[2], notes[3]},
			wantErr:   false,
		},
		{
			name: "SuccessWithoutDate",
			args: args{
				limit:  5,
				offset: 5,
				status: entity.StatusDone,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "status"}).
					AddRow(notes[1].ID, notes[1].Title, notes[1].Description, notes[1].Date, notes[1].Status).
					AddRow(notes[3].ID, notes[3].Title, notes[3].Description, notes[3].Date, notes[3].Status)

				expectedQuery := "SELECT id, title, description, date, status FROM notes WHERE status = $1 ORDER BY id ASC LIMIT 5 OFFSET 5"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.status).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantNotes: []entity.Note{notes[1], notes[3]},
			wantErr:   false,
		},
		{
			name: "SuccessWithStatusAndDate",
			args: args{
				limit:  5,
				offset: 5,
				date:   dateFormatted,
				status: entity.StatusDone,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "status"}).
					AddRow(notes[1].ID, notes[1].Title, notes[1].Description, dateFormatted, notes[1].Status).
					AddRow(notes[3].ID, notes[3].Title, notes[3].Description, dateFormatted, notes[3].Status)

				expectedQuery := "SELECT id, title, description, date, status FROM notes WHERE status = $1 AND date = $2 ORDER BY id ASC LIMIT 5 OFFSET 5"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.status, args.date).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantNotes: []entity.Note{
				{
					ID:          notes[1].ID,
					Title:       notes[1].Title,
					Description: notes[1].Description,
					Date:        dateFormatted,
					Status:      notes[1].Status,
				},
				{
					ID:          notes[3].ID,
					Title:       notes[3].Title,
					Description: notes[3].Description,
					Date:        dateFormatted,
					Status:      notes[3].Status,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			got, err := r.GetNotesExtended(context.Background(), tt.args.limit, tt.args.offset, tt.args.status, tt.args.date)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantNotes, got)
			}
		})
	}
}
