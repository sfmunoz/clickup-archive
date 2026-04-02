package archive

import (
	"encoding/json"
	"os"
)

func jsonSave(v any, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(indexFile(dir), data, 0o644)
}
