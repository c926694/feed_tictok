# Feed Video Frontend

## Stack
- Vue 3 + Vite
- Vue Router
- Axios
- Composition API

## API Base
- API requests use `baseURL = /api` (for your nginx reverse proxy)
- Media URL fields are used directly from backend response (no extra host/path concatenation)

## Run
```bash
npm install
npm run dev
```

Build:
```bash
npm run build
```

## Implemented Pages
- `/login` 登录
- `/register` 注册页占位
- `/feed` 短视频流
- `/upload` 发布页
- `/profile` 个人主页

## Implemented Components
- `VideoFeed`
- `VideoCard`
- `ActionSidebar`
- `BottomNav`
- `CommentDrawer`
- `LoginForm`
- `UploadForm`
- `ProfileHeader`
- `UserVideoGrid`

## API Coverage (from current Postman collection)
- `POST /users/login`
- `GET /users/me`
- `DELETE /users/logout`
- `GET /videos/feed`
- `GET /videos/feed/hot`
- `POST /videos/create`
- `GET /videos/{id}` (currently used as delete action exactly as Postman defines)
- `POST /follows/switchFollow/{userId}`
- `POST /comments`
- `DELETE /comments/{id}`
- `GET /comments/list/{videoId}`
- `GET /videos/swicthLike/{videoId}`
- `GET /likes/comment/switchLike/{commentId}`

## Missing / To Confirm With Backend
1. Register API is incomplete in Postman (`register` request has no URL/method body).  
   - Frontend currently keeps register page UI and explicit error prompt.
2. Upload file API is not provided.  
   - `POST /videos/create` currently uses text fields `title/cover/play` as in Postman.
3. Profile video list API is not provided.  
   - Frontend temporarily filters current feed by `author.id` for "我的视频".
4. Response schema examples are missing in Postman.  
   - Frontend includes tolerant field adapters (`play`/`play_url`, `cover`/`cover_url`, etc.).

