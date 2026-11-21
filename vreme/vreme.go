// Package vreme provides functions for working with weather stations
// and their temperature measurements.
//
// Stations are stored in a map that associates a station identifier
// (for example "LJ01") with a [Postaja]. Because the fields of [Postaja]
// are not exported, new stations are created with [DodajPostajo] and
// measurements are added with [DodajMeritev]. A summary of all stations
// is printed with [IzpisPorocila].
package vreme

import (
	"fmt"
	"sort"
)

// Postaja represents a weather station with its name, location
// and a list of measured temperatures in degrees Celsius.
type Postaja struct {
	ime 	string
	kraj 	string
	meritve	[]float64
}

// DodajPostajo adds a new station with the given identifier, name and
// location to the map postaje. If a station with the same identifier
// already exists, it is left unchanged and a message is printed.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) {
	if _, obstaja := postaje[oznaka]; obstaja {
		fmt.Println("Postaja z oznako", oznaka, "ze obstaja.");
		return;
	}
	postaje[oznaka] = Postaja{ime: ime, kraj: kraj};
}

// DodajMeritev adds the temperature measurement temperatura to the station
// with the given identifier. A measurement is accepted only if it lies
// within the interval [minTemp, maxTemp] and the station exists;
// otherwise a message is printed and the measurement is discarded.
func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura, minTemp, maxTemp float64) {
	postaja, ok := postaje[oznaka];

	if temperatura < minTemp || temperatura > maxTemp {
		fmt.Println("Neveljavna temperatura.");
	}
	if !ok {
		fmt.Println("Neveljavna oznaka.");
		return;
	}

	postaja.meritve = append(postaja.meritve, temperatura);
	postaje[oznaka] = postaja;
}

// statistika returns the lowest, highest and average temperature of the
// station with the given identifier. The result ok is false if the station
// does not exist or has fewer than stMeritev measurements.
func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (najnizja, najvisja, povprecna float64, ok bool) {
	postaja, ok := postaje[oznaka];
	if !ok || len(postaja.meritve) == 0 || len(postaja.meritve) < stMeritev {
		return 0, 0, 0, false;
	}

	najnizja = postaja.meritve[0];
	najvisja = postaja.meritve[0];
	var vsota float64;

	for _, v := range postaja.meritve {
		if v < najnizja {
			najnizja = v;
		}
		if v > najvisja {
			najvisja = v;
		}
		vsota += v;
	}

	return najnizja, najvisja, vsota / float64(len(postaja.meritve)), true;
}

// IzpisMeritev prints all stations with their names, locations and
// measurements, sorted by station identifier.
func IzpisMeritev(postaje map[string]Postaja) {
	fmt.Println("Meritve:")
	for _, oznaka := range urejeneOznake(postaje) {
		postaja := postaje[oznaka]
		fmt.Printf("%s - %s (%s): %v\n", oznaka, postaja.ime, postaja.kraj, postaja.meritve)
	}
}

// IzpisPorocila prints a report for every station: the lowest, highest
// and average temperature and a short description based on the average
// ("Mrzlo", "Hladno", "Toplo" or "Vroče"). Stations with fewer than
// stMeritev measurements are reported as having too few measurements.
func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	for _, oznaka := range urejeneOznake(postaje) {
		postaja := postaje[oznaka];
		min, max, avg, ok := statistika(postaje, oznaka, stMeritev);

		if !ok {
			fmt.Printf("%s: premalo meritev\n", postaja.ime);
			continue;
		}

		opis := "Vroče";
		switch {
		case avg < 5:
			opis = "Mrzlo"

		case avg < 15:
			opis = "Hladno"

		case avg < 25:
			opis = "Toplo"
		}

		fmt.Printf("%s: min %.2f °C, max %.2f °C, povprečje %.2f °C -> %s\n", postaja.ime, min, max, avg, opis);
	}
}


func urejeneOznake(postaje map[string]Postaja) []string {
	oznake := make([]string, 0, len(postaje));
	for oznaka := range postaje {
		oznake = append(oznake, oznaka);
	}
	sort.Strings(oznake);
	return oznake;
}