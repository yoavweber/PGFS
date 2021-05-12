package main

import (
	"PGFS/content"
	"PGFS/global"
	"PGFS/node"
	"PGFS/peers"
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"log"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Spawn a node using a temporary path, creating a temporary repo for the run
	log.Println("Spawning node on " + global.RepoPath)
	node, err := spawnNode(ctx, true)
	if err != nil {
		panic(err)
	}

	// Node identity information
	log.Println("Node spawned on " + global.RepoPath + "\nIdentity information:")
	key, _ := node.Key().Self(ctx)
	log.Println(" PeerID: " + key.ID().Pretty() + "\n Path: " + key.Path().String())

	var bootstrapNodes = []string {
		//"/ip4/10.22.201.110/tcp/4001/ipfs/Qme5aJduvtCwQ4Hojdvez9DA7FnqABo6wTeKTGtmEWGHDo",
	}

	time.Sleep(20*time.Second)

	go peers.ConnectToPeers(ctx, node, bootstrapNodes)

	time.Sleep(20*time.Second)

	cid := "QmS98pgfsLTc91kjHDzb5V9nCwXbs9pe2h2Kj8zsVCXEmR"

	filePath, err := content.GetContent(cid, node, ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Content with CID: " + cid + "\nreceived and written to " + filePath)

	time.Sleep(10*time.Second)



}

// Spawns a node
func spawnNode(ctx context.Context, isServer bool) (icore.CoreAPI, error) {
	if err := node.SetupPlugins(""); err != nil {
		return nil, err
	}

	// Checks if repo is initialized
	if !fsrepo.IsInitialized(global.RepoPath) {
		// Initializes repo in repoPath
		if err := node.RepoInit(); err != nil {
			return nil, err
		}
	}

	// Opens the repo
	nodeRepo, err := fsrepo.Open(global.RepoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo on node creation: %s", err)
	}

	// Spawns an IPFS node
	return node.CreateNode(ctx, nodeRepo, isServer)

}





