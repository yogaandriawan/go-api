package main

import (
	"flag"
	"fmt"
	"strings"
	
  "go-api/backend-api/config"
	"go-api/backend-api/database"
	"go-api/backend-api/models"
	"go-api/backend-api/routes"
)

func main() {
	// load environment variables from .env file
	config.LoadEnv()

	// connect to database
	database.ConnectDB()

	db := database.ConnectDB()
    
    // 1. Definisikan flag -migrate
	shouldMigrate := flag.Bool("migrate", false, "Tambahkan flag ini untuk menjalankan migrasi")
	flag.Parse()

	// 2. Cek apakah flag digunakan
	if *shouldMigrate {
		fmt.Println("Menjalankan migrasi database...")
		err := db.AutoMigrate(&models.User{})
		if err != nil {
			fmt.Println("Gagal migrasi:", err)
		} else {
			fmt.Println("Migrasi berhasil selesai!")
		}
		return 
	}

	//setup router
	r := routes.SetupRouter()

	// Ambil list proxy dari ENV, jika kosong biarkan nil
  proxies := config.GetEnv("TRUSTED_PROXIES", "")
		if proxies != "" {
				r.SetTrustedProxies(strings.Split(proxies, ","))
    } else {
        r.SetTrustedProxies(nil)
    }

	// jalankan server pada port 8080
	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}