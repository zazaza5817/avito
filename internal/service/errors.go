// Package service содержит бизнес-логику и реализации сервисного слоя.
package service

import (
	"errors"
)

// Общие ошибки сервисов
var (
	ErrTeamExists       = errors.New("team already exists")
	ErrTeamNotFound     = errors.New("team not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrPRExists         = errors.New("pull request already exists")
	ErrPRNotFound       = errors.New("pull request not found")
	ErrPRMerged         = errors.New("cannot modify merged pull request")
	ErrReviewerNotFound = errors.New("reviewer not assigned to this PR")
	ErrNoCandidate      = errors.New("no active replacement candidate in team")
)
