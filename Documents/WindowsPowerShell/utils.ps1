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