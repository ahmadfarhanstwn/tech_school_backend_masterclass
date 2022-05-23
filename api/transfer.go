package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/ahmadfarhanstwn/simple_bank/db/sqlc"
	"github.com/ahmadfarhanstwn/simple_bank/token"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(c *gin.Context) {
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	fromAccountId, valid := s.validAccount(c, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	// get payload from header
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccountId.Owner {
		err := fmt.Errorf("the account doesn't belong to the authorized header")
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	_, valid = s.validAccount(c, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTransactionsParams{
		FromAccountId: req.FromAccountID,
		ToAccountId: req.ToAccountID,
		Amount: req.Amount,
	}

	account, err := s.store.TransferTransaction(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)
}

func (s *Server) validAccount(c *gin.Context, accountID int64, currency string) (db.Account,bool) {
	account, err := s.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return account,false
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency missmatch: %s vs %s", account.ID, account.Currency, currency)
		c.JSON(http.StatusBadRequest, errResponse(err))
		return account, false
	}
	return account,true
}