package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/ctx"
	"github.com/reyhardy/exs-prescription-app/internal/template"
)

type patientFetcher interface {
	FetchAllPatients(context.Context, string) ([]*account.Patient, error)
	FetchPatient(ctx context.Context, id string) (*account.Patient, error)
}

type Patient struct {
	pf patientFetcher
	tr template.Renderer
}

func NewPatient(pf patientFetcher, tr template.Renderer) *Patient {
	return &Patient{
		pf: pf,
		tr: tr,
	}
}

func (h *Patient) ListAllPatients(w http.ResponseWriter, r *http.Request) {
	uCtx := ctx.GetContext(r.Context())
	if uCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	pts, err := h.pf.FetchAllPatients(r.Context(), uCtx.Id)
	if err != nil {
		http.Error(w, "failed to fetch patients", http.StatusInternalServerError)
		return
	}
	if err := h.tr.Render(r.Context(), w, template.PatientsList, pts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Patient) ViewPatientInfo(w http.ResponseWriter, r *http.Request) {
	ptId := r.PathValue("ptId")
	pt, err := h.pf.FetchPatient(r.Context(), ptId)
	if err != nil {
		slog.Error(
			"failed to fetch patient",
			"id", ptId,
			"error", err.Error(),
		)
		http.Error(w, "failed to fetch patient", http.StatusInternalServerError)
		return
	}
	if err := h.tr.Render(r.Context(), w, template.PatientView, pt); err != nil {
		slog.Error(
			"failed render template",
			"template", template.PatientView,
			"error", err.Error(),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
