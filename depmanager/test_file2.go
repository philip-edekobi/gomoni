package depmanager

import (
	"fmt"
	n "net"
	"net/http"
	"time"
)

func fooo() {
	fmt.Println("wohoo")
	fmt.Println(time.Now())
	fmt.Print(http.StateClosed)
	fmt.Println(n.IPv4len)
}
