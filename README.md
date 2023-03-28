# Conduit Connector for <SPIRE-AIS GraphQL API>
[Conduit](https://conduit.io) for [Spire Martime 2.0 API](https://documentation.spire.com/maritime-2-0/).

## How to build?
Run `make build` to build the connector.

## Testing
Run `make test` to run all the unit tests. Run `make test-integration` to run the integration tests.

The Docker compose file at `test/docker-compose.yml` can be used to run the required resource locally.

## Source
The source connector pulls data from Spire's Maritime 2.0 GraphQL API

### Configuration
A Spire API token is required to use this connector 

| name                  | description                           | required | default value |
|-----------------------|---------------------------------------|----------|---------------|
| `api_url` | Spire API URL to use for accessing the Maritime 2.0 GraphQL API. | false     | https://api.spire.com/graphql          |
| `token` | Access token to use when accessing the Spire GraphQL API. | true     |           |
| `query` | The query to send to the Spire GraphQL API. | true     |           |
| `batch_size` | The maximum number of results to retrieve from the Spire GraphQL API for each request. | false     |     100      |

