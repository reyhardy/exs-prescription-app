package service

import (
	"context"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/complaints"
)

type PatientRepo interface {
	SelectAllPatientsByTherapistId(ctx context.Context, therapistId string) ([]*account.Patient, error)
	GetPatientById(ctx context.Context, id string) (*account.Patient, error)
}

type ComplaintRepo interface {
	GetComplaintByPatientId(ctx context.Context, patientId string) ([]*complaints.Detail, error)
}

type Patient struct {
	ptr PatientRepo
	cr  ComplaintRepo
}

func NewPatient(ptr PatientRepo, cr ComplaintRepo) *Patient {
	return &Patient{
		ptr: ptr,
		cr:  cr,
	}
}

func (s *Patient) FetchAllPatients(ctx context.Context, therapistId string) ([]*account.Patient, error) {
	pts, err := s.ptr.SelectAllPatientsByTherapistId(ctx, therapistId)
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return nil, nil
	}
	return pts, nil
}

func (s *Patient) FetchPatient(ctx context.Context, id string) (*account.Patient, error) {
	pt, err := s.ptr.GetPatientById(ctx, id)
	if err != nil {
		return nil, err
	}
	c, err := s.cr.GetComplaintByPatientId(ctx, pt.Id)
	if err != nil {
		return nil, err
	}
	pt.Complaints = c
	return pt, nil
}
