package account

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

var (
	ErrInvalidParams = errors.New("invalid params")
	ErrEmptyRoles    = errors.New("roles should not be empty")
)

type Detail struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Birthday  time.Time
	Age       int
	Gender    Gender
	Roles     Roles
	Status    Status

	LastLoginAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *Detail) IsActive() bool {
	if time.Since(u.LastLoginAt) > 90*24*time.Hour {
		return false
	}
	return true
}

func (u *Detail) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

type DetailParams struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	BirthDay  time.Time
	Gender    string
	Roles     []string
}

func New(p *DetailParams) (*Detail, error) {
	id := ksuid.New()
	p.Id = id.String()
	return Map(p)
}

func Map(p *DetailParams) (*Detail, error) {
	if !p.validateParams() {
		return nil, ErrInvalidParams
	}
	var roles []Role
	for _, role := range p.Roles {
		roles = append(roles, Role(role))
	}
	age := int(time.Since(p.BirthDay).Hours() / 24 / 365)
	return &Detail{
		Id:        p.Id,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
		Password:  p.Password,
		Birthday:  p.BirthDay,
		Age:       age,
		Gender:    Gender(p.Gender),
		Roles:     roles,
		Status:    StatusActive,
	}, nil
}

func (p *DetailParams) validateParams() bool {
	if p.Id == "" || p.FirstName == "" || p.Email == "" || p.BirthDay.IsZero() {
		return false
	}
	return true
}

type Role string

const (
	RolePatient   Role = "patient"
	RoleTherapist Role = "therapist"
	RoleAdmin     Role = "admin"
)

func (r Role) String() string {
	return string(r)
}

type Roles []Role

func (r *Roles) Scan(src any) error {
	if src == nil {
		return ErrEmptyRoles
	}

	switch v := src.(type) {
	case string:
		roleSlc := strings.Split(string(v), ",")
		for _, roleStr := range roleSlc {
			*r = append(*r, Role(roleStr))
		}
	case []byte:
		roleBits := bytes.Split(v, []byte(","))
		for _, roleBit := range roleBits {
			*r = append(*r, Role(string(roleBit)))
		}
	default:
		return ErrEmptyRoles
	}

	return nil
}

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Gender string

func (g Gender) String() string {
	return string(g)
}

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "femal"
)
