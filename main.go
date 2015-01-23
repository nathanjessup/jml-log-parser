package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//I'm lazy and don't want to figure out why the gps puts an extra line return in, lets filter it out.
		if len(line) == 0 {
			continue
		}
		//gps sentences in the log say why type of sentence they are
		//		time     fix  lat         long         knots course fix	date  checksum
		//$GPRMC,194509.000,A,4042.6142,N,07400.4168,W,2.03,221.11,160412,,,A*77 
		//  0 		1 		2     3  	4      5	 6   7 	   8  	 9       indexs in the split

		if strings.Contains(line, "$GPRMC") {
			data := strings.Split(line, ",")
			latDegMin := getFloatFromString(data[3])
			longDegMin := getFloatFromString(data[5])

			//S and W neg; N and E are pos
			lat := convertDMSToDD(latDegMin)
			if data[4] == string('S') {
				lat = lat * -1 //flip the sign
			}

			long := convertDMSToDD(longDegMin)
			if data[6] == string('W') {
				long = long * -1 //flip the sign
			}

			fmt.Println(fmt.Sprintf("%f", long) + "," + fmt.Sprintf("%f", lat)) //the kml format wants cords described as long,lat,alt
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func getFloatFromString(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func convertDMSToDD(degMin float64) float64 {
	deg := math.Floor(degMin / 100)
	mins := degMin - deg*100

	return deg + mins/60
}
