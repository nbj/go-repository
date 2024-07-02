package Tests

import (
	"github.com/google/uuid"
	"time"
)

type TestCaseModel struct {
	Id                     uuid.UUID                `json:"id" gorm:"type:uuid;primaryKey;uniqueIndex"`
	Value                  string                   `json:"value"`
	TestCaseRelationModels []*TestCaseRelationModel `json:"test_case_relation_models"`
	CreatedAt              time.Time                `json:"created_at" gorm:"index;not null"`
	UpdatedAt              time.Time                `json:"updated_at" gorm:"not null"`
}

func (model *TestCaseModel) With() []string {
	return []string{
		"TestCaseRelationModels",
	}
}
