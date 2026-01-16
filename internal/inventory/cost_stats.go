package inventory

type CostStats struct {
	// Total cost (US Dollars)
	TotalCost float64

	// Cost Last 15d
	Last15DaysCost float64

	// Last month cost
	LastMonthCost float64

	// Current month so far cost
	CurrentMonthSoFarCost float64
}
