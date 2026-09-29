package repo

import (
	"context"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/db"
)

type Patient struct {
	db *db.DB
}

func NewPatient(db *db.DB) *Patient {
	return &Patient{
		db: db,
	}
}

func (r *Patient) SelectAllPatientsByTherapistId(ctx context.Context, therapistId string) ([]*account.Patient, error) {
	q := `
		SELECT a.id, a.first_name, a.last_name, a.age, a.gender, p.diagnosis
		FROM accounts AS a
		LEFT JOIN patients AS p
		ON a.id = p.id
		WHERE p.therapist_id = ?
	`
	var pts []*account.Patient
	if err := r.db.Query(ctx, q, &pts, therapistId); err != nil {
		return nil, err
	}
	return pts, nil
}

func (r *Patient) GetPatientById(ctx context.Context, id string) (*account.Patient, error) {
	q := `
		SELECT a.id, a.first_name, a.last_name, a.age, a.gender, p.diagnosis
		FROM accounts AS a
		LEFT JOIN patients AS p
		ON a.id = p.id
		WHERE a.id = ?
	`
	var pt account.Patient
	if err := r.db.QueryRow(ctx, q, &pt, id); err != nil {
		return nil, err
	}
	return &pt, nil
}
