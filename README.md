## Workflow
```mermaid
sequenceDiagram
Client ->> Server: open a connection
Server ->> Client: send a PoW challenge <rand:difficulty:solution_length> (ex. e1rn50N5cOm6EeInw5fbW/ZGkTEwScC6DjOOVhsHITc=:20:8)
Note left of Client: solve a PoW challenge
Client ->> Server: send solution
Note right of Server: verify a solution
Server ->> Client: send <random word-of-wisdom quote>
```

## Chosen algorithm
* bear some resemblance to the ones utilized in numerous cryptocurrencies, including Bitcoin;
* simplicity and comprehensibility are notable. Its implementation only necessitates fundamental procedures such as concatenation and hashing;
* the basis of this algorithm lies in cryptographic hash functions that are intended to withstand different types of attacks. To make it challenging for attackers to generate valid values, the algorithm requires a particular number of zero bytes at the beginning of the hash.

## How to run
### Server
```sh
docker-compose up [--build] server
```

### Client
```sh
docker-compose up [--build] client
```
**Note**: please run `Client` after `Server` have started.