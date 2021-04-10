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

var enchantableIndexes = [...]bool{
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

type AnalysisResults struct {
	Consumables     map[string]map[string]int
	WorldBuffs      map[string][]string
	ItemsReport     map[string]map[string]int
	MissingEnchants map[string][]string
	MissingItems    map[string][]string
}

func analyze(p *ParseResults) *AnalysisResults {
	a := &AnalysisResults{}

	if conf.WorldBuffs != nil {
		log.Debug().Msg("Checking encounters buffs")
		worldBuffs := make(map[string][]string)

		// Only check first encounter for each player
		for _, e := range p.Encounters {
			for _, p := range e.Players {
				if worldBuffs[p.Name] != nil {
					continue
				}

				worldBuffs[p.Name] = make([]string, 0)
				for _, s := range p.WorldBuffs {
					if conf.WorldBuffs[s] {
						worldBuffs[p.Name] = append(worldBuffs[p.Name], s)
					}
				}
			}
		}

		a.WorldBuffs = worldBuffs
	}

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
				continue
			}

			for _, p := range e.Players {
				if missingEnchants[p.Name] != nil {
					continue
				}

				missingEnchants[p.Name] = make([]string, 0)
				// Check every item for a permanent enchant
				for idx, i := range p.Items {
					if i.ID != "0" && enchantableIndexes[idx] && (len(i.Enchants) == 0 || i.Enchants[0] == "0" || conf.IgnoredEnchants[i.Enchants[0]]) {
						missingEnchants[p.Name] = append(missingEnchants[p.Name], itemIndexToSlot[idx])
					}
				}
			}
		}

		a.MissingEnchants = missingEnchants
	}

	if conf.TrackedItems != nil {
		log.Debug().Msg("Checking tracked items")
		trackedItems := make(map[string]map[string]int)

		trackedItemsList := make(map[string]map[string]bool)

		for _, e := range p.Encounters {
			for _, p := range e.Players {
				if trackedItemsList[p.Name] == nil {
					trackedItemsList[p.Name] = make(map[string]bool)
				}

				for _, i := range p.Items {
					if conf.TrackedItems[i.ID] != "" {
						trackedItemsList[p.Name][i.ID] = true
					}
				}
			}
		}

		for name, items := range trackedItemsList {
			trackedItems[name] = make(map[string]int)
			for id := range items {
				rarity := conf.TrackedItems[id]
				trackedItems[name][rarity]++
			}
		}

		a.ItemsReport = trackedItems
	}

	if p.Casts != nil {
		log.Debug().Msg("Checking consummables")
		consumables := make(map[string]map[string]int)

		for name, casts := range p.Casts {
			for spellId, count := range casts {
				if consumables[spellId] == nil {
					consumables[spellId] = make(map[string]int)
				}
				consumables[spellId][name] = count
			}
		}

		a.Consumables = consumables
	}

	return a
}
