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
	if info != nil {
		return errors.New(fmt.Sprintf("%s already exists", info.Name()))
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
