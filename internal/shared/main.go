package shared

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns NullString) IsValid(maxLength int) bool {
	if !ns.Valid {
		return true
	}
	return len(ns.String) <= maxLength
}
