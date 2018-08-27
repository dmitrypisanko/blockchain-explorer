drop TABLE sonm.transactions;
CREATE TABLE sonm.transactions (
    date Date,
    timestamp UInt64,
    hash String,
    blockNumber UInt64,
    value String,
    gasUsed UInt64,
    gasPrice UInt64,
    nonce UInt64,
    to String,
    from String
) ENGINE = MergeTree(date, hash, 8192);

drop TABLE sonm.blocks;
CREATE TABLE sonm.blocks (
    date Date,
    number UInt64,
    timestamp UInt64,
    hash String,
    parentHash String,
    uncleHash String,
    minedBy String,
    gasUsed UInt64,
    gasLimit UInt64,
    nonce UInt64,
    size Float64,
    transactionsCount UInt64,
    difficulty UInt64,
    extra String
) ENGINE = MergeTree(date, number, 8192);
