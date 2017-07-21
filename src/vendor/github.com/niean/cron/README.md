# cron
crontab for golang

## peference
```bash
TODO	
```

## usage
```go
import (
	"fmt"
	"time"
	
	ncron "github.com/niean/cron"
)

func main() {
	// init cron
	c := ncron.New()

	// add cron job
	c.AddFunc("* * * * * *", func() { fmt.Println("Every second") })
	c.AddFuncCC("* * * * * *", func() { fmt.Println("Every second, with max Concurrrent 2"); time.Sleep(10 * time.Second)}, 2)

	// start cron
	c.Start()

	// keep alive
	select {}
}
```

## reference

