package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"defi-lending/defi" // Go binding package for your DeFiLending contract
	"defi-lending/usdc"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Configuration constants; update these with your values.
const (
	contractAddress     = "0x6b338b0ab70B08ABEf6F4344F8dB3Bd3e42591Cc" // Your deployed DeFiLending contract address
	usdcContractAddress = "0xae624D2005c193aA546e29Ecc3346307A3dDfdD2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'deposit', 'total', or 'user' subcommand")
		os.Exit(1)
	}

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		log.Fatal("RPC_URL environment variable not set")
	}

	// Connect to Ethereum client.
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to Ethereum client:", err)
	}

	// Create an instance of the lending contract.
	lending, err := defi.NewDefi(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal("Failed to load DeFiLending contract:", err)
	}

	// Instantiate the custom uSDC token contract.
	usdcToken, err := usdc.NewUsdc(common.HexToAddress(usdcContractAddress), client)
	if err != nil {
		log.Fatal("Failed to load uSDC token contract:", err)
	}

	switch os.Args[1] {

	// Deposit subcommand: send a deposit transaction. Requires amount and private-key.
	case "deposit":
		depositCmd := flag.NewFlagSet("deposit", flag.ExitOnError)
		amountFlag := depositCmd.String("amount", "", "Amount to deposit (e.g., '10' for 10 tokens)")
		privateKeyFlag := depositCmd.String("private-key", "", "Private key for signing the transaction")
		depositCmd.Parse(os.Args[2:])

		if *amountFlag == "" || *privateKeyFlag == "" {
			fmt.Println("Usage: deposit --amount <amount> --private-key <private-key>")
			os.Exit(1)
		}

		// Convert input amount (assumed to be a whole number) to the token's smallest unit.
		amt, ok := new(big.Int).SetString(*amountFlag, 10)
		if !ok {
			log.Fatal("Invalid amount provided")
		}
		// For a token with 6 decimals (like USDC), multiply by 1e6.
		multiplier := big.NewInt(1e6)
		depositAmount := new(big.Int).Mul(amt, multiplier)

		// Create an authorized transactor using the provided private key.
		privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(*privateKeyFlag, "0x"))
		if err != nil {
			log.Fatal("Invalid private key:", err)
		}

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal("Failed to get chain ID:", err)
		}
		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			log.Fatal("Failed to create transactor:", err)
		}

		// Approve the DeFiLending contract to spend the deposit amount of uSDC tokens.
		txApprove, err := usdcToken.Approve(auth, common.HexToAddress(contractAddress), depositAmount)
		if err != nil {
			log.Fatal("Failed to approve uSDC spend:", err)
		}
		fmt.Println("Approve transaction sent, tx hash:", txApprove.Hash().Hex())

		// Set up an event subscription to wait for the Approval event from the USDC token.
		// The Approval event is standard in ERC20, so our generated binding should include a WatchApproval method.
		approvalEventCh := make(chan *usdc.UsdcApproval)
		watchOpts := &bind.WatchOpts{Context: context.Background()}
		// Filter for events where owner is the caller (auth.From) and spender is the lending contract.
		sub, err := usdcToken.WatchApproval(watchOpts, approvalEventCh, []common.Address{auth.From}, []common.Address{common.HexToAddress(contractAddress)})
		if err != nil {
			log.Fatal("Failed to subscribe to Approval events:", err)
		}
		defer sub.Unsubscribe()

		// Wait for the Approval event to be received.
		fmt.Println("Waiting for Approval event...")
		select {
		case evt := <-approvalEventCh:
			fmt.Printf("Received Approval event: owner=%s, spender=%s, value=%s\n", evt.Owner.Hex(), evt.Spender.Hex(), evt.Value.String())
		case err := <-sub.Err():
			log.Fatal("Subscription error:", err)
		case <-time.After(60 * time.Second):
			log.Fatal("Timed out waiting for Approval event")
		}

		// Once the Approval event is confirmed, call the deposit function on the DeFiLending contract.
		txDeposit, err := lending.Deposit(auth, depositAmount)
		if err != nil {
			log.Fatal("Failed to deposit uSDC:", err)
		}
		fmt.Println("Deposit transaction sent, tx hash:", txDeposit.Hash().Hex())

	// Total subcommand: read the total deposits in the contract.
	case "total":
		total, err := lending.TotalDeposits(&bind.CallOpts{})
		if err != nil {
			log.Fatal("Failed to get total deposits:", err)
		}
		fmt.Println("Total Deposits:", total)

	// User subcommand: read the deposit amount for a specific user.
	case "user":
		userCmd := flag.NewFlagSet("user", flag.ExitOnError)
		addressFlag := userCmd.String("address", "", "User address (e.g., 0x...)")
		userCmd.Parse(os.Args[2:])

		if *addressFlag == "" {
			fmt.Println("Please specify --address")
			os.Exit(1)
		}
		userAddr := common.HexToAddress(*addressFlag)
		userDeposit, err := lending.Deposits(&bind.CallOpts{}, userAddr)
		if err != nil {
			log.Fatal("Failed to get deposit for user:", err)
		}
		fmt.Printf("Deposit for user %s: %s\n", userAddr.Hex(), userDeposit.String())

	default:
		fmt.Println("Expected 'deposit', 'total', or 'user' subcommand")
		os.Exit(1)
	}
}
