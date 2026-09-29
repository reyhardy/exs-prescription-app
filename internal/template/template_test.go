package template_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/exercise"
	"github.com/reyhardy/exs-prescription-app/internal/ctx"
	"github.com/reyhardy/exs-prescription-app/internal/template"
)

func TestRenderPage(t *testing.T) {
	context := ctx.WithContext(t.Context(), ctx.MapUserCtx(&account.Detail{
		Id:    "test12345",
		Roles: []account.Role{account.RoleTherapist},
	}))
	tests := map[string]struct {
		page     template.PageName
		data     any
		wantBody string
		wantErr  error
	}{
		"render login page": {
			page:     template.Login,
			data:     "Invalid Credentials",
			wantBody: "Invalid Credentials",
			wantErr:  nil,
		},
		"render profile page": {
			page: template.UserProfile,
			data: &account.Detail{
				FirstName: "test",
				LastName:  "ping",
			},
			wantBody: "test",
			wantErr:  nil,
		},
		"render patients list page": {
			page: template.PatientsList,
			data: []*account.Patient{
				{
					Detail: &account.Detail{
						FirstName: "test",
					},
					Diagnosis: "test diagnosis 1, test diagnosis 2",
				},
			},
			wantBody: "test diagnosis 1, test diagnosis 2",
			wantErr:  nil,
		},
		"render patients list page while patients is nil": {
			page:     template.PatientsList,
			data:     nil,
			wantBody: "No patients referred yet",
			wantErr:  nil,
		},
		"render patient view page": {
			page: template.PatientView,
			data: &account.Patient{
				Detail: &account.Detail{
					FirstName: "test",
					LastName:  "user",
				},
				Diagnosis: "test diagnosis 1",
			},
			wantBody: "test user",
			wantErr:  nil,
		},
		"render patient prescription page": {
			page: template.PrescriptionDetail,
			data: &exercise.Prescription{
				Detail: &exercise.Detail{
					Name: "test exs",
				},
			},
			wantBody: "test exs",
			wantErr:  nil,
		},
		"not render page that not exists": {
			page:    "not exist",
			wantErr: template.ErrPageNotExists,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tm, err := template.NewManager()
			if err != nil {
				t.Fatalf("err creating template manager: %v", err)
			}
			err = tm.Render(context, buf, template.PageName(tt.page), tt.data)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if err != nil {
				return
			}
			if got, want := strings.Contains(buf.String(), tt.wantBody), true; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}
