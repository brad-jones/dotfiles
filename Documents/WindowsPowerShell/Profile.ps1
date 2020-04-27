function SetEnv {
    param($Key, $Value);
    Set-Item env:\$Key -Value $Value;
    [Environment]::SetEnvironmentVariable($Key, $Value, "User");
}

function AddToPath {
    param($Path);
    if (-Not ($env:PATH -like "*$Path*")) {
        SetEnv -Key "PATH" -Value "$env:PATH;$Path";
    }
}

# Define our path as early as possible
AddToPath -Path "$env:USERPROFILE\.local\bin";
AddToPath -Path "$env:USERPROFILE\.local\sbin\bin";
AddToPath -Path "$env:USERPROFILE\AppData\Roaming\Pub\Cache\bin";

# Tell goenv where to install go, personally I prefer all my tools and
# config to be hidden (ie: start with a dot) and all my actual data
# folders/files to unhidden.
SetEnv -Key "GOPATH" -Value "$env:USERPROFILE\.go";

# Configure aws-vault to use gopass to store our idenities
SetEnv -Key "AWS_VAULT_BACKEND" -Value "pass";
SetEnv -Key "AWS_VAULT_PASS_CMD" -Value "gopass";
SetEnv -Key "AWS_VAULT_PASS_PREFIX" -Value "aws-vault";
SetEnv -Key "AWS_VAULT_PASS_PASSWORD_STORE_DIR" -Value "$env:USERPROFILE\.password-store";