# DeFiLending CLI
This project provides a Command-Line Interface (CLI) written in Go to interact with a deployed **DeFiLending Smart Contract**. This application facilitates interacting with the smart contract for depositing tokens and checking deposit balances.
## Prerequisites
Before running the application, ensure you have the following:
1. **Deployed DeFiLending smart contract**: Obtain the contract address after deployment.
2. **USDC Contract Address**: Identify the address of the compatible **ERC20 token** (such as USDC) used for deposits in the lending contract.
3. **Ethereum RPC URL**: The URL of an Ethereum node (for example: Infura, Alchemy, or a self-hosted Ethereum client).
4. **Private Key**: Required for signing transactions during deposit operations.
5. **Go Binding for DeFiLending Contract**: These must be generated using [abigen]().

## Features
This CLI tool supports the following operations:
1. **Deposit tokens** into the DeFiLending contract.
2. Retrieve the **total amount of deposits** held in the contract.
3. Check the **deposit balance** of a specific user.

## Installation
### 1. Clone the repository
``` bash
git clone <repository-url>
cd <repository-folder>
```
### 2. Install dependencies
Make sure you have [Go installed]() (1.23 or later). Install project dependencies:
``` bash
go mod tidy
```
### 3. Set up environment variables
The application requires an `RPC_URL` environment variable pointing to your Ethereum RPC URL. You can export it in your terminal session or add it to your environment configuration file.
``` bash
export RPC_URL="<your-ethereum-rpc-url>"
```
Alternatively, you can set it directly in your shell script or runtime environment.
## Command Usage
Run the binary or `go run` the program followed by the appropriate commands and flags.
### 1. **Deposit Tokens**
Deposits tokens (such as **USDC**) into the **DeFiLending contract**. You must have a valid private key and the amount to deposit in the token's **smallest unit (e.g., Wei for ERC20)**.
``` bash
go run main.go deposit --amount <amount> --private-key <private-key>
```
#### Arguments:
- `--amount`: Amount of tokens to deposit (e.g., `10` for 10 tokens). The script multiplies this value by `1e6` (based on the token decimals set in the USDC contract).
- `--private-key`: The user's Ethereum private key used to sign the transaction. It must be in hexadecimal format.

Example:
``` bash
go run main.go deposit --amount 100 --private-key 0x<private-key>
```
#### Steps Within the Process:
1. **Approval**: Approves the DeFiLending smart contract to spend the token amount on your behalf.
2. **Event Monitoring**: The program waits for the Approval event confirmation from the token contract.
3. **Deposit Execution**: Calls the `deposit` function from the DeFiLending contract.

#### Expected Output:
- Approval transaction hash
- Deposit transaction hash

### 2. **Check Total Deposits**
Retrieve the total amount of deposits stored in the DeFiLending contract.
``` bash
go run main.go total
```
#### Expected Output:
The total amount of tokens deposited in the contract.
Example:
``` bash
Total Deposits: 10000000
```
### 3. **Check Deposit for a User**
Retrieve the deposit balance of a specific user.
``` bash
go run main.go user --address <user-address>
```
#### Arguments:
- `--address`: The Ethereum address of the user for whom you want to check the deposit balance.

Example:
``` bash
go run main.go user --address 0x123456789ABCDEF123456789ABCDEF123456789A
```
#### Expected Output:
The deposit balance of the specified user.
## Environment Variables
The following environment variables must be set before running the CLI:
- **`RPC_URL` **: Ethereum node RPC URL for interactions with the blockchain (e.g., `https://mainnet.infura.io/v3/<your-project-id>`).

## Project Configuration
Modify the following constants in the code for your specific deployment:
``` go
const (
	contractAddress     = "0x<your-DeFiLending-address>" // Deployed smart contract address
	usdcContractAddress = "0x<your-usdc-token-address>"  // ERC20 token (e.g., USDC) contract address
)
```
Update them to match your own smart contract and token configuration.
## Dependencies
This CLI tool leverages the following Go libraries:
1. **Geth (go-ethereum)**:
    - Provides Ethereum client connection and ABI bindings.
    - Install via: `go get github.com/ethereum/go-ethereum`

2. **Big Integer Support**:
    - Used for handling large numbers like token amounts.
    - Go's native `math/big` package.

3. **DeFiLending Binding**:
    - Generated Go bindings for your DeFiLending smart contract using `abigen`.
    - Command: `abigen --sol DeFiLending.sol --pkg defi --out defi.go`

4. **USDC Binding**:
    - Similarly, generate bindings for your ERC20 token contract.

## Notes
1. Ensure your account has sufficient funds (ETH) to pay gas fees for deposit operations.
2. Make sure the DeFiLending contract follows the expected interface for Deposit functionality.
3. Event monitoring may timeout (default: 60 seconds), so ensure network connectivity and block confirmation speed.

## License
This project is licensed under the MIT License. Feel free to modify and use it as you see fit.
This README provides a comprehensive guide to understanding, installing, and using the tool effectively.
