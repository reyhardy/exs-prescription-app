package repo

import (
	"context"

	"github.com/reyhardy/exs-prescription-app/appointment"
	"github.com/reyhardy/exs-prescription-app/internal/db"
)

type Appointment struct {
	db *db.DB
}

func NewAppointment(db *db.DB) *Appointment {
	return &Appointment{
		db: db,
	}
}

func (r *Appointment) SelectApptByPtId(ctx context.Context, ptId string) ([]*appointment.Appointment, error) {
	q := `
		SELECT id, time_slot, created_by, created_at
		FROM appointments
		WHERE patient_id = ?;
	`
	var appts []*appointment.Appointment
	if err := r.db.Query(ctx, q, &appts, ptId); err != nil {
		return nil, err
	}
	return appts, nil
}
