package Tests

import (
	"github.com/google/uuid"
	"time"
)

type TestCaseRelationModel struct {
	Id              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;uniqueIndex"`
	TestCaseModelId uuid.UUID      `json:"test_case_model_id" gorm:"type:uuid"`
	TestCaseModel   *TestCaseModel `json:"test_case_model" gorm:"foreignKey:id;references:test_case_model_id"`
	Value           string         `json:"value"`
	CreatedAt       time.Time      `json:"created_at" gorm:"index;not null"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"not null"`
}
