Function Test-CommandExists {
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
		$ErrorActionPreference=$oldPreference;
	}
}

if (!(Test-CommandExists -Command "scoop")) {
	iwr -useb get.scoop.sh | iex;
}