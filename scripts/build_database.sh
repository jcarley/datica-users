#!/bin/sh
PGPASSWORD=password psql -U admin -h db datica_users_dev <<EOF
BEGIN;

\i ./sql/drop_tables.sql

\i ./sql/users.sql

COMMIT;
EOF
