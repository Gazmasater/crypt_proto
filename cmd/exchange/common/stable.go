package common

import "strings"

// Список стабильных монет (стейблкоинов)
var stableCoins = map[string]bool{
	"USDT": true,
	"USDC": true,
	"USD1": true,
	"USDD": true,
	"USDG": true,
}

// IsStable проверяет, является ли монета стейблкоином
func IsStable(s string) bool {
	return stableCoins[strings.ToUpper(s)]
}

// AddStableCoin позволяет динамически добавить стейблкоин в список
func AddStableCoin(s string) {
	stableCoins[strings.ToUpper(s)] = true
}

// RemoveStableCoin позволяет удалить стейблкоин из списка
func RemoveStableCoin(s string) {
	delete(stableCoins, strings.ToUpper(s))
}
