package persist

import (
	"encoding/json"
	"os"
)

func Save[T any](filepath string, data T) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func Load[T any](filepath string) (data T, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return
}
