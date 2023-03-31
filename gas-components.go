package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	inputPath := flag.String("i", "", "путь к файлу с исходными данными")
	outputPath := flag.String("o", "", "путь к файлу для вывода. Необязательно, по умолчанию используется стандартный поток вывода")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		var sb strings.Builder
		sb.WriteString("\nФормат исходного файла:\nКаждая строка состоит из имени компонента или параметра и его значения, разделенных пробелом.\nНапример: Метан 89,8211\n")
		sb.WriteString("Названия параметров и компонентов:\n")
		for _, c := range allComponents {
			fmt.Fprintf(&sb, "\t%s\n", c.name)
		}
		sb.WriteString("\n\tt - температура в °С\n\tp - давление в МПа\n\n")
		sb.WriteString("Доли компонентов указываются в процентах.\nБольшие/маленькие буквы, точка или запятая в дробях - без разницы.\n\n")
		sb.WriteString("Gas Components - made by Sleepy Plov with ♥\n")
		fmt.Fprint(flag.CommandLine.Output(), sb.String())
	}
	flag.Parse()
	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "missing input path")
		flag.Usage()
		os.Exit(1)
	}
	ctx, err := readInput(*inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out := newOutput(*outputPath)
	defer out.close()
	kx := getKx(ctx)
	out.writeKx(kx)
	p0m := getP0m(kx)
	out.writeP0m(p0m)
	mm := getMm(ctx)
	out.writeMm(mm)
	d, u := getDU(ctx, kx)
	out.writeDU(d, u)
	sigma := getInitialSigma(ctx, kx, d, u)
	out.writeInitialSigma(sigma)
	pi := getPi(ctx, p0m)
	out.writePi(pi)
	tau := getTau(ctx)
	out.writeTau(tau)
	iters := getSigma(ctx, kx, pi, tau, sigma, d, u)
	out.writeSigmaIterations(iters)
	sigma = iters[len(iters)-1].sigma
	p := getP(ctx, kx, mm, sigma)
	out.writeP(p)
	z := getZ(ctx, sigma, tau, d, u)
	out.writeZ(z)
}
