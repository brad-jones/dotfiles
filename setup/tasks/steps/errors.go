package steps

import "github.com/brad-jones/goerr/v2"

// UnSupportedOsError is returned when any of our tasks/steps can not run on the detected os
var UnSupportedOsError = goerr.New("bootstrap does not support your operating system")

// Pkg404Error is returned when InstallGithubPkg could not pick an asset URL to download
var Pkg404Error = goerr.New("could not locate download url")

// ScheduledTaskError is returned when the RunAtLogin installer fails
var ScheduledTaskError = goerr.New("failed to create RunAtLogin task")

// ScoopError is returned when ever a scoop failure happens
var ScoopError = goerr.New("scoop failed")
