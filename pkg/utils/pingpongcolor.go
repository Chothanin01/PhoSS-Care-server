package utils

func CalculateColorStatus(sugar float64, pressure float64) string {
	
	if sugar <= 0 && pressure <= 0 {
		return "none"
	}

	sugarLevel := 1
	if sugar > 0 {
		if sugar < 100 {
			sugarLevel = 1 
		} else if sugar <= 125 {
			sugarLevel = 3 
		} else if sugar <= 154 {
			sugarLevel = 4 
		} else if sugar <= 182 {
			sugarLevel = 5 
		} else {
			sugarLevel = 6 
		}
	}

	pressureLevel := 1
	if pressure > 0 {
		if pressure < 120 {
			pressureLevel = 1 
		} else if pressure <= 139 {
			pressureLevel = 3 
		} else if pressure <= 159 {
			pressureLevel = 4 
		} else if pressure <= 179 {
			pressureLevel = 5 
		} else {
			pressureLevel = 6 
		}
	}

	maxLevel := sugarLevel
	if pressureLevel > maxLevel {
		maxLevel = pressureLevel
	}

	switch maxLevel {
	case 1:
		return "white"
	case 3:
		return "dark_green"
	case 4:
		return "yellow"
	case 5:
		return "orange"
	case 6:
		return "red"
	default:
		return "none"
	}
}