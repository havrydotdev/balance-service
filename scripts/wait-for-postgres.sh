#!/bin/sh
# wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"

function parse_yaml {
    local prefix=$2
    local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
    sed -ne "s|^\($s\):|\1|" \
        -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
    awk -F$fs '{
        indent = length($1)/2;
        vname[indent] = $2;
        for (i in vname) {if (i > indent) {delete vname[i]}}
        if (length($3) > 0) {
            vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
            printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
        }
    }'
}

# will put CONFIG_ in front of each variable finded in config.yml
eval $(parse_yaml /build/configs/config.yml "CONFIG_")

# sleep until db is initialized
until PGPASSWORD=$DB_PASSWORD psql -h "$host" -U "postgres" -c '\q'; do
    >&2 echo "Postgres is unavailable - sleeping"
    sleep 1
done

>&2 echo "Postgres is up - executing command"

host=$CONFIG_pg_host

if [ $COMPOSE = 'true' ] 
    then
    host=$CONFIG_pg_compose_host
fi

# run migrations
migrate -path /build/schema -database postgres://$CONFIG_pg_username:$DB_PASSWORD@$host:$CONFIG_pg_port/$CONFIG_pg_name?sslmode=disable up
>&2 echo "Migrations applied"

# run go service
exec $cmd