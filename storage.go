package main

import (
	"domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Keyword struct {
	gorm.Model
	Name                    string `gorm:"unique"`
	KeywordAnalyticsResults []KeywordAnalyticsResult
}

type KeywordAnalyticsResult struct {
	gorm.Model
	KeywordID uint
	Keyword   Keyword
	Data      domain.AnalyticsData `gorm:"serializer:json"`
}

func SaveAnalyticsForQuery(keyword Keyword, items domain.AnalyticsData) {
	db := connectToDatabase()

	// Create
	keyword = Keyword{Name: keyword.Name}
	db.Where(Keyword{Name: keyword.Name}).Attrs(Keyword{Name: keyword.Name}).FirstOrCreate(&keyword)
	result := db.Where("created_at > ? AND keyword_id = ?", time.Now().Add(-time.Hour*24), keyword.ID).Find(&KeywordAnalyticsResult{})
	if result.RowsAffected > 0 {
		return
	}

	var keywordAnalyticsResult = KeywordAnalyticsResult{Data: items, KeywordID: keyword.ID}

	db.Create(&keywordAnalyticsResult)
}

func connectToDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open("postgresql://postgres@localhost/ozon-parser?password=postgres"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Keyword{}, &KeywordAnalyticsResult{})
	return db
}

func ExtractPositionsList() []domain.KeywordAnalyticsResult {
	var results []domain.KeywordAnalyticsResult
	db := connectToDatabase()
	var dbResults []KeywordAnalyticsResult
	db.Preload("Keyword").Find(&dbResults)
	for _, analyticsModel := range dbResults {
		results = append(results, domain.KeywordAnalyticsResult{Keyword: domain.Keyword{Name: analyticsModel.Keyword.Name}, Data: analyticsModel.Data, Time: analyticsModel.CreatedAt})
	}

	return results
}