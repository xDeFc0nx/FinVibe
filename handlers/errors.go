package handlers

import (
	"fmt"
	"github.com/pkg/errors"
)

const (
	MsgInvalidData        = "Invalid form data" // You already have this as a var, using const for consistency
	MsgInvalidInput       = "Invalid input provided"
	MsgInvalidEmail       = "Invalid email address"
	MsgInvalidCredentials = "Invalid email or password"
	MsgPasswordTooShort   = "Password requires at least 8 characters"
	MsgMissingID          = "ID is required"
	MsgMissingAccountID   = "Account ID is required"
	MsgMissingUserID      = "User ID is required"
	MsgMissingFieldFmt    = "%s is required" // For dynamic fields like "First Name is required"
	MsgMissingToken       = "Token is missing"
	MsgInvalidToken       = "Token is not valid or has expired"
	MsgInvalidTokenFormat = "Invalid token format"
	MsgMissingAction      = "WebSocket action is required"
	MsgInvalidAction      = "Unknown WebSocket action"
	MsgInvalidFrequency   = "Invalid recurring frequency"
	MsgMissingFrequency   = "Recurring frequency is required"
	MsgDateRangeRequired  = "Date range is required"
)

var stackTrace = fmt.Sprintf("%+v", errors.Wrap(err, ""))

// --- Not Found Errors ---
const (
	MsgNotFound             = "Resource not found"
	MsgUserNotFound         = "User not found"
	MsgAccountNotFound      = "Account not found"
	MsgTransactionNotFound  = "Transaction not found"
	MsgBudgetNotFound       = "Budget not found"
	MsgGoalNotFound         = "Goal not found"
	MsgConnectionIDNotFound = "Connection ID not found"
)

// --- Operation Failure Errors ---
const (
	MsgDBConnectionFailed      = "Failed to connect to database"
	MsgDBPingFailed            = "Failed to ping database"
	MsgDBOperationFailed       = "Database operation failed"
	MsgCreateFailedFmt         = "Failed to create %s" // e.g., "Failed to create account"
	MsgUpdateFailedFmt         = "Failed to update %s" // e.g., "Failed to update user"
	MsgDeleteFailedFmt         = "Failed to delete %s" // e.g., "Failed to delete budget"
	MsgFetchFailedFmt          = "Failed to fetch %s"  // e.g., "Failed to fetch transactions"
	MsgCollectRowsFailed       = "Failed to collect database rows"
	MsgGenerateResponseFailed  = "Failed to generate response"
	MsgPasswordHashFailed      = "Failed to hash password"
	MsgTokenGenerationFailed   = "Failed to create token"
	MsgTokenDecodeFailed       = "Failed to decode token"
	MsgWebSocketSendFailed     = "Failed to send WebSocket message"
	MsgWebSocketCreationFailed = "Failed to create WebSocket connection"
	MsgWebSocketUpdateFailed   = "Failed to update WebSocket state"
)

// --- Authentication / Authorization Errors ---
const (
	MsgUnauthorized        = "Unauthorized"
	MsgPermissionDeniedFmt = "%s does not belong to you" // e.g., "Transaction does not belong to you"
	MsgEmailExists         = "Email already exists"
)

// --- Configuration Errors ---
const (
	MsgEnvVarNotSetFmt  = "%s environment variable not set" // e.g., "DB_CONFIG environment variable not set"
	MsgLoadConfigFailed = "Error loading configuration"
	MsgOpenFileFailed   = "Error opening file"
)
