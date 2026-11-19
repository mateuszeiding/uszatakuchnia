# Przegląd użycia — 18.11.2026

Żywa strona https://uszatakuchnia.vercel.app, oglądana jak przez znajomego, który dostał link do przepisu i otwiera go pierwszy raz. Bez logowania, z pustym `localStorage` na start. Przeglądarka wbudowana: telefon 375×812 (light i dark), dodatkowo 375×640, czyli telefon z widocznymi paskami przeglądarki. Potem desktop 1280×800.

Opublikowane są dwa przepisy: **Ragù alla bolognese** (`/recipes/6`, ma zamienniki, sprzęt i wskazówki) i **Bułki do burgerów** (`/recipes/5`, ma sprzęt i wskazówki, bez zamienników). Żaden nie ma zdjęcia ani linku do innego przepisu.

Waga: **blokuje**: znajomy nie zrobi tego, po co przyszedł, albo zrobi źle. **Myli**: zrobi, ale się zatnie, zgubi albo źle zrozumie. **Drobiazg**: zauważy, ale nie przeszkodzi.

## Znaleziska (od najpoważniejszych)

### 1. Tryb gotowania ucina początek długiego kroku i nie da się do niego przewinąć: blokuje
- **Gdzie:** tryb gotowania na telefonie, oba motywy. Bułki krok 3 i krok 8 przy 375×812. Ragù krok 3 i krok 7 przy 375×640.
- **Co się dzieje:** treść kroku jest wyśrodkowana w pionie. Kiedy krok ze wskazówkami jest wyższy niż ekran, wystaje ponad górną krawędź obszaru kroku i chowa się pod nagłówkiem („Krok 8 z 14”). Przewijanie działa tylko w dół. Bułki, krok 8, przy 375×812: ukryte ~280 px, czyli ok. 8 pierwszych linii polecenia. Ekran zaczyna się w połowie zdania („przetaczaj ją po blacie…”). W Ragù przy 375×640 znika pierwsza linia kroku 3 („Na osobnej patelni rozgrzej oliwę i smaż…”) i kroku 7 („Duś na małym ogniu co najmniej 2 godziny…”). Przy 375×812 w Ragù znika etykieta grupy „MIĘSO”.
- **Dlaczego to problem:** tryb gotowania jest od tego, żeby przy kuchence czytać krok w całości. Tu znajomy dostaje polecenie bez początku i nie wie, że czegoś brakuje. W Ragù traci „na osobnej patelni”, a w kroku 7 czas duszenia. Im więcej wskazówek pod krokiem, tym gorzej, więc wskazówki z 18.11 pogłębiają problem.

### 2. Ragù: makaron jest w kroku 9, ale nie ma go w składnikach ani na liście zakupów: blokuje (treść)
- **Gdzie:** Ragù, strona przepisu i lista zakupów, wszystkie szerokości.
- **Co się dzieje:** krok 9 każe wymieszać ugotowany makaron z sosem, a wskazówka podaje rodzaje makaronu. Na liście składników i na liście zakupów makaronu nie ma, nie ma też jego ilości na 6 porcji.
- **Dlaczego to problem:** znajomy robi zakupy z listy i wraca bez makaronu. Nawet jeśli go ma, nie wie, ile ugotować.

### 3. Lista zakupów gubi doprecyzowanie składnika („mielona”, „wędzony”, „duże”): myli
- **Gdzie:** lista zakupów `/recipes/:id/zakupy`, telefon i desktop.
- **Co się dzieje:** na stronie przepisu pod ilością stoi „mielona” (wołowina), „wędzony” (boczek), „duże” (jajka). Na liście zakupów zostaje sama nazwa: „Wołowina 450 g”, „Boczek 150 g”.
- **Dlaczego to problem:** lista zakupów jest używana w sklepie, bez przepisu przed oczami. Znajomy kupi wołowinę w kawałku albo boczek surowy, bo informacja, która decyduje o zakupie, została na innej stronie.

### 4. Link wklejony do komunikatora nic nie mówi: myli
- **Gdzie:** HTML każdej strony (`/recipes/6`, `/recipes/6/zakupy`, lista).
- **Co się dzieje:** serwer zwraca ten sam szkielet dla każdego adresu: `<title>Uszata Kuchnia</title>`, `<html lang="">`, bez `meta description`, bez `og:title`, `og:description`, `og:image` i `twitter:*`. Tytuł przepisu („Ragù alla bolognese · uszatakuchnia”) ustawia dopiero JavaScript, a komunikatory go nie uruchamiają. Nieistniejące adresy też zwracają 200.
- **Dlaczego to problem:** Messenger czy WhatsApp pokaże przy każdym linku to samo „Uszata Kuchnia” albo goły adres `…/recipes/6`, bez nazwy dania i zdjęcia. Dwa linki od właściciela (przepis i lista zakupów) wyglądają identycznie. Wysyłanie linków to dziś główny sposób, w jaki aplikacja trafia do ludzi, więc pozycja „meta OG — później” z backlogu dotyczy już obecnych odbiorców.

### 5. Pierwszy ekran przepisu na telefonie to fioletowy prostokąt: myli
- **Gdzie:** strona przepisu, 375 px, oba przepisy (żaden nie ma zdjęcia). W dark jeszcze mocniej.
- **Co się dzieje:** zastępczy obrazek 1:1 zajmuje ~670 z 812 px pierwszego ekranu. Nazwa przepisu wypada na samym dole, a składniki zaczynają się mniej więcej dwa ekrany niżej, po metadanych i karcie „Ikony przy składnikach”. W trybie ciemnym obrazek jest jasnolawendowy, czyli największa i najjaśniejsza plama na ciemnej stronie.
- **Dlaczego to problem:** znajomy otwiera link i przez pierwsze sekundy widzi pusty, ozdobny prostokąt, bez żadnej treści. Zanim zobaczy, co to za danie, musi przewinąć, a do składników jeszcze dalej. „Przepis bez zdjęcia” jest w backlogu jako „później”, ale dziś dotyczy 100% przepisów.

### 6. „Zacznij od nowa” kasuje postęp jednym dotknięciem: myli
- **Gdzie:** nagłówek „Kroki” na stronie przepisu, telefon i desktop.
- **Co się dzieje:** przycisk pojawia się po pierwszym odhaczeniu, tuż obok „Tryb gotowania”. Jedno dotknięcie bez pytania czyści odhaczone kroki i składniki oraz wraca do domyślnych porcji. Wybrane zamienniki zostają.
- **Dlaczego to problem:** przy kuchence, mokrymi palcami, łatwo trafić obok „Tryb gotowania” i stracić cały postęp bez możliwości cofnięcia. Do tego reset działa połowicznie: porcje wracają, a zamienniki nie, więc znajomy nie wie, co właściwie wyczyścił.

### 7. Przełącznik motywu przy ciemnym telefonie reaguje dopiero za drugim razem: myli
- **Gdzie:** górny pasek, każda strona, telefon z systemowym trybem ciemnym.
- **Co się dzieje:** strona jest ciemna (za systemem), ale przycisk ma `aria-pressed="false"` i etykietę „Tryb ciemny”. Pierwsze dotknięcie zapisuje motyw ciemny, więc na ekranie nic się nie zmienia. Dopiero drugie przełącza na jasny. Etykieta brzmi „Tryb ciemny” w obu stanach. Na desktopie obok ikony stoi angielskie „light”.
- **Dlaczego to problem:** kto chce jasny ekran w kuchni, dotyka ikony, nic się nie dzieje i uznaje, że przycisk nie działa. Rozjazd z decyzjami: „ręczny przełącznik motywu” jest w „Później (nie MVP)”, a na żywej stronie już jest.

### 8. Wyszukiwarka nie znajduje bez polskich znaków: myli
- **Gdzie:** lista przepisów, pole „Przepis albo składnik…”.
- **Co się dzieje:** „bulki” → 0 przepisów („Bułki” → 1). „ragu” → 0 („Ragù”). „bulka” → 0 (nie łapie innej formy słowa). Zamienniki też nie są przeszukiwane („pancetta”, „wieprzowina” → 0).
- **Dlaczego to problem:** na telefonie większość ludzi pisze bez ogonków. Znajomy wpisuje „ragu”, widzi „Nic nie pasuje” i uznaje, że przepisu nie ma, chociaż przed chwilą dostał do niego link.

### 9. Prawie każdy filtr daje 0 wyników: myli
- **Gdzie:** lista przepisów, arkusz „Filtry” (telefon) i panel (desktop).
- **Co się dzieje:** 54 opcje, z czego 45 daje „Pokaż 0 przepisów”. Filtr czasu kończy się na „do 2 h”, a oba przepisy trwają 3–5 h i 3–4 h, więc każda opcja czasu daje zero. Liczba na przycisku „Pokaż N przepisów” uprzedza o zerze, zanim się go zatwierdzi.
- **Dlaczego to problem:** przy dwóch przepisach panel z 54 opcjami obiecuje bazę, której nie ma. Znajomy klika „zupa”, „polska”, „do 1 h” i za każdym razem trafia w pustkę. Łatwo o wrażenie, że strona jest pusta albo zepsuta.

### 10. Przepis, którego nie ma, wygląda jak awaria: myli
- **Gdzie:** `/recipes/999`, a także każdy wysłany wcześniej link do przepisu, który potem usunięto albo cofnięto do szkicu.
- **Co się dzieje:** „Nie udało się wczytać przepisu.” i „← Wróć do listy przepisów”. Tytuł karty: „Uszata Kuchnia”. Nieznany adres (np. `/przepis-od-kuby`) bez słowa przekierowuje na listę.
- **Dlaczego to problem:** komunikat brzmi jak błąd sieci, więc znajomy odświeża stronę i czeka, zamiast zrozumieć, że przepisu już nie ma.

### 11. Kroki nie zmieniają się razem z zamiennikiem: drobiazg
- **Gdzie:** Ragù, kroki i tryb gotowania, po wybraniu wieprzowiny zamiast wołowiny.
- **Co się dzieje:** wiersz składnika mówi „Wieprzowina, zamiast: wołowina”, a krok 3 (też w trybie gotowania) nadal „smaż wołowinę”.
- **Dlaczego to problem:** chwila zawahania, czy zamiennik się zapisał. Przy jednym składniku łatwo się domyślić, przy kilku zamiennikach będzie trudniej.

### 12. Skalowanie daje „2,7 łodygi” i „2⅔ łyżki”: drobiazg
- **Gdzie:** Ragù, składniki i tryb gotowania, 8 porcji zamiast 6.
- **Co się dzieje:** seler naciowy „2,7 łodygi” (dziesiętnie, nie zaokrąglony), oliwa „2⅔ łyżki” i „1⅓ łyżki”, passata „665 g”. Cebula i marchew dobrze: „1–2 sztuki”.
- **Dlaczego to problem:** 0,7 łodygi i ⅔ łyżki trudno odmierzyć. Rozjazd z decyzją „jednostki sztukowe zaokrąglamy do całości lub zakresu”: łodyga zachowuje się jak waga, nie jak sztuka.

### 13. Lista zakupów wysłana linkiem nie zawiera wybranych zamienników: drobiazg
- **Gdzie:** przycisk „Lista zakupów” → adres `/recipes/6/zakupy?porcje=6`.
- **Co się dzieje:** w adresie są tylko porcje. Zamienniki siedzą w pamięci urządzenia, więc u odbiorcy lista pokaże wołowinę, choć wysyłający wybrał wieprzowinę albo „pół na pół”.
- **Dlaczego to problem:** typowa sytuacja „wyślij mi, co kupić” daje drugiej osobie inną listę niż ta, którą widzi wysyłający, a żadne z nich o tym nie wie.

### 14. Powrót z listy zakupów ląduje na górze przepisu: drobiazg
- **Gdzie:** lista zakupów, link „‹ Przepis”, telefon.
- **Co się dzieje:** link otwiera przepis od góry, czyli od fioletowego prostokąta (pkt 5), a nie przy składnikach, skąd przyszedł znajomy. Dodaje też nowy wpis w historii, więc systemowe „wstecz” wraca potem na listę zakupów.
- **Dlaczego to problem:** po przejrzeniu listy trzeba znowu przewijać przez dwa ekrany. Przyciski „wstecz” zaczynają krążyć między stronami.

### 15. Arkusz „Sprzęt i składniki” w trybie gotowania: drobiazg
- **Gdzie:** tryb gotowania → „Składniki”, telefon.
- **Co się dzieje:** okno na środku ekranu bez przycisku zamknięcia. Zamyka się tylko dotknięciem w tło, czego nic nie sugeruje. Nie ma w nim licznika porcji ani informacji „Z zamiennikami: n”, tylko liczba porcji w tytule. Pod przyciskami na telefonie stoi podpowiedź „Strzałki ← → zmieniają krok”, choć telefon nie ma strzałek.
- **Dlaczego to problem:** przy kuchence trzeba zgadywać, jak wrócić do kroku. Podpowiedź o klawiaturze to na telefonie szum.

### 16. „Zaloguj” i brak słowa o tym, czyja to strona: drobiazg
- **Gdzie:** górny pasek na każdej stronie, lista przepisów.
- **Co się dzieje:** jedynym przyciskiem w pasku (obok motywu) jest „Zaloguj”. Strona nigdzie nie mówi, czyje to przepisy i że do czytania nie trzeba konta. Lista zaczyna się od samego „Przepisy”.
- **Dlaczego to problem:** znajomy, który wszedł z linku, może uznać, że czegoś mu brakuje, póki się nie zaloguje. Nie wie też, że to przepisy kolegi, a nie kolejny serwis z przepisami.

### 17. Treść przepisów: drobiazg (do poprawy w panelu)
- **Ragù:** bazylia jest „(opcjonalnie)”, ale krok 6 każe dodać jej łodyżki, a krok 8 je wyjąć, bez „jeśli masz”. Krok 9: „z sosem na patelni”, a sos cały czas jest w garnku.
- **Bułki:** „1 sztuka / duże / Jajka” (liczba mnoga, zgłoszone 17.11, nadal jest). Literówki we wskazówkach: „moga”, „skeiło”, „Miske”. „Sezam: do smaku, do posypania”. Krok 1 podaje wodę „~35-40°C”, a składnik „ok. 40°C”.
- **Lista zakupów:** „do: do podania” (podwójne „do”). Passata w dziale „Warzywa i grzyby”, choć w sklepie stoi przy przetworach.

### 18. Drobne niespójności: drobiazg
- Tytuł karty przeglądarki: „Uszata Kuchnia” na liście, „… · uszatakuchnia” na przepisie.
- Zamiennik „Bulion warzywny — zawiera: seler”, choć przepis i tak ma seler. Decyzja mówi o dopisku, „gdy zamiennik **dodaje** alergen”.

## Co działa dobrze

- **Zamienniki** robią dokładnie to, co obiecują. Arkusz od dołu na telefonie, popover na desktopie, opisy czytelne („w przepisie”, „po połowie ilości”, „zawiera: jaja”). Wybór od razu zmienia kcal (~390 → ~420), alergeny (Grana Padano dopisuje „jaja”) i listę zakupów. „Pół na pół” daje na liście dwie pozycje po 225 g. Linia „Z zamiennikami: 2 · Przywróć przepis” i zdanie w karcie makro jasno mówią, w jakim wariancie jesteś.
- **Karta „Ikony przy składnikach”** tłumaczy ⇄ jednym zdaniem i można ją ukryć. Bez niej ikona byłaby zagadką.
- **Licznik porcji** jest przyklejony pod paskiem na telefonie, a nagłówki grup przyklejają się pod nim. „Porcje” w metadanych podążają za licznikiem. Ułamki i zakresy („1–2 sztuki”) czytają się dobrze.
- **Tryb gotowania** (poza pkt 1): wraca do pierwszego nieodhaczonego kroku (po odhaczeniu 1–4 otwiera „Krok 5 z 9”), „Dalej” odhacza krok, na końcu jest „Gotowe”. Przyciski są duże, wskazówki duże pod krokiem, sprzęt na górze arkusza składników.
- **Wskazówki pod krokami** czytają się naturalnie jako „dlaczego / co jeśli” i szarzeją po odhaczeniu kroku. Problem „rady po ugotowaniu” z feedbacku jest rozwiązany.
- **Sprzęt** to krótka karta tuż przed krokami. „(opcjonalnie)” i doprecyzowania („do rozbijania mięsa”) od razu odpowiadają na pytanie „czy mam czym to zrobić”.
- **Lista zakupów:** działy sklepu, sumowanie (oliwa „3 łyżki, do: baza, mięso”), opcjonalne osobno („jeśli chcesz”), bez wody, odhaczone zostaje na miejscu, licznik „1 z 12”, na desktopie „Drukuj”.
- **Lista przepisów:** szukanie po składnikach z dopiskiem „zawiera: Seler naciowy” na karcie. Pusty wynik wyjaśnia, co zrobić („Szukamy też po składnikach”). Filtry pokazują liczbę wyników przed zatwierdzeniem, a stan trafia do adresu.
- **Ciemny motyw** jest spójny na przepisie, w arkuszach i na liście zakupów. Kontrast przycisku głównego jest poprawiony (5,4:1, punkt 3 z przeglądu 17.11).
- **Desktop:** przyklejony panel składników przewija się w środku i licznik porcji zostaje na miejscu. Karta makro ma układ dwukolumnowy (punkt 4 z 17.11). Popover zamienników nie zasłania składników.
- Z przeglądu 17.11 poprawione: woda zniknęła z listy zakupów, tytuł zwiniętej karty makro mieści się w jednej linii.

## Funkcje „wdrożone, niesprawdzone” a decyzje w Claude Design

Porównanie z `project/Decyzje i backlog.md` (stan z 18.11).

| Funkcja | Zgodne z decyzjami | Rozjazdy / uwagi |
|---|---|---|
| Linki do przepisów, „Przydatne przepisy” | nie dało się sprawdzić | Nadal żaden opublikowany przepis nie ma linku (jak 17.11). Ragù aż prosi o link w kroku 9, ale makaronu nie ma nawet w składnikach (pkt 2). |
| Tryb gotowania | jeden krok naraz, duże Wstecz / Dalej, grupy i nagłówki w „Składnikach”, tytuł „Sprzęt i składniki” ze sprzętem na górze, ta sama ikona ⇄ i wybór wspólny ze stroną | Ucinanie długich kroków (pkt 1): błąd wdrożenia, nie decyzja. Arkusz składników to okno na środku bez zamknięcia (pkt 15), czego decyzje nie rozstrzygają. Wake Lock (ekran nie gaśnie): nie da się sprawdzić w przeglądarce wbudowanej. |
| Wskazówki przy krokach | linia pod treścią, żarówka w akcencie, bez tła i etykiety, szarzeją po odhaczeniu, duże w trybie gotowania; „Zanim zaczniesz” nie występuje (wszystko przeniesione) | brak |
| Sprzęt | karta między składnikami (z makro) a krokami na telefonie; na desktopie prawa kolumna w kolejności Sprzęt → Kroki → Wartości odżywcze; „(opcjonalnie)” kursywą, doprecyzowanie pod nazwą; bez odhaczania | Nie sprawdziłem układu dwóch kolumn od 4 pozycji (Bułki, 7 pozycji, desktop). |
| Zamienniki | goła ikona ⇄, arkusz / popover, radio z ilością, „zamiast: …”, „pół na pół”, karta objaśnień nad składnikami z „Ukryj”, „Z zamiennikami: n · Przywróć przepis”, zdanie w karcie makro, alergeny i lista zakupów według wariantu | „zawiera: seler” przy alergenie, który przepis już ma (pkt 18). Kroki nie podążają za zamiennikiem (pkt 11), czego decyzje nie obejmują. Link do listy zakupów nie niesie zamienników (pkt 13): decyzja zapisuje wybór „na urządzeniu”, więc to zgodne z nią, ale kłóci się z „link można wysłać”. Dopisku „bez danych o wartościach odżywczych” nie widziałem. |
| Poza listą | | Ręczny przełącznik motywu jest na żywo, a w decyzjach figuruje jako „Później (nie MVP)” (pkt 7). |

## Czego nie sprawdziłem

- Wake Lock w trybie gotowania (wymaga prawdziwego telefonu).
- Wydruk listy zakupów (podgląd wydruku).
- Desktop w motywie ciemnym (sprawdzony tylko telefon w dark).
- Link do innego przepisu i blok „Przydatne przepisy” (brak danych).
