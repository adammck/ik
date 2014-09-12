package ik

import (
	"fmt"
	"math"
)

type EulerAngles struct {
	Heading float64 // x
	Pitch   float64 // y
	Bank    float64 // z
}

var (
	IdentityOrientation = EulerAngles{}
)

func Euler(h float64, p float64, b float64) EulerAngles {
	return EulerAngles{rad(h), rad(p), rad(b)}
}

// TODO: GTFO?
func MakeEulerAngles(h float64, p float64, b float64) *EulerAngles {
	return &EulerAngles{h, p, b}
}

func (ea EulerAngles) String() string {
	return fmt.Sprintf("&Euler{h=%+.2f° p=%+.2f° b=%+.2f°}", deg(ea.Heading), deg(ea.Pitch), deg(ea.Bank))
}

func rad(degrees float64) float64 {
	return (math.Pi / 180) * degrees
}

func deg(rads float64) float64 {
	return rads / (math.Pi / 180)
}
