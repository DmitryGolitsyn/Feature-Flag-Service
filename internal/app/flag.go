package app

import (
	"errors"
	"strings"
)

// Базовые ошибки домена (типизированные — пригодятся для MapErrorToStatus)
var (
	ErrFlagInvalid = errors.New("flag_invalid")
)

type Percentage uint8 // 0..100

type FlagKey string

func (k FlagKey) Valid() bool { return k != "" && !strings.ContainsAny(string(k), " \t\n\r") }

type TenantID string

func (t TenantID) Valid() bool { return t != "" }

type Rule struct {
	// простейшая версия: процент rollout
	Rollout Percentage // 0..100
}

type Flag struct {
	Tenant TenantID
	Key    FlagKey
	Rule   Rule
}

// Минимальная валидация домена
func (f *Flag) Validate() error {
	if !f.Tenant.Valid() || !f.Key.Valid() {
		return ErrInvalid // уже есть у вас
	}
	if f.Rule.Rollout > 100 {
		return ErrFlagInvalid
	}
	return nil
}
