#!/usr/bin/env bash
# Kopia całej bazy: pg_dump (format custom) → gpg (AES256, hasło) → jeden plik .dump.gpg.
# Używa go .github/workflows/db-backup.yml; lokalnie działa w Git Bash (pg_dump i gpg w PATH).
#
# Wejście (zmienne środowiskowe):
#   DATABASE_URL  connection string; na Neonie wersja BEZ "-pooler" w hoście
#   PASSPHRASE    hasło szyfrowania, min. 20 znaków
#   OUT_DIR       katalog na plik (domyślnie ./backup)
# Wyjście: $OUT_DIR/uszatakuchnia-<data>_<godzina>Z.dump.gpg (czas UTC);
#          w GitHub Actions także `name` i `file` w $GITHUB_OUTPUT.
#
# Repo jest publiczne, logi też: do logu nie trafia connection string, hasło ani dane,
# tylko wersje, rozmiar, suma kontrolna i wynik kontroli.
set -euo pipefail
umask 077

: "${DATABASE_URL:?brak DATABASE_URL}"
: "${PASSPHRASE:?brak PASSPHRASE}"
if [ "${#PASSPHRASE}" -lt 20 ]; then
  echo "::error::PASSPHRASE ma mniej niż 20 znaków" >&2
  exit 1
fi

# W Actions maskujemy części connection stringa (host, użytkownik, hasło, baza),
# gdyby pojawiły się w komunikacie błędu pg_dump. Sam sekret GitHub maskuje tylko w całości.
if [ -n "${GITHUB_ACTIONS:-}" ]; then
  python3 - <<'PY'
import os, urllib.parse
u = urllib.parse.urlsplit(os.environ["DATABASE_URL"])
vals = {u.hostname, u.username, u.password, u.path.lstrip("/")}
if u.password:
    vals.add(urllib.parse.unquote(u.password))
if u.hostname:
    vals.add(u.hostname.split(".")[0])
for v in vals:
    if v and len(v) >= 4:
        print(f"::add-mask::{v}")
PY
fi

OUT_DIR="${OUT_DIR:-backup}"
mkdir -p "$OUT_DIR"
NAME="uszatakuchnia-$(date -u +%Y-%m-%d_%H%M)Z"
OUT="$OUT_DIR/$NAME.dump.gpg"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
mkdir -m 700 "$WORK/gnupg"
printf '%s' "$PASSPHRASE" >"$WORK/pass"
GPG=(gpg --homedir "$WORK/gnupg" --batch --yes --quiet --no-symkey-cache
  --pinentry-mode loopback --passphrase-file "$WORK/pass")

echo "$(pg_dump --version), $(gpg --version | head -1)"

# Jedno połączenie, jedna transakcja tylko do odczytu (snapshot) na czas zrzutu.
pg_dump --format=custom --compress=9 --no-password --dbname="$DATABASE_URL" |
  "${GPG[@]}" --symmetric --cipher-algo AES256 --s2k-digest-algo SHA512 \
    --s2k-count 65011712 --compress-algo none --output "$OUT"

# Kontrola 1: plik się odszyfrowuje i jest kompletnym archiwum pg_dump z kluczowymi tabelami.
"${GPG[@]}" --decrypt "$OUT" | pg_restore --list >"$WORK/toc"
grep -E '^; +(Dumped from|Dumped by|TOC Entries)' "$WORK/toc" | sed 's/^; *//'
for t in recipes recipe_steps recipe_ingredients ingredients photos schema_migrations schema_state; do
  if ! grep -qE "TABLE DATA public $t " "$WORK/toc"; then
    echo "::error::W kopii brakuje danych tabeli $t" >&2
    exit 1
  fi
done
echo "Tabel z danymi w kopii: $(grep -c ' TABLE DATA ' "$WORK/toc")"

# Kontrola 2: bez hasła pliku nie da się odczytać.
if pg_restore --list "$OUT" >/dev/null 2>&1; then
  echo "::error::Plik kopii da się odczytać bez hasła" >&2
  exit 1
fi
echo "Bez hasła: pg_restore odrzuca plik (OK)"

echo "Plik: $NAME.dump.gpg, $(du -h "$OUT" | cut -f1), sha256 $(sha256sum "$OUT" | cut -d' ' -f1)"

if [ -n "${GITHUB_OUTPUT:-}" ]; then
  {
    echo "name=$NAME"
    echo "file=$OUT"
  } >>"$GITHUB_OUTPUT"
fi
