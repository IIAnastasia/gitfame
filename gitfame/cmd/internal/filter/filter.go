package filter

import (
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/command"
	"log"
	"path/filepath"
	"strings"
)

func checkAnyGlob(path string, globs []string) bool {
	anyMatch := false
	for _, pattern := range globs {
		if matched, _ := filepath.Match(
			pattern,
			path,
		); matched {
			anyMatch = true
			break
		}
	}
	return anyMatch
}

func checkAnyExtension(path string, endings []string) bool {
	for _, ending := range endings {
		if strings.HasSuffix(path, ending) {
			return true
		}
	}
	return false
}

func Filter(filesToFilter []string, args command.Arguments) []string {
	log.Println("filtering files")
	var filtered []string
	for _, line := range filesToFilter {
		canUse := true
		for _, globs := range args.RestrictTo {
			if !checkAnyGlob(line, globs) {
				canUse = false
				break
			}
		}

		for _, endings := range args.Extensions {
			if !checkAnyExtension(line, endings) {
				canUse = false
				break
			}
		}

		for _, globs := range args.Exclude {
			if checkAnyGlob(line, globs) {
				canUse = false
				break
			}
		}
		if canUse {
			filtered = append(filtered, line)
		}
	}
	return filtered
}
