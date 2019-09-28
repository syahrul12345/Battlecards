//Package models represents all character models that is received from SWAPI, the postgres character model, and the http response model.
package models

import (
	"strings"
	"unicode"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//State is the flag to only index when first starting
type State struct {
	gorm.Model
	Name  string
	Index bool
}

//CharacterResponse is the JSON object of a character when it is first returend by
type CharacterResponse struct {
	Name      string
	Gender    string
	StarShips []string
	Vehicles  []string
	HomeWorld string
}

//Character represents the JSON object of a StarWars character after reindexing
type Character struct {
	gorm.Model
	NameSearch string
	Name       string
	Gender     string
	StarShips  pq.ByteaArray `gorm:"type:varchar(1000)[]"`
	Vehicles   pq.ByteaArray `gorm:"type:varchar(1000)[]"`
	Home       pq.ByteaArray `gorm:"type:varchar(8000)[]"`
}

//CharacterResult is the JSON object send to the front end
type CharacterResult struct {
	Name      string
	Gender    string
	StarShips []StarShip
	Vehicles  []Vehicle
	HomeWorld Home
}

//StarShip represents the JSON object of a one StarShip after reindexing
type StarShip struct {
	Model             string
	Starship_Class    string
	Hyperdrive_Rating string
	Cost_In_Credits   string
	Manufacturer      string
}

//Vehicle represents the JSON object of one Vehicle after reindexing
type Vehicle struct {
	Name          string
	Model         string
	CostInCredits string
}

//Home represents the homeworld stats
type Home struct {
	Name       string
	Population string
	Climate    string
}

//ToLowerAndNoSpecial removes all whitespaces, diacratics and special charates such as / and -
func ToLowerAndNoSpecial(str string) string {
	return strings.Map(func(r rune) rune {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r)) {
			return -1
		}
		return r
	}, normalize(strings.ToLower(str)))
}

//Helper function for ToLowerAndNoSpecial
func normalize(str string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, str)
	return result
}

//Helper function for ToLowerAndNoSpecial
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
