package account

import "time"

type Therapist struct {
	*Detail
	IsAPCValid   bool
	APCNumber    string
	APCExpiredAt time.Time

	ExercisesPrescirbed []string
	patientInCharg      []Patient
}
