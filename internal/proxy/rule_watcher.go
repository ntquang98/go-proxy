package proxy

import (
	"log"
	"log/slog"

	"github.com/fsnotify/fsnotify"
)

func WatchRules(path string, reload func() error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	watcher.Add(path)

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				slog.Info("rules changed, reloading...")

				if err := reload(); err != nil {
					slog.Error("reload failed", "error", err)
				}
			}
		case err := <-watcher.Errors:
			slog.Error("watch error", "error", err)
		}
	}

}
