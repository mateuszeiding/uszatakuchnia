package mappers

import (
	"sort"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func MapTaxonomyValueToDto(e entities.TaxonomyValue) dtos.TaxonomyValueDto {
	return dtos.TaxonomyValueDto{ID: e.ID, Label: e.Label, Color: e.Color}
}

// MapRecipeTaxonomyToDto grupuje wartości po kodzie osi. Wymaga preloadu "Taxonomy.Value.Axis".
func MapRecipeTaxonomyToDto(a []entities.RecipeTaxonomy) dtos.RecipeTaxonomyDto {
	sorted := make([]entities.RecipeTaxonomy, len(a))
	copy(sorted, a)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Value.SortOrder < sorted[j].Value.SortOrder
	})

	out := dtos.RecipeTaxonomyDto{}
	for _, rt := range sorted {
		code := rt.Value.Axis.Code
		out[code] = append(out[code], MapTaxonomyValueToDto(rt.Value))
	}
	return out
}

func MapTaxonomyAxesToDto(axes []entities.TaxonomyAxis, values []entities.TaxonomyValue) []dtos.TaxonomyAxisDto {
	byAxis := map[uint][]dtos.TaxonomyValueDto{}
	for _, v := range values {
		byAxis[v.AxisID] = append(byAxis[v.AxisID], MapTaxonomyValueToDto(v))
	}
	out := make([]dtos.TaxonomyAxisDto, 0, len(axes))
	for _, a := range axes {
		vals := byAxis[a.ID]
		if vals == nil {
			vals = []dtos.TaxonomyValueDto{}
		}
		out = append(out, dtos.TaxonomyAxisDto{
			ID: a.ID, Code: a.Code, Label: a.Label, Multi: a.Multi, Required: a.Required,
			HasColor: a.HasColor, FilterMode: a.FilterMode, Values: vals,
		})
	}
	return out
}

func MapUnitToDto(e entities.Unit) dtos.UnitDto {
	var forms []string
	if e.FormOne != nil && e.FormFew != nil && e.FormMany != nil {
		forms = []string{*e.FormOne, *e.FormFew, *e.FormMany}
	}
	return dtos.UnitDto{ID: e.ID, Name: e.Name, Forms: forms, NoQty: e.NoQty}
}

func MapUnitArrayToDto(a []entities.Unit) []dtos.UnitDto {
	out := make([]dtos.UnitDto, 0, len(a))
	for _, u := range a {
		out = append(out, MapUnitToDto(u))
	}
	return out
}

func MapAllergensToDto() []dtos.AllergenDto {
	out := make([]dtos.AllergenDto, 0, len(entities.Allergens))
	for _, a := range entities.Allergens {
		out = append(out, dtos.AllergenDto{Code: a.Code, Label: a.Label})
	}
	return out
}

func mapIngredientAllergenCodes(a []entities.IngredientAllergen) []string {
	out := make([]string, 0, len(a))
	for _, x := range a {
		out = append(out, x.Code)
	}
	return out
}
