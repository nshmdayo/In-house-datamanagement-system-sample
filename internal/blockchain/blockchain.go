package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	ID         string                 `json:"id"`
	DocumentID uint                   `json:"document_id"`
	UserID     uint                   `json:"user_id"`
	Action     string                 `json:"action"`
	Data       map[string]interface{} `json:"data"`
	Timestamp  time.Time              `json:"timestamp"`
	Hash       string                 `json:"hash"`
}

// Block represents a block in the blockchain
type Block struct {
	Index        int64         `json:"index"`
	Timestamp    time.Time     `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PreviousHash string        `json:"previous_hash"`
	Hash         string        `json:"hash"`
	Nonce        int64         `json:"nonce"`
	MerkleRoot   string        `json:"merkle_root"`
}

// Blockchain represents the blockchain
type Blockchain struct {
	Blocks     []Block `json:"blocks"`
	Difficulty int     `json:"difficulty"`
}

// NewBlockchain creates a new blockchain with genesis block
func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Blocks:     make([]Block, 0),
		Difficulty: 4, // Number of leading zeros required in hash
	}

	// Create genesis block
	genesisBlock := bc.createGenesisBlock()
	bc.Blocks = append(bc.Blocks, genesisBlock)

	return bc
}

// createGenesisBlock creates the first block in the blockchain
func (bc *Blockchain) createGenesisBlock() Block {
	genesisTransaction := Transaction{
		ID:        "genesis",
		Action:    "genesis",
		Data:      map[string]interface{}{"message": "Genesis block"},
		Timestamp: time.Now(),
	}
	genesisTransaction.Hash = bc.calculateTransactionHash(genesisTransaction)

	block := Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: []Transaction{genesisTransaction},
		PreviousHash: "0",
		Nonce:        0,
	}

	block.MerkleRoot = bc.calculateMerkleRoot(block.Transactions)
	block.Hash = bc.mineBlock(&block)

	return block
}

// AddTransaction adds a new transaction to the blockchain
func (bc *Blockchain) AddTransaction(transaction Transaction) error {
	// Calculate transaction hash
	transaction.Hash = bc.calculateTransactionHash(transaction)

	// Get the latest block
	latestBlock := bc.getLatestBlock()

	// Create new block
	newBlock := Block{
		Index:        latestBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: []Transaction{transaction},
		PreviousHash: latestBlock.Hash,
		Nonce:        0,
	}

	// Calculate Merkle root
	newBlock.MerkleRoot = bc.calculateMerkleRoot(newBlock.Transactions)

	// Mine the block
	newBlock.Hash = bc.mineBlock(&newBlock)

	// Add block to blockchain
	bc.Blocks = append(bc.Blocks, newBlock)

	return nil
}

// mineBlock mines a block using proof of work
func (bc *Blockchain) mineBlock(block *Block) string {
	target := bc.getTarget()

	for {
		hash := bc.calculateBlockHash(block)
		if bc.isValidHash(hash, target) {
			return hash
		}
		block.Nonce++
	}
}

// calculateBlockHash calculates the hash of a block
func (bc *Blockchain) calculateBlockHash(block *Block) string {
	data := fmt.Sprintf("%d%s%s%s%d",
		block.Index,
		block.Timestamp.Format(time.RFC3339),
		block.PreviousHash,
		block.MerkleRoot,
		block.Nonce,
	)

	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// calculateTransactionHash calculates the hash of a transaction
func (bc *Blockchain) calculateTransactionHash(tx Transaction) string {
	// Create a copy without the hash field
	txCopy := Transaction{
		ID:         tx.ID,
		DocumentID: tx.DocumentID,
		UserID:     tx.UserID,
		Action:     tx.Action,
		Data:       tx.Data,
		Timestamp:  tx.Timestamp,
	}

	data, _ := json.Marshal(txCopy)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// calculateMerkleRoot calculates the Merkle root of transactions
func (bc *Blockchain) calculateMerkleRoot(transactions []Transaction) string {
	if len(transactions) == 0 {
		return ""
	}

	if len(transactions) == 1 {
		return transactions[0].Hash
	}

	// For simplicity, we'll just hash all transaction hashes together
	var allHashes string
	for _, tx := range transactions {
		allHashes += tx.Hash
	}

	hash := sha256.Sum256([]byte(allHashes))
	return fmt.Sprintf("%x", hash)
}

// getTarget returns the target for proof of work
func (bc *Blockchain) getTarget() string {
	target := ""
	for i := 0; i < bc.Difficulty; i++ {
		target += "0"
	}
	return target
}

// isValidHash checks if a hash meets the difficulty requirement
func (bc *Blockchain) isValidHash(hash, target string) bool {
	return hash[:len(target)] == target
}

// getLatestBlock returns the latest block in the blockchain
func (bc *Blockchain) getLatestBlock() Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

// ValidateChain validates the entire blockchain
func (bc *Blockchain) ValidateChain() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]

		// Validate current block hash
		if currentBlock.Hash != bc.calculateBlockHash(&currentBlock) {
			return false
		}

		// Validate link to previous block
		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}

		// Validate proof of work
		target := bc.getTarget()
		if !bc.isValidHash(currentBlock.Hash, target) {
			return false
		}

		// Validate Merkle root
		if currentBlock.MerkleRoot != bc.calculateMerkleRoot(currentBlock.Transactions) {
			return false
		}

		// Validate all transactions in the block
		for _, tx := range currentBlock.Transactions {
			if tx.Hash != bc.calculateTransactionHash(tx) {
				return false
			}
		}
	}

	return true
}

// GetTransactionHistory returns all transactions for a document
func (bc *Blockchain) GetTransactionHistory(documentID uint) []Transaction {
	var transactions []Transaction

	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			if tx.DocumentID == documentID {
				transactions = append(transactions, tx)
			}
		}
	}

	return transactions
}

// GetUserTransactions returns all transactions for a user
func (bc *Blockchain) GetUserTransactions(userID uint) []Transaction {
	var transactions []Transaction

	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			if tx.UserID == userID {
				transactions = append(transactions, tx)
			}
		}
	}

	return transactions
}

// GetBlockByIndex returns a block by its index
func (bc *Blockchain) GetBlockByIndex(index int64) (*Block, error) {
	if index < 0 || index >= int64(len(bc.Blocks)) {
		return nil, fmt.Errorf("block index out of range")
	}

	return &bc.Blocks[index], nil
}

// GetTransactionByID returns a transaction by its ID
func (bc *Blockchain) GetTransactionByID(txID string) (*Transaction, error) {
	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			if tx.ID == txID {
				return &tx, nil
			}
		}
	}

	return nil, fmt.Errorf("transaction not found")
}

// GetChainInfo returns information about the blockchain
func (bc *Blockchain) GetChainInfo() map[string]interface{} {
	totalTransactions := 0
	for _, block := range bc.Blocks {
		totalTransactions += len(block.Transactions)
	}

	return map[string]interface{}{
		"blocks":             len(bc.Blocks),
		"total_transactions": totalTransactions,
		"difficulty":         bc.Difficulty,
		"latest_block_hash":  bc.getLatestBlock().Hash,
		"is_valid":           bc.ValidateChain(),
	}
}

// CreateDocumentTransaction creates a transaction for document operations
func CreateDocumentTransaction(txID string, documentID, userID uint, action string, data map[string]interface{}) Transaction {
	return Transaction{
		ID:         txID,
		DocumentID: documentID,
		UserID:     userID,
		Action:     action,
		Data:       data,
		Timestamp:  time.Now(),
	}
}

// GenerateTransactionID generates a unique transaction ID
func GenerateTransactionID(documentID, userID uint, action string) string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%d_%d_%s_%d", documentID, userID, action, timestamp)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)[:16] // Use first 16 characters
}
