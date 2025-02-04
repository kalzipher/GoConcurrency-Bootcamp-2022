package repositories

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"GoConcurrency-Bootcamp-2022/models"
)

type LocalStorage struct{}

const filePath = "resources/pokemons.csv"

func (l LocalStorage) Write(pokemons []models.Pokemon) error {
	file, fErr := os.Create(filePath)
	defer file.Close()
	if fErr != nil {
		return fErr
	}

	w := csv.NewWriter(file)
	records := buildRecords(pokemons)
	if err := w.WriteAll(records); err != nil {
		return err
	}

	return nil
}

func (l LocalStorage) Read() ([]models.Pokemon, error) {
	file, fErr := os.Open(filePath)
	defer file.Close()
	if fErr != nil {
		return nil, fErr
	}

	r := csv.NewReader(file)
	records, rErr := r.ReadAll()
	if rErr != nil {
		return nil, rErr
	}

	pokemons, err := parseCSVData(records)
	if err != nil {
		return nil, err
	}

	return pokemons, nil
}

func parseCSVRecord(record []string) (*models.Pokemon, error) {
	var pokemon models.Pokemon

	var toConvert = map[string]int{
		"id":     0,
		"height": 2,
		"weight": 3,
	}

	for key := range toConvert {
		value, err := strconv.Atoi(record[toConvert[key]])
		if err != nil {
			return nil, err
		}
		toConvert[key] = value
	}

	pokemon = models.Pokemon{
		ID:              toConvert["id"],
		Name:            record[1],
		Height:          toConvert["height"],
		Weight:          toConvert["weight"],
		Abilities:       nil,
		FlatAbilityURLs: record[4],
		EffectEntries:   nil,
	}

	return &pokemon, nil
}

func buildRecords(pokemons []models.Pokemon) [][]string {
	headers := []string{"id", "name", "height", "weight", "flat_abilities"}
	records := [][]string{headers}
	for _, p := range pokemons {
		record := fmt.Sprintf("%d,%s,%d,%d,%s",
			p.ID,
			p.Name,
			p.Height,
			p.Weight,
			p.FlatAbilityURLs)
		records = append(records, strings.Split(record, ","))
	}

	return records
}

func parseCSVData(records [][]string) ([]models.Pokemon, error) {
	var pokemons []models.Pokemon
	for i, record := range records {
		if i == 0 {
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		height, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}

		weight, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}

		pokemon := models.Pokemon{
			ID:              id,
			Name:            record[1],
			Height:          height,
			Weight:          weight,
			Abilities:       nil,
			FlatAbilityURLs: record[4],
			EffectEntries:   nil,
		}
		pokemons = append(pokemons, pokemon)
	}

	return pokemons, nil
}
