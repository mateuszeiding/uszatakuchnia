# Prompt dla Claude Design — podgląd szkicu przepisu (mały)

## Kontekst

Uszata Kuchnia ma jednego admina, czyli właściciela. Pracuje na desktopie. Przepis ma status **Szkic** albo **Opublikowany**. Szkic nie jest widoczny publicznie: obcy dostaje „Nie ma takiego przepisu”. Zalogowany właściciel może go jednak otworzyć, bo API już mu go zwraca.

W panelu, na liście przepisów (zakładki statusów, tabela), opublikowany przepis ma przy sobie akcję „Zobacz na stronie” (ikona oka). **Szkic jej nie ma**, więc właściciel nie może zobaczyć, jak przepis będzie wyglądał dla czytelnika. Dziś obchodzi to tak, że publikuje przepis „Test”, a ten jest widoczny dla wszystkich.

## Zadanie

1. **Panel:** ikona oka „Zobacz na stronie” także przy szkicach. Czy jakoś się różni, np. podpisem „Podgląd szkicu”?
2. **Strona przepisu w trybie szkicu**, widoczna tylko dla właściciela:
   - wyraźny, ale niezawadzający pasek, np. „Szkic — widzisz go tylko ty” i „Edytuj”,
   - poza paskiem strona wygląda dokładnie tak, jak zobaczy ją czytelnik po publikacji.

   Rozstrzygnij, czy w pasku jest też „Opublikuj”.
3. **Lista zakupów i tryb gotowania szkicu:** czy pasek też tam jest? Ważne przy liście zakupów: link do listy szkicu u znajomego nie zadziała.
4. **Formularz przepisu:** czy potrzebny jest przycisk „Podgląd” (zapis i otwarcie strony w nowej karcie), czy wystarczy oko w panelu?

## Stany
- Strona szkicu: telefon 375 i desktop, light i dark. Właściciel czasem zagląda z telefonu, choć panel jest desktopowy.
- Opublikowany przepis oglądany przez zalogowanego właściciela: bez paska (albo uzasadnij inaczej).

## Do oddania
- Artboardy z punktów wyżej.
- Krótki handoff z checklistą.
- Wpis w „Decyzje i backlog.md”.
