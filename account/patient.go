package account

import (
	"time"

	"github.com/reyhardy/exs-prescription-app/appointment"
	"github.com/reyhardy/exs-prescription-app/complaints"
)

type Patient struct {
	*Detail
	TherapistId string
	Diagnosis   string
	Appointmet  []*appointment.Appointment
	Complaints  []*complaints.Detail
}

type Complaint struct {
	Id            string
	Description   string
	PatientId     string
	CreatedAt     time.Time
	Prescriptions []string
}
