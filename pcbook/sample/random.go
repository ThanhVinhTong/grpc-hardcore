package sample

import (
	"math/rand"
	"pcbook/pb"

	"github.com/google/uuid"
)

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "HP", "Lenovo", "Asus", "Acer")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("MacBook Air", "MacBook Pro")
	case "Dell":
		return randomStringFromSet("Inspiron", "XPS", "Alienware")
	case "HP":
		return randomStringFromSet("Spectre", "Pavilion", "Envy")
	case "Lenovo":
		return randomStringFromSet("ThinkPad", "Yoga", "Legion")
	case "Asus":
		return randomStringFromSet("ROG", "Predator", "Zenbook")
	case "Acer":
		return randomStringFromSet("Aspire", "Predator", "Chromebook")
	default:
		return "Unknown"
	}
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet("Core i3", "Core i5", "Core i7", "Core i9")
	}
	return randomStringFromSet("Ryzen 3", "Ryzen 5", "Ryzen 7", "Ryzen 9")
}

func randomGPUBrand() string {
	return randomStringFromSet("NVIDIA", "AMD")
}

func randomGPUName(brand string) string {
	if brand == "NVIDIA" {
		return randomStringFromSet("GeForce RTX 5050", "GeForce RTX 5060", "GeForce RTX 5070", "GeForce RTX 5080")
	}
	return randomStringFromSet("Radeon RX 6500", "Radeon RX 6600", "Radeon RX 6700", "Radeon RX 6800")
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomPanel() pb.Screen_Panel {
	switch rand.Intn(3) {
	case 1:
		return pb.Screen_IPS
	case 2:
		return pb.Screen_OLED
	default:
		return pb.Screen_UNKNOWN
	}
}

func randomScreenResolution() *pb.Screen_Resolution {
	height := uint32(randomInt(1080, 2160))
	width := height * 16 / 9
	return &pb.Screen_Resolution{
		Width:  uint32(width),
		Height: uint32(height),
	}
}

func randomID() string {
	UUID := uuid.New()
	return UUID.String()
}

func randomStringFromSet(set ...string) string {
	n := len(set)
	if n == 0 {
		return ""
	}
	return set[rand.Intn(n)]
}

func randomBool() bool {
	return rand.Intn(2) == 1
}
