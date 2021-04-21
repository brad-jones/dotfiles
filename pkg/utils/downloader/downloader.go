package downloader

import (
	"fmt"
	"time"

	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/cavaliercoder/grab"
)

func DownloadWithProgress(prefix, src, dst string) (fileName string, err error) {
	goerr.Handle(func(e error) { err = e })

	// Create a new request
	req, err := grab.NewRequest(dst, src)
	goerr.Check(err, "failed to create NewRequest", dst, src)

	// Make the request
	fmt.Println(prefix, "|", "requesting", req.URL())
	resp := grab.NewClient().Do(req)
	fmt.Println(prefix, "|", resp.HTTPResponse.Status)

	// Start outputting progress information
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			fmt.Println(prefix, "|", fmt.Sprintf(
				"transferred %v / %v bytes (%.2f%%) ETA: %vs",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress(),
				time.Until(resp.ETA()).Seconds(),
			))

		case <-resp.Done:
			break Loop
		}
	}

	// Check for any errors
	goerr.Check(resp.Err(), "download failed")

	// Return the final resting place for the downloaded file
	fmt.Println(prefix, "|", "downloaded", resp.Filename)
	return resp.Filename, nil
}

func MustDownloadWithProgress(prefix, src, dst string) string {
	fileName, err := DownloadWithProgress(prefix, src, dst)
	goerr.Check(err)
	return fileName
}

func DownloadWithProgressAsync(prefix, src, dst string) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustDownloadWithProgress(prefix, src, dst))
	})
}
