# Go Web API Design

## 功能概述

将原有Python程序重构为Go Web应用，提供RESTful API供前端调用，去除剪贴板等与操作系统相关的功能。

## API 接口设计

### 1. 获取角色列表
```
GET /api/characters
Response:
[
  {
    "id": "ema",
    "name": "樱羽艾玛",
    "emotionCount": 8,
    "font": "font3.ttf"
  },
  ...
]
```

### 2. 设置当前角色
```
POST /api/characters/current
Request:
{
  "characterId": "sherri"
}

Response:
{
  "success": true,
  "message": "已切换到角色: 橘雪莉"
}
```

### 3. 获取当前角色
```
GET /api/characters/current
Response:
{
  "id": "sherri",
  "name": "橘雪莉",
  "emotionCount": 7,
  "font": "font3.ttf"
}
```

### 4. 生成图片
```
POST /api/generate
Request:
{
  "type": "text",                 // 类型：目前只支持text
  "content": "示例文本内容",       // 文本内容
  "textInput": "输入的文本内容",    // 用户输入的文本（冗余字段，与content相同）
  "characterId": "sherri",        // 角色ID（可选，默认为当前选择的角色）
  "emotionIndex": 1,              // 表情索引（可选，默认随机）
  "backgroundIndex": 1            // 背景索引（可选，默认随机）
}

Response:
{
  "success": true,
  "imageUrl": "/images/result_xxxxxx.png"
}
```

### 5. 获取角色表情列表
```
GET /api/characters/{characterId}/emotions
Response:
[
  {
    "id": 1,
    "name": "表情1"
  },
  ...
]
```