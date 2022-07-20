package use_cases

import (
	"GoConcurrency-Bootcamp-2022/models"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type api interface {
	FetchPokemon(id int) (models.Pokemon, error)
}

type writer interface {
	Write(pokemons []models.Pokemon) error
}

type Fetcher struct {
	api     api
	storage writer
}

func NewFetcher(api api, storage writer) Fetcher {
	return Fetcher{api, storage}
}

func (f Fetcher) Fetch(from, to int) error {
	var totalWorkers = to - from
	fmt.Println(totalWorkers)
	var pokemons []models.Pokemon
	var pokemonChannels = make(chan models.Pokemon, totalWorkers)
	var wg sync.WaitGroup
	for id := from; id <= to; id++ {
		go func(index int) {
			wg.Add(1)
			pokemon, err := f.api.FetchPokemon(index)
			if err != nil {
				return
			}

			var flatAbilities []string
			for _, t := range pokemon.Abilities {
				flatAbilities = append(flatAbilities, t.Ability.URL)
			}
			pokemon.FlatAbilityURLs = strings.Join(flatAbilities, "|")
			pokemonChannels <- pokemon
		}(id)
	}

	go func(pokeChannel chan models.Pokemon) {
		defer close(pokeChannel)
		for pokemon := range pokeChannel {
			wg.Done()
			pokemons = append(pokemons, pokemon)
		}
	}(pokemonChannels)

	wg.Wait()

	sort.Slice(pokemons, func(i, j int) bool {
		return pokemons[i].ID < pokemons[j].ID
	})

	return f.storage.Write(pokemons)
}
