package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type SimulationConfig struct {
	WorkersAmount         [3]int `json:"workers_amount"`
	StonesExtractionTime  [3]int `json:"stones_extraction_time"`
	StonesMasses          [3]int `json:"stones_masses"`
	StoneMassesLimits     [3]int `json:"stone_masses_limits"`
	TimeToTravelEmpty     int    `json:"time_to_travel_empty"`
	TimeToTravelFull      int    `json:"time_to_travel_full"`
	QuarryWorkplaces      int    `json:"quarry_workplaces"`
	TimeToPlaceStone      int    `json:"time_to_place_stone"`
	TimeToPlaceInsulation int    `json:"time_to_place_insulation"`
	TimeToChangePallet    int    `json:"time_to_change_pallet"`
}

func defaultConfig() SimulationConfig {
	return SimulationConfig{
		WorkersAmount:         [3]int{3, 3, 3},
		StonesExtractionTime:  [3]int{4, 4, 4},
		StonesMasses:          [3]int{1, 3, 5},
		StoneMassesLimits:     [3]int{14, 13, 11},
		TimeToTravelEmpty:     500,
		TimeToTravelFull:      1200,
		QuarryWorkplaces:      4,
		TimeToPlaceStone:      300,
		TimeToPlaceInsulation: 500,
		TimeToChangePallet:    1000}
}

func getConfig() SimulationConfig {
	read, err := os.ReadFile("./config/config.json")
	conf := SimulationConfig{}
	err = json.Unmarshal(read, &conf)
	if err != nil {
		fmt.Println("No config found, creating one with default values")
		conf := defaultConfig()
		file, err := os.Create("./config/config.json")
		str, err := json.Marshal(conf)
		if err != nil {
			fmt.Println("Couldn't create config file")
			return conf
		}

		defer file.Close()
		file.Write(str)
		return conf
	}
	return conf

}

func (conf SimulationConfig) printConfig() {
	fmt.Println("Current config:")
	v := reflect.ValueOf(conf)
	typeOfConf := v.Type()
	for i := 0; i < typeOfConf.NumField(); i++ {
		fmt.Printf("%s: %v\n", typeOfConf.Field(i).Name, v.Field(i).Interface())
	}
	fmt.Println("")
}
