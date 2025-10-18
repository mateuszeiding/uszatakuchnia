package models

import (
	enums "uszatakuchnia/enums/taste"
)

type Taste struct {
	Id   int
	Type enums.TasteType
}
