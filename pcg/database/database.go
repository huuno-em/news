package database

import (
	"GoNews/pcg/typeStruct"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math"

	_ "github.com/lib/pq"
)

// Константы для настройки подключения к базе данных
const (
	DBHost     = "localhost"
	DBPort     = "5434"
	DBUser     = "admin"
	DBPassword = "admin"
	DBName     = "news"
)

var DB *sql.DB

// Инициализация базы данных
func InitDB() *sql.DB {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)
	var err error
	DB, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	return DB
}

func ExecuteSchemaSQL(db *sql.DB) {
	// Чтение содержимого schema.sql
	schemaSQL, err := ioutil.ReadFile("pcg/database/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Выполнение SQL-запросов из schema.sql
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		log.Fatal(err)
	}
}

// Сохранение новости в базе данных
func SaveToDB(post typeStruct.Post) (int, error) {
	query := `
		INSERT INTO news (title, description, pub_date, source)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	row := DB.QueryRow(query, post.Title, post.Content, post.PubTime, post.Link)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	fmt.Println("Saved to DB with ID:", id)
	return id, nil
}

// Чтение новости из базы данных по названию
func ReadFromDB(id int) (typeStruct.Post, error) {
	var post typeStruct.Post

	query := `
		SELECT id, title, description, pub_date, source
		FROM news
		WHERE id = $1
	`
	row := DB.QueryRow(query, id)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
	if err != nil {
		return post, err
	}

	return post, nil
}

// Получение n последних новостей из базы данных
func GetLatestPosts(n int) ([]typeStruct.Post, error) {
	query := `
		SELECT id, title, description, pub_date, source
		FROM news
		ORDER BY pub_date DESC
		LIMIT $1
	`
	rows, err := DB.Query(query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []typeStruct.Post

	for rows.Next() {
		var post typeStruct.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// Удаление новости из базы данных по названию
func DeletePost(id int) error {
	_, err := DB.Exec("DELETE FROM news WHERE id = $1", id)
	return err
}

// SearchPostsByKeyword выполняет поиск новостей по ключевому слову в заголовке
func SearchPostsByKeyword(keyword string) ([]typeStruct.Post, error) {
	query := `
        SELECT id, title, description, pub_date, source
        FROM news
        WHERE title ILIKE '%' || $1 || '%'
    `
	rows, err := DB.Query(query, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []typeStruct.Post

	for rows.Next() {
		var post typeStruct.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPosts(page, pageSize int) (typeStruct.PaginatedPosts, error) {

	var totalResults int
	countQuery := `SELECT COUNT(id) FROM news`
	err := DB.QueryRow(countQuery).Scan(&totalResults)
	if err != nil {
		return typeStruct.PaginatedPosts{}, err
	}

	pagination := CalculatePagination(totalResults, pageSize, page)

	page = pagination.Page
	pageSize = pagination.PageSize

	// Вычисление смещения (offset) на основе номера страницы и размера страницы
	offset := (page - 1) * pageSize

	// Запрос новостей с учетом смещения и ограничения количества
	query := `
		SELECT id, title, description, pub_date, source
		FROM news
		ORDER BY pub_date DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := DB.Query(query, pageSize, offset)
	if err != nil {
		return typeStruct.PaginatedPosts{}, err
	}
	defer rows.Close()

	var posts []typeStruct.Post

	for rows.Next() {
		var post typeStruct.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			return typeStruct.PaginatedPosts{}, err
		}
		posts = append(posts, post)
	}

	return typeStruct.PaginatedPosts{
		Posts:      posts,
		Pagination: pagination,
	}, nil
}

func CalculatePagination(totalResults, pageSize, page int) typeStruct.Pagination {
	totalPages := int(math.Ceil(float64(totalResults) / float64(pageSize)))

	if totalPages < 1 {
		totalPages = 1
	}
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	return typeStruct.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalResults,
	}
}
