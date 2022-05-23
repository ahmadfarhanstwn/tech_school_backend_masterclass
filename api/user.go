package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/ahmadfarhanstwn/simple_bank/db/sqlc"
	"github.com/ahmadfarhanstwn/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	ChangedPasswordAt time.Time `json:"changed_password_at"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		ChangedPasswordAt: user.ChangedPasswordAt,
		CreatedAt: user.CreatedAt,
	}
}

func (s *Server) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		HashPassword: hashedPassword,
		Email: req.Email,
	}

	user, err := s.store.CreateUser(c, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	resp := newUserResponse(user)
	c.JSON(http.StatusOK, resp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User userResponse `json:"user"`
}

func (s *Server) loginUser(c *gin.Context) {
	var req loginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := s.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	resp := loginUserResponse{
		AccessToken: accessToken,
		User: newUserResponse(user),
	}

	c.JSON(http.StatusOK, resp)
}