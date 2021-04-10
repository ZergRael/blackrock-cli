package main

import "github.com/rs/zerolog/log"

var itemIndexToSlot = [...]string{
	"head",
	"neck",
	"shoulders",
	"shirt",
	"torso",
	"waist",
	"legs",
	"feet",
	"wrists",
	"hands",
	"ring1",
	"ring2",
	"trinket1",
	"trinket2",
	"back",
	"mainhand",
	"offhand",
	"ranged",
	"tabard",
}

var enchantableIndexes = [...]bool {
	true,
	false,
	true,
	false,
	true,
	false,
	true,
	true,
	true,
	true,
	false,
	false,
	false,
	false,
	true,
	true,
	true,
	true,
	false,
}

type AnalysysResults struct {
	PreRaidBuffs map[string][]string
	MissingEnchants map[string][]string
	MissingItems map[string][]string
}

func analyze(p *ParseResults) *AnalysysResults {
	a := &AnalysysResults{}

	log.Debug().Msg("Checking missing buffs")
	preRaidBuffs := make(map[string][]string)

	// Only check first encounter for each player
	for _, e := range p.Encounters {
		for _, p := range e.Players {
			if preRaidBuffs[p.Name] != nil {
				continue
			}

			preRaidBuffs[p.Name] = make([]string, 0)
			for _, s := range p.Buffs {
				if conf.EncounterBuffs[s] {
					preRaidBuffs[p.Name] = append(preRaidBuffs[p.Name], s)
				}
			}
		}
	}

	a.PreRaidBuffs = preRaidBuffs

	log.Debug().Msg("Checking missing items")
	missingItems := make(map[string][]string)

	// Only check first encounter for each player
	for _, e := range p.Encounters {
		for _, p := range e.Players {
			if missingItems[p.Name] != nil {
				continue
			}

			missingItems[p.Name] = make([]string, 0)
			for idx, i := range p.Items {
				if i.ID == "0" && itemIndexToSlot[idx] != "tabard" && itemIndexToSlot[idx] != "shirt" {
					missingItems[p.Name] = append(missingItems[p.Name], itemIndexToSlot[idx])
				}
			}
		}
	}

	a.MissingItems = missingItems

	if conf.CheckEnchants {
		log.Debug().Msg("Checking missing enchants")
		missingEnchants := make(map[string][]string)

		// Only check first encounter for each player
		for _, e := range p.Encounters {
			if conf.IgnoredEncountersEnchants[e.ID] {
				log.Debug().Str("encounterId", e.ID).Msg("Ignore encounter")
				continue
			}

			for _, p := range e.Players {
				if missingEnchants[p.Name] != nil {
					continue
				}

				missingEnchants[p.Name] = make([]string, 0)
				// Check every item for a permanent enchant
				for idx, i := range p.Items {
					if i.ID != "0" && enchantableIndexes[idx] && (len(i.Enchants) == 0 || i.Enchants[0] == "0" || conf.IgnoredEnchants[i.Enchants[0]]){
						missingEnchants[p.Name] = append(missingEnchants[p.Name], itemIndexToSlot[idx])
					}
				}
			}
		}

		a.MissingEnchants = missingEnchants
	}

	return a
}