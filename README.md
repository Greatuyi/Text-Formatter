# Text Formatter

A powerful command-line tool written in Go that transforms raw itinerary text into beautifully formatted, human-readable documents. The formatter intelligently processes airport codes, dates, and times while cleaning up excessive whitespace.

![Go Version](https://img.shields.io/badge/Go-1.23.2-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

## ✨ Features

### 🛫 Airport Code Resolution
- **IATA Code Support**: Convert 3-letter IATA codes (e.g., `#JFK`) to full airport names
- **ICAO Code Support**: Convert 4-letter ICAO codes (e.g., `##EGLL`) to full airport names
- **City Name Extraction**: Use `*` prefix (e.g., `*#CDG`) to display city/municipality instead of airport name
- **Comprehensive Database**: Supports thousands of airports worldwide via CSV lookup

### 📅 Date & Time Formatting
- **Date Formatting**: `D(2025-03-15T14:30-04:00)` → `15 Mar 2025`
- **12-Hour Time**: `T12(2025-03-15T14:30-04:00)` → `02:30PM (-04:00)`
- **24-Hour Time**: `T24(2025-03-16T06:30+00:00)` → `06:30 (+00:00)`
- **Timezone Support**: Handles UTC (Z) and offset-based timezones

### 🧹 Whitespace Cleanup
- **Horizontal Trimming**: Removes excessive spaces between words
- **Vertical Trimming**: Reduces multiple blank lines to maximum of two
- **Escape Sequence Handling**: Converts `\r`, `\v`, `\f` to proper newlines

### 🎨 Dual Output Modes
- **Plain Text**: Clean output written to file (no formatting codes)
- **Colorized Terminal**: Beautiful ANSI-colored output displayed in console
  - 🟢 Green: Airport names
  - 🔵 Cyan: City names and times
  - 🟣 Magenta: Dates
  - 🟡 Yellow: Timezones

## 📋 Prerequisites

- **Go**: Version 1.23.2 or higher
- **Airport Database**: CSV file containing airport information

## 🚀 Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Greatuyi/Text-Formatter.git
   cd Text-Formatter
   ```

2. **Verify Go installation**:
   ```bash
   go version
   ```

3. **Ensure you have the airport lookup CSV file** (see [Airport Database Format](#airport-database-format))

## 💻 Usage

### Basic Command

```bash
go run . <input-file> <output-file> <airport-lookup-csv>
```

### Example

```bash
go run . ./input.txt ./output.txt ./airport-lookup.csv
```

### Help

```bash
go run . -h
```

## 📝 Input Syntax

### Airport Codes

| Syntax | Description | Example Input | Example Output |
|--------|-------------|---------------|----------------|
| `#ABC` | IATA code (3 letters) | `#JFK` | John F Kennedy International Airport |
| `##ABCD` | ICAO code (4 letters) | `##EGLL` | London Heathrow Airport |
| `*#ABC` | IATA code → City | `*#CDG` | Paris |
| `*##ABCD` | ICAO code → City | `*##EDDW` | Bremen |

### Date & Time Placeholders

| Syntax | Format | Example Input | Example Output |
|--------|--------|---------------|----------------|
| `D(...)` | Date | `D(2025-03-15T14:30-04:00)` | 15 Mar 2025 |
| `T12(...)` | 12-hour time | `T12(2025-03-15T14:30-04:00)` | 02:30PM (-04:00) |
| `T24(...)` | 24-hour time | `T24(2025-03-16T06:30+00:00)` | 06:30 (+00:00) |

**Supported DateTime Formats**:
- `2006-01-02T15:04Z` (UTC)
- `2006-01-02T15:04-07:00` (with timezone offset)

### Sample Input

```text
Flight Itinerary
Passenger: Marie Charlie

Flight Details:
Flight Number: AB1234
Departure: #JFK
Arrival: ##EGLL
Date: D(2025-03-15T14:30-04:00)
Departure Time: T12(2025-03-15T14:30-04:00)
Arrival Time: T24(2025-03-16T06:30+00:00)

Layover: Location: *#CDG Duration: 2 hours
```

### Sample Output

```text
Flight Itinerary
Passenger: Marie Charlie

Flight Details:
Flight Number: AB1234
Departure: John F Kennedy International Airport
Arrival: London Heathrow Airport
Date: 15 Mar 2025
Departure Time: 02:30PM (-04:00)
Arrival Time: 06:30 (+00:00)

Layover: Location: Paris Duration: 2 hours
```

## 🗂️ Airport Database Format

The airport lookup CSV file must contain the following columns (order-independent):

| Column Name | Description | Example |
|-------------|-------------|---------|
| `name` | Full airport name | John F Kennedy International Airport |
| `iso_country` | ISO country code | US |
| `municipality` | City name | New York |
| `icao_code` | 4-letter ICAO code | KJFK |
| `iata_code` | 3-letter IATA code | JFK |
| `coordinates` | Geographic coordinates | -73.7781, 40.6413 |

**Requirements**:
- Header row must be present
- Column names are case-insensitive
- Each record must have either an IATA or ICAO code (or both)
- Empty names are not allowed

**Sample CSV**:
```csv
name,iso_country,municipality,icao_code,iata_code,coordinates
John F Kennedy International Airport,US,New York,KJFK,JFK,"-73.7781, 40.6413"
London Heathrow Airport,GB,London,EGLL,LHR,"-0.461941, 51.4706"
Charles de Gaulle Airport,FR,Paris,LFPG,CDG,"2.55, 49.0097"
```

## 🏗️ Project Structure

```
Text-Formatter/
├── main.go                 # Main application logic
├── go.mod                  # Go module definition
├── airport-lookup.csv      # Airport database
├── input.txt              # Sample input file
├── output.txt             # Generated output file
└── README.md              # This file
```

## 🔧 How It Works

1. **Argument Parsing**: Validates command-line arguments (input, output, airport CSV)
2. **Airport Database Loading**: Parses CSV and builds an in-memory lookup map
3. **Content Processing**:
   - Replaces airport codes with full names/cities
   - Formats dates and times
   - Cleans up whitespace
4. **Dual Output Generation**:
   - Plain text → Written to output file
   - ANSI-colored text → Displayed in terminal

## 🎯 Use Cases

- ✈️ **Travel Itinerary Formatting**: Convert raw booking data into readable itineraries
- 📧 **Email Templates**: Generate professional travel confirmations
- 📊 **Report Generation**: Create formatted flight schedules
- 🤖 **Data Processing**: Batch process travel documents

## ⚠️ Error Handling

The application provides clear error messages for:

| Error | Description |
|-------|-------------|
| `Input file not found` | The specified input file doesn't exist |
| `Airport lookup file not found` | The CSV database file is missing |
| `Airport lookup file is malformed` | CSV format is invalid or missing required columns |
| `Error reading input file` | Permission or I/O issues with input file |
| `Error writing output file` | Permission or I/O issues with output file |

## 🧪 Testing

Create a test input file with various airport codes and date formats:

```bash
# Create test input
echo "Flight from #JFK to *#LHR on D(2025-12-25T10:00Z)" > test-input.txt

# Run formatter
go run . test-input.txt test-output.txt airport-lookup.csv

# View results
cat test-output.txt
```

## 🚀 Building for Production

### Compile Binary

```bash
# For current platform
go build -o text-formatter .

# For Windows
GOOS=windows GOARCH=amd64 go build -o text-formatter.exe .

# For Linux
GOOS=linux GOARCH=amd64 go build -o text-formatter .

# For macOS
GOOS=darwin GOARCH=amd64 go build -o text-formatter .
```

### Run Compiled Binary

```bash
./text-formatter input.txt output.txt airport-lookup.csv
```

## 🤝 Contributing

Contributions are welcome! Here are some ways you can contribute:

1. 🐛 Report bugs
2. 💡 Suggest new features
3. 📝 Improve documentation
4. 🔧 Submit pull requests

### Development Guidelines

- Follow Go best practices and conventions
- Add comments for complex logic
- Test with various input formats
- Ensure backward compatibility

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

**Greatuyi**

## 🙏 Acknowledgments

- Airport database sourced from [OurAirports](https://ourairports.com/data/)
- Built with ❤️ using Go

## 📞 Support

If you encounter any issues or have questions:

1. Check the [error handling section](#️-error-handling)
2. Verify your input file syntax
3. Ensure the airport CSV is properly formatted
4. Open an issue on GitHub

---

**Made with ✈️ for travelers and developers**
