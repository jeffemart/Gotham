package validator

import "regexp"

// EmailValidator valida o formato do email
func EmailValidator(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// PasswordValidator valida a força da senha
func PasswordValidator(password string) bool {
	// Mínimo 8 caracteres, pelo menos uma letra maiúscula, uma minúscula e um número
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(password)
}
