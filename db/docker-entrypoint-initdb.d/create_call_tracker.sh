#!/bin/bash

echo "Searching existing user";
userExists=`psql --username postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='backend'"`
if [ -z "$userExists" ]; then
    echo "creating DB and user backend..."
    createuser backend
    createdb --owner=backend call_tracker
    echo "User and DB creation completed successfully"
fi
