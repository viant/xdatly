package predicate

type StringsFilter struct {
	Include []string
	Exclude []string
}

type IntFilter struct {
	Include []int
	Exclude []int
}

type BoolFilter struct {
	Include []bool
	Exclude []bool
}
