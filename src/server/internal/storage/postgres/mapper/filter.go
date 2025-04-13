package mapper

import (
	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage/postgres/models"
	"clusterlizer/pkg/utils"
)

func GetByUserIDToFilter(rowModels []models.Filter) []entity.Filter {
	mapping := make(map[int]entity.Filter, len(rowModels))
	idxmap := make(map[int]int)
	res := make([]entity.Filter, 0, len(mapping))

	type input struct {
		str string
		id  int
	}
	keywordIdxMap := make(map[input]int)
	for _, rowModel := range rowModels {
		id := rowModel.ID
		var keywordRow string
		if rowModel.KeywordName != nil {
			keywordRow = ptrToStr(rowModel.City)
		}
		uniqKeyword := input{id: id, str: keywordRow}

		if idx, exists := idxmap[id]; exists {
			if _, exists := keywordIdxMap[uniqKeyword]; !exists {
				res[idx].Keywords = append(res[idx].Keywords, entity.Keyword{
					Name: utils.SafeNil(rowModel.KeywordName),
					Slug: utils.SafeNil(rowModel.KeywordSlug),
				})
				keywordIdxMap[uniqKeyword] = idx
			}
		} else {
			keywords := make([]entity.Keyword, 0, 1)
			if rowModel.KeywordName != nil {
				keywords = append(keywords, entity.Keyword{
					Name: utils.SafeNil(rowModel.KeywordName),
					Slug: utils.SafeNil(rowModel.KeywordSlug),
				})
			}
			newFilter := entity.Filter{
				//ID:         rowModel.ID,
				UserTgID: rowModel.UserTgID,
				City:     utils.SafeNil(rowModel.City),
				Keywords: keywords,
				Company: entity.Company{
					Slug: utils.SafeNil(rowModel.CompanySlug),
					Name: utils.SafeNil(rowModel.CompanyName),
				},
				Experience: entity.Experience{
					Slug: utils.SafeNil(rowModel.ExperienceSlug),
					Name: utils.SafeNil(rowModel.ExperienceName),
				},
				Speciality: entity.Speciality{
					Slug: utils.SafeNil(rowModel.SpecialitySlug),
					Name: utils.SafeNil(rowModel.SpecialityName),
				},
			}

			res = append(res, newFilter)
			idxmap[id] = len(res) - 1
		}
	}

	return res
}
