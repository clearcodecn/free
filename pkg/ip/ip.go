package ip

import (
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"strings"
)

var (
	region *ip2region.Ip2Region
)

func InitIPDB(db string) {
	var err error
	region, err = ip2region.New(db)
	if err != nil {
		panic(err)
	}
}

func CloseIpDB() {
	region.Close()
}

func GetAddress(ip string) string {
	info, err := region.BtreeSearch(ip)
	if err != nil {
		return ""
	}
	var s []string

	if info.Country != "" {
		s = append(s, info.Country)
	}
	if info.Province != "" {
		s = append(s, info.Province)
	}
	if info.City != "" {
		s = append(s, info.City)
	}

	return strings.Join(s, "-")
}
