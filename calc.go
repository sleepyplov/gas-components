package main

import "math"

// сумма от i=start до end
func sum(start, end int32, getEl func(i int32) float64) float64 {
	s := 0.0
	for i := start; i < end; i++ {
		s += getEl(i)
	}
	return s
}

// смесевой параметр размера
func getKx(c *context) float64 {
	n := c.length()
	return math.Pow(
		math.Pow(
			sum(0, n, func(i int32) float64 { return c.fraction(i) * math.Pow(c.component(i).k, 5.0/2) }),
			2,
		)+2*sum(0, n-1, func(i int32) float64 {
			return sum(i+1, n, func(j int32) float64 {
				return c.fraction(i) * c.fraction(j) *
					(math.Pow(c.getKBin(i, j), 5) - 1) *
					math.Pow(c.component(i).k*c.component(j).k, 5.0/2)
			})
		}),
		1.0/5,
	)
}

// давление нормировки
func getP0m(Kx float64) float64 {
	return math.Pow10(-3) * math.Pow(Kx, -3) * R * Lt
}

// молярная масса
func getMm(c *context) float64 {
	n := c.length()
	return sum(0, n, func(i int32) float64 {
		return c.fraction(i) * c.component(i).m
	})
}

// функции молярных долей компонентов Dn и Un
func getDU(ctx *context, kx float64) (d []float64, u []float64) {
	nc := ctx.length()
	d = make([]float64, 58)
	u = make([]float64, 58)
	g := sum(0, nc, func(i int32) float64 {
		return ctx.fraction(i) * ctx.component(i).g
	}) + sum(0, nc-1, func(i int32) float64 {
		return sum(i+1, nc, func(j int32) float64 {
			return ctx.fraction(i) * ctx.fraction(j) * (ctx.gBin(i, j) - 1) * (ctx.component(i).g + ctx.component(j).g)
		})
	})
	q := sum(0, nc, func(i int32) float64 {
		return ctx.fraction(i) * ctx.component(i).q
	})
	f := sum(0, nc, func(i int32) float64 {
		return math.Pow(ctx.fraction(i), 2) * ctx.component(i).f
	})
	v := math.Pow(
		math.Pow(
			sum(0, nc, func(i int32) float64 {
				return ctx.fraction(i) * math.Pow(ctx.component(i).e, 5.0/2)
			}),
			2,
		)+2*sum(0, nc-1, func(i int32) float64 {
			return sum(i+1, nc, func(j int32) float64 {
				return ctx.fraction(i) * ctx.fraction(j) *
					(math.Pow(ctx.vBin(i, j), 5) - 1) *
					math.Pow(ctx.component(i).e*ctx.component(j).e, 5.0/2)
			})
		}),
		1.0/5,
	)
	e := make([][]float64, nc)
	gBin := make([][]float64, nc)
	for i := int32(0); i < nc; i++ {
		e[i] = make([]float64, nc)
		gBin[i] = make([]float64, nc)
		for j := int32(0); j < nc; j++ {
			e[i][j] = ctx.eBin(i, j) * math.Sqrt(ctx.component(i).e*ctx.component(j).e)
			gBin[i][j] = ctx.gBin(i, j) * (ctx.component(i).g + ctx.component(j).g) / 2
		}
	}
	for n := 0; n < 58; n++ {
		bBin := make([][]float64, nc)
		for i := int32(0); i < nc; i++ {
			bBin[i] = make([]float64, nc)
			for j := int32(0); j < nc; j++ {
				bBin[i][j] = math.Pow(gBin[i][j]+1-noDimensions.g[n], noDimensions.g[n]) *
					math.Pow(ctx.component(i).q*ctx.component(j).q+1-noDimensions.q[n], noDimensions.q[n]) *
					math.Pow(math.Sqrt(ctx.component(i).f*ctx.component(j).f)+1-noDimensions.f[n], noDimensions.f[n]) *
					math.Pow(ctx.component(i).s*ctx.component(j).s+1-noDimensions.s[n], noDimensions.s[n]) *
					math.Pow(ctx.component(i).w*ctx.component(j).w+1-noDimensions.w[n], noDimensions.w[n])
			}
		}
		b := sum(0, nc, func(i int32) float64 {
			return sum(0, nc, func(j int32) float64 {
				return ctx.fraction(i) * ctx.fraction(j) * bBin[i][j] *
					math.Pow(e[i][j], noDimensions.u[n]) *
					math.Pow(ctx.component(i).k*ctx.component(j).k, 3.0/2)
			})
		})
		c := math.Pow(g+1-noDimensions.g[n], noDimensions.g[n]) *
			math.Pow(q*q+1-noDimensions.q[n], noDimensions.q[n]) *
			math.Pow(f+1-noDimensions.f[n], noDimensions.f[n]) *
			math.Pow(v, noDimensions.u[n])

		if n <= 11 {
			d[n] = b * math.Pow(kx, -3)
			u[n] = 0
		} else if n <= 17 {
			d[n] = b*math.Pow(kx, -3) - c
			u[n] = c
		} else if n <= 57 {
			d[n] = 0
			u[n] = c
		} else {
			panic("n must be in range [0, 57]")
		}
	}
	return d, u
}

// безразмерный комплекс А0
func getA0(ctx *context, sigma float64, tau float64, d, u []float64) float64 {
	return sum(0, 58, func(n int32) float64 {
		return noDimensions.a[n] *
			math.Pow(sigma, noDimensions.b[n]) *
			math.Pow(tau, -noDimensions.u[n]) *
			(noDimensions.b[n]*d[n] +
				(noDimensions.b[n]-noDimensions.c[n]*noDimensions.k[n]*math.Pow(sigma, noDimensions.k[n]))*
					u[n]*math.Exp(-noDimensions.c[n]*math.Pow(sigma, noDimensions.k[n])))
	})
}

// безразмерный комплекс А1
func getA1(ctx *context, sigma float64, tau float64, d, u []float64) float64 {
	return sum(0, 58, func(n int32) float64 {
		return noDimensions.a[n] *
			math.Pow(sigma, noDimensions.b[n]) *
			math.Pow(tau, -noDimensions.u[n]) * ((noDimensions.b[n]+1)*noDimensions.b[n]*d[n] +
			((noDimensions.b[n]-noDimensions.c[n]*noDimensions.k[n]*math.Pow(sigma, noDimensions.k[n]))*
				(noDimensions.b[n]-noDimensions.c[n]*noDimensions.k[n]*math.Pow(sigma, noDimensions.k[n])+1)-
				noDimensions.c[n]*
					math.Pow(noDimensions.k[n], 2)*
					math.Pow(sigma, noDimensions.k[n]))*
				u[n]*math.Exp(-noDimensions.c[n]*math.Pow(sigma, noDimensions.k[n])))
	})
}

// безразмерный комплекс А2
func getA2(ctx *context, sigma, tau float64, d, u [58]float64) float64 {
	return sum(0, 58, func(n int32) float64 {
		return noDimensions.a[n] * math.Pow(sigma, noDimensions.b[n]) * math.Pow(tau, -noDimensions.u[n]) *
			(1 - noDimensions.u[n]) * (noDimensions.b[n]*d[n] +
			(noDimensions.b[n]-noDimensions.c[n]*noDimensions.k[n]*math.Pow(sigma, noDimensions.k[n]))*
				u[n]*math.Exp(-noDimensions.c[n]*math.Pow(sigma, noDimensions.k[n])))
	})
}

// безразмерный комплекс А3
func getA3(ctx *context, sigma, tau float64, d, u [58]float64) float64 {
	return sum(0, 58, func(n int32) float64 {
		return noDimensions.a[n] * math.Pow(sigma, noDimensions.b[n]) * math.Pow(tau, -noDimensions.u[n]) *
			(1 - noDimensions.u[n]) * (d[n] + u[n]*math.Exp(-noDimensions.c[n]*math.Pow(sigma, noDimensions.k[n])))
	})
}

// итерация расчета приведенной плотности
type sigmaIteration struct {
	// приведенная плотность на данном шаге итерации
	sigma float64
	// прибавление к приведенной плотности относительно предыдущей итерации
	dSigma float64
	// расчетное приведенное давление
	piCalc float64
}

// приведенное давление
func getPi(ctx *context, p0m float64) float64 {
	return ctx.p / p0m
}

// приведенная температура
func getTau(ctx *context) float64 {
	return ctx.t / Lt
}

// расчет приведенной плотности в итерационном процессе
func getSigma(ctx *context, kx, pi, tau, initialSigma float64, d, u []float64) []sigmaIteration {
	sigma := initialSigma
	var iters []sigmaIteration
	for {
		dSigma := (pi/tau - (1+getA0(ctx, sigma, tau, d, u))*sigma) / (1 + getA1(ctx, sigma, tau, d, u))
		sigma += dSigma
		piCalc := sigma * tau * (1 + getA0(ctx, sigma, tau, d, u))
		end := math.Abs((piCalc-pi)/pi) < math.Pow10(-6)
		iters = append(iters, sigmaIteration{sigma, dSigma, piCalc})
		if end {
			break
		}
	}
	return iters
}

// начальное приближение приведенной плотности
func getInitialSigma(ctx *context, kx float64, d, u []float64) float64 {
	return math.Pow10(3) * ctx.p * math.Pow(kx, 3) / (R * ctx.t)
}

// плотность смеси
func getP(ctx *context, kx, mm, sigma float64) float64 {
	return mm * math.Pow(kx, -3) * sigma
}

// коэффициент сжимаемости
func getZ(ctx *context, sigma, tau float64, d, u []float64) float64 {
	return 1 + getA0(ctx, sigma, tau, d, u)
}

// безразмерная изобарная теплоемкость природного газа в идеально-газовом состоянии
func getCp0r(ctx *context) float64 {
	n := ctx.length()
	// tau^-1
	theta := 1 / (ctx.p / Lt)
	return sum(0, n, func(i int32) float64 {
		c := ctx.component(i)
		cp0ri := c.b0 +
			c.c0*math.Pow(c.d0*theta/math.Sinh(c.d0*theta), 2) +
			c.e0*math.Pow(c.f0*theta/math.Cosh(c.f0*theta), 2) +
			c.g0*math.Pow(c.h0*theta/math.Sinh(c.h0*theta), 2) +
			c.i0*math.Pow(c.j0*theta/math.Cosh(c.j0*theta), 2)
		return ctx.fraction(i) * cp0ri
	})
}

// расчет адиабаты
func getK(ctx *context, sigma, tau, a1, a2, a3, cp0r float64, d, u []float64) float64 {
	z := getZ(ctx, sigma, tau, d, u)
	return (1 + a1 + math.Pow(1+a2, 2)/(cp0r-1+a3)) / z
}

// расчет скорости звука
func getU(ctx *context, mm, a1, a2, a3, cp0r float64) float64 {
	return math.Pow(
		math.Pow10(3)*R*(1/mm)*ctx.t*
			(1+a1+math.Pow(1+a2, 2)/(cp0r-1+a3)),
		0.5)
}

// расчет молярной плотности
func getPMol(p, mm float64) float64 {
	return p / mm
}

// псевдокритическая молярная плотность
func getPMolPc(ctx *context) float64 {
	n := ctx.length()
	return 1 / (0.125 * sum(0, n, func(i int32) float64 {
		return sum(0, n, func(j int32) float64 {
			ci := ctx.component(i)
			cj := ctx.component(j)
			return ctx.fraction(i) * ctx.fraction(j) * math.Pow(
				math.Pow(ci.m/ci.pcr, 1.0/3)+math.Pow(cj.m/cj.pcr, 1.0/3), 3)
		})
	}))
}

// расчет псевдокритической температуры
func getTpc(ctx *context, pMolPc float64) float64 {
	n := ctx.length()
	return 0.125 * pMolPc * sum(0, n, func(i int32) float64 {
		return sum(0, n, func(j int32) float64 {
			ci := ctx.component(i)
			cj := ctx.component(j)
			return ctx.fraction(i) * ctx.fraction(j) * math.Pow(
				math.Pow(ci.m/ci.pcr, 1.0/3)+math.Pow(cj.m/ci.pcr, 1.0/3), 3) *
				math.Pow(ci.tcr*cj.tcr, 1.0/2)
		})
	})
}

// расчет приведенной плотности
func getOmegaM(pMol, pMolPc float64) float64 {
	return pMol / pMolPc
}

// расчет приведенной температуры
func getTauM(ctx *context, tpc float64) float64 {
	return ctx.t / tpc
}

// расчет параметров преобразований для приведенных значений плотности и температуры
func getPhi(ctx *context) (phi [6]float64) {
	sigma := [6]float64{1, 1, 0, 1, 0, 1}
	n := ctx.length()
	for i := range phi {
		phi[i] = sigma[i] + sum(0, n, func(k int32) float64 {
			return ctx.fraction(k) * ctx.component(k).d[i]
		})
	}
	return
}

// расчет избыточной составляющей вязкости
func getDeltaU(phi [6]float64, omegaM, tauM float64) float64 {
	return sum(0, 8, func(n int32) float64 {
		return c[n] * math.Pow(phi[0]*math.Pow(omegaM, phi[1])*math.Pow(tauM, phi[2]), r[n]) *
			math.Pow(phi[3]*math.Pow(omegaM, phi[4])*math.Pow(tauM, phi[5]), -t[n])
	})
}

// расчет вязкости компонентов в разреженном состоянии
func getMu0Comp(ctx *context) []float64 {
	n := ctx.length()
	u0 := make([]float64, n)
	for i := range u0 {
		u0[i] = sum(0, 3, func(k int32) float64 {
			return ctx.component(int32(i)).a[k] * math.Pow(ctx.t/100, float64(k))
		})
	}
	return u0
}

// расчет вязкости природного газа в разреженном состоянии
func getMu0(ctx *context, mu0comp []float64) float64 {
	n := ctx.length()
	return sum(0, n, func(i int32) float64 {
		ci := ctx.component(i)
		return ctx.fraction(i) * mu0comp[i] /
			sum(0, n, func(j int32) float64 {
				cj := ctx.component(j)
				return ctx.fraction(j) * math.Pow(
					1+math.Pow(mu0comp[i]/mu0comp[j], 1.0/2)*math.Pow(cj.m/ci.m, 1.0/4), 2) /
					math.Pow(8*(1+ci.m/cj.m), 1.0/2)
			})
	})
}

// расчет псевдокритического давления природного газа
func getPpc(ctx *context, pMolPc, tpc float64) float64 {
	n := ctx.length()
	return math.Pow10(-3) * R * pMolPc * tpc *
		(0.291 - 0.08*sum(0, n, func(i int32) float64 {
			return ctx.fraction(i) * ctx.component(i).omega
		}))
}

// расчет вязкости природного газа
func getMu(ctx *context, mu0, mm, ppc, tpc, deltaU float64) float64 {
	return mu0 + 2.63094*math.Pow(mm, 1.0/2)*math.Pow(ppc, 2.0/3)/math.Pow(tpc, 1.0/6)*deltaU
}
