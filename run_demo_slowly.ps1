param (
    [string]$sqlFile = "tests\demo.sql",
    [int]$delaySeconds = 3
)

# Check if the SQL file exists
if (-not (Test-Path $sqlFile)) {
    Write-Error "SQL file not found: $sqlFile"
    exit 1
}

# Read the SQL file content
$sqlContent = Get-Content -Path $sqlFile -Raw

# Split the content into individual SQL statements
$statements = @()
$currentStatement = ""
$inMultilineStatement = $false

foreach ($line in ($sqlContent -split "`n")) {
    $line = $line.Trim()
    
    # Skip empty lines and comments
    if ([string]::IsNullOrWhiteSpace($line) -or $line.StartsWith("--")) {
        continue
    }
    
    # Add the line to the current statement
    $currentStatement += "$line`n"
    
    # Check if we're in a multiline statement (like CREATE TABLE)
    if ($line -match "CREATE TABLE" -and -not $line.EndsWith(";")) {
        $inMultilineStatement = $true
    }
    
    # Check if the statement is complete
    if ($line.EndsWith(";") -and -not $inMultilineStatement) {
        $statements += $currentStatement.Trim()
        $currentStatement = ""
    } elseif ($inMultilineStatement -and $line.EndsWith(");")) {
        $statements += $currentStatement.Trim()
        $currentStatement = ""
        $inMultilineStatement = $false
    }
}

# If there's any remaining statement, add it
if (-not [string]::IsNullOrWhiteSpace($currentStatement)) {
    $statements += $currentStatement.Trim()
}

# Create a temporary directory for individual SQL files
$tempDir = New-Item -ItemType Directory -Path "temp_sql_commands" -Force

# Process each statement
for ($i = 0; $i -lt $statements.Count; $i++) {
    $statement = $statements[$i]
    
    # Skip empty statements
    if ([string]::IsNullOrWhiteSpace($statement)) {
        continue
    }
    
    # Create a temporary file for this statement
    $tempFile = Join-Path $tempDir "command_$($i.ToString('000')).sql"
    Set-Content -Path $tempFile -Value $statement
    
    # Display the statement
    Write-Host "`n`n==============================================" -ForegroundColor Cyan
    Write-Host "Executing SQL statement ($($i+1)/$($statements.Count)):" -ForegroundColor Cyan
    Write-Host "==============================================" -ForegroundColor Cyan
    Write-Host $statement -ForegroundColor Yellow
    Write-Host "==============================================" -ForegroundColor Cyan
    
    # Execute the statement and filter out the "Executing SQL file:" line
    $output = & "cmd\sqlight.exe" $tempFile | Where-Object { $_ -notmatch "Executing SQL file:" }
    $output | ForEach-Object { Write-Host $_ }
    
    # Wait for the specified delay
    if ($i -lt $statements.Count - 1) {
        Write-Host "`nWaiting $delaySeconds seconds before next command..." -ForegroundColor Gray
        Start-Sleep -Seconds $delaySeconds
    }
}

# Clean up temporary files
Remove-Item -Path $tempDir -Recurse -Force

Write-Host "`n`nDemo completed!" -ForegroundColor Green
