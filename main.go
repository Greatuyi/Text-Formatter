package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// ANSI escape codes for terminal text formatting (used only in stdout)
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	Bold         = "\033[1m"
	Italic       = "\033[3m"
	Underline    = "\033[4m"
)

// Airport represents details of an airport.
type Airport struct {
	Name         string
	ISOCountry   string
	Municipality string // city name
	ICAOCode     string
	IATACode     string
	Coordinates  string
}

// airportMap stores airport info using IATA or ICAO codes as keys.
var airportMap map[string]*Airport

func main() {
	// Define a flag for displaying help.
	helpFlag := flag.Bool("h", false, "Display usage information")
	flag.Parse()

	if *helpFlag {
		printUsage()
		return
	}

	// Get command-line arguments.
	args := flag.Args()
	if len(args) != 3 {
		printUsage()
		return
	}

	inputPath := args[0]
	outputPath := args[1]
	airportLookupPath := args[2]

	if !fileExists(inputPath) {
		printError("Input file not found")
		return
	}
	if !fileExists(airportLookupPath) {
		printError("Airport lookup file not found")
		return
	}

	if err := loadAirportData(airportLookupPath); err != nil {
		printError(fmt.Sprintf("Airport lookup file is malformed: %v", err))
		return
	}

	input, err := os.ReadFile(inputPath)
	if err != nil {
		printError(fmt.Sprintf("Error reading input file: %v", err))
		return
	}

	// Process the content in two ways:
	// 1. Plain output for the file (no ANSI codes)
	// 2. Highlighted output for the terminal
	plainOutput := plainProcessContent(string(input))
	highlightedOutput := highlightProcessContent(string(input))

	// Write plain output to file.
	if err := os.WriteFile(outputPath, []byte(plainOutput), 0644); err != nil {
		printError(fmt.Sprintf("Error writing output file: %v", err))
		return
	}

	printSuccess("Processing completed successfully!")

	// Print highlighted output to stdout.
	fmt.Printf("\n%s%s=== Processed Output ===%s\n\n", Bold, ColorBlue, ColorReset)
	fmt.Println(highlightedOutput)
}

// printUsage prints the usage information.
func printUsage() {
	fmt.Printf("%s%sItinerary usage:%s\n", Bold, Underline, ColorReset)
	fmt.Printf("%sgo run . ./input.txt ./output.txt ./airport-lookup.csv%s\n", Italic, ColorReset)
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// loadAirportData loads airport data from a CSV into airportMap.
// It supports non-standard CSV column order by using header names.
func loadAirportData(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return err
	}

	// Build a map from trimmed, lowercased header name to index.
	requiredColumns := []string{"name", "iso_country", "municipality", "icao_code", "iata_code", "coordinates"}
	columnMap := make(map[string]int)
	for i, column := range header {
		key := strings.TrimSpace(strings.ToLower(column))
		columnMap[key] = i
	}
	// Debug: print the header mapping (remove if not needed)
	// fmt.Printf("Header mapping: %#v\n", columnMap)

	// Ensure all required columns exist.
	for _, req := range requiredColumns {
		if _, exists := columnMap[req]; !exists {
			return fmt.Errorf("missing required column: %s", req)
		}
	}

	airportMap = make(map[string]*Airport)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		// Skip empty records.
		if len(record) == 0 {
			continue
		}
		if len(record) != len(header) {
			return fmt.Errorf("malformed record")
		}
		name := record[columnMap["name"]]
		iataCode := record[columnMap["iata_code"]]
		icaoCode := record[columnMap["icao_code"]]

		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("empty name in record")
		}
		if strings.TrimSpace(iataCode) == "" && strings.TrimSpace(icaoCode) == "" {
			return fmt.Errorf("record has no IATA or ICAO code")
		}

		airport := &Airport{
			Name:         name,
			ISOCountry:   record[columnMap["iso_country"]],
			Municipality: record[columnMap["municipality"]],
			ICAOCode:     icaoCode,
			IATACode:     iataCode,
			Coordinates:  record[columnMap["coordinates"]],
		}

		// Map both IATA and ICAO codes.
		if iataCode != "" {
			airportMap[iataCode] = airport
		}
		if icaoCode != "" {
			airportMap[icaoCode] = airport
		}
	}

	return nil
}


// Plain (Non-Formatted) Processing Functions
// Used for writing plain text to the output file.

func plainProcessContent(content string) string {
	content = plainProcessAirportCodes(content)
	content = plainProcessDatesAndTimes(content)
	content = trimHorizontalWhitespace(content)
	content = trimVerticalWhitespace(content)
	return content
}

// plainProcessAirportCodes replaces airport codes with plain text names or cities.
// With "*" prefix it outputs the municipality.
func plainProcessAirportCodes(content string) string {
	// IATA codes: supports *#ABC
	iataRegex := regexp.MustCompile(`(\*?)#([A-Z]{3})`)
	content = iataRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := iataRegex.FindStringSubmatch(match)
		code := groups[2]
		if airport, exists := airportMap[code]; exists {
			if groups[1] == "*" {
				return plainCity(airport)
			}
			return plainAirport(airport)
		}
		return match
	})

	// ICAO codes: supports *##ABCD
	icaoRegex := regexp.MustCompile(`(\*?)##([A-Z]{4})`)
	content = icaoRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := icaoRegex.FindStringSubmatch(match)
		code := groups[2]
		if airport, exists := airportMap[code]; exists {
			if groups[1] == "*" {
				return plainCity(airport)
			}
			return plainAirport(airport)
		}
		return match
	})
	return content
}

// plainAirport returns the airport name.
func plainAirport(airport *Airport) string {
	return airport.Name
}

// plainCity returns the municipality (city) if available.
func plainCity(airport *Airport) string {
	if strings.TrimSpace(airport.Municipality) != "" {
		return airport.Municipality
	}
	return plainAirport(airport)
}

// plainProcessDatesAndTimes replaces date/time placeholders with plain formatted dates/times.
func plainProcessDatesAndTimes(content string) string {
	// Dates: D(...)
	dateRegex := regexp.MustCompile(`D\(([0-9T:.Z+-]{16,})\)`)
	content = dateRegex.ReplaceAllStringFunc(content, func(match string) string {
		dateStr := match[2 : len(match)-1]
		t, err := time.Parse("2006-01-02T15:04Z", dateStr)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04-07:00", dateStr)
			if err != nil {
				return match
			}
		}
		return t.Format("02 Jan 2006")
	})

	// 12-hour time: T12(...)
	time12Regex := regexp.MustCompile(`T12\(([0-9T:.Z+-]{16,})\)`)
	content = time12Regex.ReplaceAllStringFunc(content, func(match string) string {
		timeStr := match[4 : len(match)-1]
		t, err := time.Parse("2006-01-02T15:04Z", timeStr)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04-07:00", timeStr)
			if err != nil {
				return match
			}
		}
		zone := t.Format("-07:00")
		if zone == "Z" {
			zone = "(+00:00)"
		} else {
			zone = fmt.Sprintf("(%s)", zone)
		}
		return fmt.Sprintf("%s %s", t.Format("03:04PM"), zone)
	})

	// 24-hour time: T24(...)
	time24Regex := regexp.MustCompile(`T24\(([0-9T:.Z+-]{16,})\)`)
	content = time24Regex.ReplaceAllStringFunc(content, func(match string) string {
		timeStr := match[4 : len(match)-1]
		t, err := time.Parse("2006-01-02T15:04Z", timeStr)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04-07:00", timeStr)
			if err != nil {
				return match
			}
		}
		zone := t.Format("-07:00")
		if zone == "Z" {
			zone = "(+00:00)"
		} else {
			zone = fmt.Sprintf("(%s)", zone)
		}
		return fmt.Sprintf("%s %s", t.Format("15:04"), zone)
	})

	return content
}


// Highlight (Colorized) Processing Functions
// Used for printing to terminal with ANSI colors.

func highlightProcessContent(content string) string {
	content = processAirportCodes(content)
	content = processDatesAndTimes(content)
	content = trimHorizontalWhitespace(content)
	content = trimVerticalWhitespace(content)
	return content
}

// processAirportCodes replaces airport codes with highlighted airport names or cities.
// With "*" prefix it outputs the municipality.
func processAirportCodes(content string) string {
	// IATA codes: supports *#ABC
	iataRegex := regexp.MustCompile(`(\*?)#([A-Z]{3})`)
	content = iataRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := iataRegex.FindStringSubmatch(match)
		code := groups[2]
		if airport, exists := airportMap[code]; exists {
			if groups[1] == "*" {
				return highlightCity(airport)
			}
			return highlightAirport(airport)
		}
		return match
	})

	// ICAO codes: supports *##ABCD
	icaoRegex := regexp.MustCompile(`(\*?)##([A-Z]{4})`)
	content = icaoRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := icaoRegex.FindStringSubmatch(match)
		code := groups[2]
		if airport, exists := airportMap[code]; exists {
			if groups[1] == "*" {
				return highlightCity(airport)
			}
			return highlightAirport(airport)
		}
		return match
	})
	return content
}

// highlightAirport returns the airport name highlighted in green.
func highlightAirport(airport *Airport) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, airport.Name, ColorReset)
}

// highlightCity returns the municipality (city) highlighted in cyan.
func highlightCity(airport *Airport) string {
	if strings.TrimSpace(airport.Municipality) != "" {
		return fmt.Sprintf("%s%s%s", ColorCyan, airport.Municipality, ColorReset)
	}
	return fmt.Sprintf("%s%s%s", ColorGreen, airport.Name, ColorReset)
}

// processDatesAndTimes replaces date/time placeholders with highlighted formatted dates/times.
func processDatesAndTimes(content string) string {
	// Dates: D(...)
	dateRegex := regexp.MustCompile(`D\(([0-9T:.Z+-]{16,})\)`)
	content = dateRegex.ReplaceAllStringFunc(content, func(match string) string {
		dateStr := match[2 : len(match)-1]
		t, err := time.Parse("2006-01-02T15:04Z", dateStr)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04-07:00", dateStr)
			if err != nil {
				return match
			}
		}
		return fmt.Sprintf("%s%s%s", ColorMagenta, t.Format("02 Jan 2006"), ColorReset)
	})

	// 12-hour time: T12(...)
	time12Regex := regexp.MustCompile(`T12\(([0-9T:.Z+-]{16,})\)`)
	content = time12Regex.ReplaceAllStringFunc(content, func(match string) string {
		timeStr := match[4 : len(match)-1]
		t, err := time.Parse("2006-01-02T15:04Z", timeStr)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04-07:00", timeStr)
			if err != nil {
				return match
			}
		}
		zone := t.Format("-07:00")
		if zone == "Z" {
			zone = "(+00:00)"
		} else {
			zone = fmt.Sprintf("(%s)", zone)
		}
		return fmt.Sprintf("%s%s%s %s%s%s", ColorCyan, t.Format("03:04PM"), ColorReset, ColorYellow, zone, ColorReset)
	})

	// 24-hour time: T24(...)
	time24Regex := regexp.MustCompile(`T24\(([0-9T:.Z+-]{16,})\)`)
	content = time24Regex.ReplaceAllStringFunc(content, func(match string) string {
		timeStr := match[4 : len(match)-1]
		t, err := time.Parse("2006-01-02T15:04Z", timeStr)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04-07:00", timeStr)
			if err != nil {
				return match
			}
		}
		zone := t.Format("-07:00")
		if zone == "Z" {
			zone = "(+00:00)"
		} else {
			zone = fmt.Sprintf("(%s)", zone)
		}
		return fmt.Sprintf("%s%s%s %s%s%s", ColorCyan, t.Format("15:04"), ColorReset, ColorYellow, zone, ColorReset)
	})

	return content
}

// trimHorizontalWhitespace removes excessive horizontal whitespace.
func trimHorizontalWhitespace(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		words := strings.Fields(line)
		lines[i] = strings.Join(words, " ")
	}
	return strings.Join(lines, "\n")
}

// trimVerticalWhitespace removes excessive vertical whitespace.
func trimVerticalWhitespace(content string) string {
	content = regexp.MustCompile(`\\[rvf]`).ReplaceAllString(content, "\n")
	content = regexp.MustCompile("[\\r\\v\\f]+").ReplaceAllString(content, "\n")
	content = regexp.MustCompile("\n{3,}").ReplaceAllString(content, "\n\n")
	return content
}

// printError prints an error message in red and bold.
func printError(message string) {
	fmt.Printf("%s%sError: %s%s\n", ColorRed, Bold, message, ColorReset)
}

// printSuccess prints a success message in green and bold.
func printSuccess(message string) {
	fmt.Printf("%s%sSuccess: %s%s\n", ColorGreen, Bold, message, ColorReset)
}
