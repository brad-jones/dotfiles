function SetEnv {
    param($Key, $Value);
    Set-Item env:\$Key -Value $Value;
    [Environment]::SetEnvironmentVariable($Key, $Value, "User");
}

function AddToPath {
    param($Path);
    if (-Not ($env:PATH -like "*$Path*")) {
        SetEnv -Key "PATH" -Value "$Path;$env:PATH";
    }
}

function Exec {
    param (
        [scriptblock]$ScriptBlock,
        [string]$ErrorAction = $ErrorActionPreference
    )
    & @ScriptBlock
    if (($lastexitcode -ne 0) -and $ErrorAction -eq "Stop") {
        exit $lastexitcode
    }
}
