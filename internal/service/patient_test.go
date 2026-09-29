package service_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/complaints"
	"github.com/reyhardy/exs-prescription-app/internal/service"
)

var (
	errPtStore = errors.New("db error occur")
)

func TestFetchAllPatients(t *testing.T) {
	tests := map[string]struct {
		therapistId    string
		dbErr          bool ``
		wantErr        error
		hasTherapistId bool
	}{
		"success": {
			therapistId:    "therapist12345",
			wantErr:        nil,
			hasTherapistId: true,
		},
		"success with nil slice": {
			therapistId:    "notexists",
			wantErr:        nil,
			hasTherapistId: false,
		},
		"fail with db error": {
			therapistId: "therapist12345",
			dbErr:       true,
			wantErr:     errPtStore,
		},
	}
	for name, tt := range tests {
		ptr := &fakePatientRepo{
			pts: []*account.Patient{
				{
					Detail: &account.Detail{
						Id:    "pt1",
						Roles: []account.Role{account.RolePatient},
					},
					TherapistId: "therapist12345",
					Diagnosis:   "test dx 1",
				},
				{
					Detail: &account.Detail{
						Id:    "pt2",
						Roles: []account.Role{account.RolePatient},
					},
					TherapistId: "therapist12345",
					Diagnosis:   "test dx 1, test dx 2",
				},
				{
					Detail: &account.Detail{
						Id:    "pt3",
						Roles: []account.Role{account.RolePatient},
					},
					TherapistId: "other_therapist",
					Diagnosis:   "test dx 1",
				},
			},
			dbErr: tt.dbErr,
		}
		t.Run(name, func(t *testing.T) {
			pts := service.NewPatient(ptr, &fakeComplaintRepo{})
			s, err := pts.FetchAllPatients(t.Context(), tt.therapistId)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			found := slices.ContainsFunc(s, func(pt *account.Patient) bool {
				return pt.TherapistId == "therapist12345"
			})
			if got, want := found, tt.hasTherapistId; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

func TestFetchPatient(t *testing.T) {
	tests := map[string]struct {
		ptDbErr      bool
		coDbErr      bool
		id           string
		wantErr      error
		wantFullName string
	}{
		"success": {
			id:           "pt1",
			wantErr:      nil,
			wantFullName: "testpt one",
		},
		"error on pt db err": {
			ptDbErr: true,
			id:      "pt2",
			wantErr: errPtStore,
		},
		"error on complaint db err": {
			coDbErr: true,
			id:      "pt2",
			wantErr: errPtStore,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ptr := &fakePatientRepo{
				pts: []*account.Patient{
					{
						Detail: &account.Detail{
							Id:        "pt1",
							FirstName: "testpt",
							LastName:  "one",
						},
					},
					{
						Detail: &account.Detail{
							Id:        "pt2",
							FirstName: "testpt",
							LastName:  "two",
						},
					},
				},
				dbErr: tt.ptDbErr,
			}
			cr := &fakeComplaintRepo{
				co:    []*complaints.Detail{},
				dbErr: tt.coDbErr,
			}
			pts := service.NewPatient(ptr, cr)
			pt, err := pts.FetchPatient(t.Context(), tt.id)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if err != nil {
				t.Log("got error")
				return
			}
			if got, want := pt.FullName(), tt.wantFullName; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

type fakePatientRepo struct {
	pts   []*account.Patient
	dbErr bool
}

func (r *fakePatientRepo) SelectAllPatientsByTherapistId(ctx context.Context, therapistId string) ([]*account.Patient, error) {
	if r.dbErr {
		return nil, errPtStore
	}
	var newPtsSlice []*account.Patient
	for _, pt := range r.pts {
		if pt.TherapistId == therapistId {
			newPtsSlice = append(newPtsSlice, pt)
		}
	}
	return newPtsSlice, nil
}

func (r *fakePatientRepo) GetPatientById(ctx context.Context, id string) (*account.Patient, error) {
	if r.dbErr {
		return nil, errPtStore
	}
	for _, pt := range r.pts {
		if pt.Id == id {
			return pt, nil
		}
	}
	return nil, nil
}

type fakeComplaintRepo struct {
	co    []*complaints.Detail
	dbErr bool
}

func (r *fakeComplaintRepo) GetComplaintByPatientId(ctx context.Context, patientId string) ([]*complaints.Detail, error) {
	if r.dbErr {
		return nil, errPtStore
	}
	return r.co, nil
}
