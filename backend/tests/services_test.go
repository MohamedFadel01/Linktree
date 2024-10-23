package tests

import (
	"linktree-mohamedfadel-backend/internal/models"
	"linktree-mohamedfadel-backend/internal/services"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ServiceTestSuite struct {
	suite.Suite
	db               *gorm.DB
	userService      *services.UserService
	linkService      *services.LinkService
	analyticsService *services.AnalyticsService
}

func (s *ServiceTestSuite) SetupSuite() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal(err)
	}
	os.Setenv("JWT_SECRET", "test-secret-key")

	s.db.AutoMigrate(&models.User{}, &models.Link{}, &models.Analytics{})

	s.userService = services.NewUserService(s.db)
	s.linkService = services.NewLinkService(s.db)
	s.analyticsService = services.NewAnalyticsService(s.db)
}

func (s *ServiceTestSuite) TearDownSuite() {
	os.Unsetenv("JWT_SECRET")
}

func (s *ServiceTestSuite) SetupTest() {
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Analytics{})
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Link{})
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestUserSignUp() {
	testCases := []struct {
		name     string
		user     models.User
		password string
		wantErr  bool
	}{
		{
			name: "Valid User",
			user: models.User{
				FullName: "Test User",
				Username: "testuser",
				Bio:      "Test Bio",
			},
			password: "password123",
			wantErr:  false,
		},
		{
			name: "Missing Required Fields",
			user: models.User{
				Bio: "Test Bio",
			},
			password: "password123",
			wantErr:  true,
		},
		{
			name: "Duplicate Username",
			user: models.User{
				FullName: "Test User 2",
				Username: "testuser",
				Bio:      "Test Bio 2",
			},
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.userService.SignUp(tc.user, tc.password)
			if tc.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				var user models.User
				err = s.db.Where("username = ?", tc.user.Username).First(&user).Error
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.user.FullName, user.FullName)
			}
		})
	}
}

func (s *ServiceTestSuite) TestUserLogin() {
	user := models.User{
		FullName: "Test User",
		Username: "testuser",
	}
	s.userService.SignUp(user, "password123")

	testCases := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid Login",
			username: "testuser",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Invalid Username",
			username: "nonexistent",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "Invalid Password",
			username: "testuser",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "Missing Credentials",
			username: "",
			password: "",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			token, err := s.userService.Login(tc.username, tc.password)
			if tc.wantErr {
				assert.Error(s.T(), err)
				assert.Empty(s.T(), token)
			} else {
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), token)
			}
		})
	}
}

func (s *ServiceTestSuite) TestGetUserProfileInfo() {
	user := models.User{
		FullName: "Test User",
		Username: "testuser",
		Bio:      "Test Bio",
	}
	s.userService.SignUp(user, "password123")

	link := models.Link{
		Title: "Test Link",
		URL:   "https://example.com",
	}
	s.linkService.CreateLink("testuser", link)

	testCases := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "Existing User",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "Non-existent User",
			username: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			profile, err := s.userService.GetUserProfileInfo(tc.username)
			if tc.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.username, profile.Username)
				assert.Equal(s.T(), "Test Bio", profile.Bio)
				assert.Empty(s.T(), profile.PasswordHash)
				assert.NotEmpty(s.T(), profile.Links)
			}
		})
	}
}

func (s *ServiceTestSuite) TestUpdateUser() {
	user := models.User{
		FullName: "Test User",
		Username: "testuser",
		Bio:      "Test Bio",
	}
	s.userService.SignUp(user, "password123")

	testCases := []struct {
		name         string
		username     string
		updatedUser  models.User
		wantErr      bool
		verifyFields func(*testing.T, models.User)
	}{
		{
			name:     "Update Full Name",
			username: "testuser",
			updatedUser: models.User{
				FullName: "Updated Name",
			},
			wantErr: false,
			verifyFields: func(t *testing.T, user models.User) {
				assert.Equal(t, "Updated Name", user.FullName)
				assert.Equal(t, "Test Bio", user.Bio)
			},
		},
		{
			name:     "Update Bio",
			username: "testuser",
			updatedUser: models.User{
				Bio: "Updated Bio",
			},
			wantErr: false,
			verifyFields: func(t *testing.T, user models.User) {
				assert.Equal(t, "Updated Bio", user.Bio)
			},
		},
		{
			name:     "Update Password",
			username: "testuser",
			updatedUser: models.User{
				PasswordHash: "newpassword123",
			},
			wantErr: false,
			verifyFields: func(t *testing.T, user models.User) {
				assert.NotEmpty(t, user.PasswordHash)
				assert.NotEqual(t, "newpassword123", user.PasswordHash)
			},
		},
		{
			name:     "Non-existent User",
			username: "nonexistent",
			updatedUser: models.User{
				FullName: "Updated Name",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.userService.UpdateUser(tc.username, tc.updatedUser)
			if tc.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				var updatedUser models.User
				s.db.Where("username = ?", tc.username).First(&updatedUser)
				tc.verifyFields(s.T(), updatedUser)
			}
		})
	}
}

func (s *ServiceTestSuite) TestDeleteUser() {
	user := models.User{
		FullName: "Test User",
		Username: "testuser",
	}
	s.userService.SignUp(user, "password123")

	testCases := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "Existing User",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "Non-existent User",
			username: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.userService.DeleteUser(tc.username)
			if tc.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				var deletedUser models.User
				err = s.db.Where("username = ?", tc.username).First(&deletedUser).Error
				assert.Error(s.T(), err)
			}
		})
	}
}

func (s *ServiceTestSuite) TestCreateLink() {
	user := models.User{
		FullName: "Test User",
		Username: "testuser",
	}
	s.userService.SignUp(user, "password123")

	testCases := []struct {
		name    string
		link    models.Link
		wantErr bool
	}{
		{
			name: "Valid Link",
			link: models.Link{
				Title: "Test Link",
				URL:   "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "Missing Title",
			link: models.Link{
				URL: "https://example.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid URL",
			link: models.Link{
				Title: "Test Link",
				URL:   "not-a-url",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.linkService.CreateLink("testuser", tc.link)
			if tc.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				var link models.Link
				err = s.db.Where("url = ?", tc.link.URL).First(&link).Error
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.link.Title, link.Title)
			}
		})
	}
}

func (s *ServiceTestSuite) TestUpdateLink() {
	user := models.User{
		FullName: "Test User",
		Username: "testuser",
	}
	s.userService.SignUp(user, "password123")

	link := models.Link{
		Title: "Test Link",
		URL:   "https://example.com",
	}
	s.linkService.CreateLink("testuser", link)

	var createdLink models.Link
	s.db.Where("url = ?", link.URL).First(&createdLink)

	testCases := []struct {
		name        string
		username    string
		linkID      uint64
		updatedLink models.Link
		wantErr     bool
	}{
		{
			name:     "Update Title",
			username: "testuser",
			linkID:   uint64(createdLink.ID),
			updatedLink: models.Link{
				Title: "Updated Title",
			},
			wantErr: false,
		},
		{
			name:     "Update URL",
			username: "testuser",
			linkID:   uint64(createdLink.ID),
			updatedLink: models.Link{
				URL: "https://updated-example.com",
			},
			wantErr: false,
		},
		{
			name:     "Invalid URL",
			username: "testuser",
			linkID:   uint64(createdLink.ID),
			updatedLink: models.Link{
				URL: "not-a-url",
			},
			wantErr: true,
		},
		{
			name:     "Non-existent Link",
			username: "testuser",
			linkID:   9999,
			updatedLink: models.Link{
				Title: "Updated Title",
			},
			wantErr: true,
		},
		{
			name:     "Wrong User",
			username: "wronguser",
			linkID:   uint64(createdLink.ID),
			updatedLink: models.Link{
				Title: "Updated Title",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.linkService.UpdateLink(tc.username, tc.linkID, tc.updatedLink)
			if tc.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				var updatedLink models.Link
				s.db.First(&updatedLink, tc.linkID)
				if tc.updatedLink.Title != "" {
					assert.Equal(s.T(), tc.updatedLink.Title, updatedLink.Title)
				}
				if tc.updatedLink.URL != "" {
					assert.Equal(s.T(), tc.updatedLink.URL, updatedLink.URL)
				}
			}
		})
	}
}

func (s *ServiceTestSuite) TestDeleteLink() {
	user := models.User{
		FullName: "Test User",
		Username: "testuser",
	}
	s.userService.SignUp(user, "password123")

	link := models.Link{
		Title: "Test Link",
		URL:   "https://example.com",
	}
	s.linkService.CreateLink("testuser", link)

	var createdLink models.Link
	s.db.Where("url = ?", link.URL).First(&createdLink)

	testCases := []struct {
		name     string
		username string
		linkID   uint64
		wantErr  bool
	}{
		{
			name:     "Valid Delete",
			username: "testuser",
			linkID:   uint64(createdLink.ID),
			wantErr:  false,
		},
		{
			name:     "Non-existent Link",
			username: "testuser",
			linkID:   9999,
			wantErr:  true,
		},
		{
			name:     "Wrong User",
			username: "wronguser",
			linkID:   uint64(createdLink.ID),
			wantErr:  true,
		},
		{
			name:     "Already Deleted Link",
			username: "testuser",
			linkID:   uint64(createdLink.ID),
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.linkService.DeleteLink(tc.username, tc.linkID)
			if tc.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				var deletedLink models.Link
				err = s.db.First(&deletedLink, tc.linkID).Error
				assert.Error(s.T(), err)
				assert.Equal(s.T(), gorm.ErrRecordNotFound, err)
			}
		})
	}
}
