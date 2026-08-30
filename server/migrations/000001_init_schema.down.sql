-- Rollback version 1: xoá toàn bộ bảng theo thứ tự nghịch đảo
DROP TABLE IF EXISTS transaction_splits;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS post_reactions;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS friendships;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS user_devices;
DROP TABLE IF EXISTS users;
