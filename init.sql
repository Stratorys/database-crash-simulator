CREATE SCHEMA briskport;
CREATE USER briskport_user WITH PASSWORD 'admin';
GRANT USAGE ON SCHEMA briskport TO briskport_user;
GRANT CREATE ON SCHEMA briskport to briskport_user;
GRANT ALL ON ALL TABLES IN SCHEMA briskport to briskport_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA briskport GRANT ALL ON TABLES TO briskport_user;
ALTER USER briskport_user CONNECTION LIMIT 1;