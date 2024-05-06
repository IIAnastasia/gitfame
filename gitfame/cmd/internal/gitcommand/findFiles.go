package gitcommand

import (
	"bytes"
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/command"
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/filter"
	"log"
	"os"
	"os/exec"
	"strings"
)

// ищет файлы и фильтрует
func FindFiles(args command.Arguments) []string {
	log.Println("looking for files")
	cmd := exec.Command(
		"git", "-C", args.Repository, "ls-tree", "-r", args.Commit,
		"--name-only",
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		os.Exit(2)
	}
	return filter.Filter(strings.Split(out.String(), "\n"), args)
}
