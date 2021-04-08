// +build windows

package steps

import (
	"fmt"
	"strings"

	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
)

// MustDeleteAllScoopBuckets will delete all scoop buckets accept for main
func MustDeleteAllScoopBuckets() {
	prefix := colorchooser.Sprint("delete-scoop-buckets")

	ps := gopwsh.MustNew()
	defer ps.Exit()

	stdout, _ := ps.MustExecute("scoop bucket list")
	for _, line := range strings.Split(stdout, "\r\n") {
		line = strings.TrimSpace(line)

		if line == "main" || line == "" {
			continue
		}

		if _, stderr := ps.MustExecute(fmt.Sprintf("scoop bucket rm %s", gopwsh.QuoteArg(line))); len(stderr) > 0 {
			goerr.Check(ScoopError, "failed to delete bucket", line)
		}

		fmt.Println(prefix, "| bucket", line, "deleted")
	}
}

func DeleteAllScoopBucketsAsync() *task.Task {
	return task.New(func() { DeleteAllScoopBucketsAsync() })
}
