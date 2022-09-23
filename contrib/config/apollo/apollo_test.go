package apollo

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
	"time"
)

func TestApollo(t *testing.T) {
	NewApollo(appId, cluster, ip).Run()

	for range time.Tick(time.Second) {
		fmt.Println(viper.AllSettings())
	}
}
