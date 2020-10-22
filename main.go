package main
import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sort"
    "context"
    "encoding/json"
    "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Pack struct {
	Size int `json:"size"`
	Qty int `json:"qty"`
}

type Packs []Pack

func main() {
    lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	items, err := strconv.Atoi(request.QueryStringParameters["items"])
	packs := removeDuplicatesAndSort(stringToIntArray(request.QueryStringParameters["packs"]))
	results := itemsToPacks(items,packs)
	
	if err != nil {
        // handle error
        fmt.Println(err)
	}

	body, err := json.Marshal(results)

	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "https://next-products.vercel.app"
	resp.Headers["Access-Control-Allow-Credentials"] = "true"

    if err != nil {
		resp.Body = string("")
		resp.StatusCode = 200
        return resp, nil
	}
	
	resp.Body = string(body)
	resp.StatusCode = 200

    return resp, nil
}

func itemsToPacks(items int, packSet []int) Packs {

	itemsOrdered := getMinimumOrderQty(items, packSet[0]) 
	sort.Sort(sort.Reverse(sort.IntSlice(packSet))) //Reverse Pack Order

	result := Packs{}
	qty := int(0)

	for i, packSize := range packSet {
		_ = i
		qty, itemsOrdered = packsRequired(packSize, itemsOrdered)
		if qty > 0 {
			result = append(result,Pack{Size:int(packSize),Qty:int(qty)})
		}
	}

	return result
}

func getMinimumOrderQty(items int, smallestPack int) int{
	itemsOrdered := (items / smallestPack) * smallestPack

	if itemsOrdered < items {
		itemsOrdered = itemsOrdered + smallestPack
	}

	return itemsOrdered;
}

func packsRequired(packSize int, items int) (int, int) {
	packs := int(0)

	if items != 0 {
		packs = (items / packSize)
		remainingItems := items - (packSize * packs)
		return packs, remainingItems
	}
	return packs, items
}

func removeDuplicatesAndSort(elements []int) []int {
    encountered := map[int]bool{}
    result := []int{}

    for v := range elements {
        if encountered[elements[v]] == true {
        } else {
            encountered[elements[v]] = true
            result = append(result, elements[v])
        }
	}
	sort.Ints(result)
    return result
}

func stringToIntArray(packSizes string) []int { 
    tmp := strings.Split(packSizes, ",")
	values := make([]int, 0, len(tmp))
	for _, raw := range tmp {
		v, err := strconv.Atoi(raw)
		if err != nil {
			log.Print(err)
			continue
		}
		values = append(values, v)
	}
	return values
}