package blobfactory

import (
	"context"

	"github.com/celestiaorg/celestia-app/testutil/namespace"
	"github.com/celestiaorg/celestia-app/testutil/testfactory"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/rand"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	coretypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
)

var defaultSigner = testfactory.RandomAddress().String()

func RandMsgPayForBlobWithSigner(singer string, size int) (*blobtypes.MsgPayForBlob, []byte) {
	blob := tmrand.Bytes(size)
	msg, err := blobtypes.NewMsgPayForBlob(
		singer,
		namespace.RandomBlobNamespace(),
		blob,
	)
	if err != nil {
		panic(err)
	}
	return msg, blob
}

func RandMsgPayForBlobWithNamespaceAndSigner(signer string, nid []byte, size int) (*blobtypes.MsgPayForBlob, []byte) {
	blob := tmrand.Bytes(size)
	msg, err := blobtypes.NewMsgPayForBlob(
		signer,
		nid,
		blob,
	)
	if err != nil {
		panic(err)
	}
	return msg, blob
}

func RandMsgPayForBlob(size int) (*blobtypes.MsgPayForBlob, []byte) {
	blob := tmrand.Bytes(size)
	msg, err := blobtypes.NewMsgPayForBlob(
		defaultSigner,
		namespace.RandomBlobNamespace(),
		blob,
	)
	if err != nil {
		panic(err)
	}
	return msg, blob
}

func RandBlobTxsRandomlySized(enc sdk.TxEncoder, count, maxSize int) []coretypes.Tx {
	const acc = "signer"
	kr := testfactory.GenerateKeyring(acc)
	signer := blobtypes.NewKeyringSigner(kr, acc, "chainid")
	addr, err := signer.GetSignerInfo().GetAddress()
	if err != nil {
		panic(err)
	}

	coin := sdk.Coin{
		Denom:  bondDenom,
		Amount: sdk.NewInt(10),
	}

	opts := []blobtypes.TxBuilderOption{
		blobtypes.SetFeeAmount(sdk.NewCoins(coin)),
		blobtypes.SetGasLimit(10000000),
	}

	txs := make([]coretypes.Tx, count)
	for i := 0; i < count; i++ {
		// pick a random non-zero size of max maxSize
		size := tmrand.Intn(maxSize)
		if size == 0 {
			size = 1
		}
		msg, blob := RandMsgPayForBlobWithSigner(addr.String(), size)
		builder := signer.NewTxBuilder(opts...)
		stx, err := signer.BuildSignedTx(builder, msg)
		if err != nil {
			panic(err)
		}
		rawTx, err := enc(stx)
		if err != nil {
			panic(err)
		}
		wblob, err := blobtypes.NewBlob(msg.NamespaceId, blob)
		if err != nil {
			panic(err)
		}
		cTx, err := coretypes.WrapBlobTx(rawTx, wblob)
		if err != nil {
			panic(err)
		}
		txs[i] = cTx
	}

	return txs
}

func RandBlobTxsWithAccounts(
	enc sdk.TxEncoder,
	kr keyring.Keyring,
	conn *grpc.ClientConn,
	size int,
	randSize bool,
	chainid string,
	accounts []string,
) []coretypes.Tx {
	coin := sdk.Coin{
		Denom:  bondDenom,
		Amount: sdk.NewInt(10),
	}

	opts := []blobtypes.TxBuilderOption{
		blobtypes.SetFeeAmount(sdk.NewCoins(coin)),
		blobtypes.SetGasLimit(100000000000000),
	}

	txs := make([]coretypes.Tx, len(accounts))
	for i := 0; i < len(accounts); i++ {
		signer := blobtypes.NewKeyringSigner(kr, accounts[i], chainid)
		err := signer.QueryAccountNumber(context.Background(), conn)
		if err != nil {
			panic(err)
		}

		addr, err := signer.GetSignerInfo().GetAddress()
		if err != nil {
			panic(err)
		}

		randomizedSize := size
		if randSize {
			randomizedSize = rand.Intn(size)
			if randomizedSize == 0 {
				randomizedSize = 1
			}
		}
		msg, blob := RandMsgPayForBlobWithSigner(addr.String(), randomizedSize)
		builder := signer.NewTxBuilder(opts...)
		stx, err := signer.BuildSignedTx(builder, msg)
		if err != nil {
			panic(err)
		}
		rawTx, err := enc(stx)
		if err != nil {
			panic(err)
		}
		wblob, err := blobtypes.NewBlob(msg.NamespaceId, blob)
		if err != nil {
			panic(err)
		}
		cTx, err := coretypes.WrapBlobTx(rawTx, wblob)
		if err != nil {
			panic(err)
		}
		txs[i] = cTx
	}

	return txs
}

func RandBlobTxs(enc sdk.TxEncoder, count, size int) []coretypes.Tx {
	const acc = "signer"
	kr := testfactory.GenerateKeyring(acc)
	signer := blobtypes.NewKeyringSigner(kr, acc, "chainid")
	addr, err := signer.GetSignerInfo().GetAddress()
	if err != nil {
		panic(err)
	}

	coin := sdk.Coin{
		Denom:  bondDenom,
		Amount: sdk.NewInt(10),
	}

	opts := []blobtypes.TxBuilderOption{
		blobtypes.SetFeeAmount(sdk.NewCoins(coin)),
		blobtypes.SetGasLimit(10000000),
	}

	txs := make([]coretypes.Tx, count)
	for i := 0; i < count; i++ {
		msg, blob := RandMsgPayForBlobWithSigner(addr.String(), size)
		builder := signer.NewTxBuilder(opts...)
		stx, err := signer.BuildSignedTx(builder, msg)
		if err != nil {
			panic(err)
		}
		rawTx, err := enc(stx)
		if err != nil {
			panic(err)
		}
		wblob, err := blobtypes.NewBlob(msg.NamespaceId, blob)
		if err != nil {
			panic(err)
		}
		cTx, err := coretypes.WrapBlobTx(rawTx, wblob)
		if err != nil {
			panic(err)
		}
		txs[i] = cTx
	}

	return txs
}

func RandBlobTxsWithNamespaces(enc sdk.TxEncoder, nIds [][]byte, sizes []int) []coretypes.Tx {
	const acc = "signer"
	kr := testfactory.GenerateKeyring(acc)
	signer := blobtypes.NewKeyringSigner(kr, acc, "chainid")
	return RandBlobTxsWithNamespacesAndSigner(enc, signer, nIds, sizes)
}

func RandBlobTxsWithNamespacesAndSigner(
	enc sdk.TxEncoder,
	signer *blobtypes.KeyringSigner,
	nIds [][]byte,
	sizes []int,
) []coretypes.Tx {
	addr, err := signer.GetSignerInfo().GetAddress()
	if err != nil {
		panic(err)
	}

	coin := sdk.Coin{
		Denom:  bondDenom,
		Amount: sdk.NewInt(10),
	}

	opts := []blobtypes.TxBuilderOption{
		blobtypes.SetFeeAmount(sdk.NewCoins(coin)),
		blobtypes.SetGasLimit(10000000),
	}

	txs := make([]coretypes.Tx, len(nIds))
	for i := 0; i < len(nIds); i++ {
		msg, blob := RandMsgPayForBlobWithNamespaceAndSigner(addr.String(), nIds[i], sizes[i])
		builder := signer.NewTxBuilder(opts...)
		stx, err := signer.BuildSignedTx(builder, msg)
		if err != nil {
			panic(err)
		}
		rawTx, err := enc(stx)
		if err != nil {
			panic(err)
		}
		wblob, err := blobtypes.NewBlob(msg.NamespaceId, blob)
		if err != nil {
			panic(err)
		}
		cTx, err := coretypes.WrapBlobTx(rawTx, wblob)
		if err != nil {
			panic(err)
		}
		txs[i] = cTx
	}

	return txs
}
