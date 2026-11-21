/*
Vremenska-postaja demonstrates the use of package vreme.
It creates a few sample weather stations, adds temperature
measurements to them and prints a report for every station.

Measurements outside the allowed temperature interval are rejected,
and stations with too few measurements are reported without statistics.
For every other station the report contains the lowest, highest and
average temperature and a short description of the average:
"Mrzlo" (below 5 °C), "Hladno" (below 15 °C), "Toplo" (below 25 °C)
or "Vroče" (25 °C or more).

Usage:

	vremenska-postaja [flags]

The flags are:

	--minTemp value
		The lowest allowed temperature of a measurement in °C.
		Measurements below this value are rejected. Default -50.
	--maxTemp value
		The highest allowed temperature of a measurement in °C.
		Measurements above this value are rejected. Default 50.
	--stMeritev value
		The minimum number of measurements a station needs for its
		statistics to be calculated. Must be at least 1. Default 3.
	--help, -h
		Print the list of flags and exit.

The program exits with an error if minTemp is greater than maxTemp
or if stMeritev is less than 1.

For example, to accept only measurements between -10 °C and 40 °C
and require at least 5 measurements per station, run:

	vremenska-postaja --minTemp -10 --maxTemp 40 --stMeritev 5
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/TheBezstSchnitzel/vremenska-postaja/vreme"
)


func main() {
	cmd := &cli.Command{
		Name: "vremenska-postaja",
		Usage: "dodaja meritve vremenskim postajam in izpiše poročilo",
		Flags: []cli.Flag{
			&cli.FloatFlag{
				Name: "minTemp",
				Usage: "najnižja ddovoljena temperatura meritve v °C",
				Value: -50,
			},
			&cli.FloatFlag{
				Name: "maxTemp",
				Usage: "najvišja dovoljena temperatura meritve v °C",
				Value: 50,
			},
			&cli.IntFlag{
				Name: "stMeritev",
				Usage: "najmanjše število meritev, potrebnih za izračun statistike",
				Value: 3,
			},
		},
		Action: zazeni,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}


func zazeni(ctx context.Context, cmd *cli.Command) error {
	minTemp := cmd.Float("minTemp");
	maxTemp := cmd.Float("maxTemp");
	stMeritev := int(cmd.Int("stMeritev"));

	if minTemp > maxTemp {
		return fmt.Errorf("minTemp (%v) ne sme biti večja od maxTemp (%v)", minTemp, maxTemp);
	}
	if stMeritev < 1 {
		return fmt.Errorf("stMeritev mora biti vsaj 1");
	}

	postaje := map[string]vreme.Postaja{}
	vreme.DodajPostajo(postaje, "LJ01", "Joža Pličnik", "Ljubljana")
	vreme.DodajPostajo(postaje, "LJ69", "Franci Prošern", "Lublana")
	vreme.DodajPostajo(postaje, "CR12", "Vladimir Potnik", "Spodnji Duplek")
	vreme.DodajPostajo(postaje, "PO02", "Postaja št. 13194", "Zgornji Duplek")
	vreme.DodajPostajo(postaje, "P0D3", "Postaja z 2 vnosoma", "Kredarica")
	vreme.DodajPostajo(postaje, "LJ01", "Podvojena postaja", "Ljubljana")

	zacetne := map[string][]float64{
		"LJ01": {1.0, 24.0, 30.0, 45.0},
		"LJ69": {14.6, 24.0, 31.2, 9.11},
		"CR12": {22.2, 28.5, 36.9, 21.54},
		"PO02": {45.2, 3.5, 1.7, 2.54},
		"P0D3": {1.2, 2.1},
	}

	for oznaka, meritve := range zacetne {
		for _, t := range meritve {
			vreme.DodajMeritev(postaje, oznaka, t, minTemp, maxTemp)
		}
	}

	fmt.Println("Dodajanje meritev: ----------------------")
	vreme.DodajMeritev(postaje, "LJ69", 44.56, minTemp, maxTemp) // veljavna
	vreme.DodajMeritev(postaje, "LJ69", 55.56, minTemp, maxTemp) // neveljavna temperatura
	vreme.DodajMeritev(postaje, "LJO2", 44.56, minTemp, maxTemp) // neveljavna oznaka
	vreme.IzpisMeritev(postaje)

	fmt.Println("Poročilo: -----------------------")
	vreme.IzpisPorocila(postaje, stMeritev)

	return nil
}