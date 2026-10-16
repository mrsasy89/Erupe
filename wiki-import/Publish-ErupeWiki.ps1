[CmdletBinding()]
param(
    [switch]$Publish,
    [string]$WorkingDirectory = (Join-Path $env:TEMP ('Erupe-wiki-' + [guid]::NewGuid().ToString('N')))
)
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
if (-not (Get-Command git -ErrorAction SilentlyContinue)) { throw 'Installa Git for Windows prima di continuare.' }
$WikiUrl = 'https://github.com/mrsasy89/Erupe.wiki.git'
$SourceDirectory = Join-Path $PSScriptRoot 'wiki'
$Pages = @(Get-ChildItem -LiteralPath $SourceDirectory -Filter '*.md' -File)
if ($Pages.Count -ne 14) { throw 'Pacchetto incompleto: attese 14 pagine Markdown.' }
if (Test-Path -LiteralPath $WorkingDirectory) { throw 'La cartella di lavoro esiste gia. Scegli una cartella nuova.' }
# Disabilita temporaneamente le preferenze PowerShell che trattano ogni stderr
# nativo come errore: Git scrive su stderr anche i normali messaggi di progresso.
$PreviousErrorPreference = $ErrorActionPreference
$ErrorActionPreference = 'Continue'
& git -c core.autocrlf=false clone -- $WikiUrl $WorkingDirectory
$CloneExit = $LASTEXITCODE
$ErrorActionPreference = $PreviousErrorPreference
if ($CloneExit -ne 0) {
    throw 'Wiki non accessibile. Apri https://github.com/mrsasy89/Erupe/wiki e crea una prima pagina Home se non esiste; verifica inoltre il login di Git for Windows. Poi riesegui lo script.'
}
Push-Location -LiteralPath $WorkingDirectory
try {
    $BackupDirectory = Join-Path $WorkingDirectory '.git/wiki-import-backup'
    New-Item -ItemType Directory -Path $BackupDirectory -Force | Out-Null
    foreach ($Page in $Pages) {
        $Target = Join-Path $WorkingDirectory $Page.Name
        if (Test-Path -LiteralPath $Target) {
            Copy-Item -LiteralPath $Target -Destination (Join-Path $BackupDirectory $Page.Name)
        }
        Copy-Item -LiteralPath $Page.FullName -Destination $Target
    }
    $PageNames = @($Pages | ForEach-Object { $_.Name })
    & git add -- @PageNames
    if ($LASTEXITCODE -ne 0) { throw 'git add non riuscito.' }
    & git diff --cached --stat
    & git diff --cached --quiet
    $DiffExit = $LASTEXITCODE
    if ($DiffExit -eq 0) { Write-Host 'La wiki contiene gia queste pagine.'; return }
    if ($DiffExit -ne 1) { throw 'Impossibile controllare le modifiche.' }
    if (-not $Publish) {
        Write-Host "Anteprima preparata in: $WorkingDirectory"
        Write-Host 'Nessun push eseguito. Usa -Publish per pubblicare.'
        return
    }
    & git var GIT_AUTHOR_IDENT *> $null
    if ($LASTEXITCODE -ne 0) {
        throw 'Configura user.name e user.email in Git sul tuo PC, poi riesegui lo script.'
    }
    & git commit -m 'Copy upstream Erupe wiki and document client-verified quest types'
    if ($LASTEXITCODE -ne 0) { throw 'Commit non riuscito. Nessun push eseguito.' }
    # Push normale: niente force, cancellazioni o sostituzione della cronologia.
    & git push origin HEAD
    if ($LASTEXITCODE -ne 0) {
        throw 'Push non riuscito. La copia locale resta disponibile; non forzare il push.'
    }
    Write-Host 'Wiki pubblicata: https://github.com/mrsasy89/Erupe/wiki'
} finally {
    Pop-Location
}
