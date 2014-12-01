package handlebars

import "path/filepath"
import "encoding/hex"
import "crypto/md5"

func resolve(from string, to string) string {
	if filepath.IsAbs(to) {
		return to
	}

	if filepath.IsAbs(from) {
		return filepath.Join(from, to)
	}

	abs, err := filepath.Abs(from)

	if err != nil {
		abs = "/"
	}

	return filepath.Join(abs, to)
}

func hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func panicError(e error) {
	if e != nil {
		panic(e)
	}
}
