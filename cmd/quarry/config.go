package main

import (
	"encoding/json"
	"fmt"
	"os"
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
		TimeToTravelEmpty:     [2]int{20, 60},
		TimeToTravelFull:      [2]int{30, 100},
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

func (cfg SimulationConfig) printConfig() string {
	out := fmt.Sprintf(
		`Current config:
	Workers:                          Stone Masses:
        A = %-4d                          A = %-4d
        B = %-4d                          B = %-4d
        C = %-4d                          C = %-4d

	Stone Extraction Times:           Stone Masses Limits:         
        A (min = %-4d max = %-4d)         1.Layer = %-4d
        B (min = %-4d max = %-4d)         2.Layer = %-4d
        C (min = %-4d max = %-4d)         3.Layer = %-4d

	Time To Travel Empty:       min = %-4d max = %-4d
	Time To Travel Full:        min = %-4d max = %-4d
	Time To Place Stone:        min = %-4d max = %-4d
	Time To Place Insulation:   min = %-4d max = %-4d
	Time To Change Pallet:      min = %-4d max = %-4d

	Quarry Workplaces = %-2d`,
		cfg.WorkersAmount[0], cfg.StonesMasses[0],
		cfg.WorkersAmount[1], cfg.StonesMasses[1],
		cfg.WorkersAmount[2], cfg.StonesMasses[2],
		cfg.StonesExtractionTime[0][0], cfg.StonesExtractionTime[0][1], cfg.StoneMassesLimits[0],
		cfg.StonesExtractionTime[1][0], cfg.StonesExtractionTime[1][1], cfg.StoneMassesLimits[1],
		cfg.StonesExtractionTime[2][0], cfg.StonesExtractionTime[2][1], cfg.StoneMassesLimits[2],
		cfg.TimeToTravelEmpty[0], cfg.TimeToTravelEmpty[1],
		cfg.TimeToTravelFull[0], cfg.TimeToTravelFull[1],
		cfg.TimeToPlaceStone[0], cfg.TimeToPlaceStone[1],
		cfg.TimeToPlaceInsulation[0], cfg.TimeToPlaceInsulation[1],
		cfg.TimeToChangePallet[0], cfg.TimeToChangePallet[1],
		cfg.QuarryWorkplaces)
	return out
}

func (cfg SimulationConfig) saveConfig() error {
	file, err := os.Create("./config.json")
	if err != nil {
		fmt.Println("Couldn't create config file")
		return err
	}
	var str []byte
	str, err = json.Marshal(cfg)
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
