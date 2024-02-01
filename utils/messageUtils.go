package utils

// Error Message
const (
	ErrInternalServer = "Terjadi kesalahan server. Silakan coba kembali"
    ErrUserAlreadyExist  = "Username telah digunakan oleh user yang telah mendaftar sebelumnya"
	ErrConfirmPassword = "Konfirmasi kata sandi tidak sama dengan kata sandi"
	ErrValidatePassword = "Kata sandi tidak boleh kurang dari 6 karakter"
	ErrHashPassword = "Error Hashing Password!"
	ErrUserNotFound = "Username tidak terdaftar"
	ErrInvalidPassword = "Kata sandi tidak sesuai"
	ErrCreateToken = "Gagal membuat token"
	ErrUnauthorizedUser = "User belum terautentikasi!"
	ErrTimeValidation = "Invalid time cook format"
	ErrInvalidTimeCook = "Invalid time cook range"
	ErrBadRequest = "Bad request"
	ErrDataNotFound = "Data tidak ditemukan!"
	ErrRecipeNotFound = "Resep masakan tidak tersedia"
	ErrReadingFile = "Error reading/parsing file from form-data"
	ErrReadingJsonFile = "Error reading/parsing JSON file from form-data"
	ErrUploadImageMinio = "Gagal mengupload gambar ke MinIO!"
	ErrOpenFileMinio = "Gagal membuka file: %v"
	ErrInitMinio = "Gagal menginisiasi MinIO Client: %v"
	ErrDetailRecipeNotFound = "Detil Resep masakan tidak tersedia"
	ErrGetImageUrl = "Gagal mengambil image URL"
	ErrDataAlreadyDeleted = "Data sudah terhapus"
	ErrParsingToken = "Error parsing token:"
	ErrInvalidToken = "Token tidak valid"
	ErrClaimToken = "Gagal mendapatkan klaim token"
	ErrGetData = "Gagal mengambil data: "
    ErrorMessage    = "An error occurred: %s"
)

// Success Message
const (
	SuccSignUp = "User %s registered successfully!"
	SuccSignIn = "Auth User Success"
	SuccGetFavRecipe = "Berhasil memuat Resep Masakan Favorit!"
	SuccGetMyRecipe = "Berhasil memuat Resep Masakan Saya"
	SuccGetAllRecipe = "Berhasil memuat Resep Masakan"
	SuccCreateRecipe = "Resep %s berhasil ditambahkan!"
	SuccUpdateRecipe = "Resep %s berhasil diubah!"
	SuccAddFavorite = "Resep %s berhasil ditambahkan ke dalam favorit"
	SuccRemoveFavorite = "Resep %s berhasil dihapus dari favorit"
	SuccDeleteRecipe = "Resep %s berhasil dihapus!"
	Success = "Pesan Sukses"
    WelcomeMessagee  = "Server dimulai di port 8080"
    ErrorMessagee   = "An error occurred: %s"
)