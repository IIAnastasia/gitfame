package gitcommand

import (
	"bytes"
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/command"
	"log"

	"os/exec"
	"strconv"
	"strings"
)

type AuthorInfo struct {
	Files       map[string]struct{}
	Commits     map[string]struct{}
	LinesNumber int
}

func getOrDefault(
	authorInfo map[string]*AuthorInfo, author string,
) *AuthorInfo {
	authorStructPtr, contains := authorInfo[author]
	if contains {
		return authorStructPtr
	}
	authorStructPtr = &AuthorInfo{
		Files:       make(map[string]struct{}),
		Commits:     make(map[string]struct{}),
		LinesNumber: 0,
	}
	authorInfo[author] = authorStructPtr
	return authorStructPtr
}

func parseEmptyFile(
	filePath string, authorsInfo map[string]*AuthorInfo,
	args command.Arguments,
) {
	cmd := exec.Command(
		"git", "-C", args.Repository, "log", args.Commit, "-1",
		"--pretty=full", "--",
		filePath,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Println("ERROR while ", cmd.String())
		panic(err)
	}
	lines := strings.Split(out.String(), "\n")
	commit := strings.Split(lines[0], " ")[1]
	var author string
	if args.UseAuthor {
		parts := strings.Split(lines[1], " ")
		author = strings.Join(parts[1:len(parts)-1], " ")
	} else {
		parts := strings.Split(lines[2], " ")
		author = strings.Join(parts[1:len(parts)-1], " ")
	}
	authorInfoPtr := getOrDefault(authorsInfo, author)
	authorInfoPtr.Commits[commit] = struct{}{}
	authorInfoPtr.Files[filePath] = struct{}{}
}

func parseNotEmptyFile(
	filePath string,
	authorInfo map[string]*AuthorInfo, lines []string,
	needAuthor bool,
) {
	index := 0
	var author string
	for index < len(lines) {
		if len(lines[index]) == 0 {
			break
		}
		parts := strings.Split(lines[index], " ")
		lineNumb, _ := strconv.Atoi(parts[3])
		commit := parts[0]
		for i := 0; i < lineNumb; i++ {
			index++
			for lines[index][0] != ' ' && lines[index][0] != '\t' {
				if lines[index][0] == 'a' {
					if needAuthor {
						author = lines[index][len("author "):]
					}
					index += 4
				} else if lines[index][0] == 'c' {
					if !needAuthor {
						author = lines[index][len("committer "):]
					}
					index += 4
				} else if lines[index][0] == 'f' {
					//filename = lines[index][len("filename "):]
					index += 1
				} else {
					index += 1
				}
			}
			index++
		}
		authorStructPrt := getOrDefault(authorInfo, author)
		authorStructPrt.Files[filePath] = struct{}{}
		authorStructPrt.Commits[commit] = struct{}{}
		authorStructPrt.LinesNumber += lineNumb
	}

}

func Analyze(
	filePath string, authorsInfo map[string]*AuthorInfo,
	args command.Arguments,
) {
	log.Println("parsing file", filePath)
	if filePath == "" {
		return
	}
	cmd := exec.Command(
		"git", "-C", args.Repository, "blame", args.Commit, "--line-porcelain",
		filePath,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Println("ERROR while processing command: ", cmd.String())
		panic(err)
	}
	outText := out.String()
	if len(outText) == 0 || (len(outText) == 1 && outText[0] == '\n') {
		parseEmptyFile(filePath, authorsInfo, args)
	} else {
		parseNotEmptyFile(
			filePath,
			authorsInfo, strings.Split(outText, "\n"), args.UseAuthor,
		)
	}
}
