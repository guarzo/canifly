package persist

import (
	"encoding/csv"
	"io"
	"log"
	"regexp"
	"strconv"

	"github.com/guarzo/canifly/internal/embed"
)

const sysPath = "static/systems.csv"
const wormHolePattern = `^J\d{6}$`
const systemCount = 8487
const nonWHsystemCount = 5885

var SysIdToName map[string]string
var SysNameToID map[string]string
var SysNameSuggestions []string

func init() {
	file, err := embed.StaticFiles.Open(sysPath)
	if err != nil {
		log.Fatalf("failed to read embedded plans: %v", err)
	}
	defer file.Close()

	pattern := regexp.MustCompile(wormHolePattern)

	reader := csv.NewReader(file)
	SysIdToName = make(map[string]string, systemCount)
	SysNameToID = make(map[string]string, systemCount)
	SysNameSuggestions = make([]string, nonWHsystemCount)

	// Read the file line by line
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break // If we've reached the end of the file, break the loop
			} else {
				log.Fatal(err) // If it's a different error, exit
			}
		}

		SysIdToName[record[0]] = record[1]
		SysNameToID[record[1]] = record[0]

		if !pattern.MatchString(record[1]) {
			SysNameSuggestions = append(SysNameSuggestions, record[1])
		}
	}
}

func GetSystemName(systemID int64) string {
	return SysIdToName[strconv.FormatInt(systemID, 10)]
}
