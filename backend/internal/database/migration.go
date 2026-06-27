package database

import (
	"log"
	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		// User Domain
		&models.User{},
		&models.UserSetting{},
		&models.DailyCalendar{},

		// Emotion & Stress
		&models.EmotionTriangleInteraction{},
		&models.StressEvent{},
		&models.BodyTensionMap{},

		// Exercises & Activities
		&models.BreathingSession{},
		&models.CognitiveErrorGame{},
		&models.MentalMust{},
		&models.NegativeThought{},
		&models.MindCourtEvidence{},
		&models.ConflictExercise{},
		&models.MoodTracker{},
		&models.RoleAndValue{},
		&models.SkyThought{},

		// Mindfulness
		&models.MindfulTimer{},
		&models.AcceptanceExercise{},
		&models.WeeklyReport{},

		// Admin Panel (with schema)
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
