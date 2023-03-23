package main

type componentFraction struct {
	component *gasComponent
	fraction  float64
}

type context struct {
	fractions []componentFraction
	p         float64
	t         float64
}

// Modifies component fractions to magically improve precision. IDK, just following the standard.
//
// This function should be called only once after creating new context
func (c *context) initFractions() {
	iHelium, iNitrogen, iHydrogen := -1, -1, -1
	for i, cf := range c.fractions {
		if cf.component == &helium {
			iHelium = i
		} else if cf.component == &nitrogen {
			iNitrogen = i
		} else if cf.component == &hydrogen {
			iHydrogen = i
		}
	}
	if iHelium != -1 && c.fractions[iHelium].fraction <= 0.0005 {
		if iNitrogen != -1 {
			c.fractions[iNitrogen].fraction += c.fractions[iHelium].fraction
		} else {
			c.fractions = append(c.fractions, componentFraction{&nitrogen, c.fractions[iHelium].fraction})
		}
		c.fractions[iHelium].fraction = 0
	}
	if iHydrogen != -1 && c.fractions[iHydrogen].fraction <= 0.0005 {
		if iNitrogen != -1 {
			c.fractions[iNitrogen].fraction += c.fractions[iHydrogen].fraction
		} else {
			c.fractions = append(c.fractions, componentFraction{&nitrogen, c.fractions[iHydrogen].fraction})
		}
		c.fractions[iHydrogen].fraction = 0
	}
}

func (c *context) fraction(i int32) float64 {
	return c.fractions[i].fraction
}

func (c *context) component(i int32) *gasComponent {
	return c.fractions[i].component
}

func (c *context) length() int32 {
	return int32(len(c.fractions))
}

func (c *context) eBin(i, j int32) float64 {
	if i == j {
		return 1
	}
	fc := c.fractions[i].component
	sc := c.fractions[j].component
	if bp, ok := fc.binaryParams[sc]; ok {
		return bp.e
	}
	if bp, ok := sc.binaryParams[fc]; ok {
		return bp.e
	}
	return 1
}

func (c *context) getKBin(i, j int32) float64 {
	fc := c.fractions[i].component
	sc := c.fractions[j].component
	if bp, ok := fc.binaryParams[sc]; ok {
		return bp.k
	}
	return 1
}

func (c *context) gBin(i, j int32) float64 {
	if i == j {
		return 1
	}
	fc := c.fractions[i].component
	sc := c.fractions[j].component
	if bp, ok := fc.binaryParams[sc]; ok {
		return bp.g
	}
	if bp, ok := sc.binaryParams[fc]; ok {
		return bp.g
	}
	return 1
}

func (c *context) vBin(i, j int32) float64 {
	fc := c.fractions[i].component
	sc := c.fractions[j].component
	if bp, ok := fc.binaryParams[sc]; ok {
		return bp.v
	}
	return 1
}
