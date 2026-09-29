package handler_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/ctx"
	"github.com/reyhardy/exs-prescription-app/internal/handler"
	"github.com/reyhardy/exs-prescription-app/internal/template"
)

func TestListAllPatients(t *testing.T) {
	tests := map[string]struct {
		user              *account.Detail
		fetchAllPtsFunc   func(ctx context.Context, therapistId string) ([]*account.Patient, error)
		renderFn          func(ctx context.Context, w io.Writer, page template.PageName, data any) error
		wantCode          int
		wantPatientIdData string
	}{
		"success": {
			user: &account.Detail{
				Id:    "test12345",
				Roles: []account.Role{account.RoleTherapist},
			},
			fetchAllPtsFunc: func(ctx context.Context, therapistId string) ([]*account.Patient, error) {
				if therapistId != "test12345" {
					return nil, fmt.Errorf("wrong therapist Id: %v", therapistId)
				}
				return []*account.Patient{
					{
						Detail:    &account.Detail{Id: "pt12345"},
						Diagnosis: "test diagnosis",
					},
				}, nil
			},
			renderFn: func(ctx context.Context, w io.Writer, page template.PageName, data any) error {
				return nil
			},
			wantCode:          http.StatusOK,
			wantPatientIdData: "pt12345",
		},
		"failure on wrong role": {
			user: &account.Detail{
				Id:    "test12345",
				Roles: []account.Role{account.RolePatient},
			},
			fetchAllPtsFunc: func(ctx context.Context, therapistId string) ([]*account.Patient, error) {
				return nil, fmt.Errorf("this id %v has wrong role", therapistId)
			},
			wantCode: http.StatusInternalServerError,
		},
		"success when zero patients return while error is nil": {
			user: &account.Detail{
				Id:    "test12345",
				Roles: []account.Role{account.RoleTherapist},
			},
			fetchAllPtsFunc: func(ctx context.Context, therapistId string) ([]*account.Patient, error) {
				if therapistId != "test12345" {
					return nil, fmt.Errorf("wrong therapist Id: %v", therapistId)
				}
				return nil, nil
			},
			renderFn: func(ctx context.Context, w io.Writer, page template.PageName, data any) error {
				return nil
			},
			wantCode: http.StatusOK,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			pfMock := &mockPtFetcher{
				fetchAllPtsFunc: tt.fetchAllPtsFunc,
			}
			renderer := &patientsListPageRenderer{
				renderFn: tt.renderFn,
			}
			h := handler.NewPatient(pfMock, renderer)
			req := httptest.NewRequest(http.MethodGet, "/patients", nil)
			rr := httptest.NewRecorder()
			uCtx := ctx.MapUserCtx(tt.user)
			ctx := ctx.WithContext(req.Context(), uCtx)
			h.ListAllPatients(rr, req.WithContext(ctx))
			if got, want := rr.Code, tt.wantCode; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			data, ok := renderer.data.([]*account.Patient)
			if !ok || data == nil {
				t.Log("no data found")
				return
			}
			if got, want := data[0].Id, tt.wantPatientIdData; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

func TestViewPatientInfo(t *testing.T) {
	tests := map[string]struct {
	}{}

	for name := range tests {
		t.Run(name, func(t *testing.T) {

		})
	}
}

type mockPtFetcher struct {
	fetchAllPtsFunc func(context.Context, string) ([]*account.Patient, error)
}

func (m *mockPtFetcher) FetchAllPatients(ctx context.Context, therapistId string) ([]*account.Patient, error) {
	return m.fetchAllPtsFunc(ctx, therapistId)
}

func (m *mockPtFetcher) FetchPatient(ctx context.Context, id string) (*account.Patient, error) {
	return nil, nil
}

type patientsListPageRenderer struct {
	renderFn func(context.Context, io.Writer, template.PageName, any) error
	data     any
}

func (r *patientsListPageRenderer) Render(ctx context.Context, w io.Writer, name template.PageName, data any) error {
	r.data = data
	return r.renderFn(ctx, w, name, r.data)
}
