package validators

import (
	"regexp"
	"strings"
)

var (
	// Регулярное выражение для логина:
	// разрешены латинские буквы в верхнем и нижнем регистре, а так же
	// символы: -, _. Длина логина от 3 до 16 символов.
	REGEXP_USERNAME = regexp.MustCompile("^[a-zA-Z0-9_-]{3,16}$")

	// Регулярное выражение для почты
	REGEXP_EMAIL = regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9-]+.+.[a-zA-Z]{2,4}$")

	// Регулярное выражение для телефона
	REGEXP_TELEPHONE = regexp.MustCompile(`^(?:\+7)\d{10}$`)

	// Регулярное выражение для числа
	REGEXP_DIGITAL = regexp.MustCompile(`^\d{1,}$`)

	// Регулярное выражение для IPv4
	REGEXP_IP_V4 = regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.){3}(?:25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)$`)

	// Регулярное выражение для IPv6
	REGEXP_IP_V6 = regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|::([0-9a-fA-F]{1,4}:){0,5}[0-9a-fA-F]{1,4}|[0-9a-fA-F]{1,4}::([0-9a-fA-F]{1,4}:){0,4}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|:(:[0-9a-fA-F]{1,4}){1,7})$`)

	// Регулярное выражение для URL адреса
	REGEXP_URL = regexp.MustCompile(`^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$`)

	// Регулярное выражение для извелечения домена из URL адреса
	REGEXP_URL_DOMAIN = regexp.MustCompile(`(?i)https?:\/\/(?:[-\w]+\.)?([-\w]+)\.\w+(?:\.\w+)?`)

	// Регулярное выражение для извелечения протокола из URL адреса
	REGEXP_URL_PROTOCOL = regexp.MustCompile(`^([a-zA-Z]+):\/\/`)

	// вырезаю все кроме числа из строки номера телефона, что бы убрать форматирование
	regexp_telephone_cut_digitals = regexp.MustCompile(`\D`)
)

func ValidAndFormattingPhone(s string) (string, bool) {
	phone := regexp_telephone_cut_digitals.ReplaceAllString(s, "")
	if strings.HasPrefix(phone, "8") {
		phone = "+7" + phone[1:]
	}
	ok := REGEXP_TELEPHONE.Match([]byte(s))
	return phone, ok
}
