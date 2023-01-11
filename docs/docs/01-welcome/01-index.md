---
slug: /
---

import ProjectsTable from '@site/src/components/ProjectsTable';

# Introduction to Ignite

```patch (obtained with diff -W to preserve func body)
diff --git a/x/mynewchain/module_simulation.go b/x/mynewchain/module_simulation.go
index b10fbc9..a545f21 100644
--- a/x/mynewchain/module_simulation.go
+++ b/x/mynewchain/module_simulation.go
@@ -2,16 +2,17 @@ package mynewchain
 
 import (
 	"math/rand"
+	"xx"
 
 	"github.com/cosmos/cosmos-sdk/baseapp"
 	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
 	sdk "github.com/cosmos/cosmos-sdk/types"
 	"github.com/cosmos/cosmos-sdk/types/module"
 	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
 	"github.com/cosmos/cosmos-sdk/x/simulation"
 	"my-new-chain/testutil/sample"
 	mynewchainsimulation "my-new-chain/x/mynewchain/simulation"
 	"my-new-chain/x/mynewchain/types"
 )
 
 // avoid unused import issue
@@ -24,20 +25,43 @@ var (
 )
 
 const (
-// this line is used by starport scaffolding # simapp/module/const
+	opWeightMsgCreateMnylist = "op_weight_msg_mnylist"
+	// TODO: Determine the simulation weight value
+	defaultWeightMsgCreateMnylist int = 100
+
+	opWeightMsgUpdateMnylist = "op_weight_msg_mnylist"
+	// TODO: Determine the simulation weight value
+	defaultWeightMsgUpdateMnylist int = 100
+
+	opWeightMsgDeleteMnylist = "op_weight_msg_mnylist"
+	// TODO: Determine the simulation weight value
+	defaultWeightMsgDeleteMnylist int = 100
+
+	// this line is used by starport scaffolding # simapp/module/const
 )
 
 // GenerateGenesisState creates a randomized GenState of the module
 func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
 	accs := make([]string, len(simState.Accounts))
 	for i, acc := range simState.Accounts {
 		accs[i] = acc.Address.String()
 	}
 	mynewchainGenesis := types.GenesisState{
 		Params: types.DefaultParams(),
+		MnylistList: []types.Mnylist{
+			{
+				Id:      0,
+				Creator: sample.AccAddress(),
+			},
+			{
+				Id:      1,
+				Creator: sample.AccAddress(),
+			},
+		},
+		MnylistCount: 2,
 		// this line is used by starport scaffolding # simapp/module/genesisState
 	}
 	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&mynewchainGenesis)
 }
 
 // ProposalContents doesn't return any content functions for governance proposals
@@ -57,8 +81,41 @@ func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}
 // WeightedOperations returns the all the gov module operations with their respective weights.
 func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
 	operations := make([]simtypes.WeightedOperation, 0)
 
+	var weightMsgCreateMnylist int
+	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateMnylist, &weightMsgCreateMnylist, nil,
+		func(_ *rand.Rand) {
+			weightMsgCreateMnylist = defaultWeightMsgCreateMnylist
+		},
+	)
+	operations = append(operations, simulation.NewWeightedOperation(
+		weightMsgCreateMnylist,
+		mynewchainsimulation.SimulateMsgCreateMnylist(am.accountKeeper, am.bankKeeper, am.keeper),
+	))
+
+	var weightMsgUpdateMnylist int
+	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateMnylist, &weightMsgUpdateMnylist, nil,
+		func(_ *rand.Rand) {
+			weightMsgUpdateMnylist = defaultWeightMsgUpdateMnylist
+		},
+	)
+	operations = append(operations, simulation.NewWeightedOperation(
+		weightMsgUpdateMnylist,
+		mynewchainsimulation.SimulateMsgUpdateMnylist(am.accountKeeper, am.bankKeeper, am.keeper),
+	))
+
+	var weightMsgDeleteMnylist int
+	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteMnylist, &weightMsgDeleteMnylist, nil,
+		func(_ *rand.Rand) {
+			weightMsgDeleteMnylist = defaultWeightMsgDeleteMnylist
+		},
+	)
+	operations = append(operations, simulation.NewWeightedOperation(
+		weightMsgDeleteMnylist,
+		mynewchainsimulation.SimulateMsgDeleteMnylist(am.accountKeeper, am.bankKeeper, am.keeper),
+	))
+
 	// this line is used by starport scaffolding # simapp/module/operation
 
 	return operations
 }
```

[Ignite CLI](https://github.com/ignite/cli) offers everything you need to build, test, and launch your blockchain with a
decentralized worldwide community. Ignite CLI is built on top of [Cosmos SDK](https://docs.cosmos.network), the world’s
most popular blockchain framework. Ignite CLI accelerates chain development by scaffolding everything you need so you
can focus on business logic.

## What is Ignite CLI?

Ignite CLI is an easy-to-use CLI tool for creating and maintaining sovereign application-specific blockchains.
Blockchains created with Ignite CLI use Cosmos SDK and Tendermint. Ignite CLI and the Cosmos SDK modules are written in
the Go programming language. The scaffolded blockchain that is created with Ignite CLI includes a command line interface
that lets you manage keys, create validators, and send tokens.

With just a few commands, you can use Ignite CLI to:

- Create a modular blockchain written in Go
- Scaffold modules, messages, types with CRUD operations, IBC packets, and more
- Start a blockchain node in development with live reloading
- Connect to other blockchains with a built-in IBC relayer
- Use generated TypeScript/Vuex clients to interact with your blockchain
- Use the Vue.js web app template with a set of components and Vuex modules

## Install Ignite CLI

To install the `ignite` binary in `/usr/local/bin` run the following command:

```
curl https://get.ignite.com/cli | bash
```

## Projects using Tendermint and Cosmos SDK

Many projects already showcase the Tendermint BFT consensus engine and the Cosmos SDK. Explore
the [Cosmos ecosystem](https://cosmos.network/ecosystem/apps) to discover a wide variety of apps, blockchains, wallets,
and explorers that are built in the Cosmos ecosystem.

## Projects building with Ignite CLI

<ProjectsTable data={[
  { name: "Stride Labs", logo: "img/logo/stride.svg"},
  { name: "KYVE Network", logo: "img/logo/kyve.svg"},
  { name: "Umee", logo: "img/logo/umee.svg"},
  { name: "MediBloc Core", logo: "img/logo/medibloc.svg"},
  { name: "Cudos", logo: "img/logo/cudos.svg"},
  { name: "Firma Chain", logo: "img/logo/firmachain.svg"},
  { name: "BitCanna", logo: "img/logo/bitcanna.svg"},
  { name: "Source Protocol", logo: "img/logo/source.svg"},
  { name: "Sonr", logo: "img/logo/sonr.svg"},
  { name: "Neutron", logo: "img/logo/neutron.svg"},
  { name: "OKP4 Blockchain", logo: "img/logo/okp4.svg"},
  { name: "Dymension Hub", logo: "img/logo/dymension.svg"},
  { name: "Electra Blockchain", logo: "img/logo/electra.svg"},
  { name: "OLLO Station", logo: "img/logo/ollostation.svg"},
  { name: "Mun", logo: "img/logo/mun.svg"},
  { name: "Aura Network", logo: "img/logo/aura.svg"},
]}/>
