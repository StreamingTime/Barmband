package color

import (
	"errors"
	"math/rand"
	"slices"
)

// List of colors picked from https://en.wikipedia.org/wiki/Help:Distinguishable_colors
var Colors = []string{
	"F0A3FF", // Amethyst
	"0075DC", // Blue
	"993F00", // Caramel
	"4C005C", // Damson
	"191919", // Ebony
	"005C31", // Forest
	"2BCE48", // Green
	"FFCC99", // Honeydew
	"808080", // Iron
	"94FFB5", // Jade
	"8F7C00", // Khaki
	"9DCC00", // Lime
	"C20088", // Mallow
	"003380", // Navy
	"FFA405", // Orpiment
	"FFA8BB", // Pink
	"426600", // Quagmire
	"FF0010", // Red
	"5EF1F2", // Sky
	"00998F", // Turquoise
	"E0FF66", // Uranium
	"740AFF", // Violet
	"990000", // Wine
	"FFFF80", // Xanthin
	"FFE100", // Yellow
	"FF5005", // Zinnia
}

var ErrNoColorAvailable = errors.New("no color available")

func GetRandomColor(ignoreColors []string) (string, error) {
	var availableColors []string

	for _, color := range Colors {
		if !slices.Contains(ignoreColors, color) {
			availableColors = append(availableColors, color)
		}
	}

	if len(availableColors) == 0 {
		return "", ErrNoColorAvailable
	}

	index := rand.Intn(len(availableColors))

	return availableColors[index], nil
}
