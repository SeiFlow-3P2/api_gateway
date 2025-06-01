package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	appMode := os.Getenv("APP_MODE")
	if appMode == "" || appMode != "dev" && appMode != "prod" {
		return fmt.Errorf("APP_MODE is not set or invalid (dev/prod)")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return fmt.Errorf("PORT is not set")
	}

	boardServiceAddr := os.Getenv("BOARD_SERVICE_ADDR")
	if boardServiceAddr == "" {
		return fmt.Errorf("BOARD_SERVICE_ADDR is not set")
	}

	paymentServiceAddr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if paymentServiceAddr == "" {
		return fmt.Errorf("PAYMENT_SERVICE_ADDR is not set")
	}

	calendarServiceAddr := os.Getenv("CALENDAR_SERVICE_ADDR")
	if calendarServiceAddr == "" {
		return fmt.Errorf("CALENDAR_SERVICE_ADDR is not set")
	}

	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authServiceAddr == "" {
		return fmt.Errorf("AUTH_SERVICE_ADDR is not set")
	}

	otelAddr := os.Getenv("OTEL_ADDR")
	if otelAddr == "" {
		return fmt.Errorf("OTEL_ADDR is not set")
	}

	return nil
}

func IsProd() bool {
	return os.Getenv("APP_MODE") == "prod"
}

func GetPort() string {
	return os.Getenv("PORT")
}

func GetBoardServiceAddr() string {
	return os.Getenv("BOARD_SERVICE_ADDR")
}

func GetPaymentServiceAddr() string {
	return os.Getenv("PAYMENT_SERVICE_ADDR")
}

func GetCalendarServiceAddr() string {
	return os.Getenv("CALENDAR_SERVICE_ADDR")
}

func GetAuthServiceAddr() string {
	return os.Getenv("AUTH_SERVICE_ADDR")
}

func GetOtelAddr() string {
	return os.Getenv("OTEL_ADDR")
}
