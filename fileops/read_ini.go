package fileops

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rohanthewiz/serr"
)

// ReadIni reads an ini file returning keys scoped by section and their values as a map
func ReadIni(filespec string) (results map[string]string, issues []serr.SErr, err error) {
	results = make(map[string]string, 16)

	file, err := os.Open(filespec)
	if err != nil {
		return results, issues, serr.Wrap(err, "Error reading: "+filespec)
	}
	scanner := bufio.NewScanner(file)
	currSection := ""

	for scanner.Scan() { // splits on lines by default
		line := strings.TrimSpace(scanner.Text())

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
			currSection = b[1:]
			// fmt.Println("New section:", currSection)
		}

		if currSection == "" {
			fmt.Printf("It seems there is no section defined before this line:\n%q\n", line)
			return results, issues, serr.NewSErr("Missing section header")
		}

		// Keys and Values
		bef, aft, fnd := strings.Cut(line, "=")
		if !fnd {
			continue
		}

		key := strings.TrimSpace(bef)
		if key == "" {
			issues = append(issues, serr.NewSErr("key is empty", "line", line))
			continue
		}

		val := strings.TrimSpace(aft)
		if val == "" {
			issues = append(issues, serr.NewSErr("Value is empty", "line", line))
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

		if val == "" {
			issues = append(issues, serr.NewSErr("Value is empty", "line", line))
			continue
		}

		results[currSection+"::"+key] = val
	}

	if err := scanner.Err(); err != nil {
		return results, issues, serr.Wrap(err, "Error while scanning: ", filespec)
	}

	return
}
