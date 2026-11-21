package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringList is a list of strings in a jsonb column (e.g. step hints). nil is stored as [].
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
