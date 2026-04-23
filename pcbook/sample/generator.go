package sample

import (
	"pcbook/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewKeyboard returns a sample keyboard
func NewKeyboard() *pb.Keyboard {
	return &pb.Keyboard{
		Layout:     randomKeyboardLayout(),
		HasBacklit: randomBool(),
	}
}

// NewCPU returns a sample CPU
func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)
	cores := uint32(randomInt(2, 16))
	threads := uint32(randomInt(int(cores), int(cores*2)))
	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)

	cpu := &pb.CPU{
		Brand:      brand,
		Name:       name,
		NumCores:   cores,
		NumThreads: threads,
		MinGhz:     minGhz,
		MaxGhz:     maxGhz,
	}

	return cpu
}

// NewGPU returns a sample GPU
func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)
	minGhz := randomFloat64(1.0, 2.0)
	maxGhz := randomFloat64(minGhz, 4.0)
	memory := &pb.Memory{
		Value: uint64(randomInt(2, 16)),
		Unit:  pb.Memory_GIGABYTE,
	}

	gpu := &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}

	return gpu
}

// NewRAM returns a sample RAM
func NewRAM() *pb.Memory {
	return &pb.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit:  pb.Memory_GIGABYTE,
	}
}

// NewSSD returns a sample SSD
func NewSSD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(256, 2048)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
}

// NewHDD returns a sample HDD
func NewHDD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 4)),
			Unit:  pb.Memory_TERABYTE,
		},
	}
}

// NewScreen returns a sample screen
func NewScreen() *pb.Screen {
	return &pb.Screen{
		SizeInch:   randomFloat32(13.0, 17.0),
		Panel:      randomPanel(),
		Resolution: randomScreenResolution(),
		Multitouch: randomBool(),
	}
}

func NewLaptop() *pb.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	laptop := &pb.Laptop{
		Id:       randomID(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCPU(),
		Ram:      NewRAM(),
		Gpus:     []*pb.GPU{NewGPU(), NewGPU()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd:    randomFloat64(500.0, 3000.0),
		ReleaseYear: uint32(randomInt(2015, 2025)),
		UpdateAt:    timestamppb.Now(),
	}

	return laptop
}

func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
