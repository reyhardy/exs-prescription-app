package exercise

import "time"

type Detail struct {
	Name             string
	Category         string
	PrimaryMuscles   []string
	SecondaryMuscles []string
	Instructions     []string
	ImagePath        string
}

type Prescription struct {
	*Detail
	Reps      int
	Sets      int
	CreatedAt time.Time
	UpdatedAt time.Time
}
