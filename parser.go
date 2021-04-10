package main

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
)

var combatantInfoRegex = regexp.MustCompile(`(.*),(\d+),(\d+),(\d+),(\d+),(\d+),0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,(\d+),0,\(\),\(0,0,0,0\),\[],\[(.*)],\[(.*)]`)
var itemRegex = regexp.MustCompile(`\((?P<itemId>\d+),(?P<itemLvl>\d+),\((?P<enchants>[^)]*)\),\((?P<bonus>[^)]*)\),\((?P<gems>[^)]*)\)\)`)

const nilGuid string = "0000000000000000"
const nilName string = "nil"

type ParseResults struct {
	LinesCount     int
	EventsCount    map[string]int
	GuidMap        map[string]string
	Encounters     []Encounter
	Buffs          map[string]map[string]int
	EncounterBuffs map[string]map[string]string
	Casts          map[string]map[string]int
}

type Encounter struct {
	Name    string
	ID      string
	Players map[string]*EncounterPlayer
}

type EncounterPlayer struct {
	Name        string
	WorldBuffs  []string
	Items       []Item
	Consumables []string
}

type Item struct {
	ID string
	//ItemLevel string
	Enchants []string
}

var currentEncounter *Encounter

func parse(scanner *bufio.Scanner) *ParseResults {
	p := &ParseResults{
		LinesCount:     0,
		EventsCount:    make(map[string]int),
		GuidMap:        make(map[string]string),
		Encounters:     make([]Encounter, 0),
		Buffs:          make(map[string]map[string]int),
		EncounterBuffs: make(map[string]map[string]string),
		Casts:          make(map[string]map[string]int),
	}

	for scanner.Scan() {
		parseLine(scanner.Text(), p)
		p.LinesCount++
	}

	return p
}

func parseLine(line string, p *ParseResults) {
	s1 := strings.Split(line, "  ")
	//timestamp := s1[0]
	s2 := strings.SplitN(s1[1], ",", 2)
	event := s2[0]

	switch event {
	case "COMBATANT_INFO":
		parsedData := combatantInfoRegex.FindStringSubmatch(s2[1])
		guid := parsedData[1]
		//str := parsedData[2]
		//agi := parsedData[3]
		//sta := parsedData[3]
		//intl := parsedData[5]
		//spi := parsedData[6]
		//armor := parsedData[7]

		playerName := p.GuidMap[guid]
		player := &EncounterPlayer{
			Name:        playerName,
			WorldBuffs:  make([]string, 0),
			Items:       make([]Item, 0),
			Consumables: make([]string, 0),
		}

		//items := parsedData[8]
		parsedItems := itemRegex.FindAllStringSubmatch(parsedData[8], -1)
		for _, item := range parsedItems {
			itemId := item[1]
			//itemLvl := item[2]
			var enchants []string
			if item[3] != "" {
				enchants = strings.Split(item[3], ",")
			}
			//bonus := item[4] // Always empty on classic
			//gems := item[5] // Always empty on classic

			item := &Item{
				ID: itemId,
				//ItemLevel: itemLvl,
				Enchants: enchants,
			}
			player.Items = append(player.Items, *item)
		}

		//auras := parsedData[9]
		auras := make([]string, 0)
		parsedAuras := strings.Split(parsedData[9], ",")
		for i, guidOrSpellId := range parsedAuras {
			if i%2 == 1 {
				auras = append(auras, guidOrSpellId)
			}
		}

		player.WorldBuffs = auras
		currentEncounter.Players[playerName] = player
	case "SPELL_AURA_APPLIED":
		data := strings.Split(s2[1], ",")
		guid := data[0]
		name := strings.Trim(data[1], "\"")
		//flags := data[2]
		//raidFlags := data[3]
		//targetGuid := data[4]
		//targetName := data[5]
		//targetFlags := data[6]
		//targetRaidFlags := data[7]
		spellId := data[8]
		if conf.TrackedBuffs[spellId] {
			if p.Buffs[name] == nil {
				p.Buffs[name] = make(map[string]int)
			}
			p.Buffs[name][spellId]++
		}
		// Let's use this event to build our guidMap
		if guid != nilGuid && name != nilName && p.GuidMap[guid] == "" {
			p.GuidMap[guid] = name
		}
	case "SPELL_AURA_REFRESH":
		data := strings.Split(s2[1], ",")
		//guid := data[0]
		name := strings.Trim(data[1], "\"")
		//flags := data[2]
		//raidFlags := data[3]
		//targetGuid := data[4]
		//targetName := data[5]
		//targetFlags := data[6]
		//targetRaidFlags := data[7]
		spellId := data[8]
		if conf.TrackedBuffs[spellId] {
			if p.Buffs[name] == nil {
				p.Buffs[name] = make(map[string]int)
			}
			p.Buffs[name][spellId]++
		}
	case "SPELL_AURA_REMOVED":
		data := strings.Split(s2[1], ",")
		//guid := data[0]
		name := strings.Trim(data[1], "\"")
		//flags := data[2]
		//raidFlags := data[3]
		//targetGuid := data[4]
		//targetName := data[5]
		//targetFlags := data[6]
		//targetRaidFlags := data[7]
		spellId := data[8]
		if conf.TrackedEncounterBuffs[spellId] {
			if currentEncounter != nil {
				currentEncounter.Players[name].Consumables = append(currentEncounter.Players[name].Consumables, spellId)
			}
		}

	case "SPELL_PERIODIC_ENERGIZE":
		//data := strings.Split(s2[1], ",")
		//guid := data[0]
		//name := data[1]
		//flags := data[2]
		//raidFlags := data[3]
		//targetGuid := data[4]
		//targetName := data[5]
		//targetFlags := data[6]
		//targetRaidFlags := data[7]
	case "SPELL_CAST_SUCCESS":
		data := strings.Split(s2[1], ",")
		//guid := data[0]
		name := strings.Trim(data[1], "\"")
		//flags := data[2]
		//raidFlags := data[3]
		//targetGuid := data[4]
		//targetName := data[5]
		//targetFlags := data[6]
		//targetRaidFlags := data[7]
		spellId := data[8]
		if conf.TrackedCasts[spellId] {
			if p.Casts[name] == nil {
				p.Casts[name] = make(map[string]int)
			}
			p.Casts[name][spellId]++
		}
	case "ENCOUNTER_START":
		data := strings.Split(s2[1], ",")
		encounterId := data[0]
		name := strings.Trim(data[1], "\"")
		//diff := data[2]
		//playerCount := data[3]
		currentEncounter = &Encounter{
			ID:      encounterId,
			Name:    name,
			Players: make(map[string]*EncounterPlayer),
		}
		//if p.WorldBuffs[name] == nil {
		//	p.WorldBuffs[name] = make(map[string]int)
		//}
	case "ENCOUNTER_END":
		//data := strings.Split(s2[1], ",")
		//encounterId := data[0]
		//name := data[1]
		if currentEncounter == nil {
			log.Error().Msg("ENCOUNTER_END without ENCOUNTER_START")
		} else {
			p.Encounters = append(p.Encounters, *currentEncounter)
		}
		currentEncounter = nil
	}

	p.EventsCount[event]++
}
