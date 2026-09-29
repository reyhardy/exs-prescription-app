package handler

import (
	"context"
	"net/http"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/ctx"
	"github.com/reyhardy/exs-prescription-app/internal/template"
)

type userGetter interface {
	GetUser(context.Context, string) (*account.Detail, error)
}

type User struct {
	us userGetter
	tr template.Renderer
}

func NewUser(us userGetter, tr template.Renderer) *User {
	return &User{
		us: us,
		tr: tr,
	}
}

func (h User) UserProfilePage(w http.ResponseWriter, r *http.Request) {
	uCtx := ctx.GetContext(r.Context())
	if uCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	u, err := h.us.GetUser(r.Context(), uCtx.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.tr.Render(r.Context(), w, template.UserProfile, u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// func (h *User) RegisterUser(w http.ResponseWriter, r *http.Request) {
// 	birthday, err := time.Parse("2006-01-02", r.FormValue("birthday"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
// 		return
// 	}
// 	params := &user.DetailParams{
// 		FirstName: r.FormValue("first_name"),
// 		LastName:  r.FormValue("last_name"),
// 		Email:     r.FormValue("email"),
// 		Password:  r.FormValue("password"),
// 		BirthDay:  birthday,
// 		Gender:    r.FormValue("gender"),
// 		Role:      r.FormValue("role"),
// 	}
// 	if err := h.s.CreateNewUser(r.Context(), params); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintln(w, "user created")
// }
