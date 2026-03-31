package usecase

import (
	"Postpartum_BackEnd/internal/domain/symptom"
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/pkg/timeutil"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func mapToEntity(userID uuid.UUID, req dto.CreateSymptomRequest) (*entity.SymptomLog, error) {
	physicalJSON, err := json.Marshal(req.Physical)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal physical condition: %w", err)
	}

	log := &entity.SymptomLog{
		UserID:       userID,
		PhysicalData: string(physicalJSON),
	}

	for _, b := range req.Bleedings {
		log.Bleedings = append(log.Bleedings, entity.BleedingLog{
			PadUsage:   b.PadUsage,
			ClotSize:   b.ClotSize,
			BloodColor: b.BloodColor,
			Smell:      b.Smell,
		})
	}

	for _, m := range req.Moods {
		log.Moods = append(log.Moods, entity.MoodLog{Type: m})
	}

	return log, nil
}

func toDomainInput(req dto.CreateSymptomRequest) symptom.SymptomInput {
	var bleedings []symptom.BleedingInput
	for _, b := range req.Bleedings {
		bleedings = append(bleedings, symptom.BleedingInput{
			PadUsage:   b.PadUsage,
			ClotSize:   b.ClotSize,
			BloodColor: b.BloodColor,
			Smell:      b.Smell,
		})
	}
	return symptom.SymptomInput{
		Bleedings: bleedings,
		Physical: symptom.PhysicalInput{
			Temperature:    req.Physical.Temperature,
			Dizziness:      safeInt(req.Physical.Dizziness),
			Headache:       safeInt(req.Physical.Headache),
			Weakness:       safeInt(req.Physical.Weakness),
			CalfPain:       safeInt(req.Physical.CalfPain),
			AbPain:         safeInt(req.Physical.AbPain),
			Wound:          req.Physical.Wound,
			UrineProblems:  req.Physical.UrineProblems,
			UrineColor:     req.Physical.UrineColor,
			BreastProblems: req.Physical.BreastProblems,
			Swelling:       req.Physical.Swelling,
			OtherSymptoms:  req.Physical.OtherSymptoms,
		},
		Moods: req.Moods,
	}
}

func safeInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func ToResponse(alert *symptom.AlertResult) dto.AlertResponse {
	issues := make([]dto.IssueResult, 0, len(alert.Issues))
	for _, i := range alert.Issues {
		syms := i.Symptoms
		if syms == nil {
			syms = []string{}
		}
		issues = append(issues, dto.IssueResult{
			Code:        i.Code,
			Disease:     i.Disease,
			Level:       string(i.Level),
			Description: i.Description,
			Symptoms:    syms,
		})
	}
	return dto.AlertResponse{
		Level:      string(alert.Level),
		Confidence: alert.Confidence,
		Issues:     issues,
	}
}

func applyAlertSnapshot(log *entity.SymptomLog, alert *symptom.AlertResult) error {
	resp := ToResponse(alert)
	raw, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal alert snapshot: %w", err)
	}

	log.LastAlertLevel = resp.Level
	log.LastAlertData = string(raw)
	return nil
}

func parseAlertSnapshot(l entity.SymptomLog) (dto.AlertResponse, error) {
	if l.LastAlertData == "" {
		return dto.AlertResponse{
			Level:      "safe",
			Confidence: 0,
			Issues:     []dto.IssueResult{},
		}, nil
	}

	var alert dto.AlertResponse
	if err := json.Unmarshal([]byte(l.LastAlertData), &alert); err != nil {
		return dto.AlertResponse{}, fmt.Errorf("failed to unmarshal alert snapshot for symptom log %s: %w", l.ID.String(), err)
	}

	if alert.Issues == nil {
		alert.Issues = []dto.IssueResult{}
	}
	return alert, nil
}

func toHistoryResponse(logs []entity.SymptomLog) ([]dto.SymptomHistoryItem, error) {
	items := make([]dto.SymptomHistoryItem, 0, len(logs))
	for _, l := range logs {
		item, err := mapLogToHistoryItem(l)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func mapLogToHistoryItem(l entity.SymptomLog) (dto.SymptomHistoryItem, error) {
	bleedings := make([]dto.BleedingDetail, 0, len(l.Bleedings))
	for _, b := range l.Bleedings {
		bleedings = append(bleedings, dto.BleedingDetail{
			PadUsage:   b.PadUsage,
			ClotSize:   b.ClotSize,
			BloodColor: b.BloodColor,
			Smell:      b.Smell,
		})
	}

	moods := make([]string, 0, len(l.Moods))
	for _, m := range l.Moods {
		moods = append(moods, m.Type)
	}

	var physical dto.PhysicalCondition
	if l.PhysicalData != "" {
		if err := json.Unmarshal([]byte(l.PhysicalData), &physical); err != nil {
			return dto.SymptomHistoryItem{}, fmt.Errorf("failed to unmarshal physical data for symptom log %s: %w", l.ID.String(), err)
		}
	}

	alert, err := parseAlertSnapshot(l)
	if err != nil {
		return dto.SymptomHistoryItem{}, err
	}

	return dto.SymptomHistoryItem{
		ID:         l.ID.String(),
		Date:       l.Date.Format(timeutil.DateOnlyFormat),
		IsBackdate: l.IsBackdate,
		Bleedings:  bleedings,
		Physical:   physical,
		Moods:      moods,
		Alert:      alert,
		CreatedAt:  l.CreatedAt.Format(timeutil.RFC3339Format),
		UpdatedAt:  l.UpdatedAt.Format(timeutil.RFC3339Format),
	}, nil
}

func toDetailResponse(l *entity.SymptomLog) (*dto.SymptomDetailResponse, error) {
	item, err := mapLogToHistoryItem(*l)
	if err != nil {
		return nil, err
	}
	return &dto.SymptomDetailResponse{SymptomHistoryItem: item}, nil
}
