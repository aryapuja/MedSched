package services

import (
	"MedSched/config"
	"MedSched/models"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// CheckIfUsernameExists
func CheckIfUsernameExists(username string) (bool, error) {
	// Setup koneksi database
	db, err := config.SetupDatabase()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	query := `SELECT COUNT(*) FROM users WHERE username=$1`
	err = db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckIfEmailExists
func CheckIfEmailExists(email string) (bool, error) {
	// Setup koneksi database
	db, err := config.SetupDatabase()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	query := `SELECT COUNT(*) FROM users WHERE email=$1`
	err = db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// RegisterUser
func RegisterUser(user *models.User) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return fmt.Errorf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Cek username
	usernameExists, err := CheckIfUsernameExists(user.Username)
	if err != nil {
		return fmt.Errorf("Error checking username existence: %v", err)
	}
	if usernameExists {
		return errors.New("username already registered")
	}

	// Cek email
	emailExists, err := CheckIfEmailExists(user.Email)
	if err != nil {
		return fmt.Errorf("Error checking email existence: %v", err)
	}
	if emailExists {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error hashing password: %v", err)
	}
	user.Password = string(hashedPassword)

	// Simpan ke database
	sqlStatement := `
		INSERT INTO users (username, email, password, full_name, gender, birth_date, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = db.QueryRow(sqlStatement, user.Username, user.Email, user.Password, user.FullName, user.Gender, user.BirthDate, user.Role).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("Error inserting user into database: %v", err)
	}

	return nil
}

// AuthenticateUser - JWT token
func AuthenticateUser(username, password string) (string, error) {
	// Setup koneksi database
	db, err := config.SetupDatabase()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var user models.User
	query := `SELECT id, username, password, role FROM users WHERE username=$1`
	err = db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByUsername
func GetUserByUsername(username string) (*models.User, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user models.User
	sqlStatement := `SELECT id, username, password, email, role, full_name, gender, birth_date FROM users WHERE username=$1`
	err = db.QueryRow(sqlStatement, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.FullName, &user.Gender, &user.BirthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// generateJWT
func generateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // exp 24 jam
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("SECRET_KEY is not set in the environment variables")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Menandatangani token dengan secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
