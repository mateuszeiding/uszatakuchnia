// Package migrations — wersjonowane migracje danych.
//
// AutoMigrate (db.AutoMigrate) tylko DODAJE tabele i kolumny. Wszystko, co przenosi, przekształca
// albo usuwa dane (zmiana nazwy kolumny, zmiana typu, przeniesienie danych do nowej tabeli,
// usunięcie starej kolumny) musi być migracją tutaj.
//
// Zasady:
//   - Jedna migracja = jeden plik NNNN_opis.go, dopisana na końcu listy All.
//   - ID raz wdrożonej migracji nigdy się nie zmienia, a jej kod nie jest już edytowany.
//   - Migracja dostaje transakcję po AutoMigrate (nowe kolumny już istnieją).
//   - Na świeżej bazie (resetdb) migracje też się wykonują — muszą działać na pustych tabelach
//     i bez starych kolumn (sprawdzaj tx.Migrator().HasColumn / HasTable).
package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	ID string
	Up func(tx *gorm.DB) error
}

type SchemaMigration struct {
	ID        string    `gorm:"type:varchar(128);primaryKey"`
	AppliedAt time.Time `gorm:"not null"`
}

func (SchemaMigration) TableName() string { return "schema_migrations" }

// All — migracje w kolejności wykonania. Nowe dopisuj na końcu.
var All = []Migration{
	recipeTimeRange,
}

// Run wykonuje migracje, których nie ma jeszcze w schema_migrations.
func Run(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&SchemaMigration{}); err != nil {
		return err
	}

	var done []string
	if err := tx.Model(&SchemaMigration{}).Pluck("id", &done).Error; err != nil {
		return err
	}
	applied := make(map[string]bool, len(done))
	for _, id := range done {
		applied[id] = true
	}

	for _, m := range All {
		if applied[m.ID] {
			continue
		}
		if err := m.Up(tx); err != nil {
			return fmt.Errorf("migration %s: %w", m.ID, err)
		}
		if err := tx.Create(&SchemaMigration{ID: m.ID, AppliedAt: time.Now()}).Error; err != nil {
			return err
		}
	}
	return nil
}
