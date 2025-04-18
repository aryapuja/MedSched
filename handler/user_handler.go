package handler

import (
	"MedSched/models"
	"MedSched/services"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// @Summary Login a user and generate a JWT token
// @Description This endpoint allows a registered user to log in with their username and password, and receive a JWT token for authentication.
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "User login data"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func LoginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind JSON request body
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := services.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Ambil data pengguna dari database
	user, err := services.GetUserByUsername(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Mengembalikan token JWT
	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"full_name": user.FullName,
		"role":      user.Role,
	})
}

// @Summary Register a new user
// @Description This endpoint allows a new user to register with their username, email, password, full name, gender, and birth date.
// @Accept json
// @Produce json
// @Param user body models.User true "User registration data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if user.Role == "" {
		user.Role = "pasien"
	}

	// Validasi kolom diisi
	if user.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Full name is required"})
		return
	}
	if user.Gender == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gender is required"})
		return
	}
	if user.BirthDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Birth date is required"})
		return
	}
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	// Validasi email format menggunakan regex
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validasi panjang password
	if len(user.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	// Cek apakah username atau email sudah terdaftar
	if err := services.RegisterUser(&user); err != nil {
		// Jika ada error saat pendaftaran, periksa jenis error
		if err.Error() == "username already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already registered"})
			return
		}
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		// Log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// @Summary Get all users
// @Description This endpoint allows an admin to view all registered users in the system.
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	// Panggil service untuk mengambil semua pengguna
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error fetching users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Get user data based on query parameters
// @Description This endpoint allows an admin to view users based on either username or role.
// @Accept json
// @Produce json
// @Param username query string false "Username of the user"
// @Param role query string false "Role of the user"
// @Success 200 {array} models.User "List of users"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 404 {object} models.ErrorResponse "User(s) not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users/by? [get]
func GetUsersByQuery(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	role := c.DefaultQuery("role", "")

	// ambil data pengguna berdasarkan query
	if username != "" {
		// Ambil user berdasarkan username
		users, err := services.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"User not found"})
			return
		}
		c.JSON(http.StatusOK, users)
	} else if role != "" {
		// Ambil user berdasarkan role
		users, err := services.GetUserByRole(role)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"Users with the specified role not found"})
			return
		}
		if len(users) == 0 {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"No users found with role: " + role})
			return
		}
		c.JSON(http.StatusOK, users)
	} else {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Username or role must be provided"})
	}
}

// @Summary Update user data
// @Description This endpoint allows an admin to update a user's data by either username or ID using query parameters.
// @Accept json
// @Produce json
// @Param username query string false "Username of the user"
// @Param id query string false "ID of the user"
// @Param user body models.User true "User update data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/by [put]
func UpdateUser(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	id := c.DefaultQuery("id", "")

	// Cek apakah salah satu parameter query ada
	if username == "" && id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Username or ID must be provided"})
		return
	}

	// Periksa ID/Username
	var user *models.User
	var err error
	if username != "" {
		user, err = services.GetUserByUsername(username)
	} else if id != "" {
		user, err = services.GetUserByID(id) // Pastikan kamu memiliki fungsi GetUserByID di service
	}

	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"User not found"})
		return
	}

	// Bind data pengguna untuk diupdate
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Invalid input"})
		return
	}

	// Update data pengguna
	var updateErr error
	if username != "" {
		updateErr = services.UpdateUser(username, &updatedUser)
	} else if id != "" {
		updateErr = services.UpdateUser(id, &updatedUser)
	}

	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error updating user"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "User updated successfully"})
}

// @Summary Get the current user's profile
// @Description This endpoint allows the user to view their own profile.
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /users/me [get]
func GetMyProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := services.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update the current user's profile
// @Description This endpoint allows the user to update their own profile.
// @Accept json
// @Produce json
// @Param user body models.User true "User update data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/me [put]
func UpdateMyProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Invalid input"})
		return
	}

	// Update data pengguna
	if err := services.UpdateUser(username.(string), &user); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error updating user"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "User updated successfully"})
}

// @Summary Delete a user
// @Description This endpoint allows an admin to delete a user by either username or ID using query parameters.
// @Accept json
// @Produce json
// @Param username query string false "Username of the user"
// @Param id query string false "ID of the user"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/by [delete]
func DeleteUser(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	id := c.DefaultQuery("id", "")

	// Cek apakah salah satu parameter query ada
	if username == "" && id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{"Username or ID must be provided"})
		return
	}

	// Periksa ID/Username
	var user *models.User
	var err error
	if username != "" {
		user, err = services.GetUserByUsername(username)
	} else if id != "" {
		user, err = services.GetUserByID(id) // Pastikan kamu memiliki fungsi GetUserByID di service
	}

	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{"User not found"})
		return
	}

	if username != "" {
		// Hapus berdasarkan username
		err = services.DeleteUser(username)
	} else if id != "" {
		// Hapus berdasarkan ID
		err = services.DeleteUser(id)
	}

	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{"User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{"Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "User deleted successfully"})
}
