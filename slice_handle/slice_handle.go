package slice_handle

func FindIndex[T any](arr []T, f func(ele T) bool) int {
	for index, element := range arr {
		if f(element) {
			return index
		}
	}
	return -1
}

func Where[T any](arr []T, f func(ele T) bool) []T {
	newArr := make([]T, 0)
	for _, item := range arr {
		if f(item) {
			newArr = append(newArr, item)
		}
	}
	return newArr
}

func Select[T any, V any](arr []T, f func(ele T) V) []V {
	newArr := make([]V, len(arr))
	for index, item := range arr {
		newArr[index] = f(item)
	}
	return newArr
}

func SelectOne[T any, V any](arr []T, f func(ele T) V) *V {
	if len(arr) > 0 {
		a := f(arr[0])
		return &a
	}
	return nil
}

func GroupBy[T any, V any](arr []T, fKey func(ele T) (string, V)) []struct {
	Key  string
	Data []V
} {
	g := make([]struct {
		Key  string
		Data []V
	}, 0)
	mapIndex := make(map[string]int)
	for _, element := range arr {
		key, v := fKey(element)
		if vg, exist := mapIndex[key]; exist {
			d := g[vg]
			d.Data = append(d.Data, v)
			g[vg] = d
		} else {
			newIndex := len(g)
			mapIndex[key] = newIndex
			g = append(g, struct {
				Key  string
				Data []V
			}{Key: key, Data: []V{v}})
		}
	}
	return g
}

func ToDictionary[T any, V any](arr []T, f func(arg T) (string, T)) map[string]T {
	dict := make(map[string]T)

	for _, element := range arr {
		key, v := f(element)
		dict[key] = v
	}

	return dict
}

func Sum[T any](arr []T, f func(arg T) float64) float64 {
	sum := 0.0

	for _, element := range arr {
		sum += f(element)
	}

	return sum
}

func Avg[T any](arr []T, f func(arg T) float64) float64 {
	avg := 0.0
	sum := 0.0
	if len(arr) > 0 {
		for _, element := range arr {
			sum += f(element)
		}
		avg = sum / float64(len(arr))
	}
	return avg
}

func Min[T any](arr []T, f func(arg T) float64) float64 {
	min := 0.0
	if len(arr) == 0 {
		return min
	}
	min = f(arr[0])
	for i := 1; i < len(arr); i++ {
		value := f(arr[i])
		if value < min {
			min = value
		}
	}
	return min
}

func Max[T any](arr []T, f func(arg T) float64) float64 {
	max := 0.0
	if len(arr) == 0 {
		return max
	}
	max = f(arr[0])
	for i := 1; i < len(arr); i++ {
		value := f(arr[i])
		if value > max {
			max = value
		}
	}
	return max
}

func RangeMinMax[T any](arr []T, f func(arg T) float64) (float64, float64) {
	max := 0.0
	min := 0.0
	if len(arr) == 0 {
		return min, max
	}
	max = f(arr[0])
	min = f(arr[0])
	for i := 1; i < len(arr); i++ {
		value := f(arr[i])
		if value > max {
			max = value
		} else if value < min {
			min = value
		}
	}
	return min, max
}
