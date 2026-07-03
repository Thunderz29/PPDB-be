package utils

import "errors"

const (
	MsgSuccess             = "Success"
	MsgInternalServerError = "Terjadi kesalahan internal pada server"
	MsgBadRequest          = "Permintaan tidak valid"
	MsgNotFound            = "Data tidak ditemukan"
	MsgUnauthorized        = "Akses tidak sah"
	MsgForbidden           = "Akses ditolak"
	MsgValidationFailed    = "Validasi data gagal"
)

var (
	ErrNotFound      = errors.New("data tidak ditemukan")
	ErrUnauthorized  = errors.New("akses tidak sah")
	ErrForbidden     = errors.New("akses ditolak")
	ErrInvalidInput  = errors.New("input tidak valid")
	ErrInternal      = errors.New("kesalahan internal server")
)
