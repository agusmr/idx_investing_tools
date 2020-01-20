package tools

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(filepath string, url string) error {
	info, err := os.Stat(filepath)
	if os.IsExist(err) {
		return errors.New(fmt.Sprintf("%v file exists", info))
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return err
}
