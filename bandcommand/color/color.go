package color

import (
	"errors"
	"math/rand"
	"slices"
)

// List of colors picked from https://en.wikipedia.org/wiki/Help:Distinguishable_colors
var Colors = []string{
	"0xF0A3FF", // Amethyst
	"0x0075DC", // Blue
	"0x993F00", // Caramel
	"0x4C005C", // Damson
	"0x191919", // Ebony
	"0x005C31", // Forest
	"0x2BCE48", // Green
	"0xFFCC99", // Honeydew
	"0x808080", // Iron
	"0x94FFB5", // Jade
	"0x8F7C00", // Khaki
	"0x9DCC00", // Lime
	"0xC20088", // Mallow
	"0x003380", // Navy
	"0xFFA405", // Orpiment
	"0xFFA8BB", // Pink
	"0x426600", // Quagmire
	"0xFF0010", // Red
	"0x5EF1F2", // Sky
	"0x00998F", // Turquoise
	"0xE0FF66", // Uranium
	"0x740AFF", // Violet
	"0x990000", // Wine
	"0xFFFF80", // Xanthin
	"0xFFE100", // Yellow
	"0xFF5005", // Zinnia
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
