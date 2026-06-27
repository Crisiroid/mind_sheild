package database

import (
	"log"
	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.DailyCalendar{},

		&models.EmotionTriangleInteraction{},
		&models.StressEvent{},
		&models.BodyTensionMap{},

		&models.BreathingSession{},
		&models.CognitiveErrorGame{},
		&models.MentalMust{},
		&models.NegativeThought{},
		&models.MindCourtEvidence{},
		&models.ConflictExercise{},
		&models.MoodTracker{},
		&models.RoleAndValue{},
		&models.SkyThought{},

		&models.MindfulTimer{},
		&models.AcceptanceExercise{},
		&models.WeeklyReport{},

		&models.WeeklyMediaContent{},

		&models.AdminUser{},
		&models.AdminRole{},
		&models.UserReport{},
		&models.SystemLog{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
