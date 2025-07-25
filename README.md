HOW TO RUN THIS PROJECT

Step 1. Run CLI `docker-compose up -d`

Step 2. Execute all of queries below to prepare the database
CREATE TABLE merchants (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    balance NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO merchants (name, balance)
VALUES ('Demo', 1000000)

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id INTEGER NOT NULL REFERENCES merchants(id),
    status VARCHAR(25),
    amount NUMERIC(15,2) NOT NULL,
    account_number VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

Step 3. Inside app directory, run CLI `go run main.go main` to run the API

Step 4. Inside app directory, run CLI `go run main.go consumer` to run the Consumer


CURL

curl --location 'localhost:8001/transfer' \
--header 'Content-Type: application/json' \
--data '{
    "merchant_id": 1,
    "amount": 10000,
    "account_number": 1,
    "simulate_success": true
}'

note: if simulate_success value is false, then the transaction status will be failed and will refund merchant balance, while true will change the transaction status to success 



Next Improvement:
1. Database Structure and Implement Indexing
2. Enhance Retry Logic
3. Fix Dockerfile and docker-compose so we can simplify API and Consumer setup
4. API for top up merchant balance
5. Auto create tables and seeding data script