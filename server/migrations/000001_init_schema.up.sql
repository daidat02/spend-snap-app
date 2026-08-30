-- ========================================================
-- SpendSnap DB — Version 1 (khởi tạo toàn bộ schema)
-- Quy ước: mọi PK dùng VARCHAR(36) lưu UUID v4 dạng string
-- (tương thích trực tiếp với kiểu string trong Go / pgx).
-- Timestamp dùng TIMESTAMP WITH TIME ZONE.
-- ========================================================

-- 1. USER & AUTH
CREATE TABLE IF NOT EXISTS users (
    id            VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid()::text,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    firstname     VARCHAR(50),
    lastname      VARCHAR(50),
    username      VARCHAR(50) UNIQUE,
    avatar_url    TEXT,
    phone_number  VARCHAR(20)  UNIQUE,
    bio           VARCHAR(255),
    status        VARCHAR(20)  NOT NULL DEFAULT 'active', -- 'active' / 'inactive'
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at    TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS user_devices (
    id            VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id       VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fcm_token     TEXT        NOT NULL,
    device_type   VARCHAR(20) NOT NULL, -- 'ios' / 'android'
    last_active_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_user_device UNIQUE(user_id, fcm_token)
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id                 VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id            VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL,
    device_id          VARCHAR(100),
    user_agent         TEXT,
    ip_address         VARCHAR(45),
    expires_at         TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at         TIMESTAMP WITH TIME ZONE,
    created_at         TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. SOCIAL & FRIENDSHIPS
CREATE TABLE IF NOT EXISTS friendships (
    id              VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    requester_id    VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    addressee_id    VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending'/'accepted'/'rejected'/'blocked'
    is_close_friend BOOLEAN     DEFAULT FALSE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_friend_request UNIQUE(requester_id, addressee_id),
    CONSTRAINT chk_self_friend CHECK (requester_id <> addressee_id)
);

-- 3. FINANCE CORE
CREATE TABLE IF NOT EXISTS wallets (
    id         VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id    VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       VARCHAR(50) NOT NULL,
    balance    BIGINT      NOT NULL DEFAULT 0 CHECK (balance >= 0), -- VND, đơn vị nhỏ nhất
    currency   VARCHAR(10) DEFAULT 'VND',
    icon_url   TEXT,
    is_default BOOLEAN     DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
    id         VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id    VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE, -- NULL = danh mục hệ thống
    name       VARCHAR(50) NOT NULL,
    type       VARCHAR(10) NOT NULL CHECK (type IN ('expense', 'income')),
    icon       VARCHAR(50),
    color      VARCHAR(20),
    parent_id  VARCHAR(36) REFERENCES categories(id) ON DELETE SET NULL,
    is_system  BOOLEAN     DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id               VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id          VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wallet_id        VARCHAR(36) NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    category_id      VARCHAR(36) NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    amount           BIGINT      NOT NULL CHECK (amount > 0),
    type             VARCHAR(10) NOT NULL DEFAULT 'expense', -- 'expense' / 'income'
    note             TEXT,
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at       TIMESTAMP WITH TIME ZONE
);

-- 4. LOCKET POSTS & INTERACTION
CREATE TABLE IF NOT EXISTS posts (
    id            VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id       VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_url     TEXT        NOT NULL,
    caption       TEXT,
    visibility    VARCHAR(20) NOT NULL DEFAULT 'friends', -- 'private'/'friends'/'close_friends'/'public'
    location_name VARCHAR(255),
    transaction_id VARCHAR(36) REFERENCES transactions(id) ON DELETE SET NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at    TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS post_reactions (
    id         VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    post_id    VARCHAR(36) NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    emoji      VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_user_post_reaction UNIQUE(post_id, user_id, emoji)
);

CREATE TABLE IF NOT EXISTS comments (
    id         VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    post_id    VARCHAR(36) NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id  VARCHAR(36) REFERENCES comments(id) ON DELETE CASCADE, -- reply
    content    TEXT        NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. SPLIT BILL (CHIA TIỀN NHÓM)
CREATE TABLE IF NOT EXISTS groups (
    id         VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    name       VARCHAR(100) NOT NULL,
    created_by VARCHAR(36) NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS group_members (
    group_id  VARCHAR(36) NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id   VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role      VARCHAR(20) DEFAULT 'member', -- 'admin' / 'member'
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE IF NOT EXISTS transaction_splits (
    id             VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    transaction_id VARCHAR(36) NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    user_id        VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- người nợ
    group_id       VARCHAR(36) REFERENCES groups(id) ON DELETE CASCADE,
    split_amount   BIGINT      NOT NULL,
    is_paid        BOOLEAN     DEFAULT FALSE,
    paid_at        TIMESTAMP WITH TIME ZONE,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ========================================================
-- INDEXES
-- ========================================================
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(firstname, lastname);

CREATE INDEX IF NOT EXISTS idx_friendships_lookup ON friendships(requester_id, addressee_id, status);
CREATE INDEX IF NOT EXISTS idx_friendships_addressee ON friendships(addressee_id);

CREATE INDEX IF NOT EXISTS idx_posts_feed ON posts(user_id, visibility, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id, created_at ASC);
CREATE INDEX IF NOT EXISTS idx_reactions_post ON post_reactions(post_id);

CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, transaction_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_transactions_report ON transactions(user_id, category_id, transaction_date DESC);
