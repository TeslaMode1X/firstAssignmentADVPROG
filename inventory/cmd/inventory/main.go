package main

import (
	"log"
	"os"
	"time"

	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/config"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/database"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/internal/model/dto"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/inventory/server"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewPostgresDatabase(cfg)

	migrateAndSeed(db)

	l := log.New(os.Stdout, "inventory-gin ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	server.NewGinServer(cfg, db, l).Start()
}

func migrateAndSeed(db database.Database) {
	if err := db.GetDb().Migrator().AutoMigrate(&dto.CategoryDTO{}, &dto.ProductDTO{}, &dto.PromotionDTO{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	var count int64
	db.GetDb().Model(&dto.CategoryDTO{}).Count(&count)
	if count == 0 {

		categories := []dto.CategoryDTO{
			{
				Name:        "Electronics",
				Description: "Devices and gadgets",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Accessories",
				Description: "Computer and device accessories",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Clothing",
				Description: "Apparel and wearables",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		if err := db.GetDb().CreateInBatches(categories, 3).Error; err != nil {
			log.Fatalf("Failed to seed categories: %v", err)
		}
		log.Println("Successfully seeded categories")

		products := []dto.ProductDTO{
			{
				Name:        "Laptop",
				Description: "High-performance laptop",
				Price:       999.99,
				StockLevel:  10,
				CategoryID:  1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "Wireless Mouse",
				Description: "Ergonomic wireless mouse",
				Price:       29.99,
				StockLevel:  50,
				CategoryID:  2,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "T-Shirt",
				Description: "Cotton graphic t-shirt",
				Price:       19.99,
				StockLevel:  100,
				CategoryID:  3,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		if err := db.GetDb().CreateInBatches(&products, 3).Error; err != nil {
			log.Fatalf("Failed to seed products: %v", err)
		}

		log.Println("Successfully seeded products")
	} else {
		log.Println("Database already seeded, skipping initial data insertion")
	}

}
