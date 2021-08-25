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
			no:      "v1.12.6",
			hashLin: "efb3b2196a7a29ae43d49dbff760c7be5af2f636206bd56a8345ebe56091b620",
			hashWin: "5b10fa700e0661f9cb50f444fc84d3d5e6852dedce3e5367830a2f01d541943b",
		},
		"fedora": {
			no:      "33.20210226.0",
			hashWin: "b623a59f936d953071b79683beff2d082b1848beeda6d645a5f01fd31226b988",
		},
		"https://dot.net/v1/dotnet-install.sh": {
			hash: "d998926533a32537fada686ac41918b67e3e97a169a43c99cfd4100711767145",
		},
		"https://dot.net/v1/dotnet-install.ps1": {
			hash: "db76db3455b522fb548ba22724f4b33e3d6a0ceffde020f619564451e3a9eb2b",
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
		"3.1.409",
		"2.1.816",
	}
}

func GetVersion(name string) *Version {
	return AllVersions()[name]
}
