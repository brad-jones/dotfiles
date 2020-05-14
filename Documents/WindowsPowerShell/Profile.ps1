. $env:USERPROFILE\Documents\WindowsPowerShell\utils.ps1;

# Define our path as early as possible
SetEnv -Key "Path" -Value (-join("C:\Windows\System32;",
"C:\Windows;",
"C:\Windows\System32\wbem;",
"C:\Windows\System32\WindowsPowerShell\v1.0;",
"C:\Program Files\Docker\Docker\resources\bin;",
"C:\ProgramData\DockerDesktop\version-bin;",
"C:\Program Files (x86)\Microsoft Visual Studio\2019\BuildTools\MSBuild\Current\bin;",
"$env:USERPROFILE\scoop\apps\7zip\current;",
"$env:USERPROFILE\scoop\apps\adoptopenjdk-hotspot\current\bin;",
"$env:USERPROFILE\scoop\apps\aws-vault\current;",
"$env:USERPROFILE\scoop\apps\aws\current\bin;",
"$env:USERPROFILE\scoop\apps\dart\current\bin;",
"$env:USERPROFILE\scoop\apps\dotnet-sdk\current;",
"$env:USERPROFILE\scoop\apps\drun\current;",
"$env:USERPROFILE\scoop\apps\git\current\bin;",
"$env:USERPROFILE\scoop\apps\go\current\bin;",
"$env:USERPROFILE\scoop\apps\gopass\current;",
"$env:USERPROFILE\scoop\apps\gpg\current\bin;",
"$env:USERPROFILE\scoop\apps\nodejs\current;",
"$env:USERPROFILE\scoop\apps\nodejs\current\bin;",
"$env:USERPROFILE\scoop\apps\nssm\current;",
"$env:USERPROFILE\scoop\apps\nuget\current;",
"$env:USERPROFILE\scoop\apps\packer\current;",
"$env:USERPROFILE\scoop\apps\putty\current;",
"$env:USERPROFILE\scoop\apps\pwsh\current;",
"$env:USERPROFILE\scoop\apps\python\current;",
"$env:USERPROFILE\scoop\apps\ruby\current\bin;",
"$env:USERPROFILE\scoop\apps\sed\current\bin;",
"$env:USERPROFILE\scoop\apps\terraform\current;",
"$env:USERPROFILE\scoop\apps\vscode\current\bin;",
"$env:USERPROFILE\scoop\apps\wget\current;",
"$env:USERPROFILE\scoop\apps\win32-openssh\current;",
"$env:USERPROFILE\scoop\shims;",
"$env:USERPROFILE\AppData\Roaming\Pub\Cache\bin;",
"$env:USERPROFILE\.local\bin;",
"$env:USERPROFILE\.go\bin;",
"$env:USERPROFILE\.local\sbin\bin;"));

SetEnv -Key "GIT_SSH" -Value "$env:USERPROFILE\scoop\apps\win32-openssh\current\ssh.exe";

# Tell goenv where to install go, personally I prefer all my tools and
# config to be hidden (ie: start with a dot) and all my actual data
# folders/files to unhidden.
SetEnv -Key "GOPATH" -Value "$env:USERPROFILE\.go";

# Configure aws-vault to use gopass to store our idenities
SetEnv -Key "AWS_VAULT_BACKEND" -Value "pass";
SetEnv -Key "AWS_VAULT_PASS_CMD" -Value "gopass";
SetEnv -Key "AWS_VAULT_PASS_PREFIX" -Value "aws-vault";
SetEnv -Key "AWS_VAULT_PASS_PASSWORD_STORE_DIR" -Value "$env:USERPROFILE\.password-store";