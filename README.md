# Electricity Prices API
[![codecov](https://codecov.io/gh/daithihearn/electricity-prices-go/graph/badge.svg?token=I50D46PMGZ)](https://codecov.io/gh/daithihearn/electricity-prices-go)

An API for electricity prices in Spain. Data is scraped from [REData Api](https://www.ree.es/en/apidatos) and exposed via a [restful API](https://elec-api.daithiapp.com/).

This API supports a [dashboard application](https://preciosdelaelectricidad.es/), an Alexa Skill and Flash briefing.

In addition to the API there is a sync job that can be run separately to sync data from REData.

This repo is a replacement for the [original implementation](https://github.com/daithihearn/electricity-prices) written in Kotlin.

## Stack

- Go
- MongoDB

## API
To run locally you will need to have a MongoDB instance running. Update the `MONGO_URL` environment variable to point to your instance.
You will also require `make` to be installed.

Then to run locally simply run:

```bash
make run
```

To build the executable binaries locally run:

```bash
make build
```
The binaries will be installed in the build folder and can be run directly.

If you want to build the docker image run:
    
```bash
make image
```

## Sync Job
To run the sync locally run:

```bash
make sync
```

The `make build` command described in the API section will build both binaries.

The `make image` command described in the API section will build a single docker image for both the API and sync job.

To run the docker image for the sync job run like so:

```bash
docker run -d --rm elec-prices-sync ./sync
```

Running without the `./sync` will run the API.

## CORs
You must configure CORs by setting an environment variable `CORS_ALLOWED_ORIGINS` to a comma separated list of origins. For example:

```bash
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://elec.daithiapp.com
```

Please ensure there are no spaces in the list.