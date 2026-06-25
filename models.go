package main

type HomePage struct {
	Nav
	Cars []CarModel
}

type CarPage struct {
	Lookup
	RecentlyViewed  []CarModel
	Recommendations []CarModel
	Car             CarModel
	InCompare       bool
	CompareFull     bool
	BackURL         string
}

type ComparePage struct {
	Cars    []CarModel
	BackURL string
}

type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CarModel struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	ManufacturerID int            `json:"manufacturerId"`
	CategoryID     int            `json:"categoryId"`
	Year           int            `json:"year"`
	Specs          Specifications `json:"specifications"`
	Image          string         `json:"image"`
	Manufacturer   *Manufacturer  `json:"manufacturer,omitempty"`
	Category       *Category      `json:"category,omitempty"`
}

type Specifications struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

type Nav struct {
	Manufacturers         []Manufacturer
	Categories            []Category
	SelectedManufacturers []int
	SelectedCategories    []int
}

type Lookup struct {
	Manufacturers []Manufacturer
	Categories    []Category
}
