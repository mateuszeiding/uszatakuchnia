package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"

	"uszatakuchnia/db/migrations"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Odcisk schematu: encje (nazwy pól, typy, tagi gorm) + ID wersjonowanych migracji.
// Zimny start funkcji porównuje go z bazą i pomija AutoMigrate + migracje, gdy nic się nie zmieniło —
// zamiast setek zapytań o schemat pod advisory lockiem jest jedno SELECT.
// Zmiana encji albo nowa migracja = inny odcisk = migracje odpalają się jak dotąd.

const fingerprintKey = "fingerprint"

type SchemaState struct {
	Key   string `gorm:"type:varchar(64);primaryKey"`
	Value string `gorm:"type:varchar(128);not null"`
}

func (SchemaState) TableName() string { return "schema_state" }

func SchemaFingerprint() string {
	h := sha256.New()
	var walk func(t reflect.Type)
	walk = func(t reflect.Type) {
		for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return
		}
		fmt.Fprintf(h, "%s.%s{", t.PkgPath(), t.Name())
		for i := range t.NumField() {
			f := t.Field(i)
			if f.Anonymous {
				walk(f.Type)
				continue
			}
			fmt.Fprintf(h, "%s %s `%s`;", f.Name, f.Type, f.Tag)
		}
		h.Write([]byte("}"))
	}
	for _, m := range models {
		walk(reflect.TypeOf(m))
	}
	for _, m := range migrations.All {
		fmt.Fprintf(h, "m:%s;", m.ID)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// schemaCurrent — czy baza ma już schemat w wersji z tego kodu. Brak tabeli / błąd = nie.
func schemaCurrent(db *gorm.DB, fp string) bool {
	if !db.Migrator().HasTable(&SchemaState{}) {
		return false
	}
	var v string
	err := db.Raw(`SELECT value FROM schema_state WHERE key = ?`, fingerprintKey).Scan(&v).Error
	return err == nil && v == fp
}

func saveFingerprint(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&SchemaState{}); err != nil {
		return err
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&SchemaState{Key: fingerprintKey, Value: SchemaFingerprint()}).Error
}

// ensureSchema — migracje przy zimnym starcie tylko wtedy, gdy odcisk się zmienił.
// Po wzięciu locka sprawdza jeszcze raz: równoległy zimny start mógł już zmigrować.
func ensureSchema(ctx context.Context, db *gorm.DB) error {
	fp := SchemaFingerprint()
	if schemaCurrent(db, fp) {
		return nil
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockTx(tx, migrateLockKey); err != nil {
			return err
		}
		if schemaCurrent(tx, fp) {
			return nil
		}
		return runMigrations(tx)
	})
}
