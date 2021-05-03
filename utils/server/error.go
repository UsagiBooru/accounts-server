package server

import "errors"

// ErrInviteNotFound is shared error for handling createAccount method
var ErrInviteNotFound = errors.New("invite code was not found")
