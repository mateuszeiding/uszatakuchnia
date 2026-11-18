# Kopia zapasowa bazy — gdzie jest i jak ją odtworzyć

Ściąga na dzień, w którym coś zniknęło. Czytaj po kolei, nie skacz.

## W skrócie

| Co | Gdzie |
|---|---|
| Kopie | GitHub → repo `uszatakuchnia` → **Actions** → workflow **Kopia bazy** → wybrany przebieg → sekcja **Artifacts** na dole. Jedna kopia na dobę (ok. 03:17–04:17 czasu polskiego), trzymane **90 dni**. |
| Format | Plik `uszatakuchnia-RRRR-MM-DD_GGMMZ.dump.gpg` (czas UTC) w zipie. Środek: `pg_dump` w formacie custom, zaszyfrowany gpg (AES256) hasłem. |
| Hasło | W menedżerze haseł: **„Uszata Kuchnia — hasło kopii bazy”**. Kopia w GitHub Secrets (`BACKUP_PASSPHRASE`) jest tylko do zapisu — z GitHuba **nie da się** jej odczytać. Bez hasła z menedżera kopie są bezużyteczne. |
| Workflow | `.github/workflows/db-backup.yml`, skrypt `scripts/db-backup.sh` |
| Przepis po przepisie | `go run ./cmd/restore-recipe` (niżej, krok 5B) |

Repo jest publiczne: plik kopii może pobrać każdy zalogowany na GitHubie, dlatego jest zaszyfrowany. Log workflow pokazuje tylko wersje, rozmiar i sumę kontrolną.

## 0. Czy wystarczy Neon?

Jeśli szkoda stała się **mniej niż 6 godzin temu**, najszybciej cofnąć bazę w Neonie: Neon Console → projekt → **Backup & Restore** → przywrócenie gałęzi `main` do wybranego momentu. Darmowy plan sięga tylko 6 godzin wstecz. Starsze szkody → kopia z GitHuba, czytaj dalej.

## 1. Narzędzia (raz, na zapas)

1. **gpg** — jest w Git for Windows: `C:\Program Files\Git\usr\bin\gpg.exe`.
2. **PostgreSQL 17** (bez instalacji): pobierz ZIP „binaries” dla Windows x64 z <https://www.enterprisedb.com/download-postgresql-binaries> (wersja 17.x — taka jak na Neonie), rozpakuj tak, żeby istniał `C:\pgsql17\bin\pg_restore.exe` (w ZIP-ie jest katalog `pgsql` — zmień jego nazwę na `pgsql17`).

Wersja główna musi się zgadzać z Neonem (dziś 17). Jeśli Neon przejdzie na 18, potrzebujesz narzędzi 18.

Wszystkie polecenia niżej są dla PowerShell 7. Najpierw ustaw ścieżki (zmień katalog na swój):

```powershell
$gpg = "C:\Program Files\Git\usr\bin\gpg.exe"; $PG = "C:\pgsql17\bin"; $K = "$HOME\Desktop\kopia"; New-Item -ItemType Directory -Force $K | Out-Null; Set-Location $K
```

## 2. Pobierz kopię

1. GitHub → **Actions** → po lewej **Kopia bazy**.
2. Kliknij przebieg z dnia, **zanim** stała się szkoda (zielony ✓).
3. Na dole, w **Artifacts**, pobierz `uszatakuchnia-….` (zip).
4. Rozpakuj zip do `$K`. W środku jest jeden plik `….dump.gpg`.

## 3. Odszyfruj

```powershell
$p = Read-Host "Hasło kopii" -MaskInput
```

```powershell
& $gpg --batch --yes --pinentry-mode loopback --passphrase $p --output kopia.dump --decrypt (Get-ChildItem *.dump.gpg)[0]
```

`decryption failed: Bad session key` = złe hasło. Sprawdzenie, że plik jest dobry:

```powershell
& "$PG\pg_restore.exe" --list kopia.dump | Select-String "Dumped"
```

## 4. Lokalna baza na kopię

Potrzebna do kroku 5B i do próby generalnej. Tworzy tymczasowy serwer na porcie 5434, tylko na tym komputerze:

```powershell
& "$PG\initdb.exe" -D "$K\pgdata" -U postgres -A trust -E UTF8 --locale=C
```

```powershell
& "$PG\pg_ctl.exe" -D "$K\pgdata" -o "-p 5434" -l "$K\pg.log" -w start
```

```powershell
& "$PG\createdb.exe" -h localhost -p 5434 -U postgres kopia
```

```powershell
& "$PG\pg_restore.exe" -h localhost -p 5434 -U postgres --no-owner --no-privileges --exit-on-error -d kopia kopia.dump
```

Szybki rzut oka:

```powershell
& "$PG\psql.exe" -h localhost -p 5434 -U postgres -d kopia -c "select id, name, status from recipes order by id"
```

Po wszystkim serwer zatrzymujesz: `& "$PG\pg_ctl.exe" -D "$K\pgdata" stop` (katalog `$K` możesz potem usunąć).

## 5A. Cała baza z powrotem na produkcję

Nadpisuje wszystkie tabele, które są w kopii, stanem z dnia kopii. Zmiany zrobione po tym dniu przepadają.

1. **Zabezpiecz obecny stan:** Neon Console → **Branches** → **Create branch** z `main` (nazwa np. `przed-odtworzeniem-RRRR-MM-DD`). Gdyby coś poszło źle, dane z chwili przed odtworzeniem są na tej gałęzi.
2. Weź connection string **bez poolera**: Neon Console → **Connect** → wyłącz „Connection pooling” (host bez `-pooler`). To samo co `DB_DATABASE_URL_UNPOOLED` w `.env.local`.
3. Odtwórz (wszystko albo nic — jedna transakcja):

```powershell
$url = Read-Host "Connection string BEZ poolera" -MaskInput
```

```powershell
& "$PG\pg_restore.exe" --clean --if-exists --no-owner --no-privileges --single-transaction -d $url kopia.dump
```

4. Otwórz aplikację i sprawdź listę przepisów i jeden przepis ze zdjęciem. Jeśli kod jest nowszy niż kopia, aplikacja sama dogra brakujące migracje przy pierwszym połączeniu (`schema_state`, `schema_migrations` też są w kopii).

Tabele, których nie było w kopii (dodane później), zostają nietknięte.

## 5B. Jeden przepis z kopii

Gdy zniknął albo został zepsuty jeden przepis, a reszta bazy ma zostać. Wymaga kroków 2–4 (kopia odtworzona lokalnie jako baza `kopia`).

1. Znajdź id przepisu w kopii (zapytanie z kroku 4 albo adres `/recipes/<id>` w aplikacji).
2. Podgląd — nic nie zapisuje, pokazuje, co by się stało:

```powershell
$env:DB_DATABASE_URL = Read-Host "Connection string produkcji" -MaskInput
```

```powershell
go run ./cmd/restore-recipe -from "postgres://postgres@localhost:5434/kopia?sslmode=disable" -id 12
```

3. Jeśli podgląd się zgadza, ten sam rozkaz z `-apply`:

```powershell
go run ./cmd/restore-recipe -from "postgres://postgres@localhost:5434/kopia?sslmode=disable" -id 12 -apply
```

(Polecenia z katalogu repo. Zamiast `$env:DB_DATABASE_URL` możesz podać `-to "<connection string>"`.)

Co robi: przepis o tym id wstawia albo nadpisuje, jego kroki, składniki, zamienniki, tagi, sprzęt i link do zdjęcia zastępuje wersją z kopii, brakujące zdjęcie dokłada. Inne przepisy zostają. Wszystko w jednej transakcji. Jeśli w produkcji brakuje składnika, jednostki, tagu albo sprzętu, którego przepis używał (też usunięty po dniu kopii), narzędzie kończy się błędem klucza obcego i niczego nie zapisuje. Przepis odwołuje się do nich po id, więc ponowne dodanie w aplikacji nie pomoże (dostaną nowe id) — wtedy poproś Claude’a o pomoc, wklejając komunikat błędu.

Uwaga: jeśli przepis był **usunięty**, inne przepisy, które go linkowały jako składnik, straciły link (`linked_recipe_id` = NULL) — to trzeba poprawić ręcznie w aplikacji.

## Kopia ręczna (np. przed ryzykowną migracją)

GitHub → **Actions** → **Kopia bazy** → **Run workflow** → gałąź `main` → **Run workflow**. Po ok. minucie gotowe; kopia jest w Artifacts tego przebiegu. Alternatywa bez czekania: snapshot albo gałąź w Neon Console.

## Gdy backup nie działa

- **Przyszedł mail „Run failed: Kopia bazy”** → otwórz przebieg, zobacz, który krok jest czerwony.
  - `Zrzut, szyfrowanie, kontrola` z `server version mismatch` → Neon dostał nowszego Postgresa. Zmień `PG_MAJOR` w `.github/workflows/db-backup.yml` na nową wersję główną.
  - Błąd połączenia / hasła → sprawdź sekret `BACKUP_DATABASE_URL` (np. po zmianie hasła roli w Neonie trzeba go podmienić).
  - Czerwony tylko `Keepalive crona` → kopia **jest** zapisana, nie udało się tylko przedłużenie crona; zwykle przejściowe.
- **Cron wyłączony** — GitHub wyłącza zaplanowane workflow w publicznym repo po 60 dniach bez aktywności. Krok `Keepalive crona` przy każdym przebiegu włącza workflow z powrotem przez API (bez pustych commitów), więc nie powinno się to zdarzyć. Gdyby jednak: **Actions** → **Kopia bazy** → przycisk **Enable workflow**. Wyłączony cron nie wysyła maila (nic się nie uruchamia) — dlatego raz w miesiącu zerknij, czy na liście są świeże przebiegi.
- **Baza urosła** (zdjęcia w tabeli `photos`): każda kopia to cała baza ściągnięta z Neona. Darmowy plan Neona ma miesięczny limit transferu — przy bazie rzędu 100–150 MB codzienna kopia zaczyna go zjadać. Wtedy rzadsza kopia albo zdjęcia poza bazą (osobny temat w `docs/produkt.md`).

## Konfiguracja (raz, przez właściciela)

Sekrety w GitHubie: repo → **Settings** → **Secrets and variables** → **Actions** → **New repository secret**:

| Sekret | Wartość |
|---|---|
| `BACKUP_DATABASE_URL` | connection string do Neona **bez poolera** (`DB_DATABASE_URL_UNPOOLED` z `.env.local`) |
| `BACKUP_PASSPHRASE` | losowe hasło, min. 20 znaków — to samo, co w menedżerze haseł |

Hasło można wygenerować tak (wynik od razu zapisz w menedżerze haseł):

```powershell
[Convert]::ToBase64String([Security.Cryptography.RandomNumberGenerator]::GetBytes(32))
```

Zmiana hasła później: nowe kopie będą szyfrowane nowym, stare dalej starym — trzymaj w menedżerze oba, dopóki stare kopie nie wygasną (90 dni).
