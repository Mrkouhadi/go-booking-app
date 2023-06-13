package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/mrkouhadi/go-booking-app/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a Reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the this operation could not succeed within 3 seconds end it immediately
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newId int
	statement := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	err := m.DB.QueryRowContext(ctx, statement,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId) // store the returned(scan) id to newId
	if err != nil {
		log.Print("Error inserting reservation: ", err)
		return 0, err
	}
	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestrictions) error {
	// if the this operation could not succeed within 3 seconds end it immediately
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := `insert into room_restrictions(start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id)
									values($1,$2,$3,$4,$5,$6,$7)`
	_, err := m.DB.ExecContext(ctx, statement,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		res.ReservationId,
		time.Now(),
		time.Now(),
		res.RestrictionId,
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomId, and false if it no availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {
	// if the this operation could not succeed within 3 seconds, end it immediately
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var numRows int
	// give me the number of rows from rooom_restrictions where my start_date is less than your end_date
	// and my end date is greater than your start date
	query := ` 
   		select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date;
   `
	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	// if the number is 0 then there is availability, otherwise there is no availability
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailablityForAllRooms returns a slice of available rooms for given date range
func (m *postgresDBRepo) SearchAvailablityForAllRooms(start, end time.Time) ([]models.Room, error) {
	// if the this operation could not succeed within 3 seconds, end it immediately
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `
			select r.id, r.room_name 
			from rooms r 
			where r.id not in 
			(select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);
		`
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// GetRoomById gets a room by id
func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		select id, room_name, created_at, updated_at from rooms where id=$1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}
