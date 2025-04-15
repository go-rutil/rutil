package fileops

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rohanthewiz/serr"
)

// ReadIniAsMapOfSections reads an ini file returning attributes as a map of sections to a map of key values.
// This is the better way to read an ini file
func ReadIniAsMapOfSections(filespec string) (AttributesBySection map[string]map[string]string, issues []serr.SErr, err error) {
	file, err := os.Open(filespec)
	if err != nil {
		return AttributesBySection, issues, serr.Wrap(err, "Error reading: "+filespec)
	}

	AttributesBySection = make(map[string]map[string]string, 4)

	currSection := ""
	var currSectionMap map[string]string

	scanner := bufio.NewScanner(file)

	lineNbr := 0
	for scanner.Scan() { // splits on lines by default
		line := strings.TrimSpace(scanner.Text())
		lineNbr++

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") { // skip lines starting with a comment
			continue
		}

		// Check for Section
		if strings.HasPrefix(line, "[") {
			b, _, f := strings.Cut(line, "]")
			if !f {
				fmt.Println("Mismatched '['  ']'")
				continue
			}
			if len(b) <= 1 {
				fmt.Println("Section empty")
				continue
			}
			section := b[1:] // newly encountered section

			// Close out old if exists; Create new section // This is a notable pattern
			if currSection != "" && currSectionMap != nil {
				AttributesBySection[currSection] = currSectionMap
			}
			// New section
			currSection = section
			currSectionMap = make(map[string]string, 4)

			continue
		}

		if currSection == "" {
			fmt.Printf("It seems there is no section defined before lineNbr: %d, line:\n%q\n", lineNbr, line)
			return AttributesBySection, issues, serr.NewSErr("Missing section header")
		}

		// Keys and Values
		bef, aft, fnd := strings.Cut(line, "=")
		if !fnd {
			continue
		}

		key := strings.TrimSpace(bef)
		if key == "" {
			issues = append(issues, serr.NewSErr("key is empty", "line", line, "lineNbr", fmt.Sprintf("%d", lineNbr)))
			continue
		}

		val := strings.TrimSpace(aft)

		// Don't make an issue of empty values

		if val == "" {
			// issues = append(issues, serr.NewSErr("Value is empty", "line", line, "lineNbr", fmt.Sprintf("%d", lineNbr)))
			continue
		}

		// Check for delimiters and comments
		if len(val) > 1 {
			// First check if value has surrounding quotes as **quotes have the highest precedence**
			// Don't trim after delimiters removed to allow spaces in values
			if strings.HasPrefix(val, `'`) {
				if idx := strings.IndexByte(val[1:], '\''); idx != -1 {
					val = val[1 : idx+1]
				}
			} else if strings.HasPrefix(val, `"`) {
				if idx := strings.IndexByte(val[1:], '"'); idx != -1 {
					val = val[1 : idx+1]
				}
				// For comments we do want to trim space
			} else if x := strings.IndexByte(val, '#'); x != -1 {
				val = strings.TrimSpace(val[:x])
			}
		}

		// Don't make an issue of empty values
		if val == "" {
			// issues = append(issues, serr.NewSErr("Value is empty", "line", line, "lineNbr", fmt.Sprintf("%d", lineNbr)))
			continue
		}

		currSectionMap[key] = val // Store the attribute
	}

	// Close out last if exists // This is a notable pattern
	if currSection != "" && currSectionMap != nil {
		AttributesBySection[currSection] = currSectionMap
	}

	if err := scanner.Err(); err != nil {
		return AttributesBySection, issues, serr.Wrap(err, "Error while scanning: ", filespec)
	}

	return
}
