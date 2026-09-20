# Data Modeling & Entity Guidelines

Conventions for designing Domain Entities and GORM models.

## 1. Primary Keys
*   Use **UUID v4** (string representation) as the primary key:
    ```go
    ID string `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
    ```
*   Generate IDs inside the entity constructor using `github.com/google/uuid`:
    ```go
    func NewProduct(name string, price int64) *Product {
        return &Product{
            ID:        uuid.NewString(),
            Name:      name,
            Price:     price,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
    }
    ```

## 2. Timestamps and Soft Delete
*   Every persistable entity should contain standard audit timestamps:
    ```go
    CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
    UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
    ```
*   `DeletedAt` enables GORM soft deletes. Use `json:"-"` so it is not leaked in responses.

## 3. GORM Column Tags
*   Use explicit lowercase `column:name` tags:
    ```go
    Email string `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
    ```
*   Specify appropriate SQL types for portability across PostgreSQL and other RDBMS.
