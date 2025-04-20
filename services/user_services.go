package services

import (
	"MedSched/config"
	"MedSched/models"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser
func RegisterUser(user *models.User) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	// Cek apakah username sudah terdaftar
	usernameExists, err := CheckIfUsernameExists(user.Username)
	if err != nil {
		return fmt.Errorf("Gagal memeriksa ketersediaan username: %v", err)
	}
	if usernameExists {
		return errors.New("Username sudah terdaftar")
	}

	// Cek apakah email sudah terdaftar
	emailExists, err := CheckIfEmailExists(user.Email)
	if err != nil {
		return fmt.Errorf("Gagal memeriksa ketersediaan email: %v", err)
	}
	if emailExists {
		return errors.New("Email sudah terdaftar")
	}

	// Enkripsi password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Gagal mengenkripsi password: %v", err)
	}
	user.Password = string(hashedPassword)

	// Simpan ke database
	sqlStatement := `
        INSERT INTO users (username, email, password, full_name, gender, birth_date, role)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = db.QueryRow(sqlStatement, user.Username, user.Email, user.Password, user.FullName, user.Gender, user.BirthDate, user.Role).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("Gagal menyimpan data pengguna ke database: %v", err)
	}

	return nil
}

// GetAllUsers
func GetAllUsers() ([]models.User, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	var users []models.User
	rows, err := db.Query(`SELECT id, username, full_name, email, gender, birth_date, role FROM users`)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengambil data pengguna: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Email, &user.Gender, &user.BirthDate, &user.Role); err != nil {
			return nil, fmt.Errorf("Gagal memindai data pengguna: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByUsername
func GetUserByUsername(username string) (*models.User, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	var user models.User
	sqlStatement := `SELECT id, username, password, email, role, full_name, gender, birth_date FROM users WHERE username=$1`
	err = db.QueryRow(sqlStatement, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.FullName, &user.Gender, &user.BirthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Pengguna tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID
func GetUserByID(id string) (*models.User, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user models.User
	sqlStatement := `SELECT id, username, password, email, role, full_name, gender, birth_date FROM users WHERE id=$1`
	err = db.QueryRow(sqlStatement, id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.FullName, &user.Gender, &user.BirthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByRole
func GetUserByRole(role string) ([]models.User, error) {
	db, err := config.SetupDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var users []models.User
	sqlStatement := `SELECT id, username, password, email, role, full_name, gender, birth_date FROM users WHERE role=$1`
	rows, err := db.Query(sqlStatement, role)
	if err != nil {
		return nil, fmt.Errorf("kesalahan saat mengambil data users: %v", err)
	}
	defer rows.Close()

	// Looping untuk mengisi slice users
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.FullName, &user.Gender, &user.BirthDate); err != nil {
			return nil, fmt.Errorf("kesalahan saat mengambil user: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func isID(identifier string) bool {
	_, err := strconv.Atoi(identifier)
	return err == nil
}

// UpdateUser
func UpdateUser(identifier string, updatedUser *models.User) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	// Tentukan query berdasarkan username atau id
	var sqlStatement string
	if isID(identifier) {
		// Update berdasarkan ID
		sqlStatement = `UPDATE users SET username=$1, email=$2, full_name=$3, gender=$4, birth_date=$5, role=$6 WHERE id=$7`
		_, err = db.Exec(sqlStatement, updatedUser.Username, updatedUser.Email, updatedUser.FullName, updatedUser.Gender, updatedUser.BirthDate, updatedUser.Role, identifier)
	} else {
		// Update berdasarkan username
		sqlStatement = `UPDATE users SET username=$1, email=$2, full_name=$3, gender=$4, birth_date=$5, role=$6 WHERE username=$7`
		_, err = db.Exec(sqlStatement, updatedUser.Username, updatedUser.Email, updatedUser.FullName, updatedUser.Gender, updatedUser.BirthDate, updatedUser.Role, identifier)
	}

	if err != nil {
		return fmt.Errorf("Gagal memperbarui pengguna: %v", err)
	}

	return nil
}

// DeleteUser
func DeleteUser(identifier string) error {
	db, err := config.SetupDatabase()
	if err != nil {
		return fmt.Errorf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	// Tentukan query berdasarkan username atau id
	var sqlStatement string
	if isID(identifier) {
		// Hapus berdasarkan ID
		sqlStatement = `DELETE FROM users WHERE id=$1`
		_, err = db.Exec(sqlStatement, identifier)
	} else {
		// Hapus berdasarkan username
		sqlStatement = `DELETE FROM users WHERE username=$1`
		_, err = db.Exec(sqlStatement, identifier)
	}

	if err != nil {
		return fmt.Errorf("Gagal menghapus pengguna: %v", err)
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
			return "", errors.New("username atau password tidak valid")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("password tidak valid")
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateJWT
func generateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"userID":   user.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // exp 24 jam
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("SECRET_KEY tidak di set di environment variables")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Menandatangani token dengan secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

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

//RENCANA LANJUTKAN KALAU ADA WAKTU: UPDATE SESUAI FIELD YANG ADA DI BODY JSON, KALAU TIDAK ADA = TIDAK TERUPDATE
// UpdateUser - Update user data based on username or id
// func UpdateUser(identifier string, updatedUser *models.User) error {
// 	db, err := config.SetupDatabase()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	var sqlStatement string
// 	var params []interface{}
// 	var count int

// 	if isID(identifier) {
// 		// Update berdasarkan ID
// 		sqlStatement = `UPDATE users SET`
// 	} else {
// 		// Update berdasarkan username
// 		sqlStatement = `UPDATE users SET`
// 	}

// 	// Hanya update kolom yang ada dalam body JSON
// 	if updatedUser.Username != "" {
// 		sqlStatement += " username=$" + fmt.Sprintf("%d", count+1) + ","
// 		params = append(params, updatedUser.Username)
// 		count++
// 	}
// 	if updatedUser.Email != "" {
// 		sqlStatement += " email=$" + fmt.Sprintf("%d", count+1) + ","
// 		params = append(params, updatedUser.Email)
// 		count++
// 	}
// 	if updatedUser.FullName != "" {
// 		sqlStatement += " full_name=$" + fmt.Sprintf("%d", count+1) + ","
// 		params = append(params, updatedUser.FullName)
// 		count++
// 	}
// 	if updatedUser.Gender != "" {
// 		sqlStatement += " gender=$" + fmt.Sprintf("%d", count+1) + ","
// 		params = append(params, updatedUser.Gender)
// 		count++
// 	}
// 	if updatedUser.BirthDate != "" {
// 		sqlStatement += " birth_date=$" + fmt.Sprintf("%d", count+1) + ","
// 		params = append(params, updatedUser.BirthDate)
// 		count++
// 	}
// 	if updatedUser.Role != "" {
// 		sqlStatement += " role=$" + fmt.Sprintf("%d", count+1) + ","
// 		params = append(params, updatedUser.Role)
// 		count++
// 	}

// 	sqlStatement = sqlStatement[:len(sqlStatement)-1] // Menghapus koma terakhir

// 	if isID(identifier) {
// 		sqlStatement += ` WHERE id=$` + fmt.Sprintf("%d", count+1)
// 		params = append(params, identifier)
// 	} else {
// 		sqlStatement += ` WHERE username=$` + fmt.Sprintf("%d", count+1)
// 		params = append(params, identifier)
// 	}

// 	_, err = db.Exec(sqlStatement, params...)
// 	if err != nil {
// 		return fmt.Errorf("Error updating user: %v", err)
// 	}

// 	return nil
// }
