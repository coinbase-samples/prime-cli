# Prime CLI

## Overview

A Coinbase Prime command-line interface (CLI) application that generates requests to and receives responses from [Coinbase Prime's](https://prime.coinbase.com/) [REST APIs](https://docs.cloud.coinbase.com/prime/reference).

## License

The Prime CLI is free and open source and released under the [Apache License, Version 2.0](LICENSE.txt).

The application and code are only available for demonstration purposes.

## Usage

To begin, navigate to your preferred directory for development and clone the Prime CLI repository and enter the directory using the following commands: 

```
git clone https://github.com/coinbase-samples/prime-cli
cd prime-cli
```

Next, pass an environment variable via your terminal called `PRIME_CREDENTIALS` with your API and portfolio information. 

Coinbase Prime API credentials can be created in the Prime web console under Settings -> APIs. Entity ID can be retrieved by calling [Get Portfolio](https://docs.cloud.coinbase.com/prime/reference/primerestapi_getportfolio). If you are not configured yet to call this endpoint, you may proceed without including the `entityId` key and value for the time being, but certain endpoints such as List Invoices and List Assets require it. 

`PRIME_CREDENTIALS` should match the following format:
```
export PRIME_CREDENTIALS='{
"accessKey":"ACCESSKEY_HERE",
"passphrase":"PASSPHRASE_HERE",
"signingKey":"SIGNINGKEY_HERE",
"portfolioId":"PORTFOLIOID_HERE",
"svcAccountId":"SVCACCOUNTID_HERE",
"entityId":"ENTITYID_HERE"
}'
```

You may also pass an environment variable called `primeCliTimeout` which will override the default request timeout of 7 seconds. This value should be an integer in seconds. 

Build the application binary and specify an output name, e.g. `primectl`:

```
go build -o primectl
```

To ensure your project's dependencies are up-to-date, run:
```
go mod tidy
```

To verify that your application is installed correctly and accessible from any location, run the following command. It will include all available requests:

```
./primectl
```

Finally, to run commands for each endpoint, use the following format to test each endpoint. Please note that many endpoints require flags, which are detailed with the `--help` flag. 

```
./primectl portfolios list
```

```
./primectl orders create --help
```

```
./primectl orders create-preview -b 0.001 -i ETH-USD -s BUY -t MARKET
```
