package database

import (
	"fmt"
	"testing"

	"GoNews/pcg/typeStruct"

	"github.com/stretchr/testify/assert"
)

func TestSaveAndReadFromDB(t *testing.T) {
	db := InitDB()

	defer db.Close()

	// Создаем тестовый пост
	testPost := typeStruct.Post{
		Title:   "Test Title 2",
		Content: "Test Content",
		PubTime: 1692645688,
		Link:    "http://example.com/test",
	}

	// Сохраняем тестовый пост в базу данных
	idTest, err := SaveToDB(testPost)
	if err != nil {
		t.Fatalf("Failed to save post to DB: %v", err)
	}

	// Читаем пост из базы данных по названию
	readPost, err := ReadFromDB(idTest)
	if err != nil {
		t.Fatalf("Failed to read post from DB: %v", err)
	}

	// Сравниваем значения
	if readPost.Title != testPost.Title ||
		readPost.Content != testPost.Content ||
		readPost.PubTime != testPost.PubTime ||
		readPost.Link != testPost.Link {
		t.Errorf("Saved data doesn't match expected data.")
	}

	err = DeletePost(idTest)
	assert.NoError(t, err, "Failed to delete post by title")

	_, err = ReadFromDB(idTest)
	assert.Error(t, err, "Expected an error when trying to read deleted post")
}

func TestSearchPostsByKeyword(t *testing.T) {
	// Инициализация базы данных
	db := InitDB()
	defer db.Close()

	// Вставка тестовых данных
	post1 := typeStruct.Post{
		Title:   "aa24 f=f2 +++ 56ty",
		Content: "Test Description 1",
		PubTime: 1234567890,
		Link:    "http://example.com/test1",
	}

	// Сохранение тестовых данных в базе данных
	if _, err := SaveToDB(post1); err != nil {
		t.Fatalf("Failed to save post to DB: %v", err)
	}

	// Вызов функции для поиска
	keyword := "F=F2"
	posts, err := SearchPostsByKeyword(keyword)
	if err != nil {
		t.Fatalf("SearchPostsByKeyword failed: %v", err)
	}

	// Проверка результатов
	if len(posts) != 1 {
		t.Fatalf("Expected 1 post, but got %d", len(posts))
	} else {
		fmt.Println(posts)
	}
	if posts[0].Title != post1.Title {
		t.Fatalf("Expected post with title '%s', but got '%s'", post1.Title, posts[0].Title)
	}

	keyword = " "
	posts, err = SearchPostsByKeyword(keyword)
	if err != nil {
		t.Fatalf("SearchPostsByKeyword failed: %v", err)
	}

	// Проверка результатов
	if len(posts) < 1 {
		t.Fatalf("Expected 1 post, but got %d", len(posts))
	}

	keyword = ""
	posts, err = SearchPostsByKeyword(keyword)
	if err != nil {
		t.Fatalf("SearchPostsByKeyword failed: %v", err)
	}

	// Проверка результатов
	if len(posts) < 1 {
		t.Fatalf("Expected 1 post, but got %d", len(posts))
	}

	// Удаление тестовых записей из базы данных
	if err := DeletePost(post1.ID); err != nil {
		t.Fatalf("Failed to delete test record: %v", err)
	}
}
