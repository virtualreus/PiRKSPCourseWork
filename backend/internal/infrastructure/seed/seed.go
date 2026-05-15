package seed

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/nikitatisenko/pirksp/internal/domain/dto"
	"github.com/nikitatisenko/pirksp/internal/domain/entities"
	"github.com/nikitatisenko/pirksp/internal/domain/ports/repository"
	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/usecase"
)

const bcryptCost = 12

type demoHackathon struct {
	title            string
	shortDescription string
	description      string
	format           string
	status           string
	maxTeam          int
	prizes           string
	daysRegOpen      int
	daysRegClose     int
	daysStart        int
	daysEnd          int
	trackTitle       string
	trackDesc        string
	caseTitle        string
	caseDesc         string
	customer         string
	resourcesURL     string
}

func Run(ctx context.Context, users repository.UsersRepository, hackathons usecase.HackathonUseCase, log *slog.Logger) {
	if os.Getenv("SEED_DISABLE") == "1" {
		return
	}

	organizer, err := ensureUser(ctx, users, "organizer@demo.local", "Демо Организатор", "organizer", "demo12345")
	if err != nil {
		log.Warn("seed: organizer", "err", err)
		return
	}

	_, _ = ensureUser(ctx, users, "participant@demo.local", "Алексей Участников", "participant", "demo12345")

	existing, err := hackathons.ListOrganizer(ctx, organizer.ID)
	if err != nil {
		log.Warn("seed: list hackathons", "err", err)
		return
	}
	if len(existing) > 0 {
		return
	}

	now := time.Now().UTC().Truncate(time.Hour)
	demos := []demoHackathon{
		{
			title:            "Цифровой город 2026",
			shortDescription: "48 часов на прототип для умного города",
			description:      "Соберите команду, выберите кейс от городских заказчиков и сдайте работающий MVP с демо и питчем. Открытые данные, urban-tech и сервисы для жителей.",
			format:           "hybrid",
			status:           "registration",
			maxTeam:          5,
			prizes:           "MacBook для команды-победителя, стажировки в IT-парке",
			daysRegOpen:      0,
			daysRegClose:     14,
			daysStart:        15,
			daysEnd:          17,
			trackTitle:       "Умный город",
			trackDesc:        "Городские сервисы и открытые данные",
			caseTitle:        "Прогноз оттока абонентов",
			caseDesc:         "ML-модель и дашборд на открытых данных оператора связи.",
			customer:         "Городской оператор связи",
			resourcesURL:     "https://data.mos.ru",
		},
		{
			title:            "FinTech Product Sprint",
			shortDescription: "Банковские API и платёжные сценарии за выходные",
			description:      "Кейсы от финтех-партнёров: мгновенные переводы, скоринг, антифрод. Нужен рабочий прототип и сценарий демо для жюри.",
			format:           "online",
			status:           "running",
			maxTeam:          4,
			prizes:           "Грант 300 000 ₽ на пилот с партнёром",
			daysRegOpen:      -10,
			daysRegClose:     -2,
			daysStart:        -1,
			daysEnd:          2,
			trackTitle:       "Open Banking",
			trackDesc:        "Интеграции и пользовательские сценарии",
			caseTitle:        "P2P-перевод за 3 клика",
			caseDesc:         "Упростите UX перевода между физлицами без потери безопасности.",
			customer:         "Необанк «Вектор»",
			resourcesURL:     "https://www.openbanking.org.uk",
		},
		{
			title:            "EcoHack: устойчивое будущее",
			shortDescription: "Экология, ESG и зелёные технологии",
			description:      "Команды ищут решения для мониторинга выбросов, переработки и энергоэффективности. Обязательны метрики impact в питче.",
			format:           "offline",
			status:           "registration",
			maxTeam:          6,
			prizes:           "Поездка на профильную конференцию + менторство 3 месяца",
			daysRegOpen:      2,
			daysRegClose:     21,
			daysStart:        22,
			daysEnd:          24,
			trackTitle:       "Climate Tech",
			trackDesc:        "Датчики, аналитика, отчётность ESG",
			caseTitle:        "Карта качества воздуха",
			caseDesc:         "Визуализация загрязнения по районам в реальном времени.",
			customer:         "Региональный экологический центр",
			resourcesURL:     "https://www.worldbank.org",
		},
		{
			title:            "MedTech Weekend",
			shortDescription: "Цифровая медицина и телемедицина",
			description:      "Хакатон для разработчиков и аналитиков: запись к врачу, triage, напоминания о приёме лекарств. Соблюдайте требования к персональным данным в решении.",
			format:           "hybrid",
			status:           "finished",
			maxTeam:          5,
			prizes:           "Пилот в клинике-партнёре",
			daysRegOpen:      -30,
			daysRegClose:     -20,
			daysStart:        -18,
			daysEnd:          -16,
			trackTitle:       "Patient Journey",
			trackDesc:        "Путь пациента от симптома до выписки",
			caseTitle:        "Умная запись к врачу",
			caseDesc:         "Чат-бот и виджет записи с учётом расписания и специализации.",
			customer:         "Сеть клиник «Здоровье+»",
			resourcesURL:     "https://www.who.int",
		},
	}

	for _, demo := range demos {
		if err := seedHackathon(ctx, hackathons, organizer.ID, now, demo); err != nil {
			log.Warn("seed: hackathon", "title", demo.title, "err", err)
		}
	}

	log.Info("seed: demo data ready", "hackathons", len(demos))
}

func ensureUser(ctx context.Context, users repository.UsersRepository, email, name, role, password string) (entities.User, error) {
	user, err := users.GetByEmail(ctx, email)
	if err == nil {
		return user, nil
	}
	if !errorsIsNotFound(err) {
		return entities.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return entities.User{}, err
	}

	return users.Create(ctx, entities.User{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     name,
		PlatformRole: role,
	})
}

func errorsIsNotFound(err error) bool {
	return errors.Is(err, errs.ErrNotFound)
}

func seedHackathon(ctx context.Context, hackathons usecase.HackathonUseCase, organizerID uuid.UUID, now time.Time, demo demoHackathon) error {
	timeline := dto.HackathonTimeline{
		RegistrationOpensAt:  now.Add(time.Duration(demo.daysRegOpen) * 24 * time.Hour).Format(time.RFC3339),
		RegistrationClosesAt: now.Add(time.Duration(demo.daysRegClose) * 24 * time.Hour).Format(time.RFC3339),
		EventStartsAt:        now.Add(time.Duration(demo.daysStart) * 24 * time.Hour).Format(time.RFC3339),
		EventEndsAt:          now.Add(time.Duration(demo.daysEnd) * 24 * time.Hour).Format(time.RFC3339),
		SubmissionDeadlineAt: now.Add(time.Duration(demo.daysEnd) * 24 * time.Hour).Format(time.RFC3339),
	}

	created, err := hackathons.Create(ctx, organizerID, dto.CreateHackathonRequest{
		Title:            demo.title,
		ShortDescription: demo.shortDescription,
		Description:      demo.description,
		Format:           demo.format,
		Timeline:         timeline,
		MaxTeamSize:      demo.maxTeam,
		PrizesInfo:       demo.prizes,
	})
	if err != nil {
		return err
	}

	hid, err := uuid.Parse(created.ID)
	if err != nil {
		return err
	}

	track, err := hackathons.CreateTrack(ctx, organizerID, hid, dto.CreateTrackRequest{
		Title:       demo.trackTitle,
		Description: demo.trackDesc,
	})
	if err != nil {
		return err
	}

	tid, err := uuid.Parse(track.ID)
	if err != nil {
		return err
	}

	_, err = hackathons.CreateCase(ctx, organizerID, tid, dto.CreateCaseRequest{
		Title:        demo.caseTitle,
		Description:  demo.caseDesc,
		CustomerName: demo.customer,
		ResourcesURL: demo.resourcesURL,
	})
	if err != nil {
		return err
	}

	if _, err := hackathons.Publish(ctx, organizerID, hid); err != nil {
		return err
	}

	if demo.status == "registration" {
		return nil
	}

	status := demo.status
	_, err = hackathons.Update(ctx, organizerID, hid, dto.UpdateHackathonRequest{
		Status: &status,
	})
	return err
}
