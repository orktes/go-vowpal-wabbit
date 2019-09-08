package vw

// PerformanceStatistics contains model performance stats
type PerformanceStatistics struct {
	// Index of current pass
	CurrentPass uint64
	// The total number of features
	NumberOfFeatures uint64
	// Total number of examples
	NumberOfExamples uint64
	// The weighted sum of examples.
	WeightedExampleSum float64
	// The weighted sum of labels.
	WeightedLabelSum float64
	// The average loss since instance creation.
	AverageLoss float64
	// The best constant since instance creation.
	BestConstant float64
	// The best constant loss since instance creation.
	BestConstantLoss float64
}
