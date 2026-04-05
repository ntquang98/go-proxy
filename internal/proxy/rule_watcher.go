package proxy

import (
	"log"
	"log/slog"
	"time"

	"github.com/fsnotify/fsnotify"
)

func WatchRules(path string, reload func() error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	watcher.Add(path)

	var debounceTimer *time.Timer
	const debounceDuration = 200 * time.Millisecond

	for {
		select {
		case event := <-watcher.Events:
			slog.Debug("fs event", "op", event.Op.String())
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				if debounceTimer != nil {
					debounceTimer.Stop()
				}

				debounceTimer = time.AfterFunc(debounceDuration, func() {
					slog.Info("rules changed, reloading...")

					if err := reload(); err != nil {
						slog.Error("reload failed", "error", err)
					}
				})

			}
		case err := <-watcher.Errors:
			slog.Error("watch error", "error", err)
		}
	}

}
