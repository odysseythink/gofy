package tokenmanager

import (
	"encoding/json"
	"fmt"
	"maps"
	"time"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/cache"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

// TokenManager handles token generation, revocation, and retrieval.
type TokenManager struct {
	config map[string]int
}

// NewTokenManager initializes a new TokenManager with a Redis client and configuration.
func NewTokenManager(config map[string]int) *TokenManager {
	return &TokenManager{
		config: config,
	}
}

// generateToken generates a new token for the given account or email.
func (tm *TokenManager) GenerateToken(token_type string, account *models.Account, email string, additionalData map[string]any) (string, error) {
	if account == nil && email == "" {
		return "", exceptions.NewValueError("account or email must be provided")
	}

	account_id := ""
	account_email := ""
	if account != nil {
		account_id = account.ID
		account_email = account.Email
	}

	if account_id != "" {
		oldToken := tm.getCurrentTokenForAccount(account_id, token_type)
		if oldToken != "" {
			tm.revokeToken(oldToken, token_type)
		}
	}

	token := uuid.NewV4().String()
	tokenData := map[string]any{
		"account_id": account_id,
		"email":      account_email,
		"token_type": token_type,
	}
	maps.Copy(tokenData, additionalData)

	expiryMinutes := confy.Get[int](fmt.Sprintf("%s_token_expiry_minutes", token_type))
	if expiryMinutes == 0 {
		return "", exceptions.NewValueError(fmt.Sprintf("expiry minutes for %s token is not set", token_type))
	}

	tokenKey := tm.getTokenKey(token, token_type)
	expiryTime := time.Duration(expiryMinutes) * time.Minute
	err := cache.Instance().SetEx(tokenKey, tokenData, expiryTime)
	if err != nil {
		mlog.Warningf("set redis(%s) failed:%v", tokenKey, err)
	}

	if account_id != "" {
		tm.setCurrentTokenForAccount(account_id, token, token_type, expiryMinutes)
	}

	return token, nil
}

// getTokenKey generates the Redis key for a token.
func (tm *TokenManager) getTokenKey(token, tokenType string) string {
	return fmt.Sprintf("%s:token:%s", tokenType, token)
}

// revokeToken revokes a token by deleting it from Redis.
func (tm *TokenManager) revokeToken(token, tokenType string) {
	tokenKey := tm.getTokenKey(token, tokenType)
	err := cache.Instance().DelKey(tokenKey)
	if err != nil {
		mlog.Warningf("Failed to revoke token: %v", err)
	}
}

// get_token_data retrieves token data from Redis.
func (tm *TokenManager) getTokenData(token, tokenType string) (map[string]any, error) {
	key := tm.getTokenKey(token, tokenType)
	var tokenDataJSON string
	err := cache.Instance().Get(key, tokenDataJSON)
	if err != nil {
		mlog.Warningf("%s token %s not found with key %s", tokenType, token, key)
		return nil, nil
	}

	var tokenData map[string]any
	err = json.Unmarshal([]byte(tokenDataJSON), &tokenData)
	if err != nil {
		return nil, err
	}

	return tokenData, nil
}

// getCurrentTokenForAccount retrieves the current token for an account from Redis.
func (tm *TokenManager) getCurrentTokenForAccount(account_id, token_type string) string {
	key := tm.getAccountTokenKey(account_id, token_type)
	var currentToken string
	err := cache.Instance().Get(key, currentToken)
	if err != nil {
		mlog.Warningf("Failed to get current token for account: %v", err)
	}
	return currentToken
}

// setCurrentTokenForAccount sets the current token for an account in Redis.
func (tm *TokenManager) setCurrentTokenForAccount(accountID, token string, tokenType string, expiryHours int) {
	key := tm.getAccountTokenKey(accountID, tokenType)
	expiryTime := time.Duration(expiryHours) * time.Hour
	err := cache.Instance().SetEx(key, token, expiryTime)
	if err != nil {
		mlog.Warningf("Failed to set current token for account: %v", err)
	}
}

// getAccountTokenKey generates the Redis key for an account's token.
func (tm *TokenManager) getAccountTokenKey(accountID, tokenType string) string {
	return fmt.Sprintf("%s:account:%s", tokenType, accountID)
}
