# Query Optimization Guidelines

Guidelines for high-performance GORM database interactions.

## 1. Avoid N+1 Queries (Batch Fetching)
*   **NEVER** query related records in a loop (N+1 queries).
*   For cross-module relations, collect foreign IDs in a slice, call `GetByIDs(ctx, ids)` once, and build an in-memory map:
    ```go
    // In Usecase:
    userIDs := make([]string, 0, len(orders))
    for _, o := range orders {
        userIDs = append(userIDs, o.UserID)
    }
    users, err := uc.userRepo.GetByIDs(ctx, unique(userIDs))
    userMap := toMap(users)
    ```

## 2. Mandatory Pagination
*   List endpoints **MUST** use pagination. Never return unbounded result sets (`SELECT * FROM table`).
*   Use `response.ParsePagination(c)` to parse and clamp `page` and `per_page` (default 20, max 100).
*   Perform separate `Count` query for total rows before fetching paged data.

## 3. Explicit Column Selection
*   Avoid `SELECT *`. Explicitly specify columns needed using `.Select("id", "name", "price", ...)` to reduce memory allocation and database I/O.

## 4. Context Propagation
*   Always chain `.WithContext(ctx)` on every GORM operation.
*   If an HTTP client disconnects or times out, the database query will be promptly aborted.
