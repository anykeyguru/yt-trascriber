package export

import "os"

func Save(path, text string) error {
	return os.WriteFile(path, []byte(text), 0644)
}
