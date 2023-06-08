package repository

import (
	"time"

	"github.com/mrkouhadi/go-booking-app/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestrictions) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error)
}
