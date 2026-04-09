package util

import (
	"pcbook/pb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DeepCopyLaptop performs a complete manual deep copy of the Laptop object.
// While proto.Clone() handles generic deep copying, a manual copy can be
// noticeably faster by avoiding the heavy use of reflection internally.
func DeepCopyLaptop(laptop *pb.Laptop) *pb.Laptop {
	if laptop == nil {
		return nil
	}

	other := &pb.Laptop{
		Id:          laptop.Id,
		Brand:       laptop.Brand,
		Name:        laptop.Name,
		PriceUsd:    laptop.PriceUsd,
		ReleaseYear: laptop.ReleaseYear,
	}

	if laptop.Cpu != nil {
		other.Cpu = &pb.CPU{
			Brand:      laptop.Cpu.Brand,
			Name:       laptop.Cpu.Name,
			NumCores:   laptop.Cpu.NumCores,
			NumThreads: laptop.Cpu.NumThreads,
			MinGhz:     laptop.Cpu.MinGhz,
			MaxGhz:     laptop.Cpu.MaxGhz,
		}
	}

	if laptop.Ram != nil {
		other.Ram = &pb.Memory{
			Value: laptop.Ram.Value,
			Unit:  laptop.Ram.Unit,
		}
	}

	if laptop.Gpus != nil {
		other.Gpus = make([]*pb.GPU, len(laptop.Gpus))
		for i, gpu := range laptop.Gpus {
			if gpu != nil {
				otherGpu := &pb.GPU{
					Brand:  gpu.Brand,
					Name:   gpu.Name,
					MinGhz: gpu.MinGhz,
					MaxGhz: gpu.MaxGhz,
				}
				if gpu.Memory != nil {
					otherGpu.Memory = &pb.Memory{
						Value: gpu.Memory.Value,
						Unit:  gpu.Memory.Unit,
					}
				}
				other.Gpus[i] = otherGpu
			}
		}
	}

	if laptop.Storages != nil {
		other.Storages = make([]*pb.Storage, len(laptop.Storages))
		for i, st := range laptop.Storages {
			if st != nil {
				otherSt := &pb.Storage{
					Driver: st.Driver,
				}
				if st.Memory != nil {
					otherSt.Memory = &pb.Memory{
						Value: st.Memory.Value,
						Unit:  st.Memory.Unit,
					}
				}
				other.Storages[i] = otherSt
			}
		}
	}

	if laptop.Screen != nil {
		other.Screen = &pb.Screen{
			SizeInch:   laptop.Screen.SizeInch,
			Panel:      laptop.Screen.Panel,
			Multitouch: laptop.Screen.Multitouch,
		}
		if laptop.Screen.Resolution != nil {
			other.Screen.Resolution = &pb.Screen_Resolution{
				Width:  laptop.Screen.Resolution.Width,
				Height: laptop.Screen.Resolution.Height,
			}
		}
	}

	if laptop.Keyboard != nil {
		other.Keyboard = &pb.Keyboard{
			Layout:     laptop.Keyboard.Layout,
			HasBacklit: laptop.Keyboard.HasBacklit,
		}
	}

	if laptop.Weight != nil {
		switch w := laptop.Weight.(type) {
		case *pb.Laptop_WeightKg:
			other.Weight = &pb.Laptop_WeightKg{
				WeightKg: w.WeightKg,
			}
		case *pb.Laptop_WeightLb:
			other.Weight = &pb.Laptop_WeightLb{
				WeightLb: w.WeightLb,
			}
		}
	}

	if laptop.UpdateAt != nil {
		other.UpdateAt = proto.Clone(laptop.UpdateAt).(*timestamppb.Timestamp)
	}

	return other
}
