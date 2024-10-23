package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"linktree-mohamedfadel-backend/internal/api/handlers"
	"linktree-mohamedfadel-backend/internal/api/middleware"
	"linktree-mohamedfadel-backend/internal/models"
	"linktree-mohamedfadel-backend/internal/services"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type HandlerTestSuite struct {
	suite.Suite
	db          *gorm.DB
	router      *gin.Engine
	userHandler *handlers.UserHandler
	linkHandler *handlers.LinkHandler
	analytics   *handlers.AnalyticsHandler
}

func (s *HandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal(err)
	}

	os.Setenv("JWT_SECRET", "test-secret-key")

	s.db.AutoMigrate(&models.User{}, &models.Link{}, &models.Analytics{})

	userService := services.NewUserService(s.db)
	linkService := services.NewLinkService(s.db)
	analyticsService := services.NewAnalyticsService(s.db)

	s.userHandler = handlers.NewUserHandler(userService)
	s.linkHandler = handlers.NewLinkHandler(linkService)
	s.analytics = handlers.NewAnalyticsHandler(analyticsService)

	s.router = gin.New()
	s.setupRoutes()
}

func (s *HandlerTestSuite) setupRoutes() {
	s.router.POST("/users/signup", s.userHandler.SignUpHandler)
	s.router.POST("/users/login", s.userHandler.LoginHandler)
	s.router.GET("/users/:username", s.userHandler.GetUserProfileInfoHandler)

	protected := s.router.Group("")
	protected.Use(middleware.ValidateJWTFromContext())
	{
		protected.PUT("/users", s.userHandler.UpdateUserHandler)
		protected.DELETE("/users", s.userHandler.DeleteUserHandler)
		protected.POST("/links", s.linkHandler.CreateLinkHandler)
		protected.PUT("/links/:id", s.linkHandler.UpdateLinkHandler)
		protected.DELETE("/links/:id", s.linkHandler.DeleteLinkHandler)
	}

	s.router.POST("/analytics/:id/click", s.analytics.TrackLinkClickHandler)
}

func (s *HandlerTestSuite) SetupTest() {
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Analytics{})
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Link{})
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
}

func (s *HandlerTestSuite) TearDownSuite() {
	os.Unsetenv("JWT_SECRET")
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) makeRequest(method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	s.router.ServeHTTP(w, req)
	return w
}

func (s *HandlerTestSuite) TestSignUpHandler() {
	testCases := []struct {
		name       string
		payload    handlers.SignUpRequest
		wantStatus int
	}{
		{
			name: "Valid Registration",
			payload: handlers.SignUpRequest{
				FullName: "Test User",
				Username: "testuser",
				Bio:      "Test Bio",
				Password: "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Missing Required Fields",
			payload: handlers.SignUpRequest{
				Bio: "Test Bio",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Duplicate Username",
			payload: handlers.SignUpRequest{
				FullName: "Test User 2",
				Username: "testuser",
				Bio:      "Test Bio 2",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			w := s.makeRequest(http.MethodPost, "/users/signup", tc.payload, nil)
			assert.Equal(s.T(), tc.wantStatus, w.Code)
		})
	}
}

func (s *HandlerTestSuite) TestLoginHandler() {
	s.makeRequest(http.MethodPost, "/users/signup", handlers.SignUpRequest{
		FullName: "Test User",
		Username: "testuser",
		Password: "password123",
	}, nil)

	testCases := []struct {
		name       string
		payload    handlers.LoginRequest
		wantStatus int
		checkToken bool
	}{
		{
			name: "Valid Login",
			payload: handlers.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			wantStatus: http.StatusOK,
			checkToken: true,
		},
		{
			name: "Invalid Username",
			payload: handlers.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid Password",
			payload: handlers.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			w := s.makeRequest(http.MethodPost, "/users/login", tc.payload, nil)
			assert.Equal(s.T(), tc.wantStatus, w.Code)

			if tc.checkToken {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), response["token"])
			}
		})
	}
}

func (s *HandlerTestSuite) TestGetUserProfileInfoHandler() {
	s.makeRequest(http.MethodPost, "/users/signup", handlers.SignUpRequest{
		FullName: "Test User",
		Username: "testuser",
		Bio:      "Test Bio",
		Password: "password123",
	}, nil)

	w := s.makeRequest(http.MethodPost, "/users/login", handlers.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}, nil)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	s.makeRequest(http.MethodPost, "/links", handlers.CreateLinkRequest{
		Title: "Test Link 1",
		URL:   "https://example1.com",
	}, map[string]string{"Authorization": "Bearer " + token})

	s.makeRequest(http.MethodPost, "/links", handlers.CreateLinkRequest{
		Title: "Test Link 2",
		URL:   "https://example2.com",
	}, map[string]string{"Authorization": "Bearer " + token})

	testCases := []struct {
		name       string
		username   string
		wantStatus int
		checkData  bool
	}{
		{
			name:       "Valid User Profile",
			username:   "testuser",
			wantStatus: http.StatusOK,
			checkData:  true,
		},
		{
			name:       "Non-existent User",
			username:   "nonexistent",
			wantStatus: http.StatusNotFound,
			checkData:  false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			w := s.makeRequest(http.MethodGet, "/users/"+tc.username, nil, nil)
			assert.Equal(s.T(), tc.wantStatus, w.Code)

			if tc.checkData {
				var response models.User
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), "Test User", response.FullName)
				assert.Equal(s.T(), "Test Bio", response.Bio)
				assert.Equal(s.T(), 2, len(response.Links))
			}
		})
	}
}

func (s *HandlerTestSuite) TestUpdateUserHandler() {
	s.makeRequest(http.MethodPost, "/users/signup", handlers.SignUpRequest{
		FullName: "Test User",
		Username: "testuser",
		Bio:      "Test Bio",
		Password: "password123",
	}, nil)

	w := s.makeRequest(http.MethodPost, "/users/login", handlers.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}, nil)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	testCases := []struct {
		name       string
		payload    models.User
		token      string
		wantStatus int
	}{
		{
			name: "Valid Update",
			payload: models.User{
				FullName: "Updated User",
				Bio:      "Updated Bio",
			},
			token:      token,
			wantStatus: http.StatusOK,
		},
		{
			name: "No Authentication",
			payload: models.User{
				FullName: "Updated User",
				Bio:      "Updated Bio",
			},
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			headers := map[string]string{}
			if tc.token != "" {
				headers["Authorization"] = "Bearer " + tc.token
			}

			w := s.makeRequest(http.MethodPut, "/users", tc.payload, headers)
			assert.Equal(s.T(), tc.wantStatus, w.Code)

			if tc.wantStatus == http.StatusOK {
				getW := s.makeRequest(http.MethodGet, "/users/testuser", nil, nil)
				var updatedUser models.User
				json.Unmarshal(getW.Body.Bytes(), &updatedUser)
				assert.Equal(s.T(), tc.payload.FullName, updatedUser.FullName)
				assert.Equal(s.T(), tc.payload.Bio, updatedUser.Bio)
			}
		})
	}
}

func (s *HandlerTestSuite) TestDeleteUserHandler() {
	s.makeRequest(http.MethodPost, "/users/signup", handlers.SignUpRequest{
		FullName: "Test User",
		Username: "testuser",
		Bio:      "Test Bio",
		Password: "password123",
	}, nil)

	w := s.makeRequest(http.MethodPost, "/users/login", handlers.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}, nil)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	testCases := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "No Authentication",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Valid Delete",
			token:      token,
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			headers := map[string]string{}
			if tc.token != "" {
				headers["Authorization"] = "Bearer " + tc.token
			}

			w := s.makeRequest(http.MethodDelete, "/users", nil, headers)
			assert.Equal(s.T(), tc.wantStatus, w.Code)

			if tc.wantStatus == http.StatusOK {
				getW := s.makeRequest(http.MethodGet, "/users/testuser", nil, nil)
				assert.Equal(s.T(), http.StatusNotFound, getW.Code)
			}
		})
	}
}

func (s *HandlerTestSuite) TestCreateLinkHandler() {
	s.makeRequest(http.MethodPost, "/users/signup", handlers.SignUpRequest{
		FullName: "Test User",
		Username: "testuser",
		Password: "password123",
	}, nil)

	w := s.makeRequest(http.MethodPost, "/users/login", handlers.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}, nil)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	testCases := []struct {
		name       string
		payload    handlers.CreateLinkRequest
		token      string
		wantStatus int
	}{
		{
			name: "Valid Link",
			payload: handlers.CreateLinkRequest{
				Title: "Test Link",
				URL:   "https://example.com",
			},
			token:      token,
			wantStatus: http.StatusCreated,
		},
		{
			name: "Missing Title",
			payload: handlers.CreateLinkRequest{
				URL: "https://example.com",
			},
			token:      token,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid URL",
			payload: handlers.CreateLinkRequest{
				Title: "Test Link",
				URL:   "not-a-url",
			},
			token:      token,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "No Authentication",
			payload: handlers.CreateLinkRequest{
				Title: "Test Link",
				URL:   "https://example.com",
			},
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			headers := map[string]string{}
			if tc.token != "" {
				headers["Authorization"] = "Bearer " + tc.token
			}

			w := s.makeRequest(http.MethodPost, "/links", tc.payload, headers)
			assert.Equal(s.T(), tc.wantStatus, w.Code)
		})
	}
}

func (s *HandlerTestSuite) TestUpdateLinkHandler() {
	s.makeRequest(http.MethodPost, "/users/signup", handlers.SignUpRequest{
		FullName: "Test User",
		Username: "testuser",
		Password: "password123",
	}, nil)

	w := s.makeRequest(http.MethodPost, "/users/login", handlers.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}, nil)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	createResp := s.makeRequest(http.MethodPost, "/links", handlers.CreateLinkRequest{
		Title: "Original Link",
		URL:   "https://example.com",
	}, map[string]string{"Authorization": "Bearer " + token})

	assert.Equal(s.T(), http.StatusCreated, createResp.Code)

	var link models.Link
	s.db.First(&link)

	testCases := []struct {
		name       string
		linkID     string
		payload    handlers.CreateLinkRequest
		token      string
		wantStatus int
	}{
		{
			name:   "Valid Update",
			linkID: fmt.Sprintf("%d", link.ID),
			payload: handlers.CreateLinkRequest{
				Title: "Updated Link",
				URL:   "https://updated-example.com",
			},
			token:      token,
			wantStatus: http.StatusOK,
		},
		{
			name:   "Invalid Link ID",
			linkID: "999999",
			payload: handlers.CreateLinkRequest{
				Title: "Updated Link",
				URL:   "https://example.com",
			},
			token:      token,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "No Authentication",
			linkID: fmt.Sprintf("%d", link.ID),
			payload: handlers.CreateLinkRequest{
				Title: "Updated Link",
				URL:   "https://example.com",
			},
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "Invalid URL",
			linkID: fmt.Sprintf("%d", link.ID),
			payload: handlers.CreateLinkRequest{
				Title: "Updated Link",
				URL:   "not-a-url",
			},
			token:      token,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			headers := map[string]string{}
			if tc.token != "" {
				headers["Authorization"] = "Bearer " + tc.token
			}

			w := s.makeRequest(http.MethodPut, "/links/"+tc.linkID, tc.payload, headers)
			assert.Equal(s.T(), tc.wantStatus, w.Code)
		})
	}
}

func (s *HandlerTestSuite) TestDeleteLinkHandler() {
	s.makeRequest(http.MethodPost, "/users/signup", handlers.SignUpRequest{
		FullName: "Test User",
		Username: "testuser",
		Password: "password123",
	}, nil)

	w := s.makeRequest(http.MethodPost, "/users/login", handlers.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}, nil)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	createResp := s.makeRequest(http.MethodPost, "/links", handlers.CreateLinkRequest{
		Title: "Test Link",
		URL:   "https://example.com",
	}, map[string]string{"Authorization": "Bearer " + token})

	assert.Equal(s.T(), http.StatusCreated, createResp.Code)

	var link models.Link
	s.db.First(&link)

	testCases := []struct {
		name       string
		linkID     string
		token      string
		wantStatus int
	}{
		{
			name:       "Valid Delete",
			linkID:     fmt.Sprintf("%d", link.ID),
			token:      token,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Link ID",
			linkID:     "999999",
			token:      token,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "No Authentication",
			linkID:     fmt.Sprintf("%d", link.ID),
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			headers := map[string]string{}
			if tc.token != "" {
				headers["Authorization"] = "Bearer " + tc.token
			}

			w := s.makeRequest(http.MethodDelete, "/links/"+tc.linkID, nil, headers)
			assert.Equal(s.T(), tc.wantStatus, w.Code)
		})
	}
}

func (s *HandlerTestSuite) TestTrackLinkClickHandler() {
	s.makeRequest(http.MethodPost, "/users/signup", handlers.SignUpRequest{
		FullName: "Test User",
		Username: "testuser",
		Password: "password123",
	}, nil)

	w := s.makeRequest(http.MethodPost, "/users/login", handlers.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}, nil)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	s.makeRequest(http.MethodPost, "/links", handlers.CreateLinkRequest{
		Title: "Test Link",
		URL:   "https://example.com",
	}, map[string]string{"Authorization": "Bearer " + token})

	var link models.Link
	s.db.First(&link)

	testCases := []struct {
		name       string
		linkID     string
		token      string
		wantStatus int
	}{
		{
			name:       "Valid Click - Authenticated",
			linkID:     fmt.Sprintf("%d", link.ID),
			token:      token,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Valid Click - Unauthenticated",
			linkID:     fmt.Sprintf("%d", link.ID),
			token:      "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Link ID",
			linkID:     "invalid",
			token:      "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Non-existent Link",
			linkID:     "999999",
			token:      "",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			headers := map[string]string{}
			if tc.token != "" {
				headers["Authorization"] = "Bearer " + tc.token
			}

			w := s.makeRequest(http.MethodPost, "/analytics/"+tc.linkID+"/click", nil, headers)
			assert.Equal(s.T(), tc.wantStatus, w.Code)

			if tc.wantStatus == http.StatusOK {
				var analytics models.Analytics
				result := s.db.Where("link_id = ?", link.ID).First(&analytics)
				assert.NoError(s.T(), result.Error)
				assert.Greater(s.T(), analytics.ClickCount, uint(0))
			}
		})
	}
}
