package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput(inputPath string) (*context, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	ctx := context{}
	for sc.Scan() {
		line := sc.Text()
		tokens := strings.Fields(line)
		// skip empty lines
		if len(tokens) == 0 {
			continue
		}
		if len(tokens) < 2 {
			return nil, fmt.Errorf("input line too short, missing name or value: %s", line)
		}
		iValue := -1
		var (
			value float64
			err   error
		)
		for i, t := range tokens {
			if value, err = strconv.ParseFloat(strings.ReplaceAll(t, ",", "."), 64); err == nil {
				iValue = i
			}
		}
		if iValue == -1 {
			return nil, fmt.Errorf("cannot parse value: %s", line)
		}
		name := strings.ToLower(strings.Join(tokens[0:iValue], " "))
		if name == "t" {
			ctx.t = value + 273.15
		} else if name == "p" {
			ctx.p = value
		} else if comp, ok := componentsByName[name]; ok {
			// divide by 100 to convert percents into fraction
			ctx.fractions = append(ctx.fractions, componentFraction{comp, value / 100})
		} else {
			return nil, fmt.Errorf("unknown component or parameter: %s", name)
		}
	}
	ctx.initFractions()
	return &ctx, sc.Err()
}

type output struct {
	file *os.File
}

func newOutput(path string) *output {
	var (
		f   *os.File
		err error
	)
	if path == "" {
		f = os.Stdout
	} else {
		f, err = os.Create(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	return &output{f}
}

func (o *output) close() error {
	return o.file.Close()
}

func (o *output) writeKx(kx float64) {
	if _, err := fmt.Fprintf(o.file, "Смесевой параметр размера: Kx = %f м/кмоль^1/3\n", kx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *output) writeP0m(p0m float64) {
	if _, err := fmt.Fprintf(o.file, "Давление нормировки: p0m = %f\n", p0m); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *output) writeMm(mm float64) {
	if _, err := fmt.Fprintf(o.file, "Молярная масса газа: Mm = %f кг/кмоль\n", mm); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *output) writeDU(d, u []float64) {
	if len(d) != len(u) {
		panic("D and U must be of the same length")
	}
	dStr := make([]string, len(d))
	uStr := make([]string, len(u))
	for i := range d {
		dStr[i] = strconv.FormatFloat(d[i], 'f', -1, 64)
		uStr[i] = strconv.FormatFloat(u[i], 'f', -1, 64)
	}
	// get max length of formatted D strings, to pretty print table
	maxDsLen := 0
	for _, s := range dStr {
		l := len(s)
		if l > maxDsLen {
			maxDsLen = l
		}
	}
	duStr := make([]string, len(d))
	for i := range dStr {
		duStr[i] = fmt.Sprintf("%2d  |  %s  |  %s\n", i+1, dStr[i]+strings.Repeat(" ", maxDsLen-len(dStr[i])), uStr[i])
	}
	titleStr := fmt.Sprintf(" n  |  D%s  |  U\n", strings.Repeat(" ", maxDsLen-1))
	if _, err := fmt.Fprint(o.file, "Функции молярных долей компонентов:\n", titleStr, strings.Join(duStr, "")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *output) writePi(pi float64) {
	if _, err := fmt.Fprintf(o.file, "Приведенное давление: %f\n", pi); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *output) writeTau(tau float64) {
	if _, err := fmt.Fprintf(o.file, "Приведенная температура: %f\n", tau); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *output) writeInitialSigma(sigma float64) {
	if _, err := fmt.Fprintf(o.file, "Начальное приближение приведенной плотности: %f\n", sigma); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *output) writeSigmaIterations(iters []sigmaIteration) {
	n := len(iters)
	dSigmaStr := make([]string, n)
	sigmaStr := make([]string, n)
	piCalcStr := make([]string, n)
	for i := 0; i < n; i++ {
		dSigmaStr[i] = strconv.FormatFloat(iters[i].dSigma, 'f', -1, 64)
		sigmaStr[i] = strconv.FormatFloat(iters[i].sigma, 'f', -1, 64)
		piCalcStr[i] = strconv.FormatFloat(iters[i].piCalc, 'f', -1, 64)
	}
	maxLenDSigma := 0
	maxLenSigma := 0
	for i := 0; i < n; i++ {
		l := len(dSigmaStr[i])
		if l > maxLenDSigma {
			maxLenDSigma = l
		}
		l = len(sigmaStr[i])
		if l > maxLenSigma {
			maxLenSigma = l
		}
	}
	var sb strings.Builder
	sb.WriteString("Итерации расчета sigma:\n")
	fmt.Fprintf(&sb, " k  |  Δσ%s  |  σ%s  |  π расч.\n",
		strings.Repeat(" ", maxLenDSigma-2), strings.Repeat(" ", maxLenSigma-1))
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%2d  |", i+1)
		fmt.Fprint(&sb, "  ", dSigmaStr[i], strings.Repeat(" ", maxLenDSigma-len(dSigmaStr[i])+2), "|")
		fmt.Fprint(&sb, "  ", sigmaStr[i], strings.Repeat(" ", maxLenSigma-len(sigmaStr[i])+2), "|")
		fmt.Fprint(&sb, "  ", piCalcStr[i], "\n")
	}
	if _, err := fmt.Fprint(o.file, sb.String()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *output) writeP(p float64) {
	if _, err := fmt.Fprintf(o.file, "Плотность газа p = %f кг/м^3\n", p); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *output) writeZ(z float64) {
	if _, err := fmt.Fprintf(o.file, "Коэффициент сжимаемости z = %f\n", z); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
