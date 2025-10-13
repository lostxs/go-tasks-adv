package main

import (
	"4-order-api/internal/models"
	"4-order-api/internal/order"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc2NTM1MzcsIm51bWJlciI6Ijg5MTQ3Nzc0NTU1Iiwic2Vzc2lvbklkIjoib3RwdGNEbmdYeTdDIn0._Ul5H2LymZR6EtGfGYZACqCOIGJeYEAeb79n_yFXEdQ"

func initDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db

}

func initData(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Create(&models.User{
			Number:    "89147774555",
			SessionID: "otptcDngXy7C",
			Code:      "0000",
		}).Error; err != nil {
			return err
		}

		if err := db.Create(&models.Product{
			Name:        "Печенька",
			Description: "Шоколадная глазурь",
			Images:      []string{"1.png"},
			Quantity:    10,
		}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func deleteData(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Delete(&models.OrderProduct{}, 1).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&models.Order{}, 1).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&models.User{}, 1).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&models.Product{}, 1).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func TestNewOrderSuccess(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()
	db := initDB()
	err := initData(db)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(&order.OrderRequest{
		UserID: 1,
		Products: []order.QuantProductID{
			{ProductID: 1, Quantity: 5},
		},
	})

	req, err := http.NewRequest("POST", ts.URL+"/order", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		t.Fatal(err)
	}

	client := ts.Client()
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Ожидали %d,получили %d", http.StatusCreated, res.StatusCode)
	}

	err = deleteData(db)
	if err != nil {
		t.Fatal(err)
	}

}
func TestNewOrderFail(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()
	db := initDB()
	err := initData(db)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(&order.OrderRequest{
		UserID: 10,
		Products: []order.QuantProductID{
			{ProductID: 20, Quantity: 5},
		},
	})

	req, err := http.NewRequest("POST", ts.URL+"/order", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}
	client := ts.Client()
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode == http.StatusCreated {
		t.Fatalf("Ожидали %d,получили %d", http.StatusCreated, res.StatusCode)
	}

	err = deleteData(db)
	if err != nil {
		t.Fatal(err)
	}

}
