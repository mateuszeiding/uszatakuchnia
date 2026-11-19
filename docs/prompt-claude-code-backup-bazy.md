# Prompt dla Claude Code — automatyczna kopia zapasowa bazy

## Kontekst

Produkcyjna baza to Neon na **darmowym planie**. Od 14.11.2026 są w niej prawdziwe przepisy, składniki z makro i zdjęcia (tabela `photos`, bytea). Darmowy plan pozwala przywrócić bazę tylko do stanu sprzed maksymalnie **6 godzin** (okno historii, limit 1 GB zmian) i daje **1 ręczny snapshot**. Nie ma żadnej kopii, która przeżyłaby dzień.

`content/` to nie backup: trzyma przepisy w stanie z dnia dodania, a poprawki zrobione w aplikacji tam nie wracają.

Stan produktu: `docs/produkt.md`, „Ryzyka i dług”.

## Cel

Codzienna, automatyczna kopia całej bazy, z której da się odtworzyć produkcję albo pojedynczy przepis z wybranego dnia z ostatnich kilku tygodni. Bez dodatkowych kosztów.

## Rekomendowane podejście (analityk)

GitHub Actions, cron raz na dobę: `pg_dump` całej bazy (format custom) → szyfrowanie → zapis jako artefakt workflow z retencją 90 dni (maksimum GitHuba). Connection string w GitHub Secrets.

Jeśli widzisz lepsze darmowe miejsce na kopie (np. darmowy bucket R2 lub B2, osobne prywatne repo), zaproponuj je z uzasadnieniem przed wdrożeniem.

## Wymagania

- **Repo `uszatakuchnia` jest publiczne.**
  - Artefakty workflow w publicznym repo może pobrać każdy zalogowany użytkownik GitHuba. Kopia musi więc być zaszyfrowana, a klucz lub hasło trzymane w Secrets i, osobno, u właściciela.
  - Logi workflow są publiczne. Connection string, hasło i zawartość bazy nie mogą trafić do logu.
- **Pełna kopia:** schemat i dane, łącznie ze zdjęciami i tabelą `schema_migrations`.
- **Wersje:** wersja `pg_dump` zgodna z wersją Postgresa na Neonie, inaczej dump się nie uda albo nie da się go odtworzyć.
- **Nieudany backup ma być widoczny.** Właściciel dostaje powiadomienie (domyślny mail GitHuba o nieudanym workflow wystarczy, jeśli do niego trafia — sprawdź).
- **Wyłączanie crona:** GitHub wyłącza zaplanowane workflow w publicznym repo po 60 dniach bez aktywności w repo. Rozwiąż to albo opisz właścicielowi wprost, co ma robić.
- **Ręczne uruchomienie** (`workflow_dispatch`), np. przed ryzykowną migracją.
- **Odtworzenie przetestowane:** odtwórz kopię do lokalnej bazy i sprawdź, że aplikacja na niej działa (lista przepisów, przepis ze zdjęciem). Kopia, której nikt nie odtworzył, nie liczy się jako backup.
- Nie zapychaj darmowych limitów Neona: jedno połączenie dziennie, bez długich transakcji.

## Do oddania

- Workflow w `.github/workflows/`.
- Instrukcja odtwarzania w `docs/` (pełna baza; jeden przepis z kopii; gdzie jest klucz) — krok po kroku, dla właściciela, który robi to raz na rok w stresie.
- Wpis w `CLAUDE.md` w sekcji o bazie: gdzie są kopie, jak uruchomić backup ręcznie, że przed ryzykowną migracją warto go uruchomić albo zrobić snapshot w Neonie.
- Lista sekretów do dodania w GitHubie i kroki ręczne dla właściciela.

## Poza zakresem

- Zmiana planu Neona.
- Przeniesienie zdjęć z bazy do Vercel Blob (osobny temat w „Ryzykach”).

## Testowanie

`.env.local` wskazuje na produkcję. Odtwarzanie testuj na lokalnej bazie, nigdy na produkcyjnej. `resetdb` jest zakazany.

### Checklista odbioru

- [ ] Workflow uruchomiony ręcznie: zielony, artefakt zaszyfrowany (nie da się go otworzyć bez klucza).
- [ ] Log workflow nie zawiera connection stringa ani danych.
- [ ] Cron ustawiony raz na dobę, retencja 90 dni (albo inna z uzasadnieniem).
- [ ] Kopia odtworzona lokalnie, aplikacja na niej działa (lista, przepis ze zdjęciem).
- [ ] Problem 60 dni bez aktywności rozwiązany albo opisany.
- [ ] Powiadomienie o nieudanym backupie dociera do właściciela.
- [ ] Instrukcja odtwarzania w `docs/`, wpis w `CLAUDE.md`.
