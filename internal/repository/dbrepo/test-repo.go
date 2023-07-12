package dbrepo

import (
	"errors"
	"time"

	"github.com/mrkouhadi/go-booking-app/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a Reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	// if the room id is 2 then fail otherwise pass
	if res.RoomId == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestrictions) error {
	if res.RoomId == 1000000 { //   // giving it imposible data to meet
		return errors.New("soem error")
	}
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
	if id > 2 {
		return room, errors.New("an error")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {

	var u models.User
	return u, nil
}

// UpdateUser edits user's data in the DB
func (m *testDBRepo) UpdateUser(u models.User) error {

	return nil
}

// authenticating a user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

// AllReservations returns a slice of all reservations

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}


// AllNewReservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}
