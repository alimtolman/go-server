package File

import "os"

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CurrentDirectory() string {
	dir, _ := os.Getwd()
	return dir
}

func CreateDirectory(path string) {
	os.MkdirAll(path, 0755)
}

func ReadBytes(path string) []byte {
	data, _ := os.ReadFile(path)
	return data
}

func WriteBytes(path string, data []byte) {
	os.WriteFile(path, data, 0644)
}

func WriteText(path string, data string) {
	WriteBytes(path, []byte(data))
}
