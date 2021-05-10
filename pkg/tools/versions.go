package tools

import "runtime"

type Version struct {
	No   string
	Hash string
}

type version struct {
	no      string
	hash    string
	hashWin string
	hashLin string
}

func mapVersion(in *version) *Version {
	h := in.hash
	if len(h) == 0 {
		if runtime.GOOS == "linux" {
			h = in.hashLin
		}
		if runtime.GOOS == "windows" {
			h = in.hashWin
		}
	}
	return &Version{
		No:   in.no,
		Hash: h,
	}
}

func AllVersions() map[string]*Version {
	versions := map[string]*Version{}

	for k, v := range map[string]*version{
		"aws-vault": {
			no:      "v6.3.1",
			hashWin: "acba56d994a5666b16c92928299aeb52d5691ddb3ddc76015ef20bfde3d29108",
			hashLin: "84cfab75012eb272add8b09cb2d295941b977cf2bf58b3fb3caa4c4adac6f17f",
		},
		"gsudo": {
			no:      "v0.7.3",
			hashWin: "136ac9437a248786a997b7a563e17383ec6779d58e01ccb9ca07fc9e2ebc70b5",
		},
		"gopass": {
			no:      "v1.9.2",
			hashLin: "a55885a50259b623dc58f734d3f2634e3564feae1d3eba0ec8b79da645ab1d55",
			hashWin: "047e4c48aa7f828102b71855f5705ded935d6559d4c5cca178220879b7fb1f5e",
		},
		"fedora": {
			no:      "33.20210226.0",
			hashWin: "b623a59f936d953071b79683beff2d082b1848beeda6d645a5f01fd31226b988",
		},
		"https://dot.net/v1/dotnet-install.sh": {
			hash: "c96360abc54d74454105df45cba5d6ac78c8d46859d9a1c2164df2a4dd09af6c",
		},
		"https://dot.net/v1/dotnet-install.ps1": {
			hash: "8b8fe64d51d2aed4ece74fadd20098cece7a82c04656da9b0841baaedd079a2c",
		},
		"https://wslstorestorage.blob.core.windows.net/wslblob/wsl_update_x64.msi": {
			hashWin: "d872c2ef8f86798daedc295c49a31fb75fb7ba7e46f0660036ff16e55f0926fd",
		},
	} {
		versions[k] = mapVersion(v)
	}

	return versions
}

func DotnetVersions() []string {
	return []string{
		"latest",
		"3.1.408",
		"2.1.815",
	}
}

func GetVersion(name string) *Version {
	return AllVersions()[name]
}
