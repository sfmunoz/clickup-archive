package fetch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func jsonDump(v any, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	for line := range strings.SplitSeq(string(data), "\n") {
		log.Debug(line)
	}
	return os.WriteFile(filepath.Join(dir, "index.json"), data, 0o644)
}
