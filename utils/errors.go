package utils

import "errors"

// Predefined general error messages
const (
	MsgSuccess             = "Success"
	MsgInternalServerError = "Terjadi kesalahan internal pada server"
	MsgBadRequest          = "Permintaan tidak valid"
	MsgNotFound            = "Data tidak ditemukan"
	MsgUnauthorized        = "Akses tidak sah"
	MsgForbidden           = "Akses ditolak"
	MsgValidationFailed    = "Validasi data gagal"
)

// Common error variables for internal logic checks
var (
	ErrNotFound      = errors.New("data tidak ditemukan")
	ErrUnauthorized  = errors.New("akses tidak sah")
	ErrForbidden     = errors.New("akses ditolak")
	ErrInvalidInput  = errors.New("input tidak valid")
	ErrInternal      = errors.New("kesalahan internal server")
)
