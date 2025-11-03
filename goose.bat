@echo off
setlocal
powershell -NoProfile -ExecutionPolicy Bypass -Command "$envFile = '.env'; if (Test-Path $envFile) { Get-Content $envFile | ForEach-Object { if ($_ -match '^([^#=]+)=(.+)$') { $key = $Matches[1].Trim(); $value = $Matches[2].Trim(); [System.Environment]::SetEnvironmentVariable($key, $value) } } }; goose -dir ./db/migrations postgres $env:DB_DATABASE_URL %*"
endlocal
