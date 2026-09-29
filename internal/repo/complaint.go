package repo

import (
	"context"

	"github.com/reyhardy/exs-prescription-app/complaints"
	"github.com/reyhardy/exs-prescription-app/internal/db"
)

type Complaint struct {
	db *db.DB
}

func NewComplaint(db *db.DB) *Complaint {
	return &Complaint{
		db: db,
	}
}

func (r *Complaint) GetComplaintByPatientId(ctx context.Context, patientId string) ([]*complaints.Detail, error) {
	q := `
		SELECT id, description, created_at
		FROM complaints
		WHERE patient_id = ?
	`
	var c []*complaints.Detail
	if err := r.db.Query(ctx, q, &c, patientId); err != nil {
		return nil, err
	}
	return c, nil
}
