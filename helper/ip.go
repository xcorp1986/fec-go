package helper

import (
    "github.com/oschwald/geoip2-golang"
    "log"
    "net"
    "errors"
)


func GetIpInfo(ipStr string) (string, string, string, string, error){
    if ipStr == "" {
        return "", "", "", "", errors.New("ip str is empty")
    }
    ipDb, err := geoip2.Open("/www/test/ip/GeoLite2-City_20180327/GeoLite2-City.mmdb")
    if err != nil {
            log.Fatal(err)
    }
    defer ipDb.Close()
    // If you are using strings that may be invalid, check that ip is not nil
    ip := net.ParseIP(ipStr)
    record, err := ipDb.City(ip)
    if err != nil {
        log.Fatal(err)
        return "", "", "", "", err
    }
    log.Println("Ip:"+ipStr)
    iosCountryCode := record.Country.IsoCode
    log.Println("iosCode:"+iosCountryCode)
    countryName := record.Country.Names["en"]
    log.Println("countryName:"+countryName)
    log.Println(record.Subdivisions)
    stateName := ""
    if len(record.Subdivisions) > 0 {
        s := record.Subdivisions[0]
        stateName = s.Names["en"]
    }
    log.Println("stateName:"+stateName)
    cityName := record.City.Names["en"]
    log.Println("cityName:"+cityName)
    return iosCountryCode, countryName, stateName, cityName, nil
    /*
    fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["en-US"])
    fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
    fmt.Printf("Russian country name: %v\n", record.Country.Names["en"])
    fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
    fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
    fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
    */
    // Output:
    // Portuguese (BR) city name: Londres
    // English subdivision name: England
    // Russian country name: Великобритания
    // ISO country code: GB
    // Time zone: Europe/London
    // Coordinates: 51.5142, -0.0931
}