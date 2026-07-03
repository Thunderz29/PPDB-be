package utils

type MessageLang struct {
	ID string
	EN string
}

const (
	MsgSuccessCreate = "SUCCESS_CREATE"
	MsgSuccessFetch  = "SUCCESS_FETCH"
	MsgSuccessUpdate = "SUCCESS_UPDATE"
	MsgSuccessDelete = "SUCCESS_DELETE"
	MsgSuccessLogin  = "SUCCESS_LOGIN"
	MsgSuccessLogout = "SUCCESS_LOGOUT"

	MsgErrNotFound   = "ERROR_NOT_FOUND"
	MsgErrInternal   = "ERROR_INTERNAL"
	MsgErrBadRequest = "ERROR_BAD_REQUEST"
	MsgErrInvalidID  = "ERROR_INVALID_ID"
	MsgErrAuth       = "ERROR_UNAUTHORIZED"
)

var Messages = map[string]MessageLang{
	MsgSuccessCreate: {
		ID: "Data berhasil dibuat",
		EN: "Data created successfully",
	},
	MsgSuccessFetch: {
		ID: "Data berhasil diambil",
		EN: "Data retrieved successfully",
	},
	MsgSuccessUpdate: {
		ID: "Data berhasil diperbarui",
		EN: "Data updated successfully",
	},
	MsgSuccessDelete: {
		ID: "Data berhasil dihapus",
		EN: "Data deleted successfully",
	},
	MsgSuccessLogin: {
		ID: "Login berhasil",
		EN: "Login successful",
	},
	MsgSuccessLogout: {
		ID: "Logout berhasil",
		EN: "Logout successful",
	},
	MsgErrNotFound: {
		ID: "Data tidak ditemukan",
		EN: "Data not found",
	},
	MsgErrInternal: {
		ID: "Terjadi kesalahan pada server",
		EN: "Internal server error",
	},
	MsgErrBadRequest: {
		ID: "Permintaan tidak valid",
		EN: "Invalid request",
	},
	MsgErrInvalidID: {
		ID: "Format ID tidak valid",
		EN: "Invalid ID format",
	},
	MsgErrAuth: {
		ID: "Tidak ada akses (Unauthorized)",
		EN: "Unauthorized",
	},
}


func GetMsg(key string, lang string) string {
	msg, exists := Messages[key]
	if !exists {
		return key
	}
	if lang == "id" {
		return msg.ID
	}
	return msg.EN
}
