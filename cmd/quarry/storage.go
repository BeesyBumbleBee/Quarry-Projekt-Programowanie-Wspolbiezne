package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"sync"
	"time"
)

type Storage struct {
	cfg              SimulationConfig
	program          *tea.Program
	massLimit        [3]int
	cells            [3][3]int
	horizontalBars   [2][3]bool
	verticalBars     [3][2]bool
	level            int
	totalMass        int
	placeMutex       sync.Mutex
	replaceMutex     sync.Mutex
	hasPalletCleared *sync.Cond
	palletTime       time.Time
}

type storageWorkerPlacing struct {
	workerId string
}

type storageWorkerCantPlace struct {
	workerId string
}

type storageWorkerFinishedPlacing struct {
	workerId string
}
type palletFullMsg struct {
	timeTook time.Duration
}
type placeInsulationMsg struct{}

func InitStorage(cfg SimulationConfig, program *tea.Program) *Storage {
	return &Storage{
		cfg:              cfg,
		program:          program,
		massLimit:        [3]int{cfg.StoneMassesLimits[0], cfg.StoneMassesLimits[1], cfg.StoneMassesLimits[2]},
		cells:            [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		horizontalBars:   [2][3]bool{{false, false, false}, {false, false, false}},
		verticalBars:     [3][2]bool{{false, false}, {false, false}, {false, false}},
		level:            0,
		totalMass:        0,
		placeMutex:       sync.Mutex{},
		replaceMutex:     sync.Mutex{},
		hasPalletCleared: sync.NewCond(&sync.Mutex{}),
		palletTime:       time.Now(),
	}
}

func (s *Storage) Place(stoneType int, workerId string, delay func()) bool {
	s.replaceMutex.Lock()
	s.placeMutex.Lock()
	s.program.Send(storageWorkerPlacing{workerId})

	stoneMass := s.cfg.StonesMasses[stoneType]
	stoneSize := stoneType + 1

	x, y, dir, can := s.findPlace(stoneMass, stoneSize)
	if !can {
		s.placeMutex.Unlock()
		s.replaceMutex.Unlock()
		s.program.Send(storageWorkerCantPlace{workerId})
		return false
	}

	delay()

	s.totalMass += stoneMass
	for i := 0; i < stoneSize; i++ {
		s.cells[y][x] = stoneSize
		if dir > 0 {
			x += 1
			if i < stoneSize-1 {
				s.verticalBars[y][x-1] = true
			}
		} else {
			y += 1
			if i < stoneSize-1 {
				s.horizontalBars[y-1][x] = true
			}
		}
	}

	if s.checkFullLevel() {
		if s.level == 2 {
			s.placeMutex.Unlock()
			s.program.Send(storageWorkerFinishedPlacing{workerId})
			go s.changePallet()
			return true
		}
		s.placeMutex.Unlock()
		s.program.Send(storageWorkerFinishedPlacing{workerId})
		go s.placeInsulation()
		return true
	}

	s.placeMutex.Unlock()
	s.program.Send(storageWorkerFinishedPlacing{workerId})
	s.replaceMutex.Unlock()
	return true
}

func (s *Storage) checkMass(stoneMass int) bool {
	if s.totalMass+stoneMass > s.massLimit[s.level] {
		return false
	}
	return true
}

func (s *Storage) findPlace(stoneMass int, stoneSize int) (int, int, int, bool) {
	x, y, dir := 0, 0, 0
	palletCopy := s.cells

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if s.cells[row][col] != 0 {
				continue
			}
			dir = 1
			x = col
			y = row
			if !(s.canPlace(palletCopy, stoneSize, stoneMass, x, y, dir) && s.totalMass+stoneMass <= s.massLimit[s.level]) {
				dir = -1
				if !(s.canPlace(palletCopy, stoneSize, stoneMass, x, y, dir) && s.totalMass+stoneMass <= s.massLimit[s.level]) {
					continue
				}
			}
			if s.hasFillPotential(palletCopy, stoneSize, stoneMass, s.massLimit[s.level]) {
				return x, y, dir, true
			}
		}
	}
	return 0, 0, 0, false
}

func (s *Storage) canPlace(pallet [3][3]int, stoneSize int, stoneMass int, col int, row int, dir int) bool {
	if dir >= 1 {
		if col+stoneSize > 3 {
			return false
		}
		for dx := 0; dx < stoneSize; dx++ {
			if pallet[row][col+dx] != 0 {
				return false
			}
		}
	} else {
		if row+stoneSize > 3 {
			return false
		}
		for dy := 0; dy < stoneSize; dy++ {
			if pallet[row+dy][col] != 0 {
				return false
			}
		}
	}
	return true
}

func (s *Storage) hasFillPotential(pallet [3][3]int, stoneSize int, stoneMass int, massLimit int) bool {
	emptyCells := -stoneSize
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if pallet[y][x] == 0 {
				emptyCells++
			}
		}
	}

	if emptyCells == 0 {
		return true
	}
	remainingMass := massLimit - stoneMass - s.totalMass
	for stoneType := range s.cfg.StonesMasses {
		if stoneType+1 > emptyCells {
			continue
		}
		if remainingMass >= s.cfg.StonesMasses[stoneType]*(emptyCells/(stoneType+1)) {
			return true
		}
	}
	return false
}

func (s *Storage) checkFullLevel() bool {
	for _, row := range s.cells {
		for _, cell := range row {
			if cell == 0 {
				return false
			}
		}
	}
	return true
}

func (s *Storage) changePallet() {
	endTime := time.Now()
	s.program.Send(palletFullMsg{
		timeTook: endTime.Sub(s.palletTime),
	})
	time.Sleep(time.Duration(s.cfg.TimeToChangePallet[0]+
		(s.cfg.TimeToChangePallet[1]-s.cfg.TimeToChangePallet[0])) * time.Millisecond)
	s.clearPallet()
	s.level = 0
	s.palletTime = time.Now()
	s.replaceMutex.Unlock()
}

func (s *Storage) placeInsulation() {
	s.program.Send(placeInsulationMsg{})
	time.Sleep(time.Duration(s.cfg.TimeToPlaceInsulation[0]+
		(s.cfg.TimeToPlaceInsulation[1]-s.cfg.TimeToPlaceInsulation[0])) * time.Millisecond)
	s.clearPallet()
	s.level++

	s.replaceMutex.Unlock()
}

func (s *Storage) clearPallet() {
	s.totalMass = 0
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			s.cells[y][x] = 0
			if y != 2 {
				s.horizontalBars[y][x] = false
			}
			if x != 2 {
				s.verticalBars[y][x] = false
			}
		}
	}
	s.hasPalletCleared.Broadcast()
}
