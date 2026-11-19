# Prompt dla Claude Code — zdjęcia przepisów

## Kontekst

Zdjęcie przepisu było w planie od samego początku. **Właściciel nie jest w stanie dodać zdjęcia do przepisu.** Żaden opublikowany przepis go nie ma, a strona przepisu na telefonie zaczyna się od zastępczego fioletowego prostokąta na prawie cały ekran (`docs/przeglad-uzycia-2026-11-18.md`, punkt 5).

W kodzie jest ścieżka uploadu z 14.10: `PhotoField` + `CropModal` (kadrowanie 1:1) → `POST /api/photos` (`api/photos/upload.go`, bajty do tabeli `photos`, bytea) → adres `/api/photos/:id` w `recipe_photos.image_url`. Nikt jej jednak nie sprawdził na produkcji. Od 18.11 zapisy wymagają autoryzacji (c65e2cb). Jest też stara encja `images` z polami autora z Unsplash (`db/entities/image.go`) — ustal, czy cokolwiek jej jeszcze używa.

**Korekta (18.11, wieczór):** właściciel ma **własne zdjęcia (pliki)** i chce je wgrywać. Problem, który widzi: aplikacja trzyma zdjęcie jako adres lub id z obcej strony (Unsplash, `image_url`), a jego pliku nie ma gdzie zapisać. Punkt 3 (zdjęcie z linku) odpada. Magazyn: Vercel Blob (decyzja właściciela, punkt 4).

Stan produktu: `docs/produkt.md`.

## Cel

Właściciel ma zdjęcie przy każdym przepisie, który chce pokazać znajomym, i dodaje je bez kombinowania.

## Zakres

### 1. Diagnoza (najpierw, z raportem)
Przejdź na produkcji ścieżkę właściciela: formularz przepisu → dodaj zdjęcie z pliku → zapisz → zdjęcie widać na stronie przepisu, na karcie na liście i w panelu. Ustal, gdzie się sypie: wybór pliku, kadrowanie, upload (autoryzacja, limit rozmiaru Vercela, typ), zapis adresu w przepisie, wyświetlanie, routing `/api/photos/:id` w `vercel.json`. Testuj na szkicu, który potem usuniesz, albo lokalnie. **Opisz właścicielowi przyczynę, zanim zaczniesz naprawiać.**

### 2. Upload z pliku działa
Naprawa tego, co wyszło w diagnozie. Zachowanie zgodne z decyzjami w Claude Design: kadrowanie do 1:1 przed wysłaniem, zapis 1200×1200 JPEG, pasek postępu, limit rozmiaru.

### 3. ~~Zdjęcie z linku przez API~~ — ODPADA (korekta wyżej). Opcjonalnie: wgranie pliku przez API z kluczem `API_ADMIN_KEY`.
Właściciel przy wpisywaniu przepisu podaje link, a sesja Claude dodaje zdjęcie przez API, z kluczem `API_ADMIN_KEY`.
- Link może prowadzić prosto do obrazka albo do strony (np. przepisu na innym blogu). W drugim przypadku bierzemy główne zdjęcie strony (`og:image`).
- Zdjęcie jest **kopiowane do naszego magazynu**, nie hotlinkowane. Cudze linki znikają i bywają blokowane na obcych domenach.
- Kadrowanie 1:1 automatycznie (środek), ten sam rozmiar i format co przy uploadzie z pliku.
- Pobieranie po stronie serwera: walidacja typu i rozmiaru, timeout, tylko http(s), bez adresów wewnętrznych (ochrona przed SSRF).
- Działa też dla istniejącego przepisu: dodanie albo podmiana zdjęcia bez przepisywania reszty.

### 4. Gdzie trzymamy zdjęcia — decyzja z uzasadnieniem
Dziś bytea w Neonie (rozwiązanie „tymczasowe” od 10.2026). Ograniczenia: darmowy Neon ma 0,5 GB na całą bazę, a codzienny backup (`docs/backup-bazy.md`) obejmuje tylko bazę — zdjęcia poza bazą (np. Vercel Blob) wypadłyby z kopii, chyba że backup je też obejmie. Policz rozmiar jednego zdjęcia i przy ilu przepisach skończy się miejsce. Zarekomenduj: zostaje bytea albo przenosimy (i co wtedy z backupem). **Decyzja należy do właściciela — przedstaw ją przed wdrożeniem zmiany magazynu.** Diagnoza, naprawa i punkt 3 nie czekają na tę decyzję, jeśli działają na obecnym magazynie.

## Poza zakresem

- **Pole „wklej link” w formularzu.** To zmiana UI, więc najpierw idzie brief do designera (zrobi go analityk). Do tego czasu zdjęcie z linku dodaje sesja Claude przez API.
- Zdjęcia z telefonu, kilka zdjęć na przepis, zdjęcia kroków.
- Wygląd przepisu bez zdjęcia (osobny brief dla designera).

## Testowanie

`.env.local` wskazuje na produkcję. Testy zapisów rób lokalnie albo na szkicu, który usuniesz. `resetdb` jest zakazany.

## Do oddania

- Raport z diagnozy (punkt 1), zanim zaczniesz naprawiać.
- Kod (commity na `main`, bez pushowania bez prośby).
- `CLAUDE.md`, sekcja o przepisach dodawanych przez Claude: jak dodać zdjęcie z linku do nowego i do istniejącego przepisu (żądanie, klucz, co zrobić przy błędzie).
- Rekomendacja magazynu z liczbami (punkt 4).

### Checklista odbioru

- [ ] Przyczyna problemu opisana właścicielowi.
- [ ] Zdjęcie z pliku: dodanie w formularzu, zapis, widoczne na stronie przepisu, na karcie i w panelu (produkcja, sprawdzone przez właściciela).
- [ ] Zdjęcie z linku do obrazka przez API → widoczne na przepisie, 1:1, nasz adres `/api/photos/:id`.
- [ ] Zdjęcie z linku do strony (`og:image`) przez API.
- [ ] Podmiana zdjęcia w istniejącym przepisie bez utraty reszty treści (kroki, wskazówki, zamienniki).
- [ ] Zły link (404, nie-obrazek, za duży, adres wewnętrzny) → czytelny błąd, nic się nie zapisuje.
- [ ] Encja `images` (Unsplash): używana albo opisana jako martwa.
- [ ] Rekomendacja magazynu z liczbami przedstawiona właścicielowi.
- [ ] `CLAUDE.md` zaktualizowany.
