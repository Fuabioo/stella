package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

var logFile *os.File

const programName = "stella"

func main() {
	defer logFile.Close()

	if len(os.Args) < 2 {
		log.Fatal("URI path required")
	}
	uri := os.Args[1]
	step := time.Millisecond
	backoff := time.Second
	fullscreen := false
	debug := false
	failureThreshold := uint(5)

	var err error

	logFile, err = GetEphemeralLogFile(programName)
	if err != nil {
		log.Fatal(err)
	}

	// Set output of standard logger to the file
	log.SetOutput(logFile)

	level := log.InfoLevel
	if debug {
		level = log.DebugLevel
	}
	log.SetLevel(level)

	data := &state{
		failureThreshold: failureThreshold,
	}
	m := newModel(data)

	log.Debug("Initializing progress listener",
		"uri", uri,
		"step", step,
	)
	client := NewClient(uri, step, backoff, failureThreshold, false)
	client.Observe(func(progress *ProgressReport, err error) {
		if err != nil {
			log.Error("Failed to fetch progress",
				"backoff", backoff,
				"retries", fmt.Sprintf("%d/%d", data.GetFailures(), failureThreshold),
				"err", err,
			)
			data.IncFailures()
		}
		data.Set(progress)
	})

	opts := []tea.ProgramOption{}
	if fullscreen {
		opts = append(opts, tea.WithAltScreen())
	}
	if _, err := tea.NewProgram(m, opts...).Run(); err != nil {
		log.Error("Tea failed", "err", err)
	}
}

func GetEphemeralLogFile(programName string) (*os.File, error) {
	logDir := filepath.Join(os.TempDir(), programName)
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		return nil, fmt.Errorf("Failed to get current user: %v", err)
	}
	today := time.Now().Local().Format(time.DateOnly)
	path := filepath.Join(logDir, today+".log")

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("Failed to open log file: %v", err)
	}

	return file, nil
}
