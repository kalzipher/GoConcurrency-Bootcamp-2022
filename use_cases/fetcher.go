package use_cases

import (
	"GoConcurrency-Bootcamp-2022/models"
	"log"
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

type ResultPokemon struct {
	Pokemon models.Pokemon
	Err     error
}

func NewFetcher(api api, storage writer) Fetcher {
	return Fetcher{api, storage}
}

func (fetcher Fetcher) Fetch(from, to int) error {
	var listPokemon []models.Pokemon
	var finish = make(chan bool, 1)
	defer close(finish)
	chResponse := GeneratorResponsePokemon(from, to, fetcher, finish)
	for id := from; id <= to; id++ {
		pokemonResponse := <-chResponse
		if pokemonResponse.Err != nil {
			finish <- true
			continue
		}
		listPokemon = append(listPokemon, pokemonResponse.Pokemon)
	}
	return fetcher.storage.Write(listPokemon)
}

func GeneratorResponsePokemon(from, to int, fetcher Fetcher, abort <-chan bool) <-chan ResultPokemon {
	response := make(chan ResultPokemon)
	wg := sync.WaitGroup{}
	go func() {
		wg.Wait()
		close(response)
	}()

	for i := from; i <= to; i++ {
		wg.Add(1)
		go FetchPokemonById(wg, i, fetcher, response, abort)
	}
	return response
}

func FetchPokemonById(wg sync.WaitGroup, idPokemon int, fetcher Fetcher, response chan<- ResultPokemon, abort <-chan bool) {
	defer wg.Done()
	pokemon, err := fetcher.api.FetchPokemon(idPokemon)
	if err != nil {
		var flatAbilities []string
		for _, itemAbility := range pokemon.Abilities {
			flatAbilities = append(flatAbilities, itemAbility.Ability.URL)
		}
		pokemon.FlatAbilityURLs = strings.Join(flatAbilities, "|")
	}
	responsePokemon := ResultPokemon{
		pokemon,
		err,
	}
	select {
	case <-abort:
		log.Println("Abort!!!")
		return
	case response <- responsePokemon:
	}
}
