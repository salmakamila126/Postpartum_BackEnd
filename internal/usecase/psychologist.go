package usecase

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/repository"
	"Postpartum_BackEnd/internal/seed"
	"Postpartum_BackEnd/pkg/cache"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PsychologistUsecase struct {
	Repo  *repository.Repository
	Cache *cache.Cache
}

func NewPsychologistUsecase(repo *repository.Repository, appCache *cache.Cache) *PsychologistUsecase {
	return &PsychologistUsecase{Repo: repo, Cache: appCache}
}

func (u *PsychologistUsecase) SeedIfEmpty() error {
	count, err := u.Repo.PsychologistRepository.CountAll()
	if err != nil {
		return fmt.Errorf("failed to count psychologists: %w", err)
	}
	if count > 0 {
		return nil
	}
	if err := u.Repo.PsychologistRepository.SeedAll(seed.PsychologistSeedData()); err != nil {
		return fmt.Errorf("failed to seed psychologists: %w", err)
	}
	return nil
}

func (u *PsychologistUsecase) GetAll() ([]dto.PsychologistListItem, error) {
	ctx := context.Background()
	const cacheKey = "psychologists:list"

	var cached []dto.PsychologistListItem
	if ok, err := u.Cache.GetJSON(ctx, cacheKey, &cached); err == nil && ok {
		return cached, nil
	}

	list, err := u.Repo.PsychologistRepository.FindAll()
	if err != nil {
		return nil, err
	}
	items := make([]dto.PsychologistListItem, 0, len(list))
	for _, p := range list {
		items = append(items, dto.PsychologistListItem{
			ID:           p.ID.String(),
			Name:         p.Name,
			Title:        p.Title,
			Job:          p.Job,
			ExperienceYr: p.ExperienceYr,
			PriceIDR:     p.PriceIDR,
			PhotoURL:     p.PhotoURL,
		})
	}

	if err := u.Cache.SetJSON(ctx, cacheKey, items, 10*time.Minute); err != nil {
		log.Printf("cache set failed for %s: %v", cacheKey, err)
	}
	return items, nil
}

func (u *PsychologistUsecase) GetByID(id uuid.UUID) (*dto.PsychologistDetail, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("psychologists:detail:%s", id.String())

	var cached dto.PsychologistDetail
	if ok, err := u.Cache.GetJSON(ctx, cacheKey, &cached); err == nil && ok {
		return &cached, nil
	}

	p, err := u.Repo.PsychologistRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPsychologistNotFound
		}
		return nil, err
	}

	schedules := make([]dto.ScheduleSlot, 0, len(p.Schedules))
	for _, s := range p.Schedules {
		schedules = append(schedules, dto.ScheduleSlot{
			DayOfWeek: s.DayOfWeek,
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
			Label:     fmt.Sprintf("%s, %s-%s", s.DayOfWeek, s.StartTime, s.EndTime),
		})
	}

	result := &dto.PsychologistDetail{
		ID:           p.ID.String(),
		Name:         p.Name,
		Title:        p.Title,
		Job:          p.Job,
		ExperienceYr: p.ExperienceYr,
		PriceIDR:     p.PriceIDR,
		PhotoURL:     p.PhotoURL,
		Schedules:    schedules,
	}

	if err := u.Cache.SetJSON(ctx, cacheKey, result, 10*time.Minute); err != nil {
		log.Printf("cache set failed for %s: %v", cacheKey, err)
	}
	return result, nil
}

func (u *PsychologistUsecase) UpdatePhotoURL(id uuid.UUID, photoURL string) error {
	if photoURL == "" {
		return errs.New(http.StatusBadRequest, "photo_url must not be empty")
	}
	_, err := u.Repo.PsychologistRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrPsychologistNotFound
		}
		return err
	}
	if err := u.Repo.PsychologistRepository.UpdatePhotoURL(id, photoURL); err != nil {
		return err
	}

	if err := u.Cache.Delete(
		context.Background(),
		"psychologists:list",
		fmt.Sprintf("psychologists:detail:%s", id.String()),
	); err != nil {
		log.Printf("cache delete failed for psychologist id %s: %v", id.String(), err)
	}

	return nil
}

func (u *PsychologistUsecase) BuildBookingWhatsApp(
	psychologistID uuid.UUID,
	userName, userEmail string,
	req dto.BookingWhatsAppRequest,
) (*dto.BookingWhatsAppResponse, error) {

	p, err := u.Repo.PsychologistRepository.FindByID(psychologistID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPsychologistNotFound
		}
		return nil, err
	}

	validSlots := buildSlotSet(p.Schedules)
	if !validSlots[req.SelectedSlot] {
		return nil, errs.New(http.StatusBadRequest, "schedule slot not available: "+req.SelectedSlot)
	}

	adminWhatsAppNumber := os.Getenv("ADMIN_WA_NUMBER")
	if adminWhatsAppNumber == "" {
		return nil, errs.New(http.StatusInternalServerError, "ADMIN_WA_NUMBER is not configured")
	}

	message := fmt.Sprintf(
		"Halo Admin, saya ingin booking konsultasi dengan dokter %s\n\n"+
			"Nama: %s\n"+
			"Email: %s\n"+
			"Keluhan utama: \n"+
			"Riwayat penyakit (jika ada): \n"+
			"Tanggal & jam yang diinginkan (pilih satu)\n- %s\n\n"+
			"Mohon info jadwal yang tersedia dan cara pembayarannya. Terima kasih.",
		p.Name, userName, userEmail, req.SelectedSlot,
	)

	waURL := fmt.Sprintf(
		"https://wa.me/%s?text=%s",
		adminWhatsAppNumber,
		url.QueryEscape(message),
	)

	return &dto.BookingWhatsAppResponse{
		WhatsAppURL:    waURL,
		MessagePreview: message,
	}, nil
}

func buildSlotSet(schedules []entity.PsychologistSchedule) map[string]bool {
	set := make(map[string]bool, len(schedules))
	for _, s := range schedules {
		label := fmt.Sprintf("%s, %s-%s", s.DayOfWeek, s.StartTime, s.EndTime)
		set[label] = true
	}
	return set
}
