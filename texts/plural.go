package texts

type Countable interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
}

// Plural возвращает корректную форму множественного числа: Plural(11, "рыба", "рыбы", "рыб") -> "рыб"
func Plural[T Countable](count T, formSingular, formPluralWeak, formPluralStrong string) string {
	if count%10 == 1 && count%100 != 11 {
		return formSingular
	}
	if count%10 >= 2 && count%10 <= 4 && (count%100 < 10 || count%100 >= 20) {
		return formPluralWeak
	}
	return formPluralStrong
}
