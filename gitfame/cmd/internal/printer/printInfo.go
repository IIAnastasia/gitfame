package printer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/command"
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/stats"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
)

func printTabular(
	authorInfo []stats.FinalAuthorInfo,
) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	tabStrings := []string{"Name\tLines\tCommits\tFiles"}
	for _, line := range authorInfo {
		tabStrings = append(
			tabStrings, line.Author+"\t"+strconv.Itoa(line.Lines)+"\t"+
				strconv.Itoa(line.Commits)+"\t"+strconv.Itoa(line.Files),
		)
	}

	for _, str := range tabStrings {
		_, err := fmt.Fprintln(w, str)
		if err != nil {
			return
		}
	}

	err := w.Flush()
	if err != nil {
		return
	}
}

func printCSV(authorInfo []stats.FinalAuthorInfo) {
	w := csv.NewWriter(os.Stdout)
	lines := [][]string{{"Name", "Lines", "Commits", "Files"}}
	for _, line := range authorInfo {
		lines = append(
			lines, []string{
				line.Author,
				strconv.Itoa(line.Lines),
				strconv.Itoa(line.Commits),
				strconv.Itoa(line.Files),
			},
		)
	}
	err := w.WriteAll(lines)
	if err != nil {
		return
	}
	w.Flush()
}

func printJSON(authorInfo []stats.FinalAuthorInfo) {
	var dicts []map[string]interface{}
	w := json.NewEncoder(os.Stdout)
	for _, line := range authorInfo {
		dicts = append(
			dicts, map[string]interface{}{
				"name":    line.Author,
				"lines":   line.Lines,
				"commits": line.Commits,
				"files":   line.Files,
			},
		)
	}
	err := w.Encode(dicts)
	if err != nil {
		return
	}
}

func printJSONLines(authorInfo []stats.FinalAuthorInfo) {
	w := json.NewEncoder(os.Stdout)
	for _, line := range authorInfo {
		err := w.Encode(
			map[string]interface{}{
				"name":    line.Author,
				"lines":   line.Lines,
				"commits": line.Commits,
				"files":   line.Files,
			},
		)
		if err != nil {
			return
		}
	}
}

func PrintInfo(
	args command.Arguments,
	authorInfo []stats.FinalAuthorInfo,
) {
	log.Println("converting output")
	switch args.Format {
	case "csv":
		printCSV(authorInfo)
	case "json":
		printJSON(authorInfo)
	case "json-lines":
		printJSONLines(authorInfo)
	case "tabular":
		printTabular(authorInfo)
	default:
		panic(4)
	}
}
