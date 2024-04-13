docker compose up       //to get the docker instance running using the yml file

docker compose up -d    //to run docker instance in a detached mode to clean up our terminal

docker compose ls       //to check for any running docker instance

docker compose stop     //to stop running instances of docker

docker exec -it lenslocked-db-1 /usr/bin/psql -U baloo -d lenslocked  //to execute the psql binary inside 
// docker so that we can write sql queries right inside the terminal. Remember that '-it' is two commands combined.

goose postgres \
"host=localhost port=5432 user=baloo password=bruh dbname=lenslocked sslmode=disable" \
status
// to connect goose to our database, we need to provide a connection string.
// to run a migration, change 'status' to 'up', to do rollback, change 'status' to 'down'