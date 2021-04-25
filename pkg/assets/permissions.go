package assets

var Permissions = []*Permission{
	{"~/.gnupg", FOLDER, 0700},
	{"~/.gnupg/**", FILE, 0600},
	{"~/.ssh", FOLDER, 0700},
	{"~/.ssh/**", FILE, 0600},
	{"~/.local/bin/*", FILE, 0755},
	{"~/.local/sbin/bin/*", FILE, 0755},
}
