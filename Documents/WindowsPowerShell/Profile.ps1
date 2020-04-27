function SetEnv {
    param($Key, $Value);
    New-Item env:\$Key -Value $Value;
    [Environment]::SetEnvironmentVariable($Key, $Value, "User");
}

function AddToPath {
    param($Path);
    if (-Not ($env:PATH -like "*$Path*")) {
        SetEnv -Key "PATH" -Value "$env:PATH;$Path";
    }
}

AddToPath -Path "$env:USERPROFILE\.local\bin";
AddToPath -Path "$env:USERPROFILE\.local\sbin\bin";
AddToPath -Path "$env:USERPROFILE\AppData\Roaming\Pub\Cache\bin";