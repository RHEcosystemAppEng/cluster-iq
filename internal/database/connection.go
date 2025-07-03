package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func Connect(dbURL string, logger *zap.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		logger.Error("Database connection failed", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("Database ping failed after connection", zap.Error(err))
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully")
	return db, nil
}
