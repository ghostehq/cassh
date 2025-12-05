// Provides critically important LSP and Flash meme rotation system
// This is the most vital component of cassh
package memes

import (
	"math/rand/v2"
)

type Character struct {
	Name       string
	Image      string // Filename in static/images
	AltText    string
	Quotes     []string
	ColorTheme string // CSS color for theming
}

// LumpySpacePrincess - the drama queen we need
var LumpySpacePrincess = Character{
	Name:       "Lumpy Space Princess",
	Image:      "lsp.png",
	AltText:    "Lumpy Space Princess floating dramatically",
	ColorTheme: "#9B59B6",
	Quotes: []string{
		"Oh my GLOB, just click the button!",
		"These lumps aren't gonna authenticate themselves!",
		"Whatever, I'm getting cheese fries... after you sign in.",
		"You want access? FINE. But only because I'm NICE.",
		"My SSH certs bring all the devs to the yard.",
		"I'm, like, totally secure and stuff.",
		"Drama bomb! Your cert expired!",
		"UGH, just sign in already! I have PLANS.",
		"I didn't spend 30 seconds on this login page for you to hesitate!",
		"Your old cert? It's DEAD to me now.",
		"Listen, I'm only gonna say this once... JK I'll say it forever: SIGN IN.",
		"I'm not bossy, I just have better ideas. Like 12-hour certs.",
		"Oh, you need access? How MATHEMATICAL.",
		"Authenticate me! I mean... authenticate yourself!",
		"I'm too beautiful for permanent SSH keys.",
	},
}

// DMVSloth - patience is a virtue, especially at the DMV
var DMVSloth = Character{
	Name:       "Flash Slothmore",
	Image:      "sloth.png",
	AltText:    "Flash from Zootopia, smiling slowly",
	ColorTheme: "#27AE60",
	Quotes: []string{
		"Ha... ha... ha... click... the... button...",
		"Let... me... just... verify... your... identity...",
		"What... do... you... call... a... three... humped... camel?",
		"Your... cert... will... be... ready... in... *checks watch* ...12 hours.",
		"Wel... come... to... the... cassh... authentication... portal...",
		"I'm... working... as... fast... as... I... can...",
		"Sign... in... please... I'll... wait...",
		"Ha... ha... ha... security... is... no... joke...",
		"Take... your... time... I... have... all... day...",
		"Pre... gnant... pause... just... kidding... click... SSO...",
		"Your... patience... is... appreciated...",
		"Almost... there... just... 12... more... hours...",
		"I... love... ephemeral... certs... they're... so... fleeting...",
		"Flash... is... my... name... authentication... is... my... game...",
		"Need... to... go... faster?... Too... bad...",
	},
}

var Characters = []Character{
	LumpySpacePrincess,
	DMVSloth,
}

func GetRandomCharacter() Character {
	return Characters[rand.IntN(len(Characters))]
}

func GetCharacterByName(name string) Character {
	switch name {
	case "lsp":
		return LumpySpacePrincess
	case "sloth":
		return DMVSloth
	default:
		return GetRandomCharacter()
	}
}

func GetRandomQuote(c Character) string {
	return c.Quotes[rand.IntN(len(c.Quotes))]
}

// MemeData is passed to templates for rendering
type MemeData struct {
	Character  Character
	Quote      string
	ColorTheme string
}

// GetMemeData returns randomized meme data for template rendering
func GetMemeData(preference string) MemeData {
	var char Character
	if preference == "random" || preference == "" {
		char = GetRandomCharacter()
	} else {
		char = GetCharacterByName(preference)
	}

	return MemeData{
		Character:  char,
		Quote:      GetRandomQuote(char),
		ColorTheme: char.ColorTheme,
	}
}
