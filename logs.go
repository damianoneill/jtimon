package main

import (
	"fmt"
	"log"
	"os"
)

// Log it, each go routine collecting JTI data has its own log file so its concurrency safe.
// In case of print, caller has to guarantee safety.
func l(safe bool, jctx *JCtx, s string) {
	if *print {
		switch safe {
		case true:
			gmutex.Lock()
			fmt.Printf(s)
			gmutex.Unlock()

		case false:
			fmt.Printf(s)
		}
	} else if jctx.cfg.Log.loger != nil {
		jctx.cfg.Log.loger.Printf(s)
	}
}

func logInit(jctx *JCtx) {
	file := jctx.cfg.Log.File
	if file != "" {
		if *print {
			fmt.Println("Both print and log options are specified, ignoring log")
		} else {
			f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Printf("Could not create log file(%s): %v\n", file, err)
			} else {
				flags := 0
				if !jctx.cfg.Log.CSVStats {
					flags = log.LstdFlags
				}
				jctx.cfg.Log.loger = log.New(f, "", flags)
				jctx.cfg.Log.handle = f
				fmt.Printf("logging in %s for %s:%d [periodic stats every %d seconds]\n",
					jctx.cfg.Log.File, jctx.cfg.Host, jctx.cfg.Port, jctx.cfg.Log.PeriodicStats)
			}
		}
	}
}
