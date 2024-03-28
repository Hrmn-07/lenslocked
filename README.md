docker compose up       //to get the docker instance running using the yml file

docker compose up -d    //to run docker instance in a detached mode to clean up our terminal

docker compose ls       //to check for any running docker instance

docker compose stop     //to stop running instances of docker

docker exec -it lenslocked-db-1 /usr/bin/psql -U baloo -d lenslocked  //to execute the psql binary inside 
// docker so that we can write sql queries right inside the terminal. Remember that '-it' is two commands combined.