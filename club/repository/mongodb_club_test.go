package repository_test

import (
	"testing"
)

func TestGetByID(t *testing.T) {

	// db, mock, err := sqlmock.New()
	// if err != nil {
	// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	// }
	// defer db.Close()
	// rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
	// 	AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	// query := "SELECT id,title,content, author_id, updated_at, created_at FROM article WHERE ID = \\?"

	// mock.ExpectQuery(query).WillReturnRows(rows)
	// a := clubRepo.NewMysqlArticleRepository(db)

	// num := int64(5)
	// anArticle, err := a.GetByID(num)
	// assert.NoError(t, err)
	// assert.NotNil(t, anArticle)
	// assert.Equal(true, true)
}
