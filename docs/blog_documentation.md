# Blog Module Documentation

## 1. Overview

The Blog module provides functionality for creating, managing, and interacting with long-form content. It supports a full publishing workflow (Draft/Publish), categorization, tagging, and user engagement features (Comments, Reactions, Bookmarks).

## 2. Architecture

The module follows a Clean Architecture pattern, separating concerns into layers:

- **Handler (Interface)**: HTTP request handling and validation.
- **UseCase (Application)**: Orchestration of business logic.
- **Service (Domain)**: Core domain logic and rules.
- **Repository (Infrastructure)**: Data persistence.

### C4 Container Diagram

```mermaid
C4Context
    title Blog System Context

    Person(author, "Author", "Creates and manages content")
    Person(reader, "Reader", "Consumes content")

    System_Boundary(blog_system, "Blog Module") {
        Container(api, "API Handler", "Go/Gin", "Handles HTTP requests")
        Container(usecase, "UseCase", "Go", "Business logic orchestration")
        Container(service, "Domain Service", "Go", "Core domain rules")
        ContainerDb(db, "Database", "PostgreSQL", "Stores blogs, comments, reactions")
    }

    Rel(author, api, "Manage Blogs", "HTTPS/JSON")
    Rel(reader, api, "Read/React", "HTTPS/JSON")
    
    Rel(api, usecase, "Calls")
    Rel(usecase, service, "Uses")
    Rel(service, db, "Persists", "GORM")
```

## 3. Data Model

The data model centers around the `Blog` entity, with relationships to Users, Categories, and Tags.

### Entity Relationship Diagram (ERD)

```mermaid
erDiagram
    User ||--o{ Blog : "authors"
    User ||--o{ BlogReaction : "reacts"
    User ||--o{ Comment : "writes"
    
    Category ||--o{ Blog : "categorizes"
    
    Blog ||--o{ BlogTag : "has"
    Tag ||--o{ BlogTag : "in"
    
    Blog ||--o{ BlogReaction : "receives"
    Blog ||--o{ Comment : "has"
    
    Blog {
        uuid id PK
        uuid author_id FK
        uuid category_id FK
        string title
        string slug
        text content
        enum status "draft, published"
        enum visibility "public, subscribers_only"
        timestamp published_at
        int upvote_count
        int downvote_count
    }

    BlogReaction {
        uuid id PK
        uuid blog_id FK
        uuid user_id FK
        enum type "upvote, downvote"
    }

    Comment {
        uuid id PK
        uuid blog_id FK
        uuid user_id FK
        uuid parent_id FK "Self-referencing for threads"
        text content
    }
    
    Category {
        uuid id PK
        string name
        string slug
    }

    Tag {
        uuid id PK
        string name
        string slug
    }
```

## 4. User Flows

### 4.1 Author Publishing Flow

An author creates a draft, updates it, and finally publishes it to make it visible.

```mermaid
sequenceDiagram
    participant Author
    participant API as Blog Handler
    participant UC as UseCase
    participant DB as Database

    %% Create Draft
    Author->>API: POST /blogs (Draft Content)
    API->>UC: Create(Draft)
    UC->>DB: Insert Blog (Status: Draft)
    DB-->>UC: Blog ID
    UC-->>API: Blog Object
    API-->>Author: 201 Created

    %% Update Content
    Author->>API: PUT /blogs/:id (Refined Content)
    API->>UC: Update(Content)
    UC->>DB: Update Fields
    DB-->>UC: Success
    UC-->>API: Updated Object
    API-->>Author: 200 OK

    %% Publish
    Author->>API: POST /blogs/:id/publish
    API->>UC: Publish(Visibility, Date)
    UC->>DB: Update Status='published', PublishedAt=Now
    DB-->>UC: Success
    UC-->>API: Published Object
    API-->>Author: 200 OK
```

### 4.2 Reader Interaction Flow

A reader views a blog post and reacts to it.

```mermaid
sequenceDiagram
    participant Reader
    participant API as Blog Handler
    participant UC as UseCase
    participant DB as Database

    %% Read Blog
    Reader->>API: GET /blogs/:id
    API->>UC: GetByID(id, viewerID)
    UC->>DB: Fetch Blog
    DB-->>UC: Blog Data
    UC->>DB: Fetch Reaction Counts (Denormalized)
    DB-->>UC: Counts
    UC-->>API: Blog Response
    API-->>Reader: 200 OK (Content)

    %% React (Upvote)
    Reader->>API: POST /blogs/:id/reaction {reaction: "upvote"}
    API->>UC: React(id, user, upvote)
    UC->>DB: Upsert Reaction
    UC->>DB: Increment Upvote Count (Atomic)
    DB-->>UC: Success
    UC-->>API: New Counts
    API-->>Reader: 200 OK
```

## 5. API Reference

### Base URL: `/api/v1`

| Method | Endpoint | Description | Auth Required |
|:-------|:---------|:------------|:--------------|
| `GET` | `/blogs` | List blogs with filters (author, category, status) | No |
| `GET` | `/blogs/feed` | Get personalized feed | Yes |
| `GET` | `/blogs/:id` | Get single blog details | Optional |
| `POST` | `/blogs` | Create a new blog (Draft) | Yes |
| `PUT` | `/blogs/:id` | Update blog content | Yes (Author) |
| `DELETE` | `/blogs/:id` | Soft delete blog | Yes (Author) |
| `POST` | `/blogs/:id/publish` | Publish a blog | Yes (Author) |
| `POST` | `/blogs/:id/unpublish`| Revert to draft | Yes (Author) |
| `POST` | `/blogs/:id/reaction` | Upvote, Downvote, or Remove | Yes |
| `POST` | `/blogs/:id/comments` | Add a comment | Yes |

### Reaction Endpoint
**POST** `/blogs/:id/reaction`

**Payload**
```json
{
  "reaction": "upvote"
}
```
*   **reaction**: `upvote`, `downvote`, or `none`. Use `none` to remove an existing reaction.

### Key Data Structures

#### Blog Response
```json
{
  "id": "uuid",
  "authorId": "uuid",
  "title": "My First Blog",
  "slug": "my-first-blog",
  "content": "...",
  "status": "published",
  "author": {
    "id": "uuid",
    "name": "Jane Doe"
  },
  "upvoteCount": 10,
  "createdAt": "2023-10-01T12:00:00Z"
}
```
