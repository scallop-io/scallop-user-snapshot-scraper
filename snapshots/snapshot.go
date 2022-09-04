package snapshots

import (
	"context"
	"fmt"
	"github.com/scallop-io/scallop-user-snapshot-scraper/common/utils"
	"github.com/scallop-io/scallop-user-snapshot-scraper/types"
	"net/http"
	"strings"
)

func FetchSnapshot(ctx context.Context, poolBase string, periodNumber uint) (types.Snapshot, error) {
	result := struct {
		Message   string         `json:"message"`
		Timestamp uint64         `json:"timestamp"`
		Data      types.Snapshot `json:"data"`
	}{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.scallop.io/v1/snapshot/%s/%d", poolBase, periodNumber), nil)
	if err != nil {
		return result.Data, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return result.Data, err
	}

	err = utils.PopulateModelFromBody(resp.Body, &result)
	if err != nil {
		return result.Data, err
	}

	if result.Message != "OK" {
		return result.Data, fmt.Errorf("returned message is %s", result.Message)
	}
	return result.Data, nil
}

func FetchMultipleSnapshot(ctx context.Context, poolBase string, periodNumbers []uint) ([]types.Snapshot, error) {
	result := struct {
		Message   string           `json:"message"`
		Timestamp uint64           `json:"timestamp"`
		Data      []types.Snapshot `json:"data"`
	}{}

	periodNumbersParams := ""
	for _, periodNumber := range periodNumbers {
		periodNumbersParams += fmt.Sprintf("%v,", periodNumber)
	}
	periodNumbersParams = strings.Trim(periodNumbersParams, ",")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.scallop.io/v1/snapshot/%s?period_numbers=%s", poolBase, periodNumbersParams), nil)
	if err != nil {
		return result.Data, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return result.Data, err
	}

	err = utils.PopulateModelFromBody(resp.Body, &result)
	if err != nil {
		return result.Data, err
	}

	if result.Message != "OK" {
		return result.Data, fmt.Errorf("returned message is %s", result.Message)
	}
	return result.Data, nil
}
