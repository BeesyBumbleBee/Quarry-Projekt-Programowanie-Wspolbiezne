package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type SimulationConfig struct {
	WorkersAmount         [3]int    `json:"workers_amount"`
	StonesExtractionTime  [3][2]int `json:"stones_extraction_time"`
	StonesMasses          [3]int    `json:"stones_masses"`
	StoneMassesLimits     [3]int    `json:"stone_masses_limits"`
	TimeToTravelEmpty     [2]int    `json:"time_to_travel_empty"`
	TimeToTravelFull      [2]int    `json:"time_to_travel_full"`
	QuarryWorkplaces      int       `json:"quarry_workplaces"`
	TimeToPlaceStone      [2]int    `json:"time_to_place_stone"`
	TimeToPlaceInsulation [2]int    `json:"time_to_place_insulation"`
	TimeToChangePallet    [2]int    `json:"time_to_change_pallet"`
}

func defaultConfig() SimulationConfig {
	return SimulationConfig{
		WorkersAmount:         [3]int{2, 2, 3},
		StonesExtractionTime:  [3][2]int{{500, 2000}, {1000, 2500}, {2000, 3500}},
		StonesMasses:          [3]int{1, 3, 5},
		StoneMassesLimits:     [3]int{14, 13, 11},
		TimeToTravelEmpty:     [2]int{50, 100},
		TimeToTravelFull:      [2]int{100, 200},
		QuarryWorkplaces:      3,
		TimeToPlaceStone:      [2]int{300, 1500},
		TimeToPlaceInsulation: [2]int{500, 1000},
		TimeToChangePallet:    [2]int{1000, 2000}}
}

func getConfig() (SimulationConfig, error) {

	read, err := os.ReadFile("./config.json")
	conf := SimulationConfig{}
	if err != nil {
		//fmt.Println("No config found, creating one with default values")
		conf := defaultConfig()
		file, err := os.Create("./config.json")
		var str []byte
		str, err = json.Marshal(conf)
		if err != nil {
			fmt.Println("Couldn't create config file")
			return conf, err
		}

		defer func(file *os.File) {
			err = file.Close()
			if err != nil {

			}
		}(file)
		_, err = file.Write(str)
		if err != nil {
			return defaultConfig(), err
		}
		return conf, err
	}
	err = json.Unmarshal(read, &conf)
	if err != nil {
		return defaultConfig(), err
	}
	return conf, err

}

func (conf SimulationConfig) printConfig() string {
	out := "Current config:\n"
	v := reflect.ValueOf(conf)
	typeOfConf := v.Type()
	for i := 0; i < typeOfConf.NumField(); i++ {
		out += fmt.Sprintf("%s: %v\n", typeOfConf.Field(i).Name, v.Field(i).Interface())
	}
	out += "\n"
	return out
}

func (conf SimulationConfig) saveConfig() error {
	file, err := os.Create("./config.json")
	if err != nil {
		fmt.Println("Couldn't create config file")
		return err
	}
	var str []byte
	str, err = json.Marshal(conf)
	if err != nil {
		fmt.Println("Couldn't create json")
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)
	_, err = file.Write(str)
	if err != nil {
		return err
	}
	return nil
}
