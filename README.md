# Food Swipe Recommender

A swipe-based food recommendation app that helps you decide what to eat.

## How It Works

The app shows you food suggestions one at a time. Your swipes teach it what you're in the mood for:

- **Swipe Left** - "Nah, not feeling it" - The app learns to avoid similar foods
- **Swipe Right** - "Yeah, I like this direction!" - The app shows more foods like this
- **Super Swipe** - "This is it!" - You've made your decision, session ends

The more you swipe, the better the suggestions get. It's like having a friend who actually remembers what you said five minutes ago.

---

## Architecture

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────┐
│   Frontend  │  HTTP   │   Backend        │  API    │   OpenAI    │
│  (Next.js)  │ ◄─────► │  (Go + Gin)      │ ◄─────► │  Embeddings │
│             │         │                  │         │             │
│  - Shows UI │         │ - Sessions       │         │ - text-emb- │
│  - Swipes   │         │ - Vector math    │         │   3-small   │
│  - No logic │         │ - Cosine scoring │         └─────────────┘
└─────────────┘         └──────────────────┘
```

**Frontend**: A thin client that just renders and sends swipes. Zero business logic.

**Backend**: The brain. Maintains sessions, updates intent vectors, scores foods using cosine similarity.

**OpenAI**: Generates 1536-dimensional embeddings that capture semantic meaning of each food.

---

## How Vector Embeddings Work

### The Magic of Embeddings

Instead of manually tagging foods, we let OpenAI read the food descriptions and create a "semantic fingerprint" for each one. Foods with similar vibes end up close together in 1536-dimensional space.

**Example:**

```
"Butter Chicken: A rich and creamy North Indian curry..."
   ↓ OpenAI embedding ↓
[0.82, 0.31, 0.15, 0.09, -0.21, 0.44, ...] (1536 numbers)

"Paneer Tikka: Grilled Indian cottage cheese with spices..."
   ↓ OpenAI embedding ↓
[0.78, 0.29, 0.18, 0.12, -0.19, 0.41, ...] (1536 numbers)
```

These vectors are _close_ because both are Indian, both have bold flavors. The AI learned this from millions of texts, no manual tagging needed!

---

## The Recommendation Flow

### 1. Session Starts (Neutral Intent)

```
User's Intent Vector: [0, 0, 0, 0, ...] (all zeros - no preferences yet)
```

We show the first food in the list.

### 2. User Swipes Right on "Butter Chicken"

```go
// Butter Chicken's embedding
food := []float64{0.82, 0.31, 0.15, 0.09, ...}

// Update intent with +0.2 weight (weak positive)
intent = addWeightedVector(intent, food, 0.2)
intent = normalize(intent) // [0.164, 0.062, 0.030, ...]

// Intent now "points" toward Butter Chicken's direction!
```

### 3. User Swipes Left on "Sushi"

```go
// Sushi's embedding (different flavor profile)
food := []float64{0.21, 0.87, 0.05, 0.93, ...}

// Update intent with -0.5 weight (strong negative)
intent = addWeightedVector(intent, food, -0.5)
intent = normalize(intent) // [0.15, -0.38, 0.03, -0.41, ...]

// Intent now AVOIDS sushi direction!
```

### 4. Finding Next Recommendation

For each unseen food, we calculate **cosine similarity** with the intent vector:

```
Paneer Tikka:    similarity = 0.87  ← Very close to intent!
Caesar Salad:    similarity = 0.12  ← Meh
Tom Yum Soup:    similarity = -0.23 ← Opposite direction
```

**Winner: Paneer Tikka** gets shown next! It's the most similar to what you've liked.

---

## The Math Behind It

### Cosine Similarity

Measures the angle between two vectors (not the distance). Ranges from -1 to +1:

```
similarity = (A · B) / (||A|| × ||B||)

Where:
  A · B     = dot product (sum of element-wise multiplication)
  ||A||, ||B|| = magnitudes (vector lengths)
```

**Why cosine and not distance?**

- Direction matters more than magnitude
- "Spicy Indian" and "VERY spicy Indian" should be similar
- Normalized vectors = fair comparisons

### Swipe Weights

```go
const (
    LeftSwipeWeight  = -0.5  // Strong negative signal
    RightSwipeWeight = 0.2   // Weak positive (exploratory)
    SuperSwipeWeight = 1.0   // Strong positive + ends session
)
```

**Why asymmetric weights?**

- Dislikes should be strong (avoid bad suggestions)
- Likes should be gentle (allow exploration)
- Super swipe is definitive (you found it!)

## Example Session

**Session State:**

```
Intent: [0, 0, 0, ...]  (neutral)
Seen: []
```

### Swipe 1: RIGHT on "Butter Chicken"

```
Intent: [0.82×0.2, 0.31×0.2, ...] → normalized → [0.93, 0.35, ...]
Seen: ["Butter Chicken"]
```

### Swipe 2: LEFT on "Sushi Platter"

```
Intent: [0.93, 0.35, ...] + [0.21×-0.5, 0.87×-0.5, ...]
      → [0.78, -0.21, ...] (normalized)
Seen: ["Butter Chicken", "Sushi Platter"]
```

### Next Recommendation

```
Scoring all unseen foods:
  Paneer Tikka:     cos(intent, food) = 0.94  ← WINNER
  Margherita Pizza: cos(intent, food) = 0.31
  Spicy Ramen:      cos(intent, food) = -0.18
```

**Shows: Paneer Tikka** (most aligned with intent!)

### Swipe 3: SUPER on "Paneer Tikka"

```
Intent: Updated strongly toward Paneer Tikka
Session: COMPLETED
Final Choice: "Paneer Tikka"
```

---

## What Makes This Approach Cool

### 1. **Zero Manual Tagging**

Just write food descriptions naturally. OpenAI figures out the semantics.

### 2. **Learns Relationships**

"Masala Dosa" and "Idli Sambar" are close even without "South Indian" tags.

### 3. **Handles Nuance**

"Spicy Indian" and "Mild Italian" can both be recommended based on context, not just tags.

### 4. **Explainable**

Every recommendation is a math operation. No black-box ML models.

### 5. **Stateless**

Each session is independent. No long-term user profiles to manage.

---

## Food Data

50 foods with natural descriptions in `server2/data/food.json`:

```json
{
  "id": "1",
  "name": "Butter Chicken",
  "description": "A rich and creamy North Indian curry made with tender chicken in a tomato-based sauce with butter and spices. Mildly spiced and comforting, it's perfect for those who enjoy hearty, warming dishes."
}
```

At startup, the backend:

1. Loads all 50 foods
2. Sends each description to OpenAI
3. Gets back 1536-dimensional embedding vectors
4. Stores them in memory for fast lookups

---

## Running the App

### Prerequisites

- Go 1.21+ (or latest)
- Node.js 18+
- OpenAI API key
- Docker (optional, for containerized deployment)

### Start Backend (Local)

```bash
cd server2
export OPENAI_API_KEY="your-key-here"
go run main.go
```

Server runs on http://localhost:8000

_Note: On first start, it generates embeddings for all 50 foods (~10 seconds)_

### Start Backend (Docker)

```bash
cd server2

# Build for local testing
docker build -t bite-decide-backend .

# Run with environment variable
docker run -p 8000:8000 -e OPENAI_API_KEY="your-key-here" bite-decide-backend
```


### Start Frontend

```bash
cd client
npm install
npm run dev
```

Frontend runs on http://localhost:3000

---

## API Endpoints

### `POST /session`

Creates a new session with neutral intent vector.

**Response:**

```json
{ "session_id": "abc-123-def" }
```

### `GET /recommendation?session_id=<id>`

Returns the best unseen food for this session.

**Response:**

```json
{
  "name": "Butter Chicken",
  "description": "A rich and creamy North Indian curry..."
}
```

### `POST /swipe`

Processes a swipe action and updates intent vector.

**Request:**

```json
{
  "session_id": "abc-123-def",
  "food_name": "Butter Chicken",
  "action": "right" // "left" | "right" | "super"
}
```

**Response:**

```json
{ "status": "ok" }
```

---


## Testing

Run unit tests:

```bash
cd server2
go test ./...
```

## Project Structure

```
.
├── client/              # Next.js frontend
│   ├── app/
│   │   ├── page.tsx           # Landing page
│   │   ├── swipe/page.tsx     # Main swipe interface
│   │   ├── complete/page.tsx  # Success screen
│   │   └── lib/api.ts         # API client
│   └── ...
│
└── server2/             # Go backend
    ├── Dockerfile             # Docker build
    ├── go.mod                 # Go dependencies
    ├── go.sum                 # Dependency checksums
    ├── main.go                # Entry point + Gin router
    ├── data/
    │   └── food.json          # 50 foods with descriptions
    ├── handlers/
    │   └── handlers.go        # HTTP request handlers
    ├── openai/
    │   └── client.go          # OpenAI embedding API client
    ├── store/
    │   ├── food.go            # Food storage + embeddings
    │   └── session.go         # In-memory session store
    ├── engine/
    │   └── recommender.go     # Cosine similarity logic
    └── models/
        ├── food.go            # Food data structures
        ├── session.go         # Session data structures
        └── swipe.go           # Swipe action types
```

