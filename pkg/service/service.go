package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
	"mwz.com/data"
)

type TestService struct {
	d            *data.Data
	requestGroup singleflight.Group
}

func NewTestService(d *data.Data) TestService {
	var testService TestService
	testService.d = d

	return testService
}

type Article struct {
	Id      int32  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

func generateData() []Article {
	return []Article{
		{
			Id:      1,
			Title:   "abc",
			Content: "c",
		},
		{
			Id:      2,
			Title:   "def",
			Content: "c",
		},
		{
			Id:      3,
			Title:   "ghi",
			Content: "c",
		},
	}
}

func (s *TestService) queryFromDb(ctx context.Context) ([]Article, error) {
	fmt.Println("query_from_db")
	db := s.d.GetDb()

	var articles []Article = make([]Article, 0)

	rows, err := db.QueryContext(ctx, "SELECT id, title, content FROM articles")
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var article Article
		if err := rows.Scan(&article.Id, &article.Title, &article.Content); err != nil {
			return nil, fmt.Errorf("albumsByArtist %v", err)
		}
		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %v", err)
	}

	return articles, nil
}

func (s *TestService) GetRedis(ctx context.Context) string {
	key := "k1"

	val := s.d.Get(ctx, key)
	if val != "" {
		return val
	}

	v, err, _ := s.requestGroup.Do(key, func() (interface{}, error) {
		// dataList := generateData()
		dataList, err := s.queryFromDb(ctx)
		if err != nil {
			fmt.Printf("%v", err)
			dataList = generateData()
		}

		s.d.SetWithExpir(ctx, key, dataList, time.Second*2)

		rsVal, _ := json.Marshal(dataList)
		return string(rsVal), nil
	})

	if err != nil {
		fmt.Printf("%v", err)
	}

	return v.(string)
}
