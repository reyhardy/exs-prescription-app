package appointment

import (
	"errors"
	"time"

	"github.com/segmentio/ksuid"
)

type Appointment struct {
	Id        string
	PatientId string
	TimeSlot  time.Time
	CreatedBy string
	CreatedAt time.Time

	Prescriptions []string
}

type Params struct {
	Id        string
	PatientId string
	CreatedBy string
	TimeSlot  time.Time
}

func New(p *Params) (*Appointment, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		return nil, err
	}
	p.Id = id.String()
	return Map(p)
}

func Map(p *Params) (*Appointment, error) {
	if p.PatientId == "" {
		return nil, errors.New("patient id not provided")
	}
	var prescriptions []string
	return &Appointment{
		Id:            p.Id,
		PatientId:     p.PatientId,
		CreatedBy:     p.CreatedBy,
		TimeSlot:      p.TimeSlot,
		Prescriptions: prescriptions,
	}, nil
}
