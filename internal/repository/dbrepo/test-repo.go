package dbrepo

import (
	"time"

	"github.com/mrkouhadi/go-booking-app/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a Reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestrictions) error {
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomId, and false if it no availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {
	return false, nil
}

// SearchAvailablityForAllRooms returns a slice of available rooms for given date range
func (m *testDBRepo) SearchAvailablityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRoomById gets a room by id
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {

	var room models.Room

	return room, nil
}
