# CSV → JSON Lines Converter (Go)

This command-line program reads a csv and converts each row into a JSON Lines (.jl) record while maintaining column order.
Each output line is a valid JSON object with field names taken from the CSV header.

Example output:

{"value":452600,"income":8.3252,"age":41,"rooms":880,"bedrooms":129,"pop":322,"hh":126}
{"value":358500,"income":8.3014,"age":21,"rooms":7099,"bedrooms":1106,"pop":2401,"hh":1138}
{"value":352100,"income":7.2574,"age":52,"rooms":1467,"bedrooms":190,"pop":496,"hh":177}

### Requirements

* Go 1.18+

### How to Build and Run

1. Clone the repo
2. Initialize a Go module: ```go mod init csvToJsonCLI```
4. Build it: ```go build -o main .```
5. Run it: ```./main <inputfilename.csv> <outputfilename.jl>```

### Command Line Usage
Usage: ```./main <input.csv> <output.jl>```

Example: ```./main housesInput.csv housesOutput.jl```

In this project I used AI to help create verbose and frequent error checks since I'm still getting used to go style error checking and norms.
