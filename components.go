package main

type gasComponent struct {
	// название
	name string
	// молярная масса M
	m float64
	// коэффициент сжимаемости при стандартных условиях Zc
	zc float64
	// энергетический параметр
	e float64
	// параметр размера K, (m^3/kMole)^1/3
	k float64
	// ориентационный параметр
	g float64
	// квадрупольный параметр
	q float64
	// высокотемпературный параметр
	f float64
	// дипольный параметр
	s float64
	// параметр ассоциации
	w float64
	// b0 - j0 - коэффициенты для расчета безразмерных изобарных теплоемкостей
	b0 float64
	c0 float64
	d0 float64
	e0 float64
	f0 float64
	g0 float64
	h0 float64
	i0 float64
	j0 float64
	// критическая температура
	tcr float64
	// критическое давление
	pcr float64
	// фактор Питцера
	omega float64
	// коэффициент для расчета параметров преобразований phi
	d [6]float64
	// коэффициент для расчета вязкости
	a [4]float64
	// параметры бинарного взаимодействия с другими компонентами
	binaryParams map[*gasComponent]binaryInteractionParams
}

// параметры бинарного взаимодействия компонентов
type binaryInteractionParams struct {
	e float64
	v float64
	k float64
	g float64
}

var methane = gasComponent{
	name:  "метан",
	m:     16.043,
	zc:    0.9981,
	e:     151.318300,
	k:     0.4619255,
	g:     0.0,
	q:     0.0,
	f:     0.0,
	s:     0.0,
	w:     0.0,
	b0:    4.00088,
	c0:    0.76315,
	d0:    820.659,
	e0:    0.00460,
	f0:    178.410,
	g0:    8.74432,
	h0:    1062.82,
	i0:    -4.46921,
	j0:    1090.53,
	tcr:   190.564,
	pcr:   162.66,
	omega: 0.064294,
	d:     [6]float64{0, 0, 0, 0, 0, 0},
	a:     [4]float64{-0.838029104, 4.88406903, -0.344504244, 0.0151593109},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&propane: {
			e: 0.994635,
			v: 0.990877,
			k: 1.007619,
			g: 1,
		},
		&iButane: {
			e: 1.019530,
			v: 1,
			k: 1,
			g: 1,
		},
		&nButane: {
			e: 0.989844,
			v: 0.992291,
			k: 0.997596,
			g: 1,
		},
		&iPentane: {
			e: 1.002350,
			v: 1,
			k: 1,
			g: 1,
		},
		&nPentane: {
			e: 0.999268,
			v: 1.003670,
			k: 1.002529,
			g: 1,
		},
		&nHexane: {
			e: 1.107274,
			v: 1.302576,
			k: 0.982962,
			g: 1,
		},
		&nitrogen: {
			e: 0.971640,
			v: 0.886106,
			k: 1.003630,
			g: 1,
		},
		&carbonDioxide: {
			e: 0.960644,
			v: 0.963827,
			k: 0.995933,
			g: 0.807653,
		},
		&hydrogen: {
			e: 1.170520,
			v: 1.156390,
			k: 1.023260,
			g: 1.957310,
		},
	},
}

var ethane = gasComponent{
	name:  "этан",
	m:     30.070,
	zc:    0.992,
	e:     244.166700,
	k:     0.5279209,
	g:     0.079300,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    4.00263,
	c0:    4.33939,
	d0:    559.314,
	e0:    1.23722,
	f0:    223.284,
	g0:    13.1974,
	h0:    1031.38,
	i0:    -6.01989,
	j0:    1071.29,
	tcr:   305.32,
	pcr:   206.58,
	omega: 0.10958,
	d:     [6]float64{0.04156931, 0, 0.06408111, 0.04763455, -0.1889656, 0.1533738},
	a:     [4]float64{-1.21924490, 4.05145591, -0.200150993, 0.00662746099},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&propane: {
			e: 1.022560,
			v: 1.065173,
			k: 0.986893,
			g: 1,
		},
		&iButane: {
			e: 1,
			v: 1.250000,
			k: 1,
			g: 1,
		},
		&nButane: {
			e: 1.013060,
			v: 1.250000,
			k: 1,
			g: 1,
		},
		&iPentane: {
			e: 1,
			v: 1.250000,
			k: 1,
			g: 1,
		},
		&nPentane: {
			e: 1.005320,
			v: 1.250000,
			k: 1,
			g: 1,
		},
		&nitrogen: {
			e: 0.970120,
			v: 0.816431,
			k: 1.007960,
			g: 1,
		},
		&carbonDioxide: {
			e: 0.925053,
			v: 0.969870,
			k: 1.008510,
			g: 0.370296,
		},
		&hydrogen: {
			e: 1.164460,
			v: 1.616660,
			k: 1.020340,
			g: 1,
		},
	},
}

var propane = gasComponent{
	name:  "пропан",
	m:     44.097,
	zc:    0.9834,
	e:     298.118300,
	k:     0.5837490,
	g:     0.141239,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    4.02939,
	c0:    6.60569,
	d0:    479.856,
	e0:    3.19700,
	f0:    200.893,
	g0:    19.1921,
	h0:    955.312,
	i0:    -8.37267,
	j0:    1027.29,
	tcr:   369.825,
	pcr:   220.49,
	omega: 0.18426,
	d:     [6]float64{0.03976538, 0.08375624, 0.1747180, 1.250272, 0.5283498, 0.2458511},
	a:     [4]float64{0.254518256, 2.54779249, 0.0683095277, 0.0114348793},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&nButane: {
			e: 1.004900,
			v: 1,
			k: 1,
			g: 1,
		},
		&nitrogen: {
			e: 0.945939,
			v: 0.915502,
			k: 1,
			g: 1,
		},
		&carbonDioxide: {
			e: 0.960237,
			v: 1,
			k: 1,
			g: 1,
		},
		&hydrogen: {
			e: 1.034787,
			v: 1,
			k: 1,
			g: 1,
		},
	},
}

var iButane = gasComponent{
	name:  "и-бутан",
	m:     58.123,
	zc:    0.971,
	e:     324.068900,
	k:     0.6406937,
	g:     0.256692,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    4.06714,
	c0:    8.97575,
	d0:    438.270,
	e0:    5.25156,
	f0:    198.018,
	g0:    25.1423,
	h0:    1905.02,
	i0:    16.1388,
	j0:    893.765,
	tcr:   407.85,
	pcr:   224.36,
	omega: 0.16157,
	d:     [6]float64{0.07234927, 0.009435210, -0.03673568, 0.4516722, 0.3272680, -0.6135352},
	a:     [4]float64{1.04273843, 1.69220741, 0.194077419, -0.0159867334},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&nitrogen: {
			e: 0.946914,
			v: 1,
			k: 1,
			g: 1,
		},
		&carbonDioxide: {
			e: 0.906849,
			v: 1,
			k: 1,
			g: 1,
		},
		&hydrogen: {
			e: 1.300000,
			v: 1,
			k: 1,
			g: 1,
		},
	},
}

var nButane = gasComponent{
	name:  "н-бутан",
	m:     58.123,
	zc:    0.9682,
	e:     337.638900,
	k:     0.6341423,
	g:     0.281835,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    4.33944,
	c0:    9.44893,
	d0:    468.270,
	e0:    6.89406,
	f0:    183.636,
	g0:    24.4618,
	h0:    1914.10,
	i0:    14.7824,
	j0:    903.185,
	tcr:   425.16,
	pcr:   227.85,
	omega: 0.21340,
	d:     [6]float64{-0.06667775, 0.2100174, 0.06330205, 0.3182660, 0.1474434, -1.113935},
	a:     [4]float64{-0.524058048, 2.81260308, -0.0496574363, 0},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&nitrogen: {
			e: 0.973384,
			v: 0.993556,
			k: 1,
			g: 1,
		},
		&carbonDioxide: {
			e: 0.897362,
			v: 1,
			k: 1,
			g: 1,
		},
		&hydrogen: {
			e: 1.300000,
			v: 1,
			k: 1,
			g: 1,
		},
	},
}

var iPentane = gasComponent{
	name:  "и-пентан",
	m:     72.150,
	zc:    0.953,
	e:     365.599900,
	k:     0.6738577,
	g:     0.332267,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    4,
	c0:    11.7618,
	d0:    292.503,
	e0:    20.1101,
	f0:    910.237,
	g0:    33.1688,
	h0:    1919.37,
	i0:    0,
	j0:    0,
	tcr:   460.39,
	pcr:   236.0,
	omega: 0.26196,
	d:     [6]float64{0.02229787, 0.08380246, 0.04639638, -0.1450583, 0.03725585, -0.4106772},
	a:     [4]float64{0.550744125, 1.75702204, 0.173363456, -0.0167839786},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&nitrogen: {
			e: 0.959340,
			v: 1,
			k: 1,
			g: 1,
		},
		&carbonDioxide: {
			e: 0.726255,
			v: 1,
			k: 1,
			g: 1,
		},
	},
}

var nPentane = gasComponent{
	name:  "н-пентан",
	m:     72.150,
	zc:    0.945,
	e:     370.682300,
	k:     0.6798307,
	g:     0.366911,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    4,
	c0:    8.95043,
	d0:    178.670,
	e0:    21.8360,
	f0:    840.538,
	g0:    33.4032,
	h0:    1774.25,
	i0:    0,
	j0:    0,
	tcr:   469.65,
	pcr:   232.0,
	omega: 0.29556,
	d:     [6]float64{0, 0.1651156, -0.07126922, 0.06698673, -0.5283166, -0.7803174},
	a:     [4]float64{0.452603096, 1.79775689, 0.157002776, -0.0158057627},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&nitrogen: {
			e: 0.945520,
			v: 1,
			k: 1,
			g: 1,
		},
		&carbonDioxide: {
			e: 0.859764,
			v: 1,
			k: 1,
			g: 1,
		},
	},
}

var nHexane = gasComponent{
	name:  "н-гексан",
	m:     86.177,
	zc:    0.919,
	e:     402.636293,
	k:     0.7175118,
	g:     0.289731,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    4,
	c0:    11.6977,
	d0:    182.326,
	e0:    26.8142,
	f0:    859.207,
	g0:    38.6164,
	h0:    1826.59,
	i0:    0,
	j0:    0,
	tcr:   507.85,
	pcr:   233.6,
	omega: 0.29965,
	d:     [6]float64{0.1753529, -0.08018375, -0.03543316, -0.09677546, -0.2015218, -1.206562},
	a:     [4]float64{0.658064311, 1.50818329, 0.178280027, -0.0161050134},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&carbonDioxide: {
			e: 0.855134,
			v: 1.066638,
			k: 0.910183,
			g: 1,
		},
	},
}

var nitrogen = gasComponent{
	name:  "азот",
	m:     28.0135,
	zc:    0.9997,
	e:     99.737780,
	k:     0.4479153,
	g:     0.027815,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    3.50031,
	c0:    0.13732,
	d0:    662.738,
	e0:    0.14660,
	f0:    680.562,
	g0:    0.90066,
	h0:    1740.06,
	i0:    0,
	j0:    0,
	tcr:   126.2,
	pcr:   313.1,
	omega: 0.013592,
	d:     [6]float64{-0.005352690, 0.09101896, 0.01501200, 0.2640642, -0.1032012, -0.1078872},
	a:     [4]float64{-0.279070091, 7.81221301, -0.699863421, 0.0378831186},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&carbonDioxide: {
			e: 1.022740,
			v: 0.835058,
			k: 0.982361,
			g: 0.982746,
		},
		&hydrogen: {
			e: 1.086320,
			v: 0.408838,
			k: 1.032270,
			g: 1,
		},
	},
}

var carbonDioxide = gasComponent{
	name:  "диоксид углерода",
	m:     44.010,
	zc:    0.9947,
	e:     241.960600,
	k:     0.4557489,
	g:     0.189065,
	q:     0.690000,
	f:     0,
	s:     0,
	w:     0,
	b0:    3.50002,
	c0:    2.04452,
	d0:    919.306,
	e0:    -1.06044,
	f0:    865.070,
	g0:    2.03366,
	h0:    483.553,
	i0:    0.01393,
	j0:    341.109,
	tcr:   304.2,
	pcr:   468.0,
	omega: 0.20625,
	d:     [6]float64{-0.03468202, 0.1130498, 0.05811886, 0.05767935, -0.1814105, -0.5971794},
	a:     [4]float64{-0.468233636, 5.37907799, -0.0349633355, -0.0126198032},
	binaryParams: map[*gasComponent]binaryInteractionParams{
		&hydrogen: {
			e: 1.281790,
			v: 1,
			k: 1,
			g: 1,
		},
	},
}

var helium = gasComponent{
	name:  "гелий",
	m:     4.0026,
	zc:    1.0005,
	e:     2.610111,
	k:     0.3589888,
	g:     0,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    2.5,
	tcr:   5.19,
	pcr:   69.64,
	omega: -0.14949,
	d:     [6]float64{0.299249, -0.1490941, -0.1577329, -0.225324, -0.2731058, -0.8827831},
	a:     [4]float64{2.95929817, 7.1775132, -0.641191946, 0.0451852767},
}

var hydrogen = gasComponent{
	name:  "водород",
	m:     2.0159,
	zc:    1.0006,
	e:     26.957940,
	k:     0.3514916,
	g:     0.034369,
	q:     0,
	f:     0,
	s:     0,
	w:     0,
	b0:    2.47906,
	c0:    0.95806,
	d0:    228.734,
	e0:    0.45444,
	f0:    326.843,
	g0:    1.56039,
	h0:    1651.71,
	i0:    -1.3756,
	j0:    1671.69,
	tcr:   32.938,
	pcr:   31.36,
	omega: -0.12916,
	d:     [6]float64{-0.03937273, 0.01532106, -0.03423876, -0.1399209, -0.06955475, -1.049055},
	a:     [4]float64{1.42410895, 3.03739469, -0.203048737, 0.0106137856},
}

var componentsByName = map[string]*gasComponent{
	methane.name:       &methane,
	ethane.name:        &ethane,
	propane.name:       &propane,
	iButane.name:       &iButane,
	nButane.name:       &nButane,
	iPentane.name:      &iPentane,
	nPentane.name:      &nPentane,
	nHexane.name:       &nHexane,
	nitrogen.name:      &nitrogen,
	carbonDioxide.name: &carbonDioxide,
	helium.name:        &helium,
	hydrogen.name:      &hydrogen,
}

var allComponents = []*gasComponent{
	&methane,
	&ethane,
	&propane,
	&propane,
	&iButane,
	&nButane,
	&iPentane,
	&nPentane,
	&nHexane,
	&nitrogen,
	&carbonDioxide,
	&helium,
	&hydrogen,
}
