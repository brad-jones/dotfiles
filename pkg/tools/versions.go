package tools

type Version struct {
	No   string
	Hash string
}

func AllVersions() map[string]*Version {
	return map[string]*Version{
		"aws-vault": {No: "5.4.4"},
		"gsudo":     {No: "v0.7.3", Hash: "136ac9437a248786a997b7a563e17383ec6779d58e01ccb9ca07fc9e2ebc70b5"},
		"gopass":    {No: "v1.9.2", Hash: "047e4c48aa7f828102b71855f5705ded935d6559d4c5cca178220879b7fb1f5e"},
	}
}

func GetVersion(name string) *Version {
	return AllVersions()[name]
}
