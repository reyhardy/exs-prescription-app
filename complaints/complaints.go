package complaints

import (
	"errors"
	"time"

	"github.com/segmentio/ksuid"
)

var (
	ErrPatientIdNotProvided = errors.New("patient id not provided")
)

type Detail struct {
	Id          string
	Description string
	PatientId   string
	CreatedAt   time.Time
}

type Params struct {
	Id          string
	PatientId   string
	Description string
}

func New(p *Params) (*Detail, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		return nil, err
	}
	p.Id = id.String()
	return Map(p)
}

func Map(p *Params) (*Detail, error) {
	if p.PatientId == "" {
		return nil, ErrPatientIdNotProvided
	}
	return &Detail{
		Id:          p.Id,
		PatientId:   p.PatientId,
		Description: p.Description,
	}, nil
}
