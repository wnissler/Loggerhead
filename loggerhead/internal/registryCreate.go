package loggerhead

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/bi-zone/etw"
	"golang.org/x/sys/windows"
)

func RegETW() {
	guid, _ := windows.GUIDFromString("{70EB4F03-C1DE-4F73-A051-33D13D5413BD}")
	session, err := etw.NewSession(guid)
	if err != nil {
		log.Fatalf("Failed to create ETW session: %s", err)
	}

	cb := func(e *etw.Event) {
		if data, err := e.EventProperties(); err == nil {
			if err != nil {
				log.Fatalf(err.Error())
			}
			log.Printf("%s", data["RelativeName"])
		}
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := session.Process(cb); err != nil {
			log.Printf("[ERR] Got error processing events: %s", err)
		}
		wg.Done()
	}()

	// Trap cancellation.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	if err := session.Close(); err != nil {
		log.Printf("[ERR] Got error closing the session: %s", err)
	}
	wg.Wait()
}
