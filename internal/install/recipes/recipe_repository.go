package recipes

import (
	"math"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeRepository struct {
	RecipeLoaderFunc func() ([]types.OpenInstallationRecipe, error)
	recipes          []types.OpenInstallationRecipe
}

type recipeMatch struct {
	matchCount int
	recipe     types.OpenInstallationRecipe
}

// NewRecipeRepository returns a new instance of types.RecipeRepository.
func NewRecipeRepository(loaderFunc func() ([]types.OpenInstallationRecipe, error)) *RecipeRepository {
	rr := RecipeRepository{
		RecipeLoaderFunc: loaderFunc,
		recipes:          nil,
	}

	return &rr
}

func (rf *RecipeRepository) FindAll(m types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	results := []types.OpenInstallationRecipe{}
	matchRecipes := make(map[string][]recipeMatch)
	hostMap := getHostMap(m)

	if rf.recipes == nil {
		recipes, err := rf.RecipeLoaderFunc()
		if err != nil {
			return nil, err
		}
		log.Debugf("Loaded %d recipes", len(recipes))

		rf.recipes = recipes
	}

	log.Debugf("Find all available out of %d recipes for host %+v", len(rf.recipes), hostMap)
	log.Debugf("All recipes: %+v", rf.recipes)

	for _, recipe := range rf.recipes {
		matchTargetCount := []int{}

		for _, rit := range recipe.InstallTargets {
			matchCount := 0
			for k, v := range getRecipeTargetMap(rit) {
				if v == "" {
					continue
				}
				isValueMatching := matchRecipeCriteria(hostMap, k, v)
				if isValueMatching {
					log.Debugf("matching recipe %s field name %s and value %s", recipe.Name, k, v)
					matchCount++
				} else {
					log.Debugf("recipe %s defines %s but input did not provide a match", recipe.Name, k)
					matchCount = -1
					break
				}
			}
			if matchCount >= 0 {
				matchTargetCount = append(matchTargetCount, matchCount)
			}
		}

		if len(recipe.InstallTargets) == 0 || len(matchTargetCount) > 0 {
			maxMatchTargetCount := 0
			if len(matchTargetCount) > 0 {
				maxMatchTargetCount = mathMax(matchTargetCount)
			}
			log.Debugf("Recipe InstallTargetsCount %d and maxMatchCount %d", len(recipe.InstallTargets), maxMatchTargetCount)

			match := recipeMatch{
				recipe:     recipe,
				matchCount: maxMatchTargetCount,
			}
			if _, ok := matchRecipes[recipe.Name]; !ok {
				matches := []recipeMatch{match}
				matchRecipes[recipe.Name] = matches
			} else {
				matchRecipes[recipe.Name] = append(matchRecipes[recipe.Name], match)
			}
		}
	}

	for _, matches := range matchRecipes {
		if len(matches) > 0 {
			match := findMaxMatch(matches)
			singleRecipe := match.recipe
			results = append(results, singleRecipe)
		}
	}

	return results, nil
}

func findMaxMatch(matches []recipeMatch) recipeMatch {
	var result *recipeMatch

	for _, match := range matches {
		if result == nil {
			result = &match
		} else {
			if match.matchCount > result.matchCount {
				result = &match
			}
		}
	}

	return *result
}

func mathMax(numbers []int) int {
	result := math.MinInt32
	for _, number := range numbers {
		if number > result {
			result = number
		}
	}
	return result
}

func matchRecipeCriteria(hostMap map[string]string, rkey string, rvalue string) bool {
	if val, ok := hostMap[rkey]; ok {
		return strings.EqualFold(val, rvalue)
	}

	return false
}

func getHostMap(m types.DiscoveryManifest) map[string]string {
	hostMap := map[string]string{
		"KernelArch":      m.KernelArch,
		"KernelVersion":   m.KernelVersion,
		"OS":              m.OS,
		"Platform":        m.Platform,
		"PlatformFamily":  m.PlatformFamily,
		"PlatformVersion": m.PlatformVersion,
	}
	return hostMap
}

func getRecipeTargetMap(rit types.OpenInstallationRecipeInstallTarget) map[string]string {
	targetMap := map[string]string{
		"KernelArch":      rit.KernelArch,
		"KernelVersion":   rit.KernelVersion,
		"OS":              string(rit.Os),
		"Platform":        string(rit.Platform),
		"PlatformFamily":  string(rit.PlatformFamily),
		"PlatformVersion": rit.PlatformVersion,
	}
	return targetMap
}
