# CodeQuest CLI

A command-line tool for practicing coding challenges locally. Write solutions in Go, JavaScript, TypeScript, or Python and test them against the same test cases used in the online platform.

## Features

- **Multi-language support**: Go, JavaScript, TypeScript, and Python
- **Local testing**: Run tests without Docker - just needs Go, Node.js, and/or Python
- **Offline practice**: Download challenges and work on them locally
- **Fast execution**: Native runtime execution for quick feedback
- **Same test cases**: Uses identical test cases as the web platform

## Installation

### Quick install (recommended)

**One-line installer (Linux/macOS):**
```bash
curl -fsSL https://raw.githubusercontent.com/crisecheverria/codequest/main/install.sh | bash
```

### Manual install from releases

Download the latest binary from the [releases page](https://github.com/crisecheverria/codequest/releases):

**For macOS (Intel):**
```bash
curl -L https://github.com/crisecheverria/codequest/releases/latest/download/codequest-darwin-amd64 -o codequest
chmod +x codequest
sudo mv codequest /usr/local/bin/
```

**For macOS (M1/M2):**
```bash
curl -L https://github.com/crisecheverria/codequest/releases/latest/download/codequest-darwin-arm64 -o codequest
chmod +x codequest
sudo mv codequest /usr/local/bin/
```

**For Linux (AMD64):**
```bash
curl -L https://github.com/crisecheverria/codequest/releases/latest/download/codequest-linux-amd64 -o codequest
chmod +x codequest
sudo mv codequest /usr/local/bin/
```

**For Linux (ARM64):**
```bash
curl -L https://github.com/crisecheverria/codequest/releases/latest/download/codequest-linux-arm64 -o codequest
chmod +x codequest
sudo mv codequest /usr/local/bin/
```

**For Windows:**
Download [codequest-windows-amd64.exe](https://github.com/crisecheverria/codequest/releases/latest/download/codequest-windows-amd64.exe) and add it to your PATH.

### Install from source

```bash
go install github.com/crisecheverria/codequest@latest
```

## Prerequisites

- **Go 1.23+** (for Go challenges)
- **Node.js 18+** (for JavaScript/TypeScript challenges)
- **Python 3.8+** (for Python challenges)

## Usage

### List available challenges

```bash
codequest list
codequest list --language python --difficulty easy
codequest list --language go --difficulty medium
```

### Download a challenge

```bash
codequest fetch challenge-slug
cd challenge-slug
```

### Test your solution

```bash
codequest test
```

### Example workflow

```bash
# List challenges
codequest list --language python

# Download a challenge
codequest fetch python-list-comprehension

# Move to the challenge directory
cd challenge-python-list-comprehension

# Edit the solution file
vim solution.py

# Test your solution
codequest test
```

## Supported Languages

- **Go**: Full support with native `go run` execution
- **JavaScript**: Full support with Node.js execution
- **TypeScript**: Basic support with type annotation removal
- **Python**: Full support with native `python3` execution

## Challenge Structure

Each challenge creates a workspace with:
- `solution.go` / `solution.js` / `solution.py` - Your solution file
- `README.md` - Challenge description and examples
- `.challenge.json` - Challenge metadata (don't modify)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details
