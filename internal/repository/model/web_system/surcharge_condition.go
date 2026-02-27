package websystem

import (
	"fmt"

	"github.com/google/uuid"
)

type SurchargeCondition struct {
	ID            uuid.UUID `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	ConditionType string    `json:"condition_type"`
	Config        []byte    `json:"config"`
	IsActive      bool      `json:"is_active"`
}

func (surchargeCondition *SurchargeCondition) ValidateConfig() error {
	if surchargeCondition == nil {
		return fmt.Errorf("surcharge condition is nil")
	}
	return ValidateConditionConfig(surchargeCondition.ConditionType, surchargeCondition.Config)
}
