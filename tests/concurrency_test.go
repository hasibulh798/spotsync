package tests

import (
	"errors"
	"log"
	"sync"
	"testing"

	"spotsync/config"
	"spotsync/models"
	"spotsync/repository"
	"spotsync/utils"

	"github.com/joho/godotenv"
)

// TestReservationConcurrency checks if parallel reservation requests are transactionally locked
// and exactly one is allowed to book the spot.
func TestReservationConcurrency(t *testing.T) {
	// Load env from project root (one level up) since test runs in package directory
	_ = godotenv.Load("../.env")

	// 1. Load config and connect to DB
	cfg := config.LoadConfig()
	if cfg.DBURL == "" {
		t.Skip("Skipping concurrency integration test as DB_URL is not configured")
	}

	db := config.InitDB(cfg)

	// Clean up database tables for a clean test run
	db.Exec("DELETE FROM reservations")
	db.Exec("DELETE FROM parking_zones")
	db.Exec("DELETE FROM users")

	// 2. Create test zone with capacity = 1
	zone := models.ParkingZone{
		Name:          "EV Charging Spot A",
		Type:          "ev_charging",
		TotalCapacity: 1,
		PricePerHour:  120.00,
	}
	if err := db.Create(&zone).Error; err != nil {
		t.Fatalf("Failed to create test zone: %v", err)
	}

	// Create two driver users
	user1 := models.User{Name: "Driver A", Email: "driver.a@example.com", Password: "hashedpassword", Role: "driver"}
	user2 := models.User{Name: "Driver B", Email: "driver.b@example.com", Password: "hashedpassword", Role: "driver"}
	if err := db.Create(&user1).Error; err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}
	if err := db.Create(&user2).Error; err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	// 3. Initialize repository
	resRepo := repository.NewReservationRepository(db)

	var wg sync.WaitGroup
	wg.Add(2)

	var successCount int
	var conflictCount int
	var mu sync.Mutex

	// Launch Goroutine 1 (Driver A trying to reserve)
	go func() {
		defer wg.Done()
		_, err := resRepo.CreateReservationWithLock(zone.ID, user1.ID, "PLATE-CONC-1")
		mu.Lock()
		defer mu.Unlock()
		if err == nil {
			successCount++
			log.Println("Goroutine 1 (Driver A): Reservation Succeeded (201 Created)")
		} else if errors.Is(err, utils.ErrZoneFull) {
			conflictCount++
			log.Println("Goroutine 1 (Driver A): Reservation Failed - Zone Full (409 Conflict)")
		} else {
			t.Errorf("Goroutine 1 unexpected error: %v", err)
		}
	}()

	// Launch Goroutine 2 (Driver B trying to reserve)
	go func() {
		defer wg.Done()
		_, err := resRepo.CreateReservationWithLock(zone.ID, user2.ID, "PLATE-CONC-2")
		mu.Lock()
		defer mu.Unlock()
		if err == nil {
			successCount++
			log.Println("Goroutine 2 (Driver B): Reservation Succeeded (201 Created)")
		} else if errors.Is(err, utils.ErrZoneFull) {
			conflictCount++
			log.Println("Goroutine 2 (Driver B): Reservation Failed - Zone Full (409 Conflict)")
		} else {
			t.Errorf("Goroutine 2 unexpected error: %v", err)
		}
	}()

	wg.Wait()

	// 4. Assert concurrency constraints
	if successCount != 1 {
		t.Errorf("Expected exactly 1 successful reservation, got %d", successCount)
	}
	if conflictCount != 1 {
		t.Errorf("Expected exactly 1 conflict error, got %d", conflictCount)
	}

	log.Printf("[SUCCESS] Concurrency Test Passed. Successes: %d, Conflicts: %d", successCount, conflictCount)
}
