# assets

This packages makes use of <https://golang.org/pkg/embed/> to embed all files
below `./files` into the final compiled binary. We do this so our binaries are
truly self-contained, unlike say a tool like <https://www.chezmoi.io/>.
