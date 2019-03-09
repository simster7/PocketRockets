# PocketRockets
An open-source poker platform.

All commands below should be run on the project directory (`PocketRockets/`).

## Installing
1. Install [Docker](https://docs.docker.com/install/)

2. Run

```
docker-compose build
```

## Running the server
To run the server:

1. Run

```
docker-compose up
```

2. Visit

```
localhost:3000
```

## Testing
### Backend testing

Run tests:
```
docker run -ti -v `pwd`/backend/:/app/backend/ pr-backend pytest
```

Run type checker:
```
docker run -ti -v `pwd`/backend/:/app/backend/ pr-backend pyre --source-directory ./poker/engine/ check
```

Frontend testing coming soon.

## Contributing
The `master` branch is protected, to submit changes open a pull request. At least one person other than the submitter must approve the request for it to be merged.
