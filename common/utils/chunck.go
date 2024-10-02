package utils

import (
	"encoding/json"

	"github.com/samber/lo"
)

func GetChunkProductsForPrint(data interface{}, productField string, options map[string]int, firstPageRows, pageRows int) map[string]interface{} {
	jsonByte, _ := json.Marshal(data)
	dataMap := map[string]interface{}{}
	_ = json.Unmarshal(jsonByte, &dataMap)
	productsInterface, ok := dataMap[productField].([]interface{})
	if !ok {
		return dataMap
	}
	products := []map[string]interface{}{}
	for _, item := range productsInterface {
		productItem := item.(map[string]interface{})
		products = append(products, productItem)
	}
	dataMap["productsChunk"] = ChunkProductsForPrint(products, options, firstPageRows, pageRows)
	return dataMap
}

func ChunkProductsForPrint(products []map[string]interface{}, options map[string]int, firstPageRows, pageRows int) [][]map[string]interface{} {
	optionsKeys := lo.Keys(options)
	heightSlice := [][]int{}
	heightTotal := 0
	keyStart := 0
	productListChunk := [][]map[string]interface{}{}
	for pIndex, pItem := range products {
		subHeightSlice := []int{}
		for _, oItem := range optionsKeys {
			resName, height := FormatNameForPrint(pItem[oItem], options[oItem])
			subHeightSlice = append(subHeightSlice, height)
			products[pIndex][oItem] = resName

		}
		heightSlice = append(heightSlice, subHeightSlice)
	}
	for index, hItemSlice := range heightSlice {
		heightTotal += lo.Max(hItemSlice)
		limitRows := pageRows
		if keyStart == 0 {
			limitRows = firstPageRows
		}
		if heightTotal > limitRows {
			productListChunk = append(productListChunk, lo.Slice(products, keyStart, index+1))
			keyStart = index + 1
			heightTotal = 0
		}
	}
	if heightTotal > 0 {
		productListChunk = append(productListChunk, lo.Slice(products, keyStart, len(products)))
	}
	return productListChunk
}

func FormatNameForPrint(nameData interface{}, length int) (string, int) {
	name, ok := nameData.(string)
	if !ok {
		return "", 0
	}
	height := 0
	resName := ""
	runeSlice := []rune(name)
	nameLen := len(runeSlice)
	if nameLen > length {
		count := nameLen / length
		if nameLen%length != 0 {
			count++
		}
		for i := 0; i < count; i++ {
			height++
			resName += lo.Substring(name, i*length, uint(length)) + "<br>"
		}
	} else {
		height++
		resName = name
	}
	return resName, height
}
