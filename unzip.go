package main

import (
	"fmt"
	"github.com/takoyaki-3/goc"
)

func main(){
	fmt.Println("start")

	paths,_ := goc.Dirwalk("./docomo-cycle-tokyo_station_status")

	for _,v:=range paths{
		fmt.Println(v)

		goc.Unzip(v,"./gbfs-stationstatus")
	}
}
