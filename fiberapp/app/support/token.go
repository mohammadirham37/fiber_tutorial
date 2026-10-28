package support

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"fiberapp/app/models"
	"fiberapp/database"
)

// CreateToken membuat token API baru untuk user.
// Setara: $user->createToken('nama')->plainTextToken di Laravel Sanctum.
//
// Token asli hanya dikembalikan SEKALI ke client. Di database kita simpan
// hash-nya saja, jadi kalau database bocor, token tidak bisa dipakai.
func CreateToken(user *models.User, name string) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	plain := hex.EncodeToString(bytes)

	token := models.PersonalAccessToken{
		UserID: user.ID,
		Name:   name,
		Token:  hashToken(plain),
	}
	if err := database.DB.Create(&token).Error; err != nil {
		return "", err
	}
	return plain, nil
}

// FindToken mencari token berdasarkan token asli yang dikirim client.
func FindToken(plain string) (*models.PersonalAccessToken, error) {
	var token models.PersonalAccessToken
	err := database.DB.Preload("User").Where("token = ?", hashToken(plain)).First(&token).Error
	if err != nil {
		return nil, err
	}

	now := time.Now()
	database.DB.Model(&token).Update("last_used_at", &now)
	return &token, nil
}

func hashToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
