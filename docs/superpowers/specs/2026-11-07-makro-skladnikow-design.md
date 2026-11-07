# Makro składników i dań — design

Status: przed designem UI (czeka na input z Claude Design).

## Po co

Użytkownik chce widzieć orientacyjne makro przepisu. Wartości tego samego produktu
różnią się między markami (np. fusilli z czerwonej soczewicy — ~2 g białka różnicy
między producentami), więc składnik trzyma makro jako zakres, a nie jedną liczbę.
Chodzi o orientację, nie o liczenie co do grama.

## Cele

- Makro składnika na 100 g jako zakres min–max (max opcjonalny → pojedyncza wartość).
- Makro przepisu wyliczane automatycznie ze składników.
- Wyświetlanie: `min – max (~średnia)`; przy braku zakresu sama wartość.

## Poza zakresem (na teraz)

- Warianty wielkości (małe / duże jajko) — zakres w przeliczniku na razie to pokrywa.
- Cukry (są częścią węglowodanów), mikroelementy, witaminy.
- Uwzględnianie zmiany masy przy obróbce (gotowanie, smażenie).
- Integracja z zewnętrznymi bazami produktów.

## Anty-cele

- Fałszywa precyzja: liczby nie mogą udawać dokładniejszych niż są.
- Makro „na 100 g dania” — liczone z surowej masy wprowadzałoby w błąd, odrzucone.

## Model danych

### Składnik — makro na 100 g

Pięć wartości, każda jako para `min` (wymagane, jeśli podane w ogóle) + `max` (opcjonalne):
kcal, białko, tłuszcz, węglowodany, błonnik.

Średnia nie jest przechowywana — liczona jako (min + max) / 2.

### Dziedziczenie z drzewka

Składnik bez własnych danych makro bierze zakres z dzieci: min = najniższe min wśród
dzieci, max = najwyższe max. Np. ogólny „makaron” wyjdzie szeroki — to celowe, pokazuje,
że przepis jest niedookreślony. Własne dane rodzica mają pierwszeństwo.

### Przeliczniki jednostek per składnik

Waga jednostki zależy od pary składnik + jednostka (łyżka oliwy ≈ 10 g, łyżka mąki ≈ 8 g),
więc przelicznik siedzi na składniku: `jednostka → gramy (min, opcjonalnie max)`.
Przykład: jajko — 1 szt ≈ 50–65 g.

- `g` — przeliczenie wprost. `ml` — 1:1 domyślnie, chyba że składnik ma własny przelicznik
  (decyzja Claude'a, do potwierdzenia).
- Jednostki „do smaku” (`NoQty`) — nie wliczane, bez ostrzeżenia (pomijalne makro).

## Wyliczanie makro przepisu

Dla każdej pozycji: zakres ilości (ilość bywa już zakresem „2–3”) × zakres gramów
jednostki × zakres makro / 100 → zakres makro pozycji. Suma po pozycjach = makro całego
dania; dzielone przez `Servings` = makro na porcję.

- Składniki opcjonalne — nie wliczane do sumy (decyzja Claude'a).
- Składnik bez danych makro lub bez przelicznika dla użytej jednostki — pomijany, a przy
  makro pokazuje się dopisek „nie wliczono: X, Y” (decyzja Claude'a — zawsze pokazujemy,
  z informacją o brakach, zamiast ukrywać).
- Makro ukryte tylko wtedy, gdy nic nie da się policzyć.
- Ręczne pole `Recipe.KcalPerServing` zostaje usunięte — zastępuje je wyliczenie.
  Wymaga migracji wersjonowanej (usunięcie kolumny) albo resetdb, jeśli w bazie są tylko seedy.

## Gdzie widać

- **Strona przepisu:** makro na porcję jako główna wartość, pod spodem na całe danie.
  Podpięte pod istniejące skalowanie porcji.
- **Składnik** (panel admina / widok składnika): makro na 100 g.
- **Panel admina — składnik:** edycja zakresów makro i listy przeliczników jednostek.
- Lista przepisów — do decyzji przy designie (np. kcal na porcję na karcie).

## UI i doprecyzowania z Claude Design (11.2026)

Źródło: `.design/Decyzje i backlog.md` → „Makroskładniki” oraz `.design/Handoff - Claude Code.md`
→ „Makroskładniki — prompt dla Claude Code”. Prototyp: `.design/Makro.html` → `.design/mk-main.jsx`
(`rnd`, `fmt`, `fmtText` przenosimy 1:1). Te pliki są źródłem prawdy dla UI, zaokrągleń,
walidacji, dziedziczenia przeliczników (w górę drzewa) i cache'u `recipes.nutrition`.

Najważniejsze, co dochodzi względem sekcji wyżej:
- Licznik porcji nie zmienia wartości na porcję — przelicza tylko „Całe danie”.
- Na liście przepisów: ` · ~520 kcal` w metadanych karty.
- Edycja składnika w panelu jako osobna strona (nie modal), odziedziczone wartości kursywą.
- Przelicznik jednostki dziedziczony z rodzica, jeśli składnik nie ma własnego.
- Wynik cache'owany w `recipes.nutrition` (JSON), przeliczany przy zapisie przepisu i przy
  zmianie makro / przelicznika / rodzica składnika.

Rozbieżności handoffu z kodem — rozstrzygnięte:
- Handoff: `ingredient_unit_conversions.unit` jako tekst. Kod ma słownik `units` →
  używamy `unit_id` (FK), unikalne `(ingredient_id, unit_id)`.
- Handoff: „dodaj `parent_id`” — już istnieje.
- `ml`/`l`: handoff przelicza w kodzie tylko `g`/`kg`. Decyzja Claude'a: `ml` = 1 g i
  `l` = 1000 g domyślnie (bulion, śmietana, mleko — błąd pomijalny dla orientacyjnych
  wartości), własny przelicznik składnika nadpisuje (np. oliwa 1 ml ≈ 0,92 g).

## Zostawione dla buildera

- Kształt tabel (kolumny na `ingredients` vs osobna tabela makro; tabela przeliczników).
- Gdzie liczyć (backend w mapperze / serwisie — preferowane, frontend tylko formatuje).
- Zaokrąglenia: kcal do całości, gramy do 1 miejsca po przecinku (sugestia).
