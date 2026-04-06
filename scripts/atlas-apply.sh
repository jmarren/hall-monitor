#!/bin/bash
atlas schema apply -u "postgresql://postgres:nothing-exciting@localhost:5436/postgres" --to "file://internal/db/sql/schema.sql" --dev-url "docker://postgres/18/dev"
