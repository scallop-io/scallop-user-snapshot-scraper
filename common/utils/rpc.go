package utils

import (
	"context"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetAccountParseWithLimit(ctx context.Context, rpcClient *rpc.Client, address solana.PublicKey, data interface{}, limit int) error {
	resp, err := rpcClient.GetAccountInfo(
		context.TODO(),
		address,
	)
	if err != nil {
		return err
	}

	var borshDec *bin.Decoder
	if limit == -1 {
		borshDec = bin.NewBorshDecoder(resp.Value.Data.GetBinary())
	} else {
		borshDec = bin.NewBorshDecoder(resp.Value.Data.GetBinary()[:limit])
	}
	err = borshDec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func GetAccountParse(ctx context.Context, rpcClient *rpc.Client, address solana.PublicKey, data interface{}) error {
	return GetAccountParseWithLimit(ctx, rpcClient, address, data, -1)
}
