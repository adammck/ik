package three

import "fmt"

type EulerAngles struct {
	Heading float64 // heading
	Pitch   float64 // pitch
	Bank    float64 // bank
}

func MakeEulerAngles(h float64, p float64, b float64) *EulerAngles {
	return &EulerAngles{h, p, b}
}

func (ea *EulerAngles) String() string {
	return fmt.Sprintf("&Euler{h=%g p=%g b=%g}", ea.Heading, ea.Pitch, ea.Bank)
}
