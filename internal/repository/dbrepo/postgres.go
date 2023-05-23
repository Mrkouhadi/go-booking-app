package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mrkouhadi/go-booking-app/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	fmt.Println("CONTACT has been hit")
	return true
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the this operation could not succeed within 3 seconds end it immediately
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newId int
	statement := `insert into reservations (first_name, last_name, email,phone,start_date,end_date,room_id, created_at, updated_at)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	err := m.DB.QueryRowContext(ctx, statement,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		time.Now(),
		time.Now(),
	).Scan(&newId) // store the returned(scan) id to newId
	if err != nil {
		log.Print("Error inserting reservation: ", err)
		return 0, err
	}
	return newId, nil
}
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestrictions) error {
	return nil
}
