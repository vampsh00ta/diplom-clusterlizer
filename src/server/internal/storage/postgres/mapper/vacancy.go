package mapper

import (
	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage/postgres/models"
	"clusterlizer/pkg/utils"
)

func GetByFiltersToVacancies(rowModels []models.GetAllWithFilter) []entity.Vacancy {
	res := make([]entity.Vacancy, 0, len(rowModels))
	idxMap := make(map[int]int)

	type input struct {
		str string
		id  int
	}
	citiesIdxMap := make(map[input]int)
	keywordsIdxMap := make(map[input]int)

	for _, rowModel := range rowModels {
		var keywordRow, cityRow string
		if rowModel.KeywordName != nil {
			keywordRow = *rowModel.KeywordName
		}
		if rowModel.City != "" {
			cityRow = rowModel.City
		}
		id := rowModel.ID
		uniqCity := input{id: id, str: cityRow}
		uniqKeyword := input{id: id, str: keywordRow}

		if idx, exists := idxMap[id]; exists {
			if _, exists := citiesIdxMap[uniqCity]; !exists {
				res[idx].Cities = append(res[idx].Cities, cityRow)
				citiesIdxMap[uniqCity] = idx
			}
			if _, exists := keywordsIdxMap[uniqKeyword]; !exists {
				res[idx].Keywords = append(res[idx].Keywords, keywordRow)
				keywordsIdxMap[uniqKeyword] = idx
			}
		} else {
			keywords := make([]string, 0)
			cities := make([]string, 0)
			if rowModel.KeywordName != nil {
				keywords = append(keywords, *rowModel.KeywordName)
			}
			if rowModel.City != "" {
				cities = append(cities, rowModel.City)
			}
			newVacancy := entity.Vacancy{
				ID:       rowModel.ID,
				SlugID:   rowModel.SlugID,
				Name:     rowModel.Name,
				Cities:   cities,
				Keywords: keywords,

				Category:       rowModel.Category,
				Description:    rowModel.Description,
				ExperienceName: rowModel.ExperienceName,
				ExperienceSlug: rowModel.ExperienceSlug,
				CompanyName:    rowModel.CompanyName,
				CompanySlug:    rowModel.CompanySlug,

				SpecialityName: utils.SafeNil(rowModel.SpecialityName),
				SpecialitySlug: utils.SafeNil(rowModel.SpecialitySlug),
				Link:           rowModel.Link,
				Rank:           rowModel.Rank,
			}
			citiesIdxMap[uniqCity] = idx
			keywordsIdxMap[uniqKeyword] = idx

			res = append(res, newVacancy)
			idxMap[id] = len(res) - 1
		}
	}
	return res
}

func GetAllToVacancies(rowModels []models.GetAll) []entity.Vacancy {
	res := make([]entity.Vacancy, 0, len(rowModels))
	idxmap := make(map[int]int)

	type input struct {
		str string
		id  int
	}
	citiesIdxMap := make(map[input]int)

	for _, rowModel := range rowModels {
		var cityRow string
		if rowModel.City != "" {
			cityRow = rowModel.City
		}
		id := rowModel.ID
		uniqCity := input{id: id, str: cityRow}

		if idx, exists := idxmap[id]; exists {
			if _, exists := citiesIdxMap[uniqCity]; !exists {
				res[idx].Cities = append(res[idx].Cities, cityRow)
				citiesIdxMap[uniqCity] = idx
			}
		} else {
			keywords := make([]string, 0)
			cities := make([]string, 0)
			if rowModel.City != "" {
				cities = append(cities, rowModel.City)
			}
			newVacancy := entity.Vacancy{
				ID:             rowModel.ID,
				SlugID:         rowModel.SlugID,
				Name:           rowModel.Name,
				Cities:         cities,
				Keywords:       keywords,
				Category:       rowModel.Category,
				Description:    rowModel.Description,
				ExperienceName: rowModel.ExperienceName,
				ExperienceSlug: rowModel.ExperienceSlug,
				CompanyName:    rowModel.CompanyName,
				CompanySlug:    rowModel.CompanySlug,
				Link:           rowModel.Link,
			}
			citiesIdxMap[uniqCity] = idx

			res = append(res, newVacancy)
			idxmap[id] = len(res) - 1
		}
	}
	return res
}
