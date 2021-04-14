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

function RetryCommand {
	[CmdletBinding()]
	param (
		[parameter(Mandatory, ValueFromPipeline)] 
		[ValidateNotNullOrEmpty()]
		[scriptblock] $ScriptBlock,
		[int] $RetryCount = 5,
		[int] $TimeoutInSecs = 1,
		[string] $SuccessMessage = "Command executed successfuly!",
		[string] $FailureMessage = "Failed to execute the command"
	)
		
	process {
		$Attempt = 1
		$Flag = $true
		
		do {
			try {
				$PreviousPreference = $ErrorActionPreference
				$ErrorActionPreference = 'Stop'
				Invoke-Command -ScriptBlock $ScriptBlock -OutVariable Result              
				$ErrorActionPreference = $PreviousPreference
				Write-Verbose "$SuccessMessage `n"
				$Flag = $false
			}
			catch {
				if ($Attempt -gt $RetryCount) {
					Write-Verbose "$FailureMessage! Total retry attempts: $RetryCount"
					Write-Verbose "[Error Message] $($_.exception.message) `n"
					$Flag = $false
				}
				else {
					Write-Verbose "[$Attempt/$RetryCount] $FailureMessage. Retrying in $TimeoutInSecs seconds..."
					Start-Sleep -Seconds $TimeoutInSecs
					$Attempt = $Attempt + 1
				}
			}
		}
		While ($Flag)
		
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

Function RmIfExists {
	Param ($Path);
	if (Test-Path -Path $Path) {
		Remove-Item -Path $Path -Recurse -Force;
	}
}

Function CommandExists {
	Param ($Command);
	$oldPreference = $ErrorActionPreference;
	$ErrorActionPreference = 'stop';
	try {
		if (Get-Command $Command) {
			return $true;
		}
	} catch {
		return $false;
	} finally {
		$ErrorActionPreference = $oldPreference;
	}
}

Function ModuleExists {
	Param ($Module);
	$oldPreference = $ErrorActionPreference;
	$ErrorActionPreference = 'stop';
	try {
		if (Get-InstalledModule -Name $Module) {
			return $true;
		}
	} catch {
		return $false;
	} finally {
		$ErrorActionPreference = $oldPreference;
	}
}