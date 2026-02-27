package web_system

import dto_common "go-structure/internal/dto/common"

type (
	CreateSurchargeRuleRequestDto struct {
		ServiceID     string   `json:"service_id" binding:"required,uuid"`
		ZoneID        string   `json:"zone_id" binding:"required,uuid"`
		Amount        float64  `json:"amount" binding:"required,gte=0"`
		Unit          string   `json:"unit" binding:"required"` // 'percent' | 'fixed'
		Priority      int      `json:"priority" binding:"gte=0"`
		ConditionIDs  []string `json:"condition_ids" binding:"omitempty,dive,uuid"`
		IsActive      bool     `json:"is_active"`
	}

	UpdateSurchargeRuleRequestDto struct {
		ServiceID    string   `json:"service_id" binding:"required,uuid"`
		ZoneID       string   `json:"zone_id" binding:"required,uuid"`
		Amount       float64  `json:"amount" binding:"required,gte=0"`
		Unit         string   `json:"unit" binding:"required"`
		Priority     int      `json:"priority" binding:"gte=0"`
		ConditionIDs []string `json:"condition_ids" binding:"omitempty,dive,uuid"`
		IsActive     bool     `json:"is_active"`
	}

	SurchargeRuleServiceDto struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	}

	SurchargeRuleZoneDto struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	}

	SurchargeRuleItemDto struct {
		ID           string                      `json:"id"`
		Amount       float64                     `json:"amount"`
		Unit         string                      `json:"unit"`
		Priority     int                         `json:"priority"`
		ConditionIDs []string                    `json:"condition_ids"`
		Conditions   []SurchargeConditionItemDto `json:"conditions"`
		Service      *SurchargeRuleServiceDto    `json:"service,omitempty"`
		Zone         *SurchargeRuleZoneDto       `json:"zone,omitempty"`
		IsActive     bool                        `json:"is_active"`
	}

	ListSurchargeRulesResponseDto struct {
		Items      []SurchargeRuleItemDto  `json:"items"`
		Pagination dto_common.PaginationMeta `json:"pagination"`
	}
)
