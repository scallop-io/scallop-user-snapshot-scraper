package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gagliardetto/solana-go"
	splToken "github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/scallop-io/scallop-user-snapshot-scraper/common/utils"
	"github.com/scallop-io/scallop-user-snapshot-scraper/config"
	"github.com/scallop-io/scallop-user-snapshot-scraper/snapshots"
	"github.com/scallop-io/scallop-user-snapshot-scraper/types"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Print("Please enter the pool base of the pool: ")
	var poolBase string
	fmt.Scanf("%s", &poolBase)

	fmt.Print("Please enter the start date (yyyy-mm-dd): ")
	var startDate string
	fmt.Scanf("%s", &startDate)

	fmt.Print("Please enter the end date (yyyy-mm-dd): ")
	var endDate string
	fmt.Scanf("%s", &endDate)

	fmt.Print("Please enter minimal balance of token that user should had: ")
	var minimalBalance float64
	fmt.Scanln(&minimalBalance)

	conf := config.LoadConfig("./config")

	programId, err := solana.PublicKeyFromBase58(conf.ProgramId)
	if err != nil {
		panic(err)
	}

	poolBasePubkey, err := solana.PublicKeyFromBase58(poolBase)
	if err != nil {
		panic(err)
	}

	poolAddress, poolBump, err := solana.FindProgramAddress([][]byte{
		[]byte(types.PoolSeed),
		poolBasePubkey.Bytes(),
	}, programId)
	_ = poolBump
	if err != nil {
		panic(err)
	}

	rpcClient := rpc.New(conf.Endpoint)
	ctx := context.Background()

	var pool types.Pool
	err = utils.GetAccountParseWithLimit(ctx, rpcClient, poolAddress, &pool, types.PoolDataSize)
	if err != nil {
		panic(err)
	}

	var mint splToken.Mint
	err = utils.GetAccountParse(ctx, rpcClient, pool.PoolMint, &mint)
	if err != nil {
		panic(err)
	}

	tStartDate, err := time.Parse(utils.RFC3339FullDate, startDate)
	if err != nil {
		panic(err)
	}

	tEndDate, err := time.Parse(utils.RFC3339FullDate, endDate)
	if err != nil {
		panic(err)
	}
	tEndDate = tEndDate.AddDate(0, 0, 1)

	users := []string{}
	initialized := false

	for periodNumber := uint(1); periodNumber <= uint(pool.PoolPeriodNumber); periodNumber++ {
		result, err := snapshots.FetchSnapshot(ctx, poolBase, periodNumber)
		if err != nil {
			panic(err)
		}

		if tStartDate.Unix() > result.Timestamp {
			continue
		}

		if result.Timestamp > tEndDate.Unix() {
			break
		}

		if !initialized {
			for _, user := range result.Users {
				users = append(users, user.User)
			}
			initialized = true
			continue
		}

		temp := users
		for i, userAddress := range temp {
			ok := false
			for _, user := range result.Users {
				if user.User == userAddress {
					amount, err := strconv.ParseUint(user.Amount, 10, 64)
					if err != nil {
						panic(err)
					}

					if amount >= uint64(minimalBalance*math.Pow(10, float64(mint.Decimals))) {
						ok = true
					} else {
						break
					}
				}
			}
			if !ok {
				users[i] = users[len(users)-1]
				users = users[:len(users)-1]
			}
		}
	}

	if !initialized {
		panic(errors.New("there's no valid snapshot in this range of date"))
		return
	}

	result := struct {
		PoolBase       string   `json:"poolBase"`
		StartDate      string   `json:"startDate"`
		EndDate        string   `json:"endDate"`
		MinimalBalance float64  `json:"minimalBalance"`
		TakenOn        string   `json:"takenOn"`
		Result         []string `json:"result"`
	}{
		PoolBase:       poolBase,
		StartDate:      startDate,
		EndDate:        endDate,
		MinimalBalance: minimalBalance,
		TakenOn:        time.Now().String(),
		Result:         users,
	}

	buff, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("scrapper_%s_%d.json", poolBase, time.Now().Unix())

	err = os.WriteFile(fileName, buff, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully write to file: " + fileName)
}
