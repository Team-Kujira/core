# kujira

**kujira** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

Please refer to [docs.kujira.app](https://docs.kujira.app/) and join our [Discord](https://t.co/kur923FTZk) for guidance in getting set up.

## Get started

Ensure your ignite version is 0.26.1

```
ignite version
```

If not, you can download the latest. See [the docs](https://docs.ignite.com/welcome/install#verify-your-ignite-cli-version) for more info

```
curl https://get.ignite.com/cli | bash
sudo mv ignite /usr/local/bin/
```

Then start the chain

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.
