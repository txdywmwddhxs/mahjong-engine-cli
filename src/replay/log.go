package replay

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const LogStartPrefix = "REPLAY START "
const LogIncompletePrefix = "REPLAY INCOMPLETE "
const LogBlockPrefix = "-----BEGIN MAHJONG REPLAY V1 "
const LogBlockEnd = "-----END MAHJONG REPLAY-----"

const (
	logBlockLineWidth         = 76
	maxLogIDLength            = 128
	initialLogScanBufferBytes = 64 * 1024
	// Bound base64 text independently of the decompressed JSON limit.
	maxEncodedLogBytes = MaxFileBytes * 2
)

// LogBlock keeps complete, verified snapshots compact enough to copy from a log.
// Human-readable game messages remain outside this gzip/base64 transport block.
func (f *File) LogBlock() (string, error) {
	if err := f.Validate(); err != nil {
		return "", err
	}
	if !validLogID(f.GameID) {
		return "", fmt.Errorf("invalid game ID for log transport")
	}
	data, err := json.Marshal(f)
	if err != nil {
		return "", err
	}
	if len(data) > MaxFileBytes {
		return "", fmt.Errorf("replay exceeds 32 MiB")
	}
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err = writer.Write(data); err != nil {
		return "", err
	}
	if err = writer.Close(); err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(compressed.Bytes())
	var block strings.Builder
	fmt.Fprintf(&block, "%s%s-----\n", LogBlockPrefix, f.GameID)
	for len(encoded) > logBlockLineWidth {
		block.WriteString(encoded[:logBlockLineWidth])
		block.WriteByte('\n')
		encoded = encoded[logBlockLineWidth:]
	}
	block.WriteString(encoded)
	block.WriteByte('\n')
	block.WriteString(LogBlockEnd)
	block.WriteByte('\n')
	return block.String(), nil
}

func validLogID(id string) bool {
	if len(id) == 0 || len(id) > maxLogIDLength {
		return false
	}
	for _, c := range id {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_' || c == '.') {
			return false
		}
	}
	return true
}

// ReadSource accepts native JSON, a copied log block or a complete text log.
// With no game ID it selects the latest recorded round, including its failure
// status: an unfinished latest round never silently falls back to an older one.
func ReadSource(in io.Reader, gameID string) (*File, error) {
	reader := bufio.NewReader(in)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("empty replay input")
			}
			return nil, err
		}
		if b == ' ' || b == '\t' || b == '\r' || b == '\n' {
			continue
		}
		if err = reader.UnreadByte(); err != nil {
			return nil, err
		}
		if b == '{' {
			file, err := decodeJSON(reader)
			if err != nil {
				return nil, err
			}
			if gameID != "" && file.GameID != gameID {
				return nil, fmt.Errorf("game %s not found", gameID)
			}
			return file, nil
		}
		break
	}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, initialLogScanBufferBytes), maxEncodedLogBytes)
	var payload strings.Builder
	var selected, selectedID, activeID string
	var selectedErr error
	selectedRound := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, LogStartPrefix) || strings.HasPrefix(line, LogIncompletePrefix) {
			prefix := LogStartPrefix
			if strings.HasPrefix(line, LogIncompletePrefix) {
				prefix = LogIncompletePrefix
			}
			fields := strings.Fields(strings.TrimPrefix(line, prefix))
			if len(fields) == 0 {
				return nil, fmt.Errorf("malformed replay status line")
			}
			id := fields[0]
			if gameID == "" || id == gameID {
				selectedRound = true
				selected = ""
				selectedID = id
				selectedErr = fmt.Errorf("game %s has an incomplete recording; select an earlier completed game with --game ID", id)
			}
			continue
		}
		if strings.HasPrefix(line, "-----BEGIN MAHJONG REPLAY ") {
			if !strings.HasPrefix(line, LogBlockPrefix) || !strings.HasSuffix(line, "-----") {
				return nil, fmt.Errorf("unsupported or malformed replay block header")
			}
			id := strings.TrimSuffix(strings.TrimPrefix(line, LogBlockPrefix), "-----")
			if !validLogID(id) {
				return nil, fmt.Errorf("invalid replay block game ID")
			}
			activeID = id
			payload.Reset()
			if gameID == "" || gameID == id {
				selectedRound = true
				selected = ""
				selectedID = id
				selectedErr = fmt.Errorf("game %s replay block is truncated (missing end marker)", id)
			}
			continue
		}
		if line == LogBlockEnd {
			if activeID == "" {
				return nil, fmt.Errorf("replay end marker without begin marker")
			}
			if gameID == "" || activeID == gameID {
				selected = payload.String()
				selectedID = activeID
				selectedErr = nil
			}
			activeID = ""
			payload.Reset()
			continue
		}
		if activeID != "" {
			if payload.Len()+len(line) > maxEncodedLogBytes {
				return nil, fmt.Errorf("encoded replay block exceeds size limit")
			}
			payload.WriteString(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read replay log: %w", err)
	}
	if selectedErr != nil {
		return nil, selectedErr
	}
	if !selectedRound {
		if gameID != "" {
			return nil, fmt.Errorf("game %s not found in log", gameID)
		}
		return nil, fmt.Errorf("no replay block found; old text-only logs cannot be replayed")
	}
	decoded, err := base64.StdEncoding.Strict().DecodeString(selected)
	if err != nil {
		return nil, fmt.Errorf("corrupt replay block for %s: %w", selectedID, err)
	}
	zipped, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return nil, fmt.Errorf("corrupt replay block for %s: %w", selectedID, err)
	}
	defer zipped.Close()
	file, err := decodeJSON(zipped)
	if err != nil {
		return nil, fmt.Errorf("replay block for %s: %w", selectedID, err)
	}
	if file.GameID != selectedID {
		return nil, fmt.Errorf("replay block game ID mismatch")
	}
	return file, nil
}

func decodeJSON(in io.Reader) (*File, error) {
	data, err := io.ReadAll(io.LimitReader(in, MaxFileBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read replay: %w", err)
	}
	if len(data) > MaxFileBytes {
		return nil, fmt.Errorf("replay exceeds 32 MiB")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var file File
	if err = decoder.Decode(&file); err != nil {
		return nil, fmt.Errorf("read replay: %w", err)
	}
	var trailing any
	if err = decoder.Decode(&trailing); err != io.EOF {
		return nil, fmt.Errorf("unexpected trailing replay data")
	}
	if err = file.Validate(); err != nil {
		return nil, err
	}
	return &file, nil
}
