package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringList — lista tekstów w kolumnie jsonb (np. wskazówki kroku). nil zapisuje się jako [].
type StringList []string

func (l StringList) Value() (driver.Value, error) {
	if l == nil {
		return "[]", nil
	}
	b, err := json.Marshal([]string(l))
	return string(b), err
}

func (l *StringList) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*l = StringList{}
		return nil
	case []byte:
		return json.Unmarshal(v, l)
	case string:
		return json.Unmarshal([]byte(v), l)
	default:
		return fmt.Errorf("StringList: nieobsługiwany typ %T", src)
	}
}
