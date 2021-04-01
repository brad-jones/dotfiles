package steps

import (
	"fmt"
	"strings"

	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
)

// MustInstallScoopBucket will install the given scoop bucket
func MustInstallScoopBucket(bucketName, repo string) {
	prefix := colorchooser.Sprint("install-scoop-bucket-" + bucketName)

	ps := gopwsh.MustNew()
	defer ps.Exit()

	stdout, _ := ps.MustExecute("scoop bucket list")
	if strings.Contains(stdout, bucketName) {
		fmt.Println(prefix, "bucket already exists")
		return
	}

	if len(repo) > 0 {
		goexec.MustRunPrefixed(prefix, "powershell", "-Command",
			fmt.Sprintf("scoop bucket add %s %s",
				gopwsh.QuoteArg(bucketName),
				gopwsh.QuoteArg(repo),
			),
		)
		return
	}

	goexec.MustRunPrefixed(prefix, "powershell", "-Command",
		fmt.Sprintf("scoop bucket add %s",
			gopwsh.QuoteArg(bucketName),
		),
	)
}

func InstallScoopBucketAsync(bucketName, repo string) *task.Task {
	return task.New(func() { MustInstallScoopBucket(bucketName, repo) })
}
