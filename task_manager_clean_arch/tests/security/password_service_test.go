package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/infrastructure/security"
	"golang.org/x/crypto/bcrypt"
)

// PasswordSecuritySuite defines the test suite for password security functions
type PasswordSecuritySuite struct {
	suite.Suite
}

// TestPasswordSecuritySuite runs the test suite
func TestPasswordSecuritySuite(t *testing.T) {
	suite.Run(t, new(PasswordSecuritySuite))
}

// TestHashPassword tests the HashPassword function
func (s *PasswordSecuritySuite) TestHashPassword() {
	s.Run("Success", func() {
		password := "Abebe123"
		hashedPassword, err := security.HashPassword(password)

		s.NoError(err)
		s.NotEmpty(hashedPassword)
		// Verify the hash is valid by checking it with bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		s.NoError(err)
	})

	s.Run("EmptyPassword", func() {
		hashedPassword, err := security.HashPassword("")

		s.NoError(err) // bcrypt allows empty passwords
		s.NotEmpty(hashedPassword)
		// Verify the hash is valid
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(""))
		s.NoError(err)
	})

	s.Run("LongPassword", func() {
		// bcrypt has a max input length of 72 bytes; test with a long password
		password := strings.Repeat("a", 100)
		truncatedPassword := password[:72]
		hashedPassword, err := security.HashPassword(password)

		s.Error(err)
		s.Empty(hashedPassword)

		hashedPassword, err = security.HashPassword(truncatedPassword)

		s.NoError(err)
		s.NotEmpty(hashedPassword)
		// bcrypt truncates passwords > 72 bytes; verify with first 72 bytes
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password[:72]))
		s.NoError(err)
	})
}

// TestValidatePassword tests the ValidatePassword function
func (s *PasswordSecuritySuite) TestValidatePassword() {
	s.Run("Success", func() {
		password := "Abebe123"
		hashedPassword, err := security.HashPassword(password)
		s.NoError(err)

		result := security.ValidatePassword(hashedPassword, password)

		s.True(result)
	})

	s.Run("WrongPassword", func() {
		password := "Abebe123"
		wrongPassword := "Kebede123"
		hashedPassword, err := security.HashPassword(password)
		s.NoError(err)

		result := security.ValidatePassword(hashedPassword, wrongPassword)

		s.False(result)
	})

	s.Run("EmptyPassword", func() {
		password := ""
		hashedPassword, err := security.HashPassword(password)
		s.NoError(err)

		result := security.ValidatePassword(hashedPassword, "")

		s.True(result)
	})

	s.Run("EmptyHashedPassword", func() {
		result := security.ValidatePassword("", "Abebe123")

		s.False(result) // Invalid hash format
	})

	s.Run("InvalidHashFormat", func() {
		invalidHash := "invalid_hash_format"
		result := security.ValidatePassword(invalidHash, "Abebe123")

		s.False(result) // bcrypt will return an error for invalid hash
	})
}
