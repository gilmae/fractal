package fractal

import (
	"encoding/json"
	"fmt"
	"os"
)

func (p *Point) MarshalJSON() ([]byte, error) {
	type Alias Point

	return json.Marshal(&struct {
		X      int
		Y      int
		Escape float64
	}{
		X:      p.X,
		Y:      p.Y,
		Escape: p.Escape,
	})
}

func Write_Raw(points map[Key]CalculatedPoint, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range points {
		json.NewEncoder(file).Encode(&v)
	}

	file.Close()
}
