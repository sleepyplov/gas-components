package main

import (
	"math"
	"testing"
)

type gasTestCase struct {
	// температура
	t float64
	// давление
	p float64
	// ожидаемая вязкость, не путать с p - давление
	pe float64
	// ожидаемый коэффициент сжимаемости
	ze float64
}

var n1TestCases = []gasTestCase{
	{
		t:  250,
		p:  0.1,
		pe: 0.8112,
		ze: 0.9966,
	},
	{
		t:  300,
		p:  0.1,
		pe: 0.6749,
		ze: 0.9982,
	},
	{
		t:  350,
		p:  0.1,
		pe: 0.5780,
		ze: 0.9990,
	},
	{
		t:  250,
		p:  5,
		pe: 49.295,
		ze: 0.8200,
	},
	{
		t:  300,
		p:  5,
		pe: 36.949,
		ze: 0.9116,
	},
	{
		t:  350,
		p:  5,
		pe: 30.253,
		ze: 0.9543,
	},
	{
		t:  250,
		p:  15,
		pe: 196.15,
		ze: 0.6182,
	},
	{
		t:  300,
		p:  15,
		pe: 125.53,
		ze: 0.8050,
	},
	{
		t:  350,
		p:  15,
		pe: 95.519,
		ze: 0.9068,
	},
	{
		t:  250,
		p:  30,
		pe: 285.18,
		ze: 0.8504,
	},
	{
		t:  300,
		p:  30,
		pe: 223.21,
		ze: 0.9054,
	},
	{
		t:  350,
		p:  30,
		pe: 178.53,
		ze: 0.9703,
	},
}

var n2TestCases = []gasTestCase{
	{
		t:  250,
		p:  0.1,
		pe: 0.9577,
		ze: 0.9963,
	},
	{
		t:  300,
		p:  0.1,
		pe: 0.7967,
		ze: 0.9980,
	},
	{
		t:  350,
		p:  0.1,
		pe: 0.6823,
		ze: 0.9989,
	},
	{
		t:  250,
		p:  5,
		pe: 59.396,
		ze: 0.8032,
	},
	{
		t:  300,
		p:  5,
		pe: 43.980,
		ze: 0.9039,
	},
	{
		t:  350,
		p:  5,
		pe: 35.869,
		ze: 0.9500,
	},
	{
		t:  250,
		p:  15,
		pe: 241.91,
		ze: 0.5916,
	},
	{
		t:  300,
		p:  15,
		pe: 151.67,
		ze: 0.7864,
	},
	{
		t:  350,
		p:  15,
		pe: 114.10,
		ze: 0.8960,
	},
	{
		t:  250,
		p:  30,
		pe: 342.04,
		ze: 0.8369,
	},
	{
		t:  300,
		p:  30,
		pe: 267.56,
		ze: 0.8915,
	},
	{
		t:  350,
		p:  30,
		pe: 213.16,
		ze: 0.9592,
	},
}

var n3TestCases = []gasTestCase{
	{
		t:  250,
		p:  0.1,
		pe: 0.7454,
		ze: 0.9972,
	},
	{
		t:  300,
		p:  0.1,
		pe: 0.6203,
		ze: 0.9986,
	},
	{
		t:  350,
		p:  0.1,
		pe: 0.5313,
		ze: 0.9993,
	},
	{
		t:  250,
		p:  5,
		pe: 43.206,
		ze: 0.8602,
	},
	{
		t:  300,
		p:  5,
		pe: 33.217,
		ze: 0.9324,
	},
	{
		t:  350,
		p:  5,
		pe: 27.454,
		ze: 0.9670,
	},
	{
		t:  250,
		p:  15,
		pe: 158.30,
		ze: 0.7044,
	},
	{
		t:  300,
		p:  15,
		pe: 108.18,
		ze: 0.8589,
	},
	{
		t:  350,
		p:  15,
		pe: 84.803,
		ze: 0.9391,
	},
	{
		t:  250,
		p:  30,
		pe: 253.14,
		ze: 0.8809,
	},
	{
		t:  300,
		p:  30,
		pe: 196.78,
		ze: 0.9443,
	},
	{
		t:  350,
		p:  30,
		pe: 158.80,
		ze: 1.0030,
	},
}

func getPZ(ctx *context) (p float64, z float64) {
	kx := getKx(ctx)
	p0m := getP0m(kx)
	mm := getMm(ctx)
	d, u := getDU(ctx, kx)
	sigma := getInitialSigma(ctx, kx, d, u)
	pi := getPi(ctx, p0m)
	tau := getTau(ctx)
	iters := getSigma(ctx, kx, pi, tau, sigma, d, u)
	sigma = iters[len(iters)-1].sigma
	p = getP(ctx, kx, mm, sigma)
	z = getZ(ctx, sigma, tau, d, u)
	return
}

// returns the nearest number with the specified number of fraction digits
//
// e.g. roundDecimals(10.0267, 2) = 10.03
func roundDecimals(value float64, fractions int) float64 {
	return math.Round(value*math.Pow10(fractions)) / math.Pow10(fractions)
}

func almostEqual(a, b, threshold float64) bool {
	return math.Abs(a-b) <= threshold
}

func TestGasN1(t *testing.T) {
	ctx := &context{}
	ctx.fractions = []componentFraction{
		{&methane, 0.965},
		{&ethane, 0.018},
		{&propane, 0.0045},
		{&iButane, 0.001},
		{&nButane, 0.001},
		{&iPentane, 0.0005},
		{&nPentane, 0.0003},
		{&nHexane, 0.0007},
		{&nitrogen, 0.003},
		{&carbonDioxide, 0.006},
	}
	ctx.initFractions()
	for _, tc := range n1TestCases {
		ctx.t = tc.t
		ctx.p = tc.p
		p, z := getPZ(ctx)
		p, z = roundDecimals(p, 4), roundDecimals(z, 4)
		// t.Log(math.Abs(p-tc.pe) <= 0.1, z, tc.ze)
		if !almostEqual(p, tc.pe, 0.1) {
			t.Errorf("Wrong p for N1; actual = %f, expected = %f", p, tc.pe)
		}
		if !almostEqual(z, tc.ze, 0.0001) {
			t.Errorf("Wrong z for N1; actual = %f, expected = %f", z, tc.ze)
		}
	}
}

func TestGasN2(t *testing.T) {
	ctx := context{}
	ctx.fractions = []componentFraction{
		{&methane, 0.812},
		{&ethane, 0.043},
		{&propane, 0.009},
		{&iButane, 0.0015},
		{&nButane, 0.0015},
		{&nitrogen, 0.057},
		{&carbonDioxide, 0.076},
	}
	ctx.initFractions()
	for _, tc := range n2TestCases {
		ctx.t = tc.t
		ctx.p = tc.p
		p, z := getPZ(&ctx)
		p, z = roundDecimals(p, 4), roundDecimals(z, 4)
		if !almostEqual(p, tc.pe, 0.1) {
			t.Errorf("Wrong p for N2; actual = %f, expected = %f", p, tc.pe)
		}
		if !almostEqual(z, tc.ze, 0.0001) {
			t.Errorf("Wrong z for N2; actual = %f, expected = %f", z, tc.ze)
		}
	}
}

func TestGasN3(t *testing.T) {
	ctx := context{}
	ctx.fractions = []componentFraction{
		{&methane, 0.8641},
		{&ethane, 0.018},
		{&propane, 0.0045},
		{&iButane, 0.001},
		{&nButane, 0.001},
		{&iPentane, 0.0003},
		{&nPentane, 0.0005},
		{&nHexane, 0.0012},
		{&nitrogen, 0.0034},
		{&carbonDioxide, 0.006},
		{&helium, 0.005},
		{&hydrogen, 0.095},
	}
	ctx.initFractions()
	for _, tc := range n3TestCases {
		ctx.t = tc.t
		ctx.p = tc.p
		p, z := getPZ(&ctx)
		p, z = roundDecimals(p, 4), roundDecimals(z, 4)
		if !almostEqual(p, tc.pe, 0.1) {
			t.Errorf("Wrong p for N3; actual = %f, expected = %f", p, tc.pe)
		}
		if !almostEqual(z, tc.ze, 0.0001) {
			t.Errorf("Wrong z for N3; actual = %f, expected = %f", z, tc.ze)
		}
	}
}
