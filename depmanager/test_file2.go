package depmanager

import (
	"fmt"
	"net/http"
	"time"
)

func fooo() {
	fmt.Println("wohoo")
	fmt.Println(time.Now())
	fmt.Print(http.StateClosed)
}
