package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	bc *BlockChain
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	interator := BlockchainInterator{cli.bc.tip, cli.bc.db}
	for {

		block := interator.Readblock()
		fmt.Printf("\n")
		//fmt.Printf("The %dth block:\n",)
		fmt.Printf("Prev. hash: %x\n", block.PreBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		if len(interator.hashofBlocktoRead) == 0 {
			break
		}
	}

}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")

}

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("no parameter")
		os.Exit(1)
	}

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "block data") // return *string, 这里面的data不加-!
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:]) // parse是用来分析输入参数的
		if err != nil {
			log.Panic(err)
		}

	case "printchain":

		err := printChainCmd.Parse(os.Args[2:]) // parse是用来分析输入参数的
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)

	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}
