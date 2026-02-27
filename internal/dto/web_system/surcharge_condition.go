package web_system

import (
	"encoding/json"

	dto_common "go-structure/internal/dto/common"
)

type (
	CreateSurchargeConditionRequestDto struct {
		Code          string          `json:"code" binding:"required"`
		Name          string          `json:"name" binding:"required"`
		ConditionType string          `json:"condition_type" binding:"required,oneof=time_window weather traffic holiday"`
		Config        json.RawMessage `json:"config" binding:"required"`
		IsActive      bool            `json:"is_active"`
	}

	UpdateSurchargeConditionRequestDto struct {
		Code          string          `json:"code" binding:"required"`
		Name          string          `json:"name" binding:"required"`
		ConditionType string          `json:"condition_type" binding:"required,oneof=time_window weather traffic holiday"`
		Config        json.RawMessage `json:"config" binding:"required"`
		IsActive      bool            `json:"is_active"`
	}

	SurchargeConditionItemDto struct {
		ID            string          `json:"id"`
		Code          string          `json:"code"`
		Name          string          `json:"name"`
		ConditionType string          `json:"condition_type"`
		Config        json.RawMessage `json:"config"`
		IsActive      bool            `json:"is_active"`
	}

	ListSurchargeConditionsResponseDto struct {
		Items      []SurchargeConditionItemDto `json:"items"`
		Pagination dto_common.PaginationMeta   `json:"pagination"`
	}
)

// Template config for API docs / client (only for reference, validation is in model).
const (
	TimeWindowTemplate    = `{"from": "17:00", "to": "19:00", "days": ["mon","tue","wed","thu","fri"]}`
	WeatherRainTemplate   = `{"rain": true}`
	WeatherRainMMTemplate = `{"rain_mm": {"operator": ">=", "value": 10}}`
	TrafficTemplate       = `{"level": "high"}`
	HolidayTemplate       = `{"dates": ["2025-12-29", "2025-12-30", "2026-01-01"]}`
)
