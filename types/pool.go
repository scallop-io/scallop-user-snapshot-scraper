package types

import "github.com/gagliardetto/solana-go"

type Pool struct {
	Discriminator    [8]byte
	AccountType      [1]byte
	PoolBase         solana.PublicKey
	PoolType         [1]byte
	PoolMarket       solana.PublicKey
	PoolAuthority    solana.PublicKey
	PoolVaultToken   solana.PublicKey
	PoolRewardToken  solana.PublicKey
	PoolMint         solana.PublicKey
	PoolPeriodNumber uint16
}

const PoolDataSize = 8 + // discriminator
	1 + // accountType
	32 + // poolBase
	1 + // poolType
	32*5 + // poolMarket - poolMint
	4 // poolPeriodNumber

const PoolSeed = "pool_seed"
