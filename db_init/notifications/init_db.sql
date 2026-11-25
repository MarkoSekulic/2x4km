CREATE KEYSPACE IF NOT EXISTS notifications_db
  WITH replication = {
    'class': 'NetworkTopologyStrategy',
    'datacenter1': 1
  }
  AND durable_writes = true;

USE notifications_db;

DROP TABLE IF EXISTS notifications_by_user;
DROP TABLE IF EXISTS notifications_by_user_type;
DROP TABLE IF EXISTS notifications_by_user_actor;
DROP TABLE IF EXISTS notifications_by_user_file;

CREATE TABLE notifications_by_user (
  user_id      text,          -- external ID
  notified_at  timeuuid,
  n_type       text,          -- 'shared' | 'added' | 'updated'
  by_user_id   text,          -- external ID
  by_user_name text,
  f_type       text,          -- 'file' | 'folder'
  f_id         text,          -- external ID
  f_name       text,
  added_f_id   text,          -- external ID - only when n_type = 'added'
  added_f_name text,
  PRIMARY KEY ((user_id), notified_at)
) WITH CLUSTERING ORDER BY (notified_at DESC)
  AND default_time_to_live = 7776000;

CREATE TABLE notifications_by_user_type (
  user_id text, n_type text, notified_at timeuuid,
  by_user_id text, by_user_name text, f_type text, f_id text, f_name text,
  added_f_id text, added_f_name text,
  PRIMARY KEY ((user_id, n_type), notified_at)
) WITH CLUSTERING ORDER BY (notified_at DESC)
  AND default_time_to_live = 7776000;

CREATE TABLE notifications_by_user_actor (
  user_id text, by_user_id text, notified_at timeuuid,
  n_type text, by_user_name text, f_type text, f_id text, f_name text,
  added_f_id text, added_f_name text,
  PRIMARY KEY ((user_id, by_user_id), notified_at)
) WITH CLUSTERING ORDER BY (notified_at DESC)
  AND default_time_to_live = 7776000;

CREATE TABLE notifications_by_user_file (
  user_id text, f_id text, notified_at timeuuid,
  n_type text, by_user_id text, by_user_name text, f_type text, f_name text,
  added_f_id text, added_f_name text,
  PRIMARY KEY ((user_id, f_id), notified_at)
) WITH CLUSTERING ORDER BY (notified_at DESC)
  AND default_time_to_live = 7776000;