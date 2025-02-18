package services

import (
	"context"
	"net/http"
	"time"

	"github.com/bhbdev/jam/config"
	"github.com/bhbdev/jam/lib/logger"
	"github.com/bhbdev/jam/lib/password"
	"github.com/bhbdev/jam/lib/redis"
	"github.com/bhbdev/jam/lib/session"
	"github.com/google/uuid"

	//"bhbdev/lib/util"
	"github.com/bhbdev/jam/models"
	"github.com/bhbdev/jam/services/repo"

	"gorm.io/gorm"
)

type UserService struct {
	*repo.Query
	DB *gorm.DB // this lets us use the db directly when repo.Query is not enough
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		Query: repo.Use(db),
		DB:    db,
	}
}

func (s *UserService) PageServiceName() string {
	return "UserService"
}

func (s *UserService) List(ctx context.Context, params *models.ListParams) ([]*models.User, error) {
	query := s.User.WithContext(ctx) //.Where(s.UCType.TenantId.Eq(ctxdata.TenantId(ctx)))
	apps, err := query.
		Order(s.User.CreatedAt.Desc(), s.User.ID.Desc()).
		Limit(params.Pagination.PerPage()).
		Offset(params.Pagination.Offset()).
		Find()

	if err != nil {
		logger.Error("UserService.List error", "error", err)
		return nil, err
	}

	total, err := query.Count()
	if err != nil {
		logger.Error("UserService.List error counting records", "error", err)
		return nil, err
	}
	params.Pagination.SetTotal(int(total))

	return apps, nil
}

func (s *UserService) Get(ctx context.Context, id int) (*models.User, error) {
	app, err := s.User.WithContext(ctx).GetByID(id)
	if err != nil {
		logger.Error("UserService.Get error", "error", err)
		return nil, err
	}
	return app, nil
}

func (s *UserService) Save(ctx context.Context, user *models.User) error {
	err := s.User.WithContext(ctx).Save(user)
	if err != nil {
		logger.Error("UserService.Save error saving record", "error", err)
		return err
	}
	return nil
}

func (s *UserService) SetPassword(ctx context.Context, id int64, pwd string) error {
	hashedPassword, err := password.Hash(pwd)
	if err != nil {
		logger.Error("UserService.SetPassword error hashing password", "error", err)
		return err
	}
	_, err = s.User.WithContext(ctx).Where(s.User.ID.Eq(id)).Update(s.User.Password, hashedPassword)
	if err != nil {
		logger.Error("UserService.SetPassword error updating password", "error", err)
		return err
	}
	return nil
}

func (s *UserService) UsernameExists(ctx context.Context, username string) bool {
	user, err := s.User.WithContext(ctx).Where(s.User.Username.Eq(username)).First()
	if err != nil {
		logger.Error("UserService.UsernameExists error", "error", err)
		return false
	}
	return user != nil
}

func (s *UserService) EmailExists(ctx context.Context, email string) bool {
	user, err := s.User.WithContext(ctx).Where(s.User.Email.Eq(email)).First()
	if err != nil {
		logger.Error("UserService.EmailExists error", "error", err)
		return false
	}
	return user != nil
}

func (s *UserService) Register(ctx context.Context, user *models.User) (errs map[string]string) {
	errs = make(map[string]string)

	// validate username and email are unique
	if s.UsernameExists(ctx, user.Username) {
		errs["Username"] = "Username already exists"
	}
	if s.EmailExists(ctx, user.Email) {
		errs["Email"] = "Email already exists"
	}
	if len(errs) > 0 {
		return
	}

	user.Status = models.UserStatusActive
	user.IsAdmin = false
	user.Password = password.TemporaryPassword(10)

	err := s.Save(ctx, user)
	if err != nil {
		logger.Error("UserService.Register error saving record", "error", err)
		errs["Database"] = "Error registering account"
	}

	logger.Info("UserService.Register", "user", user)

	return
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request, user *models.User) (sessionId string, errs map[string]string) {
	ctx := r.Context()
	errs = make(map[string]string)

	// find user with username and password match
	u, err := s.User.WithContext(ctx).
		Where(
			s.User.Email.Eq(user.Email),
			s.User.Password.Eq(user.Password),
		).
		First()
	if err != nil {
		logger.Error("UserService.Login error", "error", err)
		errs["Database"] = "Login failed"
		return
	}

	if u == nil {
		errs["Database"] = "Login failed"
		return
	}

	sh := session.NewSessionHandler(redis.Client)

	uuid, err := uuid.NewRandom()
	if err != nil {
		logger.Error("UserService.Login error generating session key", "error", err)
		errs["Session"] = "Login failed"
		return
	}
	sessionId = uuid.String()
	// create a session

	err = sh.Create(ctx, sessionId, session.UserSession{Id: int(u.ID)}, config.SessionExpireTime)
	if err != nil {
		logger.Error("UserService.Login error creating session", "error", err)
		errs["Session"] = "Login failed"
		return
	}

	// Ping the user's last login time.
	u.LastLoginAt = time.Now()

	err = s.Save(ctx, u)
	if err != nil {
		logger.Error("UserService.Login error saving record", "error", err)
		errs["Database"] = "Login failed"
	}

	return
}

func (s *UserService) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sh := session.NewSessionHandler(redis.Client)

	// get the session id from the cookie
	sessionId, err := r.Cookie(config.SessionCookieName)
	if err != nil {
		logger.Error("UserService.Logout error getting session id", "error", err)
		return
	}

	// expire the session
	err = sh.Expire(ctx, sessionId.Value, 0)
	if err != nil {
		logger.Error("UserService.Logout error expiring session", "error", err)
		return
	}
}
