package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateMatrixUsername(email string) string {
	// Remove the domain part and replace @ with _
	username := strings.Split(email, "@")[0]
	username = strings.ReplaceAll(username, ".", "_")
	return fmt.Sprintf("@%s:your_matrix_domain.com", username)
}

func GenerateSecurePassword() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
