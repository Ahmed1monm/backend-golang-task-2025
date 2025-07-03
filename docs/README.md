# HOW TO RUN THE APP

1. Clone the repository
2. Navigate to the project directory

```bash
docker-compose up --build
```

3. Access the API at <http://localhost:8080>

4. Access the Swagger documentation at <http://localhost:8080/docs/>

# Tech decisions

1. Rate limiting

- I implemented the sliding window algorithm for rate limiting using Redis
2. Performance
- I used Redis for caching [I cached the response og products getters only for simplicity]
- Utilized go routines for concurrent processing in order creation not to block the main thread [notification is an example]
3. real time updates
- I used websockets for real time updates in both inventory updates and notifications
4. Race conditions
- I used transactions to prevent race conditions in inventory updates and orders

> [!CAUTION] 
> I went for super simplecity in the design and implementation in addition to using vipe coding in code writing not implementation itself(design of the code and architecture) and this due to a health issue I had so I believe this is not a production ready code or system design and we can discuss the production ready in the interview. Sorry for the latency and I will try to make it better in the next interview.

