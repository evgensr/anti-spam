package app

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

func (a *App) loadStopWords() error {

	errorResponse := a.download(
		a.cfg.App.GoogleCSV,
		a.cfg.App.LocalNameCSV,
		30000,
	)
	if errorResponse != nil {
		return fmt.Errorf("error download file: %w", errorResponse)
	}

	filter, err := a.readFromCSV()
	if err != nil {
		a.l.Error("error read from csv", zap.Error(err))
	}

	if filter != nil {
		a.filter.Mu.Lock()
		a.filter.Word = filter
		a.filter.Mu.Unlock()
	}
	return nil
}

func (a *App) download(url string, filename string, timeout int64) error {
	client := http.Client{
		Timeout: time.Duration(timeout * int64(time.Second)),
	}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("cannot download file: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("cannot download file: %s", resp.Status)
	}

	if len(resp.Header["Content-Type"]) == 0 {
		return fmt.Errorf("cannot download file: no Content-Type header")
	}

	if resp.Header["Content-Type"][0] != "text/csv" {
		return fmt.Errorf("content type not valid: %s", resp.Header["Content-Type"][0])
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read body: %w", err)
	}

	err = os.WriteFile(filename, b, 0644)
	if err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	return nil

}

func (a *App) readFromCSV() ([]string, error) {

	f, err := os.Open(a.cfg.App.LocalNameCSV)
	defer func() {
		if err = f.Close(); err != nil {
			a.l.Error("error close file", zap.Error(err))
		}
	}()

	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	var counter int
	var line []string
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		counter++
		if err == io.EOF {
			a.l.Info("load csv", zap.Int("count", counter-1))
			break
		}
		if counter == 1 {
			continue
		}

		if err != nil {
			a.l.Logger.Error("error read csv", zap.Error(err))
		}
		line = append(line, rec[0])
	}

	return line, nil
}

func (a *App) refresh() {
	for {
		err := a.loadStopWords()
		if err != nil {
			a.l.Error("error load stop words", zap.Error(err))
		}
		time.Sleep(time.Duration(a.cfg.App.RefreshIntervalSec) * time.Second)
	}
}
