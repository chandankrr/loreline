CREATE EXTENSION IF NOT EXISTS VECTOR;

CREATE TYPE book_status AS ENUM (
    'processing',
    'ready',
    'failed'
);
CREATE TYPE book_indexing_status AS ENUM (
    'pending',
    'indexing',
    'ready',
    'failed'
);

CREATE TABLE folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT folders_id_user_uniq UNIQUE (id, user_id),
    -- Prevent duplicate folder names for the same user under the same parent
    CONSTRAINT folders_user_parent_name_uniq
    UNIQUE (user_id, parent_id, name),
    -- Parent folder must belong to the same user
    CONSTRAINT folders_parent_owner_fk
        FOREIGN KEY (parent_id, user_id)
        REFERENCES folders (id, user_id)
        ON DELETE CASCADE
);

CREATE INDEX folders_user_idx ON folders(user_id);
CREATE INDEX folders_parent_user_idx ON folders(parent_id, user_id);

CREATE TRIGGER set_updated_at_folders
    BEFORE UPDATE ON folders
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    folder_id UUID,
    title TEXT NOT NULL,
    author TEXT,
    object_key TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    mime_type TEXT NOT NULL DEFAULT 'application/pdf',
    file_size BIGINT NOT NULL,
    page_count INTEGER NOT NULL DEFAULT 1,
    status book_status NOT NULL DEFAULT 'processing',
    indexing_status book_indexing_status NOT NULL DEFAULT 'pending',
    total_chunks INTEGER NOT NULL DEFAULT 0,
    indexed_chunks INTEGER NOT NULL DEFAULT 0,
    indexing_attempts INTEGER NOT NULL DEFAULT 0,
    indexing_error TEXT,
    indexing_updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_page INTEGER NOT NULL DEFAULT 1,
    progress REAL NOT NULL DEFAULT 0,
    error_message TEXT,
    last_opened_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT books_folder_owner_fk
        FOREIGN KEY (folder_id, user_id)
        REFERENCES folders (id, user_id)
        ON DELETE CASCADE
);

CREATE INDEX books_folder_user_idx ON books (folder_id, user_id);
CREATE INDEX books_user_folder_recent_idx ON books (user_id, folder_id, last_opened_at);
CREATE INDEX books_indexing_queue_idx
    ON books (indexing_status, indexing_updated_at)
    WHERE status = 'ready'
      AND indexing_status <> 'ready';

CREATE TRIGGER set_updated_at_books
    BEFORE UPDATE ON books
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE book_chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID NOT NULL,
    user_id UUID NOT NULL,
    page_start INTEGER NOT NULL,
    page_end INTEGER NOT NULL,
    content text NOT NULL,
    search_vector TSVECTOR
        GENERATED ALWAYS AS (
            to_tsvector('english', coalesce(content, ''))
        ) STORED,
    embedding VECTOR(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX book_chunks_book_user_idx ON book_chunks (book_id, user_id);
CREATE INDEX book_chunks_search_idx ON book_chunks USING GIN (search_vector);
CREATE INDEX book_chunks_embedding_idx
    ON book_chunks
    USING HNSW (embedding vector_cosine_ops)
    WHERE embedding IS NOT NULL;

CREATE TRIGGER set_updated_at_book_chunks
    BEFORE UPDATE ON book_chunks
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
