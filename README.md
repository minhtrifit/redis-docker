## SOURCE: https://www.youtube.com/watch?v=qucL1F2YEKE

## Step 1: docker-compose -f .\redis-docker-compose.yml up

## Step 2: http://localhost:8001

## Step 3: Choose "I already have a database" -> "Connect to a Redis Database"

## Step 4: Fill form with
`
    Host: Redis
    Port: 6379
    Name: redis-local
`

## Test run redis successfully:
`
    set name minhtrifit
    get name
`