package helpers

import (
	"ItDevTest/internal/models"
	"errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
	"os"
	"regexp"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidateUserInput(user *models.User) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SendEmail(toEmail, subject, content string) error {
	from := mail.NewEmail("Ваше имя", "ваш_email@example.com")
	to := mail.NewEmail("Получатель", toEmail)
	message := mail.NewSingleEmail(from, subject, to, content, content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	return err
}
